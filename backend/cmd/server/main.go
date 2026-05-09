package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pdf-parser/internal/config"
	"pdf-parser/internal/handler"
	"pdf-parser/internal/middleware"
	"pdf-parser/internal/mineru"
	"pdf-parser/internal/oauth"
	"pdf-parser/internal/repository"
	"pdf-parser/internal/service"
	"pdf-parser/pkg/database"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	if err := config.InitConfig("./configs/config.yaml"); err != nil {
		logger.Fatal("Failed to initialize config", zap.Error(err))
	}

	cfg := config.GlobalConfig

	if err := database.InitMySQL(&cfg.Database); err != nil {
		logger.Fatal("Failed to initialize MySQL", zap.Error(err))
	}
	defer database.Close()

	if err := database.InitRedis(&cfg.Redis); err != nil {
		logger.Warn("Failed to initialize Redis, some features may not work", zap.Error(err))
	}
	defer database.CloseRedis()

	if err := database.AutoMigrate(); err != nil {
		logger.Fatal("Failed to migrate database", zap.Error(err))
	}

	taskRepo := repository.NewTaskRepository()
	resultRepo := repository.NewResultRepository()
	cacheRepo := repository.NewCacheRepository()
	userRepo := repository.NewUserRepository()
	sessionRepo := repository.NewSessionRepository()

	_ = userRepo
	_ = sessionRepo

	fileStorage, err := service.NewFileStorage(cfg.Upload.StoragePath)
	if err != nil {
		logger.Fatal("Failed to initialize file storage", zap.Error(err))
	}

	mineruClient := mineru.NewClient(&cfg.Minoru)
	taskService := service.NewTaskService(
		taskRepo,
		resultRepo,
		cacheRepo,
		mineruClient,
		fileStorage,
		cfg,
	)

	oauthClient := oauth.NewOAuthClient(&cfg.OAuth)
	taskHandler := handler.NewTaskHandler(taskService)
	authHandler := handler.NewAuthHandler(oauthClient)
	h := handler.NewHandler(taskHandler, authHandler)

	if cfg.App.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(middleware.Logger(logger))
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())
	r.Use(middleware.RateLimit(cacheRepo, &cfg.RateLimit))

	h.RegisterRoutes(r)

	srv := &http.Server{
		Addr:    cfg.App.Addr(),
		Handler: r,
	}

	go func() {
		logger.Info("Starting server", zap.String("addr", cfg.App.Addr()))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}
