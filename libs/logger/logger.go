package logger

import (
	"go.uber.org/zap"
)

var (
	globalLogger  *zap.Logger
	defaultConfig = Config{
		Strict:           false,
		Production:       false,
		EnableStacktrace: true,
	}
)

type Config struct {
	Out              []string
	Strict           bool
	Production       bool
	EnableStacktrace bool
}

func NewLogger(cfg Config) error {

	builder := zap.NewProductionConfig()
	builder.DisableStacktrace = !cfg.EnableStacktrace
	builder.Development = false
	builder.Encoding = "json"

	builder.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	if cfg.Strict {
		builder.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	}

	if cfg.Production {
		builder.OutputPaths = cfg.Out
	}

	logger, err := builder.Build()
	if err != nil {
		return err
	}
	globalLogger = logger
	return nil
}

func Get() *zap.Logger {
	if globalLogger == nil {
		if err := NewLogger(defaultConfig); err != nil {
			panic(err)
		}
	}
	return globalLogger
}
