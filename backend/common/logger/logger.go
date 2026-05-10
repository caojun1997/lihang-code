package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/songquanpeng/go-api-starter/common/config"
)

var (
	infoLog  *log.Logger
	errorLog *log.Logger
	warnLog  *log.Logger
	debugLog *log.Logger
	sysLog   *log.Logger

	infoFile  *os.File
	errorFile *os.File
	warnFile  *os.File
	debugFile *os.File
	sysFile   *os.File

	mutex sync.Mutex
)

func SetupLogger() {
	flag := log.LstdFlags | log.Lshortfile

	infoLog = log.New(os.Stdout, "[INFO] ", flag)
	errorLog = log.New(os.Stderr, "[ERROR] ", flag)
	warnLog = log.New(os.Stdout, "[WARN] ", flag)
	debugLog = log.New(os.Stdout, "[DEBUG] ", flag)
	sysLog = log.New(os.Stdout, "[SYS] ", flag)

	if config.LogDir != "" {
		if err := os.MkdirAll(config.LogDir, 0755); err != nil {
			log.Fatalf("Failed to create log directory: %v", err)
		}

		now := time.Now()
		logFileName := filepath.Join(config.LogDir, fmt.Sprintf("info-%s.log", now.Format("2006-01-02")))

		var err error
		infoFile, err = os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Failed to open info log file: %v", err)
		}

		infoLog.SetOutput(infoFile)
		errorLog.SetOutput(infoFile)
		warnLog.SetOutput(infoFile)
		debugLog.SetOutput(infoFile)
		sysLog.SetOutput(infoFile)
	}
}

func CloseLogger() {
	mutex.Lock()
	defer mutex.Unlock()

	if infoFile != nil {
		infoFile.Close()
	}
	if errorFile != nil {
		errorFile.Close()
	}
	if warnFile != nil {
		warnFile.Close()
	}
	if debugFile != nil {
		debugFile.Close()
	}
	if sysFile != nil {
		sysFile.Close()
	}
}

func Info(args ...interface{}) {
	infoLog.Println(args...)
}

func Infof(format string, args ...interface{}) {
	infoLog.Printf(format, args...)
}

func Error(args ...interface{}) {
	errorLog.Println(args...)
}

func Errorf(format string, args ...interface{}) {
	errorLog.Printf(format, args...)
}

func Warn(args ...interface{}) {
	warnLog.Println(args...)
}

func Warnf(format string, args ...interface{}) {
	warnLog.Printf(format, args...)
}

func Debug(args ...interface{}) {
	debugLog.Println(args...)
}

func Debugf(format string, args ...interface{}) {
	debugLog.Printf(format, args...)
}

func SysLog(args ...interface{}) {
	sysLog.Println(args...)
}

func SysLogf(format string, args ...interface{}) {
	sysLog.Printf(format, args...)
}

func FatalLog(args ...interface{}) {
	Error(args...)
	os.Exit(1)
}

func FatalLogf(format string, args ...interface{}) {
	Errorf(format, args...)
	os.Exit(1)
}
