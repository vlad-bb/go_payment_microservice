package auth

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"go_payment_microservice/cmd/server/middlewares/logger"
	"go_payment_microservice/internal/config"
	"go_payment_microservice/internal/services"
	us "go_payment_microservice/internal/services/user"
	"log/slog"
	//"time"
)

const userTokenName = "X-User-Token"
const ctxUserKey = "contextUserKey"

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
	//startedAt := time.Now()
	token := ctx.Get(userTokenName)

	userID, err := m.svc.Auth.VerifyAuthToken(token)
	if err != nil {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	u, err := m.svc.User.GetUserByID(ctx.Context(), userID)
	if err != nil {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	ctx.Locals(ctxUserKey, u)
	logger.WithLoggerAttrs(ctx, slog.String("user_id", userID))

	//telemetry.MeasureAuthTime(m.cfg.InstanceName, time.Since(startedAt))

	return ctx.Next()
}

func MustGetUser(ctx *fiber.Ctx) *us.User {
	return mustBeUser(ctx.Locals(ctxUserKey))
}

func MustGetUserWs(conn *websocket.Conn) *us.User {
	return mustBeUser(conn.Locals(ctxUserKey))
}

func mustBeUser(u interface{}) *us.User {
	if u == nil {
		panic(errors.New("user not found"))
	}
	result, ok := u.(*us.User)
	if !ok {
		panic(errors.New("user type assertion failed"))
	}
	return result
}
