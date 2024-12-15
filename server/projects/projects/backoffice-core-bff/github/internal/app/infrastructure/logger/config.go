package logger

import (
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewPismo(options ...zap.Option) (*zap.Logger, error) {
	return NewPismoConfig().Build(options...)
}

func NewPismoConfig() zap.Config {
	logLevel := env.GetEnvWithDefaultAsString(env.LogLevel, env.DefaultLogLevel)
	atomLevel := zap.NewAtomicLevel()
	_ = atomLevel.UnmarshalText([]byte(logLevel))

	return zap.Config{
		Level:       atomLevel,
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "Timestamp",
			LevelKey:       "SeverityText",
			NameKey:        "Logger",
			CallerKey:      "Caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "Body",
			StacktraceKey:  "Stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.EpochNanosTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

