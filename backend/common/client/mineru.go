package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/songquanpeng/go-api-starter/common/config"
)

type MineruClient struct {
	config     *config.MinoruConfig
	httpClient *http.Client
}

type ParseRequest struct {
	TaskType string `json:"task_type"`
	FileURL  string `json:"file_url"`
}

type ParseResponse struct {
	TaskID string `json:"task_id"`
	Status string `json:"status"`
}

func NewMineruClient(cfg *config.MinoruConfig) *MineruClient {
	return &MineruClient{
		config: cfg,
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Second,
		},
	}
}

func (c *MineruClient) Parse(filePath string, outputFormat string) (string, error) {
	taskType := "pdf_to_markdown"
	if outputFormat == "txt" {
		taskType = "pdf_to_text"
	}

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	var result string
	if c.config.Mode == "flash" {
		result, err = c.parseWithFlash(fileData, taskType)
	} else {
		result, err = c.parseWithPrecision(fileData, taskType)
	}

	if err != nil {
		return "", err
	}

	return result, nil
}

func (c *MineruClient) parseWithFlash(fileData []byte, taskType string) (string, error) {
	if c.config.APIKey == "" {
		return "", fmt.Errorf("API key is required for flash mode")
	}

	endpoint := "/api/v1/flash"
	url := c.config.BaseURL + endpoint

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "document.pdf")
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := part.Write(fileData); err != nil {
		return "", fmt.Errorf("failed to write file data: %w", err)
	}

	writer.WriteField("task_type", taskType)
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if content, ok := result["content"].(string); ok {
		return content, nil
	}

	return string(respBody), nil
}

func (c *MineruClient) parseWithPrecision(fileData []byte, taskType string) (string, error) {
	if c.config.APIKey == "" {
		return "", fmt.Errorf("API key is required for precision mode")
	}

	endpoint := "/api/v1/precision/upload"
	url := c.config.BaseURL + endpoint

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "document.pdf")
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := part.Write(fileData); err != nil {
		return "", fmt.Errorf("failed to write file data: %w", err)
	}

	writer.WriteField("task_type", taskType)
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return string(respBody), nil
}

func (c *MineruClient) DownloadResult(taskID string) (string, error) {
	url := fmt.Sprintf("%s/api/v1/task/%s/result", c.config.BaseURL, taskID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return string(respBody), nil
}

func SaveFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
