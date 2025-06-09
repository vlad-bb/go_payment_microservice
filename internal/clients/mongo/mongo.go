package app_mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go_payment_microservice/internal/config"
)

type Client struct {
	Db *mongo.Database
}

func NewMongo(ctx context.Context, cfg *config.Config) (*Client, error) {
	opts := options.Client()
	opts.ApplyURI(cfg.MongoCn)
	conn, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}
	return &Client{
		Db: conn.Database(cfg.MongoDb),
	}, nil
}
