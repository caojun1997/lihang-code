package middleware

import (
	"net/http"
	"strings"

	"pdf-parser/internal/model"
	"pdf-parser/pkg/database"

	"github.com/gin-gonic/gin"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Authorization header required",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := parts[1]
		session, err := validateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Invalid or expired token",
			})
			c.Abort()
			return
		}

		var user model.User
		if err := database.GetDB().Where("user_id = ?", session.UserID).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "User not found",
			})
			c.Abort()
			return
		}

		c.Set("user", &user)
		c.Set("session", session)
		c.Next()
	}
}

func validateToken(token string) (*model.Session, error) {
	var session model.Session
	err := database.GetDB().Where("access_token = ?", token).First(&session).Error
	if err != nil {
		return nil, err
	}

	return &session, nil
}
