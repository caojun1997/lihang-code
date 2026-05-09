package handler

import (
	"pdf-parser/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	taskHandler *TaskHandler
	authHandler *AuthHandler
}

func NewHandler(taskHandler *TaskHandler, authHandler *AuthHandler) *Handler {
	return &Handler{
		taskHandler: taskHandler,
		authHandler: authHandler,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", h.taskHandler.HealthCheck)

	v1 := r.Group("/api/v1")
	{
		tasks := v1.Group("/tasks")
		{
			tasks.POST("", h.taskHandler.CreateTask)
			tasks.GET("", h.taskHandler.ListTasks)
			tasks.GET("/:task_id", h.taskHandler.GetTask)
			tasks.GET("/:task_id/status", h.taskHandler.GetTaskStatus)
			tasks.GET("/:task_id/result", h.taskHandler.GetTaskResult)
			tasks.GET("/:task_id/download", h.taskHandler.DownloadResult)
			tasks.DELETE("/:task_id", h.taskHandler.DeleteTask)
		}

		auth := v1.Group("/auth")
		{
			auth.GET("/url", h.authHandler.GetAuthURL)
			auth.GET("/callback", h.authHandler.HandleCallback)
			auth.GET("/me", middleware.RequireAuth(), h.authHandler.GetCurrentUser)
			auth.POST("/refresh", middleware.RequireAuth(), h.authHandler.RefreshToken)
			auth.POST("/logout", middleware.RequireAuth(), h.authHandler.Logout)
			auth.GET("/apikey", middleware.RequireAuth(), h.authHandler.GetAPIKey)
		}
	}
}
