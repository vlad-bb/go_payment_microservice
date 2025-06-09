package logger

import (
	"context"
	"go_payment_microservice/internal/config"
	"log/slog"
	"os"
)

var logger *Logger

func GetLogger() *Logger {
	return logger
}

type Logger struct {
	cfg *config.Config
	log *slog.Logger
}

func InitLogger(cfg *config.Config) {

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	h := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: convertToSlogLevel(cfg.LoggerLevel),
	})
	logger = &Logger{
		cfg: cfg,
		log: slog.New(h),
	}
}

func convertToSlogLevel(level string) slog.Leveler {
	switch level {
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	case "debug":
		return slog.LevelDebug
	}
	return slog.LevelInfo
}

type CtxLoggerKey struct{}

func WithAttrs(ctx context.Context, args ...slog.Attr) context.Context {
	ctx = context.WithValue(ctx, CtxLoggerKey{}, MergeAttrs(getAttrs(ctx), args))
	return ctx
}

func MergeAttrs(left []slog.Attr, right []slog.Attr) []slog.Attr {
	return append(left, right...)
}

func getAttrs(ctx context.Context) []slog.Attr {
	attrs := ctx.Value(CtxLoggerKey{})
	if attrs == nil {
		return []slog.Attr{}
	}
	result, ok := attrs.([]slog.Attr)
	if !ok {
		return []slog.Attr{}
	}
	return result
}

func convertAttrsToAny(a []slog.Attr) []any {
	result := make([]any, len(a))
	for i, v := range a {
		result[i] = v
	}
	return result
}

func (s *Logger) Info(ctx context.Context, msg string, args ...slog.Attr) {
	s.log.InfoContext(ctx, msg, convertAttrsToAny(MergeAttrs(getAttrs(ctx), args))...)
}

func (s *Logger) Error(ctx context.Context, err error, args ...slog.Attr) {
	s.log.ErrorContext(ctx, err.Error(), convertAttrsToAny(MergeAttrs(getAttrs(ctx), args))...)
}

func (s *Logger) Panic(ctx context.Context, err error, args ...slog.Attr) {
	s.log.ErrorContext(ctx, err.Error(), args)
	panic(err)
}

func (s *Logger) Fatal(ctx context.Context, err error, args ...slog.Attr) {
	s.log.ErrorContext(ctx, err.Error(), args)
	os.Exit(1)
}
