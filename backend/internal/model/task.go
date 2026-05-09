package model

import (
	"time"
)

type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusProcessing TaskStatus = "processing"
	StatusCompleted  TaskStatus = "completed"
	StatusFailed     TaskStatus = "failed"
)

type OutputFormat string

const (
	FormatMarkdown OutputFormat = "markdown"
	FormatTxt      OutputFormat = "txt"
)

type ParseTask struct {
	ID             uint64       `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskID         string       `gorm:"column:task_id;type:varchar(64);uniqueIndex;not null" json:"task_id"`
	FileName       string       `gorm:"column:file_name;type:varchar(255);not null" json:"file_name"`
	FilePath       string       `gorm:"column:file_path;type:varchar(512);not null" json:"file_path"`
	FileSize       int64        `gorm:"column:file_size;not null" json:"file_size"`
	PageCount      int          `gorm:"column:page_count;default:0" json:"page_count"`
	OutputFormat   OutputFormat `gorm:"column:output_format;type:enum('markdown','txt');default:'markdown'" json:"output_format"`
	Status         TaskStatus   `gorm:"column:status;type:enum('pending','processing','completed','failed');default:'pending'" json:"status"`
	MinoruTaskID   string       `gorm:"column:mineru_task_id;type:varchar(128)" json:"mineru_task_id"`
	ErrorMessage   string       `gorm:"column:error_message;type:text" json:"error_message"`
	CreatedAt      time.Time    `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time    `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (ParseTask) TableName() string {
	return "parse_tasks"
}

type ParseResult struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskID    string    `gorm:"column:task_id;type:varchar(64);index;not null" json:"task_id"`
	Content   string    `gorm:"column:content;type:longtext;not null" json:"content"`
	WordCount int       `gorm:"column:word_count;default:0" json:"word_count"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (ParseResult) TableName() string {
	return "parse_results"
}
