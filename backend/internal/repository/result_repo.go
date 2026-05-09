package repository

import (
	"errors"
	"pdf-parser/internal/model"
	"pdf-parser/pkg/database"
)

type ResultRepository struct{}

func NewResultRepository() *ResultRepository {
	return &ResultRepository{}
}

func (r *ResultRepository) Create(result *model.ParseResult) error {
	return database.GetDB().Create(result).Error
}

func (r *ResultRepository) GetByTaskID(taskID string) (*model.ParseResult, error) {
	var result model.ParseResult
	err := database.GetDB().Where("task_id = ?", taskID).First(&result).Error
	if err != nil {
		if errors.Is(err, errors.New("record not found")) {
			return nil, model.ErrResultNotFound
		}
		return nil, err
	}
	return &result, nil
}

func (r *ResultRepository) Update(result *model.ParseResult) error {
	return database.GetDB().Save(result).Error
}

func (r *ResultRepository) Delete(taskID string) error {
	return database.GetDB().Where("task_id = ?", taskID).Delete(&model.ParseResult{}).Error
}
