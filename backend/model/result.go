package model

import (
	"time"
)

type ParseResult struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskID    string    `gorm:"column:task_id;type:varchar(64);uniqueIndex;not null" json:"task_id"`
	Content   string    `gorm:"column:content;type:longtext;not null" json:"content"`
	FilePath  string    `gorm:"column:file_path;type:varchar(512);not null" json:"file_path"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (ParseResult) TableName() string {
	return "parse_results"
}
