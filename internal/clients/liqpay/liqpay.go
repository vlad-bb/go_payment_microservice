package liqpay

import (
	"context"
	"github.com/liqpay/go-sdk"
	"go_payment_microservice/internal/config"
)

type Client struct {
	LP  *liqpay.Client
	cfg *config.Config
}

func NewLiqPay(ctx context.Context, cfg *config.Config) (*Client, error) {
	return &Client{
		LP:  liqpay.New(cfg.LiqPayPublicKey, cfg.LiqPaySecretKey, nil),
		cfg: cfg,
	}, nil
}
