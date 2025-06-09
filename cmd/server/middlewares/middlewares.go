package middlewares

import (
	"go_payment_microservice/cmd/server/middlewares/auth"
	"go_payment_microservice/cmd/server/middlewares/logger"
	//"module6/cmd/server/middlewares/telemetry"
	"go_payment_microservice/internal/config"
	"go_payment_microservice/internal/services"
)

type Middlewares struct {
	Logger *logger.Middleware
	Auth   *auth.Middleware
	//Telemetry *telemetry.Middleware
}

func NewMiddlewares(svc *services.Services, cfg *config.Config) *Middlewares {
	return &Middlewares{
		Logger: logger.NewMiddleware(svc, cfg),
		Auth:   auth.NewMiddleware(svc, cfg),
		//Telemetry: telemetry.NewMiddleware(cfg),
	}
}
