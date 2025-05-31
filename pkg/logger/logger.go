package logger

import (
	"github.com/crafty-ezhik/blog-api/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger(cfg *config.Config) error {
	var level zapcore.Level

	switch cfg.Log.Mode {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	logConfig := zap.NewProductionConfig()
	logConfig.Encoding = cfg.Log.Encoding
	logConfig.Level = zap.NewAtomicLevelAt(level)
	logConfig.OutputPaths = cfg.Log.OutputPath

	var err error
	Log, err = logConfig.Build()
	if err != nil {
		return err
	}
	return nil
}
