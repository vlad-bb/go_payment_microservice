package logger

import (
	"github.com/gofiber/fiber/v2"
	"go_payment_microservice/internal/logger"
	"log/slog"
)

func WithLoggerAttrs(ctx *fiber.Ctx, attrs ...slog.Attr) {
	ctx.Locals(logger.CtxLoggerKey{}, logger.MergeAttrs(getAttrs(ctx), attrs))
}

func getAttrs(ctx *fiber.Ctx) []slog.Attr {
	attrs := ctx.Locals(logger.CtxLoggerKey{})
	if attrs == nil {
		return []slog.Attr{}
	}
	result, ok := attrs.([]slog.Attr)
	if !ok {
		return []slog.Attr{}
	}
	return result
}
