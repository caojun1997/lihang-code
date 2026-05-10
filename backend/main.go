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
	"github.com/songquanpeng/go-api-starter/model"
	"github.com/songquanpeng/go-api-starter/router"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	common.Init()
	logger.SetupLogger()
	logger.SysLogf("PDF Parser %s started", common.Version)

	if os.Getenv("GIN_MODE") != gin.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	if config.DebugEnabled {
		logger.SysLog("running in debug mode")
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
			logger.SysLogf("Config loaded from: %s", path)
			break
		}
	}

	if initErr != nil {
		logger.FatalLog("Failed to initialize config: " + initErr.Error())
	}

	model.InitDB()
	logger.SysLog("Database initialized")

	err := common.InitRedisClient()
	if err != nil {
		logger.FatalLog("Failed to initialize Redis: " + err.Error())
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
	server.Use(gin.Recovery())

	router.SetRouter(server)

	port := os.Getenv("PORT")
	if port == "" {
		if config.ForcePort {
			port = strconv.Itoa(config.Port)
		} else {
			port = strconv.Itoa(config.GlobalConfig.Server.Port)
		}
	}

	logger.SysLogf("Server started on http://0.0.0.0:%s", port)
	logger.SysLog("API文档: http://localhost:" + port + "/api/v1")
	logger.SysLog("前端页面: http://localhost:" + port)
	err = server.Run(":" + port)
	if err != nil {
		logger.FatalLog("Failed to start HTTP server: " + err.Error())
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
