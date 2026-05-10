package common

import (
	"os"
	"path/filepath"

	"github.com/songquanpeng/go-api-starter/common/config"
	"github.com/songquanpeng/go-api-starter/common/logger"
)

func Init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.FatalLog("failed to get home directory: " + err.Error())
	}

	if config.LogDir == "" {
		config.LogDir = filepath.Join(homeDir, "logs")
	}
}
