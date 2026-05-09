package model

import "errors"

var (
	ErrTaskNotFound       = errors.New("task not found")
	ErrResultNotFound     = errors.New("result not found")
	ErrInvalidFileType    = errors.New("invalid file type, only PDF is allowed")
	ErrFileTooLarge       = errors.New("file size exceeds limit")
	ErrInvalidFormat      = errors.New("invalid output format")
	ErrMineruAPIError     = errors.New("mineru API error")
	ErrParseTimeout       = errors.New("parse timeout")
	ErrInvalidTaskStatus  = errors.New("invalid task status")
	ErrUploadFailed       = errors.New("file upload failed")
	ErrTaskInProgress     = errors.New("task is still in progress")
)

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

func (e *APIError) Error() string {
	return e.Message
}

func NewAPIError(code int, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}
