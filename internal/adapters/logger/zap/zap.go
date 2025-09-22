package zap

import (
	"fmt"
	"notification/pkg/logger"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	sl *zap.SugaredLogger
}

func NewLogger() (logger.Logger, error) {
	config := zap.NewProductionConfig()

	// options zap
	config.OutputPaths = []string{"stdout", "logs/app.log"}
	config.Level.SetLevel(zapcore.InfoLevel)
	config.Encoding = "json"

	rawLogger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("init zap error: %v", err)
	}

	return &zapLogger{sl: rawLogger.Sugar()}, nil
}

func (l *zapLogger) Infof(format string, args ...any) {
	l.sl.Infof(format, args...)
}

func (l *zapLogger) Debugf(format string, args ...any) {
	l.sl.Debugf(format, args...)
}

func (l *zapLogger) Warnf(format string, args ...any) {
	l.sl.Warnf(format, args...)
}

func (l *zapLogger) Errorf(format string, args ...any) {
	l.sl.Errorf(format, args...)
}

func (l *zapLogger) Fatalf(format string, args ...any) {
	l.sl.Fatalf(format, args...)
}
