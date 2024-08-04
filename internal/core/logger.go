package core

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

func NewEncoder(env Env) zapcore.Encoder {
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
	if env.LogEncoder == "console" {
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
	return zapcore.NewJSONEncoder(encoderConfig)
}

func NewWriteSyncer(env Env) zapcore.WriteSyncer {
	if len(env.LogFileName) > 0 {
		return zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
			zapcore.AddSync(&lumberjack.Logger{
				Filename:   env.LogFileName,
				MaxSize:    500, // megabytes
				MaxBackups: 3,
				MaxAge:     28, // days
			}))
	}
	return zapcore.AddSync(os.Stdout)
}

func NewLogger(
	env Env,
	encoder zapcore.Encoder,
	writeSyncer zapcore.WriteSyncer,
) *zap.Logger {
	level, _ := zapcore.ParseLevel(env.LogLevel)
	core := zapcore.NewCore(encoder, writeSyncer, level)
	return zap.New(core)
}
