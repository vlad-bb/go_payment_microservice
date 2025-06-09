package services

import (
	"go_payment_microservice/internal/clients"
	"go_payment_microservice/internal/config"
	"go_payment_microservice/internal/services/auth"
	"go_payment_microservice/internal/services/customer"
	"go_payment_microservice/internal/services/subscription"

	//"module6/internal/services/counter"
	//"module6/internal/services/message"
	"go_payment_microservice/internal/services/liqpay"
	"go_payment_microservice/internal/services/user"
)

type Services struct {
	//Counter *counter.Service
	Auth         *auth.Service
	User         *user.Service
	Liqpay       *liqpay.Service
	Customer     *customer.Service
	Subscription *subscription.Service
	//Message *message.Service
}

func NewServices(clients *clients.Clients, cfg *config.Config) *Services {
	return &Services{
		//Counter: counter.NewService(clients),
		Auth:         auth.NewService(cfg),
		User:         user.NewService(clients),
		Liqpay:       liqpay.NewService(cfg, clients),
		Customer:     customer.NewService(clients),
		Subscription: subscription.NewService(clients),
		//Message: message.NewService(clients),
	}
}
