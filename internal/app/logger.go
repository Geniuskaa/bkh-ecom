package app

import (
	"bkh-ecom/internal/logger"
	"context"
	"fmt"

	"go.uber.org/zap/zapcore"
)

// InitLogger inits logger by env log level
func InitLogger(ctx context.Context) {
	var levelString string //просто вытащили из конфига
	var level zapcore.Level
	if err := level.Set(levelString); err != nil {
		logger.WarnKV(ctx, logger.Data{
			Msg:    "parse log level  failed, will be set WARN level",
			Error:  err,
			Detail: fmt.Sprintf("log level: %v", levelString),
		})
		level = zapcore.WarnLevel
	}
	logger.SetLevel(level)
	logger.WarnKV(ctx, logger.Data{
		Msg:    "logger set succeeded",
		Detail: fmt.Sprintf("log level: %v", level),
	})
}
