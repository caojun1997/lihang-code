package model

import (
	"time"

	"gorm.io/gorm"
)

type ParseTask struct {
	ID           uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskID       string         `gorm:"column:task_id;type:varchar(64);uniqueIndex;not null" json:"task_id"`
	UserID       string         `gorm:"column:user_id;type:varchar(64);index" json:"user_id"`
	FileName     string         `gorm:"column:file_name;type:varchar(256);not null" json:"file_name"`
	FilePath     string         `gorm:"column:file_path;type:varchar(512)" json:"file_path"`
	FileSize     int64          `gorm:"column:file_size" json:"file_size"`
	OutputFormat string         `gorm:"column:output_format;type:varchar(16);not null;default:markdown" json:"output_format"`
	Status       string         `gorm:"column:status;type:varchar(32);not null;default:pending;index" json:"status"`
	Progress     int            `gorm:"column:progress;not null;default:0" json:"progress"`
	ErrorMsg     string         `gorm:"column:error_msg;type:text" json:"error_msg"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (ParseTask) TableName() string {
	return "parse_tasks"
}

type OutputFormat string

const (
	FormatMarkdown OutputFormat = "markdown"
	FormatText     OutputFormat = "txt"
	FormatJSON     OutputFormat = "json"
)

const (
	TaskStatusPending    = "pending"
	TaskStatusUploading  = "uploading"
	TaskStatusProcessing = "processing"
	TaskStatusCompleted  = "completed"
	TaskStatusPartial    = "partial_failed"
	TaskStatusFailed     = "failed"
	TaskStatusCancelled  = "cancelled"
)
