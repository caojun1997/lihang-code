package mineru

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type FlashAdapter struct {
	baseURL string
	client  *http.Client
	timeout time.Duration
}

func NewFlashAdapter(baseURL string, timeout int) *FlashAdapter {
	return &FlashAdapter{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
		timeout: time.Duration(timeout) * time.Second,
	}
}

type flashTaskResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		TaskID    string `json:"task_id"`
		Status    string `json:"status"`
		ResultURL string `json:"result_url"`
	} `json:"data"`
}

type flashStatusResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		TaskID    string `json:"task_id"`
		Status    string `json:"status"`
		Progress  int    `json:"progress"`
		ResultURL string `json:"result_url"`
		UpdatedAt string `json:"updated_at"`
	} `json:"data"`
}

func (a *FlashAdapter) CreateTask(ctx context.Context, filePath string) (*TaskInfo, error) {
	url := fmt.Sprintf("%s/mineru/api/v1/agent/parse/file", a.baseURL)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	writer.Close()

	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result flashTaskResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("mineru API error: %s", result.Msg)
	}

	return &TaskInfo{
		TaskID: result.Data.TaskID,
		Status: result.Data.Status,
	}, nil
}

func (a *FlashAdapter) GetTaskStatus(ctx context.Context, taskID string) (*TaskInfo, error) {
	url := fmt.Sprintf("%s/mineru/api/v1/agent/task/%s", a.baseURL, taskID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result flashStatusResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("mineru API error: %s", result.Msg)
	}

	updatedAt, _ := time.Parse(time.RFC3339, result.Data.UpdatedAt)

	return &TaskInfo{
		TaskID:    result.Data.TaskID,
		Status:    result.Data.Status,
		Progress:  result.Data.Progress,
		UpdatedAt: updatedAt,
	}, nil
}

func (a *FlashAdapter) GetResult(ctx context.Context, taskID string) (*ParseResult, error) {
	url := fmt.Sprintf("%s/mineru/api/v1/agent/task/%s/result", a.baseURL, taskID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get result")
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return &ParseResult{
		Content:   string(content),
		WordCount: len(string(content)),
		TaskID:    taskID,
	}, nil
}
