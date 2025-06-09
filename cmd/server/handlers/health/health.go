package health

import (
	"github.com/gofiber/fiber/v2"
	"go_payment_microservice/internal/config"
)

type Handler struct {
	cfg *config.Config
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		cfg: cfg,
	}
}

type Response struct {
	Status string `json:"status"`
	Env    string `json:"env"`
}

func (h *Handler) Health(c *fiber.Ctx) error {
	return c.JSON(&Response{
		Status: "ok",
		Env:    h.cfg.Env,
	})
}
