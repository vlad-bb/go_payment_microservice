package handlers

import (
	//"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"go_payment_microservice/cmd/server/handlers/auth"
	"go_payment_microservice/cmd/server/handlers/customer"

	//"module6/cmd/server/handlers/counter"
	"go_payment_microservice/cmd/server/handlers/health"
	"go_payment_microservice/cmd/server/handlers/liqpay"
	//"module6/cmd/server/handlers/message"
	"go_payment_microservice/cmd/server/middlewares"
	"go_payment_microservice/internal/config"
	"go_payment_microservice/internal/services"
)

type Handlers struct {
	Health *health.Handler
	//Counter *counter.Handler
	Auth     *auth.Handler
	LiqPay   *liqpay.Handler
	Customer *customer.Handler
	//Message *message.Handler

	mdlwr *middlewares.Middlewares
}

func NewHandlers(cfg *config.Config, svcs *services.Services, mdlwr *middlewares.Middlewares) *Handlers {
	return &Handlers{
		Health: health.NewHandler(cfg),
		//Counter: counter.NewHandler(svcs),
		Auth: auth.NewHandler(svcs),
		//Message: message.NewHandler(svcs, cfg),
		LiqPay:   liqpay.NewHandler(svcs, cfg),
		Customer: customer.NewHandler(svcs),
		mdlwr:    mdlwr,
	}
}

func (h *Handlers) RegisterRoutes(router fiber.Router) {

	//router.Get("/metrics", prometheusHandler())
	router.Get("/health", h.Health.Health)
	router.Get("/docs/*", swagger.HandlerDefault)

	api := router.Group("/api")
	//api.Use(h.mdlwr.Telemetry.Handle)
	api.Use(h.mdlwr.Logger.Handle)

	authGroup := api.Group("/auth")
	authGroup.Post("/signup", h.Auth.SignUp)
	authGroup.Post("/signin", h.Auth.SignIn)

	paymentGroup := api.Group("/subscription")
	paymentGroup.Use(h.mdlwr.Auth.Handle)
	paymentGroup.Post("/create-sub", h.LiqPay.CreateSub)
	paymentGroup.Delete("/delete-sub", h.LiqPay.DeleteSub)
	paymentGroup.Put("/update-sub", h.LiqPay.UpdateSub)
	paymentGroup.Get("/:telegram_id", h.LiqPay.GetSub)

	customerGroup := api.Group("/customer")
	customerGroup.Use(h.mdlwr.Auth.Handle)
	customerGroup.Post("/create", h.Customer.CreateCustomer)
	customerGroup.Put("/update", h.Customer.UpdateCustomer)
	customerGroup.Get("/list", h.Customer.ListCustomer)
	customerGroup.Get("/:telegram_id", h.Customer.GetCustomer)

	//messageGroup := api.Group("/message")
	//messageGroup.Use(h.mdlwr.Auth.Handle)
	//messageGroup.Post("/send", h.Message.SendMessage)
	//messageGroup.Post("/list", h.Message.List)
	//messageGroup.Get("/listen/:channel_name", websocket.New(h.Message.Listen))
}

func prometheusHandler() fiber.Handler {
	handler := promhttp.Handler()
	return func(c *fiber.Ctx) error {
		fasthttpadaptor.NewFastHTTPHandler(handler)(c.Context())
		return nil
	}
}
