package logger

import (
	std "notification/internal/adapters/logger/std"
	zap "notification/internal/adapters/logger/zap"
	pkgLogger "notification/pkg/logger"
)

// DIP idea on logger...

func Init() pkgLogger.Logger {
	zapLogger, err := zap.NewLogger()

	// bad init zap (fallback)
	if err != nil {
		fallbackLogger := std.NewLogger()
		fallbackLogger.Warnf("init logger, and used default logger %v", err)
		return fallbackLogger
	}

	return zapLogger
}
