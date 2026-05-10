package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/songquanpeng/go-api-starter/common"
	"github.com/songquanpeng/go-api-starter/common/logger"
)

func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID := c.GetString("request_id")
				
				logger.ErrorWithContext(c, "Panic recovered",
					"request_id", requestID,
					"error", err,
				)

				common.ErrorResponse(c, 500, common.CodeInternalError, "服务器内部错误，请稍后重试")
				c.Abort()
			}
		}()
		c.Next()
	}
}

func Logger() gin.HandlerFunc {
	return logger.RequestLogger()
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID")
		c.Header("Access-Control-Expose-Headers", "X-Request-ID")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Rate limiting logic will be handled by common.RateLimiter
		c.Next()
	}
}

func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		done := make(chan struct{})

		go func() {
			c.Next()
			close(done)
		}()

		select {
		case <-done:
			return
		case <-time.After(timeout):
			c.Header("X-Timeout", "true")
			common.ErrorResponse(c, 408, common.CodeRequestTimeout, "请求超时")
			c.Abort()
		}
	}
}

func SetUpLogger(r *gin.Engine) {
	r.Use(Logger())
}

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, userExists := c.Get("user")
		_, sessionExists := c.Get("session")

		if !userExists && !sessionExists {
			requestID := c.GetString("request_id")
			logger.WarnWithContext(c, "Unauthorized access attempt",
				"request_id", requestID,
				"path", c.Request.URL.Path,
			)
			common.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}
		c.Next()
	}
}
