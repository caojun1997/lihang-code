package handler

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}

func Error(c *gin.Context, httpCode int, code int, message string) {
	c.JSON(httpCode, Response{
		Code: code,
		Msg:  message,
	})
}

func ParamError(c *gin.Context, message string) {
	Error(c, 400, 400, message)
}

func NotFound(c *gin.Context, message string) {
	Error(c, 404, 404, message)
}

func ServerError(c *gin.Context, message string) {
	Error(c, 500, 500, message)
}

func SuccessWithMessage(c *gin.Context, message string) {
	c.JSON(200, Response{
		Code: 0,
		Msg:  message,
	})
}
