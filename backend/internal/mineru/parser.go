package mineru

import "context"

type Parser interface {
	CreateTask(ctx context.Context, filePath string) (*TaskInfo, error)
	GetTaskStatus(ctx context.Context, taskID string) (*TaskInfo, error)
	GetResult(ctx context.Context, taskID string) (*ParseResult, error)
}
