package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

type LogEntry struct {
	Timestamp  string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	RequestID string                 `json:"request_id,omitempty"`
	Method    string                 `json:"method,omitempty"`
	Path      string                 `json:"path,omitempty"`
	Status    int                    `json:"status,omitempty"`
	Latency   string                 `json:"latency,omitempty"`
	ClientIP  string                `json:"client_ip,omitempty"`
	UserID    string                 `json:"user_id,omitempty"`
	Duration  float64               `json:"duration_ms,omitempty"`
	Error     string                 `json:"error,omitempty"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	Stack     string                 `json:"stack,omitempty"`
}

type Logger struct {
	mu       sync.Mutex
	output   io.Writer
	file     *os.File
	level    Level
	jsonMode bool
}

var (
	defaultLogger *Logger
	infoLog      *Logger
	errorLog      *Logger
	warnLog       *Logger
	debugLog      *Logger
	sysLog        *Logger
)

func Init(logDir string, level string, jsonMode bool) {
	defaultLogger = &Logger{
		output:   os.Stdout,
		level:    parseLevel(level),
		jsonMode: jsonMode,
	}

	infoLog = defaultLogger
	errorLog = defaultLogger
	warnLog = defaultLogger
	debugLog = defaultLogger
	sysLog = defaultLogger

	if logDir != "" {
		setupFileOutput(logDir)
	}
}

func setupFileOutput(logDir string) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		defaultLogger.log(ERROR, "Failed to create log directory: "+err.Error(), nil)
		return
	}

	now := time.Now()
	
	infoFilePath := filepath.Join(logDir, fmt.Sprintf("info-%s.log", now.Format("2006-01-02")))
	errorFilePath := filepath.Join(logDir, fmt.Sprintf("error-%s.log", now.Format("2006-01-02")))

	infoFile, err := os.OpenFile(infoFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		defaultLogger.log(ERROR, "Failed to open info log file: "+err.Error(), nil)
		return
	}

	errorFile, err := os.OpenFile(errorFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		defaultLogger.log(ERROR, "Failed to open error log file: "+err.Error(), nil)
		return
	}

	infoLog.file = infoFile
	infoLog.output = io.MultiWriter(os.Stdout, infoFile)

	errorLog.file = errorFile
	errorLog.output = io.MultiWriter(os.Stderr, errorFile)
}

func parseLevel(level string) Level {
	switch level {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn", "warning":
		return WARN
	case "error":
		return ERROR
	case "fatal":
		return FATAL
	default:
		return INFO
	}
}

func (l *Logger) log(level Level, message string, fields map[string]interface{}) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339Nano),
		Level:     level.String(),
		Message:   message,
		Fields:    fields,
	}

	if l.jsonMode {
		data, _ := json.Marshal(entry)
		fmt.Fprintln(l.output, string(data))
	} else {
		formatStr := "[%s] [%s] %s"
		args := []interface{}{entry.Timestamp, entry.Level, entry.Message}
		
		if entry.RequestID != "" {
			formatStr += " [request_id=%s]"
			args = append(args, entry.RequestID)
		}
		if entry.Method != "" {
			formatStr += " [method=%s]"
			args = append(args, entry.Method)
		}
		if entry.Path != "" {
			formatStr += " [path=%s]"
			args = append(args, entry.Path)
		}
		if entry.Status > 0 {
			formatStr += " [status=%d]"
			args = append(args, entry.Status)
		}
		if entry.Latency != "" {
			formatStr += " [latency=%s]"
			args = append(args, entry.Latency)
		}
		if entry.ClientIP != "" {
			formatStr += " [ip=%s]"
			args = append(args, entry.ClientIP)
		}
		if entry.UserID != "" {
			formatStr += " [user_id=%s]"
			args = append(args, entry.UserID)
		}
		if entry.Error != "" {
			formatStr += " [error=%s]"
			args = append(args, entry.Error)
		}
		
		fmt.Fprintf(l.output, formatStr+"\n", args...)
	}

	if l.file != nil && level >= ERROR {
		fmt.Fprintln(l.file, formatEntry(entry))
	}
}

func formatEntry(entry LogEntry) string {
	data, _ := json.Marshal(entry)
	return string(data)
}

func (l *Logger) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		l.file.Close()
	}
}

func Info(args ...interface{}) {
	message := fmt.Sprint(args...)
	infoLog.log(INFO, message, nil)
}

func Infof(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	infoLog.log(INFO, message, nil)
}

func Error(args ...interface{}) {
	message := fmt.Sprint(args...)
	errorLog.log(ERROR, message, nil)
}

func Errorf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	errorLog.log(ERROR, message, nil)
}

func Warn(args ...interface{}) {
	message := fmt.Sprint(args...)
	warnLog.log(WARN, message, nil)
}

func Warnf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	warnLog.log(WARN, message, nil)
}

func Debug(args ...interface{}) {
	message := fmt.Sprint(args...)
	debugLog.log(DEBUG, message, nil)
}

func Debugf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	debugLog.log(DEBUG, message, nil)
}

func Fatal(args ...interface{}) {
	message := fmt.Sprint(args...)
	errorLog.log(FATAL, message, nil)
	os.Exit(1)
}

func Fatalf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	errorLog.log(FATAL, message, nil)
	os.Exit(1)
}

func SysLog(args ...interface{}) {
	message := fmt.Sprint(args...)
	sysLog.log(INFO, "[SYSTEM] "+message, nil)
}

func SysLogf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	sysLog.log(INFO, "[SYSTEM] "+message, nil)
}

func InfoWithContext(c *gin.Context, message string, keyvals ...interface{}) {
	fields := extractFields(keyvals...)
	addContextInfo(c, &fields)
	infoLog.log(INFO, message, fields)
}

func ErrorWithContext(c *gin.Context, message string, keyvals ...interface{}) {
	fields := extractFields(keyvals...)
	addContextInfo(c, &fields)

	if _, ok := getError(fields); ok {
		fields["stack"] = getStack(4)
	}

	errorLog.log(ERROR, message, fields)
}

func WarnWithContext(c *gin.Context, message string, keyvals ...interface{}) {
	fields := extractFields(keyvals...)
	addContextInfo(c, &fields)
	warnLog.log(WARN, message, fields)
}

func DebugWithContext(c *gin.Context, message string, keyvals ...interface{}) {
	fields := extractFields(keyvals...)
	addContextInfo(c, &fields)
	debugLog.log(DEBUG, message, fields)
}

func addContextInfo(c *gin.Context, fields *map[string]interface{}) {
	if c == nil {
		return
	}

	if requestID := c.GetString("request_id"); requestID != "" {
		(*fields)["request_id"] = requestID
	}

	if user, exists := c.Get("user"); exists {
		if userID, ok := user.(string); ok {
			(*fields)["user_id"] = userID
		}
	}

	(*fields)["method"] = c.Request.Method
	(*fields)["path"] = c.Request.URL.Path
	(*fields)["client_ip"] = c.ClientIP()

	if status, exists := c.Get("status"); exists {
		(*fields)["status"] = status
	}

	if latency, exists := c.Get("latency"); exists {
		(*fields)["latency"] = latency
	}
}

func extractFields(keyvals ...interface{}) map[string]interface{} {
	fields := make(map[string]interface{})
	
	for i := 0; i < len(keyvals)-1; i += 2 {
		if key, ok := keyvals[i].(string); ok {
			fields[key] = keyvals[i+1]
		}
	}
	
	return fields
}

func getError(fields map[string]interface{}) (error, bool) {
	for _, v := range fields {
		if err, ok := v.(error); ok {
			return err, true
		}
	}
	return nil, false
}

func getStack(skip int) string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		
		fields := map[string]interface{}{
			"method":      c.Request.Method,
			"path":       path,
			"query":      query,
			"status":     c.Writer.Status(),
			"latency":    latency.String(),
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
			"duration_ms": float64(latency.Nanoseconds()) / 1e6,
		}

		if requestID := c.GetString("request_id"); requestID != "" {
			fields["request_id"] = requestID
		}

		if errs := c.Errors.String(); errs != "" {
			fields["errors"] = errs
		}

		status := c.Writer.Status()
		level := INFO
		if status >= 500 {
			level = ERROR
		} else if status >= 400 {
			level = WARN
		}

		message := fmt.Sprintf("%s %s", c.Request.Method, path)
		if query != "" {
			message += "?" + query
		}
		message += fmt.Sprintf(" %d %s", status, latency.String())

		defaultLogger.log(level, message, fields)
	}
}
