package common

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Response struct {
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	RequestID string      `json:"request_id,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Error     *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
	Code      string `json:"code"`
	Detail    string `json:"detail,omitempty"`
	Stack     string `json:"stack,omitempty"`
	Timestamp string `json:"timestamp"`
}

const (
	CodeSuccess           = 0
	CodeUnknownError      = 10000
	CodeInvalidParameters = 10001
	CodeUnauthorized      = 10002
	CodeForbidden        = 10003
	CodeNotFound         = 10004
	CodeInternalError    = 10005
	CodeRateLimited      = 10006
	CodeFileTooLarge     = 10007
	CodeUnsupportedFormat = 10008
	CodeParseFailed      = 10009
	CodeDatabaseError    = 10010
	CodeCacheError       = 10011
	CodeRequestTimeout   = 10012
)

type BizError struct {
	HTTPStatus int
	Code       int
	Message    string
	Detail     string
	Err        error
}

func (e *BizError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func (e *BizError) Unwrap() error {
	return e.Err
}

func NewBizError(httpStatus, code int, message string) *BizError {
	return &BizError{
		HTTPStatus: httpStatus,
		Code:       code,
		Message:    message,
	}
}

func NewBizErrorWithDetail(httpStatus, code int, message, detail string) *BizError {
	return &BizError{
		HTTPStatus: httpStatus,
		Code:       code,
		Message:    message,
		Detail:     detail,
	}
}

func NewBizErrorWrap(err error, httpStatus, code int, message string) *BizError {
	return &BizError{
		HTTPStatus: httpStatus,
		Code:       code,
		Message:    message,
		Err:        err,
	}
}

var (
	ErrInvalidParams     = NewBizError(http.StatusBadRequest, CodeInvalidParameters, "参数错误")
	ErrUnauthorized      = NewBizError(http.StatusUnauthorized, CodeUnauthorized, "未授权")
	ErrForbidden         = NewBizError(http.StatusForbidden, CodeForbidden, "禁止访问")
	ErrNotFound          = NewBizError(http.StatusNotFound, CodeNotFound, "资源不存在")
	ErrInternalServer   = NewBizError(http.StatusInternalServerError, CodeInternalError, "服务器内部错误")
	ErrRateLimited       = NewBizError(http.StatusTooManyRequests, CodeRateLimited, "请求过于频繁")
	ErrFileTooLarge      = NewBizError(http.StatusRequestEntityTooLarge, CodeFileTooLarge, "文件过大")
	ErrUnsupportedFormat = NewBizError(http.StatusBadRequest, CodeUnsupportedFormat, "不支持的文件格式")
	ErrParseFailed       = NewBizError(http.StatusInternalServerError, CodeParseFailed, "解析失败")
	ErrDatabaseError     = NewBizError(http.StatusInternalServerError, CodeDatabaseError, "数据库错误")
	ErrCacheError        = NewBizError(http.StatusInternalServerError, CodeCacheError, "缓存错误")
)

func Success(c *gin.Context, data interface{}) {
	requestID := c.GetString("request_id")
	c.JSON(http.StatusOK, Response{
		Code:      CodeSuccess,
		Msg:       "success",
		RequestID: requestID,
		Data:      data,
	})
}

func SuccessWithMessage(c *gin.Context, msg string) {
	requestID := c.GetString("request_id")
	c.JSON(http.StatusOK, Response{
		Code:      CodeSuccess,
		Msg:       msg,
		RequestID: requestID,
	})
}

func SuccessWithData(c *gin.Context, msg string, data interface{}) {
	requestID := c.GetString("request_id")
	c.JSON(http.StatusOK, Response{
		Code:      CodeSuccess,
		Msg:       msg,
		RequestID: requestID,
		Data:      data,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, code int, msg string) {
	requestID := c.GetString("request_id")
	c.JSON(statusCode, Response{
		Code:      code,
		Msg:       msg,
		RequestID: requestID,
	})
}

func ErrorWithDetail(c *gin.Context, statusCode int, code int, msg string, detail string) {
	requestID := c.GetString("request_id")
	c.JSON(statusCode, Response{
		Code: code,
		Msg:  msg,
		Error: &ErrorInfo{
			Code:      fmt.Sprintf("ERR_%d", code),
			Detail:    detail,
			Timestamp: time.Now().Format(time.RFC3339),
		},
		RequestID: requestID,
	})
}

func ErrorWithStack(c *gin.Context, err *BizError) {
	requestID := c.GetString("request_id")

	c.JSON(err.HTTPStatus, Response{
		Code: err.Code,
		Msg:  err.Message,
		Error: &ErrorInfo{
			Code:      fmt.Sprintf("ERR_%d", err.Code),
			Detail:    err.Detail,
			Timestamp: time.Now().Format(time.RFC3339),
		},
		RequestID: requestID,
	})
}

func ParamError(c *gin.Context, msg string) {
	ErrorResponse(c, http.StatusBadRequest, CodeInvalidParameters, msg)
}

func ParamErrorWithDetail(c *gin.Context, msg, detail string) {
	ErrorWithDetail(c, http.StatusBadRequest, CodeInvalidParameters, msg, detail)
}

func ServerError(c *gin.Context, msg string) {
	ErrorResponse(c, http.StatusInternalServerError, CodeInternalError, msg)
}

func ServerErrorWithStack(c *gin.Context, err error) {
	requestID := c.GetString("request_id")

	c.JSON(http.StatusInternalServerError, Response{
		Code: CodeInternalError,
		Msg:  "服务器内部错误",
		Error: &ErrorInfo{
			Code:      "ERR_10005",
			Detail:    err.Error(),
			Stack:     string(debug.Stack()),
			Timestamp: time.Now().Format(time.RFC3339),
		},
		RequestID: requestID,
	})
}

func NotFound(c *gin.Context, msg string) {
	ErrorResponse(c, http.StatusNotFound, CodeNotFound, msg)
}

func Unauthorized(c *gin.Context, msg string) {
	ErrorResponse(c, http.StatusUnauthorized, CodeUnauthorized, msg)
}

func Forbidden(c *gin.Context, msg string) {
	ErrorResponse(c, http.StatusForbidden, CodeForbidden, msg)
}

func Fail(c *gin.Context, err *BizError) {
	ErrorWithStack(c, err)
}

func FailIfError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	if bizErr, ok := err.(*BizError); ok {
		ErrorWithStack(c, bizErr)
		return
	}

	ServerErrorWithStack(c, err)
}

func getRequestID(c *gin.Context) string {
	requestID := c.GetString("request_id")
	if requestID == "" {
		requestID = uuid.New().String()
	}
	return requestID
}
