package log

import (
	"github.com/golang_fiber_base/internal/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

func NewEncoder(config config.Config) zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,    // Capitalize the log level names
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC timestamp format
		EncodeDuration: zapcore.SecondsDurationEncoder, // Duration in seconds
		EncodeCaller:   zapcore.ShortCallerEncoder,     // Short caller (file and line)
	}
	if config.LogEncoder == "console" {
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
	return zapcore.NewJSONEncoder(encoderConfig)
}

func NewWriteSyncer(config config.Config) zapcore.WriteSyncer {
	if len(config.LogFileName) > 0 {
		return zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
			zapcore.AddSync(&lumberjack.Logger{
				Filename:   config.LogFileName,
				MaxSize:    500, // megabytes
				MaxBackups: 3,
				MaxAge:     28, // days
			}))
	}
	return zapcore.AddSync(os.Stdout)
}

func NewLogger(
	config config.Config,
	encoder zapcore.Encoder,
	writeSyncer zapcore.WriteSyncer,
) *zap.Logger {
	level, _ := zapcore.ParseLevel(config.LogLevel)
	core := zapcore.NewCore(encoder, writeSyncer, level)
	return zap.New(core)
}

func Module() fx.Option {
	return fx.Module(
		"ModuleLog",
		fx.Provide(
			NewLogger,
			NewEncoder,
			NewWriteSyncer,
		),
	)
}
