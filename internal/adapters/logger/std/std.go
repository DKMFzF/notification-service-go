package std

import (
	"log"
	logger "notification/pkg/logger"
)

// default (or fallback) logger

type stdLogger struct{}

func NewLogger() logger.Logger {
	return &stdLogger{}
}

func (l *stdLogger) Debugf(format string, args ...any) {
	log.Printf("[DEBUG] "+format, args...)
}

func (l *stdLogger) Infof(format string, args ...any) {
	log.Printf("[INFO]: "+format, args...)
}

func (l *stdLogger) Warnf(format string, args ...any) {
	log.Printf("[WARN]: "+format, args...)
}

func (l *stdLogger) Errorf(format string, args ...any) {
	log.Printf("[ERROR]: "+format, args...)
}

func (l *stdLogger) Fatalf(format string, args ...any) {
	log.Printf("[FATAL]: "+format, args...)
}
