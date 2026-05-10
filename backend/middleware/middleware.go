package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/songquanpeng/go-api-starter/common"
	"github.com/songquanpeng/go-api-starter/common/ctxkey"
	"github.com/songquanpeng/go-api-starter/common/helper"
	"github.com/songquanpeng/go-api-starter/model"
)

func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := c.GetHeader("X-Request-Id")
		if requestId == "" {
			requestId = helper.GetRandomCode(32)
		}
		c.Set(string(ctxkey.KeyRequestID), requestId)
		c.Writer.Header().Set("X-Request-Id", requestId)
		c.Next()
	}
}

func GetRequestId(c *gin.Context) string {
	requestId, exists := c.Get(string(ctxkey.KeyRequestID))
	if !exists {
		return ""
	}
	return requestId.(string)
}

func GetUsername(c *gin.Context) string {
	username, exists := c.Get(string(ctxkey.KeyUsername))
	if !exists {
		return ""
	}
	return username.(string)
}

func GetUserId(c *gin.Context) string {
	userId, exists := c.Get(string(ctxkey.KeyUserId))
	if !exists {
		return ""
	}
	return userId.(string)
}

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			common.Unauthorized(c, "Authorization header required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			common.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		token := parts[1]
		session, err := validateToken(token)
		if err != nil {
			common.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		var user model.User
		if err := common.DB.Where("user_id = ?", session.UserID).First(&user).Error; err != nil {
			common.Unauthorized(c, "User not found")
			c.Abort()
			return
		}

		c.Set(string(ctxkey.KeyUserId), user.UserID)
		c.Set(string(ctxkey.KeyUsername), user.Username)
		c.Set("user", &user)
		c.Set("session", session)
		c.Next()
	}
}

func validateToken(token string) (*model.Session, error) {
	var session model.Session
	err := common.DB.Where("access_token = ?", token).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func SetUpLogger(r *gin.Engine) {
	r.Use(gin.Logger())
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
