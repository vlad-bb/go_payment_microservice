package subscription

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go_payment_microservice/internal/clients"
	"go_payment_microservice/internal/constants"
	"time"
)

var (
	OrderIdNotFound  = errors.New("order Id not found")
	ErrorSubNotFound = errors.New("subscription not found")
)

type Service struct {
	subColl *mongo.Collection
}

func NewService(clients *clients.Clients) *Service {
	return &Service{
		subColl: clients.Mongo.Db.Collection(constants.CollectionSubscriptions),
	}
}

// Customer represents a subscription in app.
// swagger:model
type Subscription struct {
	Amount               float64   `json:"amount"`
	Description          string    `json:"description"`
	OrderId              string    `json:"order_id" bson:"order_id"`
	CustomerId           string    `json:"customer_id" bson:"customer_id"`
	Status               string    `json:"status" bson:"status"`
	Subscribe            string    `json:"subscribe"`
	SubscribeDateStart   string    `json:"subscribe_date_start"`
	SubscribeDateEnd     string    `json:"subscribe_date_end"`
	SubscribePeriodicity string    `json:"subscribe_periodicity"`
	CreatedAt            time.Time `json:"created_at" bson:"created_at"`
}

func (s *Service) SaveNewSub(ctx context.Context, data Subscription) error {
	existingSub := &Subscription{}
	err := s.subColl.FindOne(ctx, bson.M{"order_id": data.OrderId}).Decode(existingSub)
	if err == nil {
		return fmt.Errorf("subscription with order_id %s already exists", data.OrderId)
	}
	_, err = s.subColl.InsertOne(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to insert new subscription: %v", err)
	}
	return nil
}

func (s *Service) UpdateStatusSub(ctx context.Context, OrderId string, Status string) error {
	sub := &Subscription{}
	err := s.subColl.FindOne(ctx, bson.M{"order_id": OrderId}).Decode(sub)
	if err != nil {
		return OrderIdNotFound
	}
	sub.Status = Status
	update := bson.M{
		"$set": bson.M{
			"status": sub.Status,
		},
	}
	_, err = s.subColl.UpdateOne(ctx, bson.M{"order_id": OrderId}, update)
	return nil
}

func (s *Service) UpdateSub(ctx context.Context, OrderId string, Description string, Amount float64) error {
	sub := &Subscription{}
	err := s.subColl.FindOne(ctx, bson.M{"order_id": OrderId}).Decode(sub)
	if err != nil {
		return OrderIdNotFound
	}
	sub.Description = Description
	sub.Amount = Amount
	update := bson.M{
		"$set": bson.M{
			"description": sub.Description,
			"amount":      sub.Amount,
		},
	}
	_, err = s.subColl.UpdateOne(ctx, bson.M{"order_id": OrderId}, update)
	return nil
}

func (s *Service) CalculateSubscribeDateEnd(startDate string, subscribe int, periodicity string) (string, error) {
	start, err := time.Parse("2006-01-02 15:04:05", startDate)
	if err != nil {
		return "", fmt.Errorf("invalid start date format: %v", err)
	}
	var end time.Time
	switch periodicity {
	case "day":
		end = start.AddDate(0, 0, subscribe)
	case "week":
		end = start.AddDate(0, 0, subscribe*7)
	case "month":
		end = start.AddDate(0, subscribe, 0)
	case "year":
		end = start.AddDate(subscribe, 0, 0)
	default:
		return "", fmt.Errorf("invalid periodicity value: %v", periodicity)
	}
	return end.Format("2006-01-02 15:04:05"), nil
}

func (s *Service) GetSubscriptionByID(ctx context.Context, telegramId string) (*Subscription, error) {
	var result *Subscription
	err := s.subColl.FindOne(ctx, bson.M{"customer_id": telegramId, "status": "subscribed"}).Decode(&result)
	if err != nil {
		return nil, ErrorSubNotFound
	}
	return result, nil
}
