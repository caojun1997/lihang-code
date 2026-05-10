package router

import (
	"github.com/gin-gonic/gin"
	"github.com/songquanpeng/go-api-starter/common/config"
	"github.com/songquanpeng/go-api-starter/controller"
	"github.com/songquanpeng/go-api-starter/middleware"
)

var (
	taskController *controller.TaskController
	authController *controller.AuthController
)

func InitRouter() {
	taskController = controller.NewTaskController()
	authController = controller.NewAuthController()
}

func SetRouter(r *gin.Engine) {
	InitRouter()

	r.Use(middleware.RequestId())
	middleware.SetUpLogger(r)
	middleware.CORS()

	r.GET("/health", taskController.HealthCheck)

	v1 := r.Group("/api/v1")
	{
		tasks := v1.Group("/tasks")
		{
			tasks.POST("", middleware.RequireAuth(), taskController.CreateTask)
			tasks.GET("", taskController.ListTasks)
			tasks.GET("/:task_id", taskController.GetTask)
			tasks.GET("/:task_id/status", taskController.GetTaskStatus)
			tasks.GET("/:task_id/result", taskController.GetTaskResult)
			tasks.GET("/:task_id/download", taskController.DownloadResult)
			tasks.DELETE("/:task_id", middleware.RequireAuth(), taskController.DeleteTask)
		}

		authGroup := v1.Group("/auth")
		{
			authGroup.GET("/url", authController.GetAuthURL)
			authGroup.GET("/callback", authController.HandleCallback)
			authGroup.GET("/me", middleware.RequireAuth(), authController.GetCurrentUser)
			authGroup.POST("/refresh", middleware.RequireAuth(), authController.RefreshToken)
			authGroup.POST("/logout", middleware.RequireAuth(), authController.Logout)
			authGroup.GET("/apikey", middleware.RequireAuth(), authController.GetAPIKey)
		}
	}

	_ = config.Version
}
