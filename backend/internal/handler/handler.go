package handler

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
	taskHandler *TaskHandler
}

func NewHandler(taskHandler *TaskHandler) *Handler {
	return &Handler{
		taskHandler: taskHandler,
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
	}
}
