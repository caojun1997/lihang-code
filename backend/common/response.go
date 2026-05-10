package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

const (
	CodeSuccess           = 0
	CodeUnknownError      = 10000
	CodeInvalidParameters = 10001
	CodeUnauthorized      = 10002
	CodeForbidden         = 10003
	CodeNotFound          = 10004
	CodeInternalError     = 10005
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  "success",
		Data: data,
	})
}

func SuccessWithMessage(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  msg,
	})
}

func Error(c *gin.Context, statusCode int, code int, msg string) {
	c.JSON(statusCode, Response{
		Code: code,
		Msg:  msg,
	})
}

func ParamError(c *gin.Context, msg string) {
	Error(c, http.StatusBadRequest, CodeInvalidParameters, msg)
}

func ServerError(c *gin.Context, msg string) {
	Error(c, http.StatusInternalServerError, CodeInternalError, msg)
}

func NotFound(c *gin.Context, msg string) {
	Error(c, http.StatusNotFound, CodeNotFound, msg)
}

func Unauthorized(c *gin.Context, msg string) {
	Error(c, http.StatusUnauthorized, CodeUnauthorized, msg)
}

func Forbidden(c *gin.Context, msg string) {
	Error(c, http.StatusForbidden, CodeForbidden, msg)
}
