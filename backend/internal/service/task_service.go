package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"pdf-parser/internal/config"
	"pdf-parser/internal/mineru"
	"pdf-parser/internal/model"
	"pdf-parser/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TaskService struct {
	taskRepo     *repository.TaskRepository
	resultRepo   *repository.ResultRepository
	cacheRepo    *repository.CacheRepository
	mineruClient *mineru.MineruClient
	storage      *FileStorage
	config       *config.Config
}

func NewTaskService(
	taskRepo *repository.TaskRepository,
	resultRepo *repository.ResultRepository,
	cacheRepo *repository.CacheRepository,
	mineruClient *mineru.MineruClient,
	storage *FileStorage,
	cfg *config.Config,
) *TaskService {
	return &TaskService{
		taskRepo:     taskRepo,
		resultRepo:   resultRepo,
		cacheRepo:    cacheRepo,
		mineruClient: mineruClient,
		storage:      storage,
		config:       cfg,
	}
}

type CreateTaskResult struct {
	TaskID    string `json:"task_id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

func (s *TaskService) CreateTask(ctx context.Context, file *multipart.FileHeader, outputFormat model.OutputFormat) (*CreateTaskResult, error) {
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".pdf") {
		return nil, model.ErrInvalidFileType
	}

	if file.Size > s.config.Upload.MaxSize {
		return nil, model.ErrFileTooLarge
	}

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	taskID := fmt.Sprintf("task_%s", uuid.New().String()[:8])
	filePath, err := s.storage.Save(src, taskID, file.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	task := &model.ParseTask{
		TaskID:       taskID,
		FileName:     file.Filename,
		FilePath:     filePath,
		FileSize:     file.Size,
		OutputFormat: outputFormat,
		Status:       model.StatusPending,
	}

	if err := s.taskRepo.Create(task); err != nil {
		s.storage.Delete(filePath)
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	go s.processTask(taskID, filePath)

	return &CreateTaskResult{
		TaskID:    taskID,
		Status:    string(model.StatusPending),
		CreatedAt: task.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *TaskService) processTask(taskID, filePath string) {
	ctx := context.Background()

	task, err := s.taskRepo.GetByTaskID(ctx, taskID)
	if err != nil {
		return
	}

	s.taskRepo.UpdateStatus(taskID, model.StatusProcessing, "")
	s.cacheRepo.SetTaskStatus(ctx, taskID, string(model.StatusProcessing))

	taskInfo, err := s.mineruClient.CreateTask(ctx, filePath)
	if err != nil {
		s.taskRepo.UpdateStatus(taskID, model.StatusFailed, "")
		s.cacheRepo.SetTaskStatus(ctx, taskID, string(model.StatusFailed))
		task.ErrorMessage = err.Error()
		s.taskRepo.Update(task)
		return
	}

	s.taskRepo.UpdateStatus(taskID, model.StatusProcessing, taskInfo.TaskID)

	maxAttempts := 60
	for i := 0; i < maxAttempts; i++ {
		time.Sleep(5 * time.Second)

		status, err := s.mineruClient.GetTaskStatus(ctx, taskInfo.TaskID)
		if err != nil {
			continue
		}

		s.cacheRepo.SetTaskProgress(ctx, taskID, status.Progress)

		if status.Status == "completed" {
			result, err := s.mineruClient.GetResult(ctx, taskInfo.TaskID)
			if err != nil {
				s.taskRepo.UpdateStatus(taskID, model.StatusFailed, "")
				s.cacheRepo.SetTaskStatus(ctx, taskID, string(model.StatusFailed))
				task.ErrorMessage = err.Error()
				s.taskRepo.Update(task)
				return
			}

			parseResult := &model.ParseResult{
				TaskID:    taskID,
				Content:   result.Content,
				WordCount: result.WordCount,
			}

			s.resultRepo.Create(parseResult)

			s.taskRepo.UpdateStatus(taskID, model.StatusCompleted, "")
			s.cacheRepo.SetTaskStatus(ctx, taskID, string(model.StatusCompleted))
			s.cacheRepo.SetTaskProgress(ctx, taskID, 100)
			return
		}

		if status.Status == "failed" {
			s.taskRepo.UpdateStatus(taskID, model.StatusFailed, "")
			s.cacheRepo.SetTaskStatus(ctx, taskID, string(model.StatusFailed))
			task.ErrorMessage = "MinerU解析失败"
			s.taskRepo.Update(task)
			return
		}
	}

	s.taskRepo.UpdateStatus(taskID, model.StatusFailed, "")
	s.cacheRepo.SetTaskStatus(ctx, taskID, string(model.StatusFailed))
	task.ErrorMessage = "解析超时"
	s.taskRepo.Update(task)
}

type TaskListResult struct {
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
	List     []model.ParseTask `json:"list"`
}

func (s *TaskService) ListTasks(ctx context.Context, page, pageSize int, status string) (*TaskListResult, error) {
	tasks, total, err := s.taskRepo.List(page, pageSize, status)
	if err != nil {
		return nil, err
	}

	return &TaskListResult{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		List:     tasks,
	}, nil
}

func (s *TaskService) GetTask(ctx context.Context, taskID string) (*model.ParseTask, error) {
	return s.taskRepo.GetByTaskID(ctx, taskID)
}

type TaskStatusResult struct {
	TaskID    string `json:"task_id"`
	Status    string `json:"status"`
	Progress  int    `json:"progress"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

func (s *TaskService) GetTaskStatus(ctx context.Context, taskID string) (*TaskStatusResult, error) {
	task, err := s.taskRepo.GetByTaskID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	progress, _ := s.cacheRepo.GetTaskProgress(ctx, taskID)
	if progress == 0 {
		switch task.Status {
		case model.StatusPending:
			progress = 10
		case model.StatusProcessing:
			progress = 50
		case model.StatusCompleted:
			progress = 100
		case model.StatusFailed:
			progress = 0
		}
	}

	return &TaskStatusResult{
		TaskID:    task.TaskID,
		Status:    string(task.Status),
		Progress:  progress,
		UpdatedAt: task.UpdatedAt.Format(time.RFC3339),
	}, nil
}

type TaskResultResponse struct {
	TaskID    string `json:"task_id"`
	Content   string `json:"content"`
	WordCount int    `json:"word_count"`
}

func (s *TaskService) GetTaskResult(ctx context.Context, taskID string) (*TaskResultResponse, error) {
	task, err := s.taskRepo.GetByTaskID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	if task.Status != model.StatusCompleted {
		return nil, model.ErrTaskInProgress
	}

	result, err := s.resultRepo.GetByTaskID(taskID)
	if err != nil {
		return nil, err
	}

	return &TaskResultResponse{
		TaskID:    taskID,
		Content:   result.Content,
		WordCount: result.WordCount,
	}, nil
}

func (s *TaskService) DownloadResult(ctx context.Context, taskID string, c *gin.Context) error {
	task, err := s.taskRepo.GetByTaskID(ctx, taskID)
	if err != nil {
		return err
	}

	if task.Status != model.StatusCompleted {
		return model.ErrTaskInProgress
	}

	result, err := s.resultRepo.GetByTaskID(taskID)
	if err != nil {
		return err
	}

	ext := ".md"
	if task.OutputFormat == model.FormatTxt {
		ext = ".txt"
	}
	filename := task.FileName[:len(task.FileName)-4] + "_parsed" + ext

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", filename))
	c.Header("Content-Type", "text/plain; charset=utf-8")

	content := result.Content
	if task.OutputFormat == model.FormatTxt {
		content = stripMarkdown(result.Content)
	}

	c.String(200, content)
	return nil
}

func (s *TaskService) DeleteTask(ctx context.Context, taskID string) error {
	task, err := s.taskRepo.GetByTaskID(ctx, taskID)
	if err != nil {
		return err
	}

	s.storage.Delete(task.FilePath)

	s.resultRepo.Delete(taskID)
	s.cacheRepo.DeleteTaskCache(ctx, taskID)

	return s.taskRepo.Delete(taskID)
}

func stripMarkdown(content string) string {
	lines := strings.Split(content, "\n")
	var result []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "```") {
			continue
		}
		if strings.HasPrefix(line, "|") {
			continue
		}
		if strings.HasPrefix(line, "![]") {
			continue
		}
		line = strings.ReplaceAll(line, "**", "")
		line = strings.ReplaceAll(line, "__", "")
		line = strings.ReplaceAll(line, "*", "")
		line = strings.ReplaceAll(line, "_", "")
		line = strings.ReplaceAll(line, "[", "")
		line = strings.ReplaceAll(line, "]", "")
		if strings.HasPrefix(line, "(") && strings.HasSuffix(line, ")") {
			continue
		}
		if line != "" {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}
