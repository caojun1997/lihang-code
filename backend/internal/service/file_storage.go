package service

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type FileStorage struct {
	basePath string
}

func NewFileStorage(basePath string) (*FileStorage, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}
	return &FileStorage{basePath: basePath}, nil
}

func (s *FileStorage) Save(reader io.Reader, taskID, filename string) (string, error) {
	ext := filepath.Ext(filename)
	newFilename := fmt.Sprintf("%s%s", taskID, ext)
	filePath := filepath.Join(s.basePath, newFilename)

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return filePath, nil
}

func (s *FileStorage) Read(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

func (s *FileStorage) Delete(filePath string) error {
	if _, err := os.Stat(filePath); err == nil {
		return os.Remove(filePath)
	}
	return nil
}

func (s *FileStorage) Exists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func (s *FileStorage) GetBasePath() string {
	return s.basePath
}
