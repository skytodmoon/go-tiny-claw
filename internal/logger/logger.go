package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

var (
	levelNames = map[LogLevel]string{
		DEBUG: "DEBUG",
		INFO:  "INFO",
		WARN:  "WARN",
		ERROR: "ERROR",
		FATAL: "FATAL",
	}

	levelColors = map[LogLevel]string{
		DEBUG: "\033[36m",
		INFO:  "\033[32m",
		WARN:  "\033[33m",
		ERROR: "\033[31m",
		FATAL: "\033[35m",
	}

	resetColor = "\033[0m"
)

type Logger struct {
	mu       sync.Mutex
	level    LogLevel
	module   string
	outputs  []io.Writer
	loggers  []*log.Logger
	useColor bool
}

var (
	globalLogger *Logger
	once         sync.Once
)

func Init(level LogLevel, logDir string, useColor bool) error {
	var err error
	once.Do(func() {
		globalLogger, err = NewLogger(level, "main", logDir, useColor)
	})
	return err
}

func NewLogger(level LogLevel, module string, logDir string, useColor bool) (*Logger, error) {
	l := &Logger{
		level:    level,
		module:   module,
		useColor: useColor,
	}

	l.outputs = append(l.outputs, os.Stdout)

	if logDir != "" {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		logFile := filepath.Join(logDir, fmt.Sprintf("claw-%s.log", time.Now().Format("2006-01-02")))
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		l.outputs = append(l.outputs, file)
	}

	for _, output := range l.outputs {
		l.loggers = append(l.loggers, log.New(output, "", 0))
	}

	return l, nil
}

func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	levelName := levelNames[level]
	message := fmt.Sprintf(format, args...)

	// Always use plain text for log files
	plainLine := fmt.Sprintf("%s [%s] [%s] %s", timestamp, levelName, l.module, message)

	for i, logger := range l.loggers {
		// Apply color only for stdout (first output) when useColor is enabled
		if l.useColor && i == 0 {
			color := levelColors[level]
			coloredLine := fmt.Sprintf("%s [%s%s%s] [%s] %s%s",
				timestamp, color, levelName, resetColor, l.module, message, resetColor)
			logger.Println(coloredLine)
		} else {
			logger.Println(plainLine)
		}
	}

	if level == FATAL {
		os.Exit(1)
	}
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
}

func (l *Logger) WithModule(module string) *Logger {
	return &Logger{
		level:    l.level,
		module:   module,
		outputs:  l.outputs,
		loggers:  l.loggers,
		useColor: l.useColor,
	}
}

func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

func GetLogger() *Logger {
	if globalLogger == nil {
		Init(INFO, "", true)
	}
	return globalLogger
}

func WithModule(module string) *Logger {
	return GetLogger().WithModule(module)
}

func SetLevel(level LogLevel) {
	GetLogger().SetLevel(level)
}

func Debug(format string, args ...interface{}) {
	GetLogger().Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	GetLogger().Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	GetLogger().Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	GetLogger().Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
	GetLogger().Fatal(format, args...)
}

func ParseLevel(levelStr string) (LogLevel, error) {
	switch levelStr {
	case "debug", "DEBUG":
		return DEBUG, nil
	case "info", "INFO":
		return INFO, nil
	case "warn", "WARN":
		return WARN, nil
	case "error", "ERROR":
		return ERROR, nil
	case "fatal", "FATAL":
		return FATAL, nil
	default:
		return INFO, fmt.Errorf("unknown log level: %s", levelStr)
	}
}
