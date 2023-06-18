package log

import (
	"errors"
	"github.com/iyurev/tfmirror/pkg/types"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DevLogger = iota
	ProdLogger
)

type LoggerType int32

func NewLogger(loggerType LoggerType, level zapcore.Level, options ...zap.Option) (*zap.Logger, error) {
	var loggerCfg zap.Config
	switch loggerType {
	case DevLogger:
		conf := zap.NewDevelopmentConfig()
		conf.Level = zap.NewAtomicLevelAt(level)
		loggerCfg = conf
	case ProdLogger:
		conf := zap.NewProductionConfig()
		conf.Level = zap.NewAtomicLevelAt(level)
		loggerCfg = conf
	default:
		return nil, errors.New("unsupported logger type")
	}
	return loggerCfg.Build(options...)
}

func FieldProviderSrc(providerSrc string) zap.Field {
	return zap.String("provider_src", providerSrc)
}

func FieldProviderVersion(version string) zap.Field {
	return zap.String("provider_version", version)
}

func FieldPlatform(platform *types.Platform) []zap.Field {
	return []zap.Field{zap.String("os", platform.Os), zap.String("arch", platform.Arch)}
}

func LevelFromString(level string) zapcore.Level {
	switch level {
	case "info":
		return zap.InfoLevel
	case "debug":
		return zap.DebugLevel
	default:
		return zap.InfoLevel
	}
}
