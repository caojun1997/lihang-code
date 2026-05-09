package mineru

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type PrecisionAdapter struct {
	baseURL string
	apiKey  string
	client  *http.Client
	timeout time.Duration
}

func NewPrecisionAdapter(baseURL, apiKey string, timeout int) *PrecisionAdapter {
	return &PrecisionAdapter{
		baseURL: baseURL,
		apiKey:  apiKey,
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
		timeout: time.Duration(timeout) * time.Second,
	}
}

type precisionStatusResponse struct {
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

func (a *PrecisionAdapter) CreateTask(ctx context.Context, filePath string) (*TaskInfo, error) {
	return nil, fmt.Errorf("precision mode not implemented in this version")
}

func (a *PrecisionAdapter) GetTaskStatus(ctx context.Context, taskID string) (*TaskInfo, error) {
	url := fmt.Sprintf("%s/mineru/api/v4/extract/task/%s", a.baseURL, taskID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.apiKey))

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result precisionStatusResponse
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

func (a *PrecisionAdapter) GetResult(ctx context.Context, taskID string) (*ParseResult, error) {
	return nil, fmt.Errorf("precision mode not implemented in this version")
}
