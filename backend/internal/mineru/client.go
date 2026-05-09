package mineru

import (
	"context"
	"time"

	"pdf-parser/internal/config"
)

type MineruClient struct {
	parser Parser
}

type TaskInfo struct {
	TaskID    string
	Status    string
	Progress  int
	UpdatedAt time.Time
}

type ParseResult struct {
	Content   string
	WordCount int
	TaskID    string
}

func NewClient(cfg *config.MinoruConfig) *MineruClient {
	var parser Parser

	if cfg.Mode == "flash" {
		parser = NewFlashAdapter(cfg.BaseURL, cfg.Timeout)
	} else {
		parser = NewPrecisionAdapter(cfg.BaseURL, cfg.APIKey, cfg.Timeout)
	}

	return &MineruClient{
		parser: parser,
	}
}

func (c *MineruClient) CreateTask(ctx context.Context, filePath string) (*TaskInfo, error) {
	return c.parser.CreateTask(ctx, filePath)
}

func (c *MineruClient) GetTaskStatus(ctx context.Context, taskID string) (*TaskInfo, error) {
	return c.parser.GetTaskStatus(ctx, taskID)
}

func (c *MineruClient) GetResult(ctx context.Context, taskID string) (*ParseResult, error) {
	return c.parser.GetResult(ctx, taskID)
}
