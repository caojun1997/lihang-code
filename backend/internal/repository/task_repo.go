package repository

import (
	"context"
	"errors"
	"pdf-parser/internal/model"
	"pdf-parser/pkg/database"
)

type TaskRepository struct{}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{}
}

func (r *TaskRepository) Create(task *model.ParseTask) error {
	return database.GetDB().Create(task).Error
}

func (r *TaskRepository) GetByTaskID(ctx context.Context, taskID string) (*model.ParseTask, error) {
	var task model.ParseTask
	err := database.GetDB().Where("task_id = ?", taskID).First(&task).Error
	if err != nil {
		if errors.Is(err, errors.New("record not found")) {
			return nil, model.ErrTaskNotFound
		}
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) Update(task *model.ParseTask) error {
	return database.GetDB().Save(task).Error
}

func (r *TaskRepository) Delete(taskID string) error {
	return database.GetDB().Where("task_id = ?", taskID).Delete(&model.ParseTask{}).Error
}

func (r *TaskRepository) List(page, pageSize int, status string) ([]model.ParseTask, int64, error) {
	var tasks []model.ParseTask
	var total int64

	query := database.GetDB().Model(&model.ParseTask{})

	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&tasks).Error
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *TaskRepository) UpdateStatus(taskID string, status model.TaskStatus, mineruTaskID string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if mineruTaskID != "" {
		updates["mineru_task_id"] = mineruTaskID
	}
	return database.GetDB().Model(&model.ParseTask{}).Where("task_id = ?", taskID).Updates(updates).Error
}

func (r *TaskRepository) GetPendingTasks(limit int) ([]model.ParseTask, error) {
	var tasks []model.ParseTask
	err := database.GetDB().Where("status = ?", model.StatusPending).Order("created_at ASC").Limit(limit).Find(&tasks).Error
	return tasks, err
}
