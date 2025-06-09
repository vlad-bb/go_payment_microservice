package logger

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
	"github.com/pkg/errors"
	"go_payment_microservice/internal/config"
	"go_payment_microservice/internal/logger"
	"go_payment_microservice/internal/services"
	"log/slog"
	"time"
)

type Middleware struct {
	svc *services.Services
	cfg *config.Config
}

func NewMiddleware(svc *services.Services, cfg *config.Config) *Middleware {
	return &Middleware{
		svc: svc,
		cfg: cfg,
	}
}

func (m *Middleware) Handle(ctx *fiber.Ctx) error {
	now := time.Now()
	rqId := ulid.MustNew(ulid.Timestamp(now), ulid.DefaultEntropy()).String()
	path := ctx.Path()

	WithLoggerAttrs(ctx, slog.String("instance", m.cfg.InstanceName))
	WithLoggerAttrs(ctx, slog.String("rqId", rqId), slog.String("path", path))

	logger.GetLogger().Info(
		ctx.Context(),
		"request start",
	)

	err := ctx.Next()

	duration := time.Since(now)
	WithLoggerAttrs(ctx, slog.String("duration", duration.String()))
	var fe *fiber.Error

	if errors.As(err, &fe) {
		WithLoggerAttrs(ctx, slog.String("resp_message", fe.Message))
	} else if err != nil {
		WithLoggerAttrs(ctx, slog.String("resp_message", err.Error()))
	}

	logger.GetLogger().Info(
		ctx.Context(),
		"request end",
		slog.Int("status_code", ctx.Response().StatusCode()),
	)

	return err
}
