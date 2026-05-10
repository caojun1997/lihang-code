package main

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/songquanpeng/go-api-starter/common"
	"github.com/songquanpeng/go-api-starter/common/client"
	"github.com/songquanpeng/go-api-starter/common/config"
	"github.com/songquanpeng/go-api-starter/common/logger"
	"github.com/songquanpeng/go-api-starter/controller"
	"github.com/songquanpeng/go-api-starter/model"
	"github.com/songquanpeng/go-api-starter/router"
)

func main() {
	common.Init()
	logger.SetupLogger()
	logger.SysLogf("PDF Parser %s started", common.Version)

	if os.Getenv("GIN_MODE") != gin.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	if config.DebugEnabled {
		logger.SysLog("running in debug mode")
	}

	if err := config.InitConfig("./configs/config.yaml"); err != nil {
		logger.FatalLog("Failed to initialize config: " + err.Error())
	}

	model.InitDB()
	logger.SysLog("Database initialized")

	var err error
	err = common.InitRedisClient()
	if err != nil {
		logger.FatalLog("Failed to initialize Redis: " + err.Error())
	}

	common.RateLimiter.Init(60)

	controller.InitOAuth(&client.OAuthConfig{
		ClientID:    config.GlobalConfig.OAuth.ClientID,
		ClientSecret: config.GlobalConfig.OAuth.ClientSecret,
		AuthURL:     config.GlobalConfig.OAuth.AuthURL,
		TokenURL:    config.GlobalConfig.OAuth.TokenURL,
		UserInfoURL: config.GlobalConfig.OAuth.UserInfoURL,
		RedirectURI: config.GlobalConfig.OAuth.RedirectURI,
		APIKeyURL:   config.GlobalConfig.OAuth.APIKeyURL,
		LogoutURL:   config.GlobalConfig.OAuth.LogoutURL,
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

	logger.SysLogf("Server started on http://localhost:%s", port)
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
