package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/songquanpeng/go-api-starter/common"
	"github.com/songquanpeng/go-api-starter/common/client"
	"github.com/songquanpeng/go-api-starter/common/config"
	"github.com/songquanpeng/go-api-starter/common/logger"
	"github.com/songquanpeng/go-api-starter/controller"
	"github.com/songquanpeng/go-api-starter/middleware"
	"github.com/songquanpeng/go-api-starter/model"
	"github.com/songquanpeng/go-api-starter/router"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	if os.Getenv("GIN_MODE") != gin.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	configPaths := []string{
		"./configs/config.yaml",
		"./config.yaml",
		"/etc/pdf-parser/config.yaml",
	}

	var initErr error
	for _, path := range configPaths {
		initErr = config.InitConfig(path)
		if initErr == nil {
			break
		}
	}

	if initErr != nil {
		log.Fatalf("Failed to initialize config: %v", initErr)
	}

	logDir := config.LogDir
	if logDir == "" {
		logDir = "./logs"
	}
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	logger.Init(logDir, logLevel, false)
	logger.SysLogf("PDF Parser started, version: %s", common.Version)

	if config.DebugEnabled {
		logger.SysLog("running in debug mode")
	}

	model.InitDB()
	logger.Info("Database initialized")

	err := common.InitRedisClient()
	if err != nil {
		logger.Fatal("Failed to initialize Redis: " + err.Error())
	}

	common.RateLimiter.Init(60)

	controller.InitOAuth(&client.OAuthConfig{
		ClientID:     config.GlobalConfig.OAuth.ClientID,
		ClientSecret: config.GlobalConfig.OAuth.ClientSecret,
		AuthURL:      config.GlobalConfig.OAuth.AuthURL,
		TokenURL:     config.GlobalConfig.OAuth.TokenURL,
		UserInfoURL:  config.GlobalConfig.OAuth.UserInfoURL,
		RedirectURI:  config.GlobalConfig.OAuth.RedirectURI,
		APIKeyURL:    config.GlobalConfig.OAuth.APIKeyURL,
		LogoutURL:    config.GlobalConfig.OAuth.LogoutURL,
		WellKnownURL: config.GlobalConfig.OAuth.WellKnownURL,
	})

	server := gin.New()
	server.Use(middleware.Recovery())
	server.Use(middleware.RequestId())
	server.Use(middleware.CORS())
	middleware.SetUpLogger(server)

	router.SetRouter(server)

	port := os.Getenv("PORT")
	if port == "" {
		if config.ForcePort {
			port = strconv.Itoa(config.Port)
		} else {
			port = strconv.Itoa(config.GlobalConfig.Server.Port)
		}
	}

	logger.SysLogf("Server starting on http://0.0.0.0:%s", port)
	logger.SysLogf("API docs: http://localhost:%s/api/v1", port)
	logger.SysLogf("Frontend: http://localhost:%s", port)

	err = server.Run(":" + port)
	if err != nil {
		logger.Fatal("Failed to start HTTP server: " + err.Error())
	}

	defer func() {
		err := model.CloseDB()
		if err != nil {
			logger.Error("Failed to close database:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")
}
