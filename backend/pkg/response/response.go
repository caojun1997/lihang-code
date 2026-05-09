package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}

func SuccessWithMessage(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
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
	Error(c, http.StatusBadRequest, 400, msg)
}

func ServerError(c *gin.Context, msg string) {
	Error(c, http.StatusInternalServerError, 500, msg)
}

func NotFound(c *gin.Context, msg string) {
	Error(c, http.StatusNotFound, 404, msg)
}

func Unauthorized(c *gin.Context, msg string) {
	Error(c, http.StatusUnauthorized, 401, msg)
}
