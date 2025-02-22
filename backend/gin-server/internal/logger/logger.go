package logger

import (
	"log"

	"gin-server/internal/utils"
)

type Logger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
}

type SimpleLogger struct {
}

func NewSimpleLogger() Logger {
	return &SimpleLogger{}
}

func (l *SimpleLogger) Info(args ...interface{}) {
	log.Println("[INFO]", args)
}

func (l *SimpleLogger) Infof(format string, args ...interface{}) {
	log.Printf("[INFO] "+format, args...)
}

func (l *SimpleLogger) Warn(args ...interface{}) {
	log.Println("[WARN]", args)
}

func (l *SimpleLogger) Warnf(format string, args ...interface{}) {
	log.Printf("[WARN] "+format, args...)
}

func (l *SimpleLogger) Error(args ...interface{}) {
	log.Println("[ERROR]", args)
}

func (l *SimpleLogger) Errorf(format string, args ...interface{}) {
	log.Printf("[ERROR] "+format, args...)
}

func (l *SimpleLogger) Debug(args ...interface{}) {
	if utils.IsDebugEnabled() {
		log.Println("[DEBUG]", args)
	}
}

func (l *SimpleLogger) Debugf(format string, args ...interface{}) {
	if utils.IsDebugEnabled() {
		log.Printf("[DEBUG] "+format, args...)
	}
}
