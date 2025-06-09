package clients

import (
	"context"
	app_liqpay "go_payment_microservice/internal/clients/liqpay"
	app_mongo "go_payment_microservice/internal/clients/mongo"
	"go_payment_microservice/internal/config"
)

type Clients struct {
	Mongo  *app_mongo.Client
	LiqPay *app_liqpay.Client
}

func NewClients(ctx context.Context, cfg *config.Config) (*Clients, error) {
	mongo, err := app_mongo.NewMongo(ctx, cfg)
	if err != nil {
		return nil, err
	}
	liqpay, err := app_liqpay.NewLiqPay(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &Clients{
		Mongo:  mongo,
		LiqPay: liqpay,
	}, nil
}
