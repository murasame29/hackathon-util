package logger

import (
	"context"
	"reflect"

	"github.com/murasame29/hackathon-util/cmd/config"
	"go.uber.org/zap"
)

type LoggerKey struct{}

func NewLogger() *zap.Logger {
	var logger *zap.Logger

	switch config.Config.Application.Env {
	case config.Dev:
		logger, _ = zap.NewDevelopment(zap.WithCaller(false))
	case config.Prod:
		logger, _ = zap.NewProduction(zap.WithCaller(false))
	default:
		logger, _ = zap.NewDevelopment(zap.WithCaller(false))
	}

	return logger
}

func NewLoggerWithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, LoggerKey{}, NewLogger())
}

func FromContext(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(LoggerKey{}).(*zap.Logger)
	if !ok {
		return NewLogger()
	}

	return logger
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx).Info(msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx).Warn(msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx).Error(msg, fields...)
}

func Field(key string, val any) zap.Field {
	switch reflect.TypeOf(val).Kind() {
	case reflect.Int:
		return zap.Int(key, val.(int))
	case reflect.Int64:
		return zap.Int64(key, val.(int64))
	case reflect.String:
		return zap.String(key, val.(string))
	case reflect.Float32:
		return zap.Float32(key, val.(float32))
	case reflect.Float64:
		return zap.Float64(key, val.(float64))
	default:
		return zap.Any(key, val)
	}
}
