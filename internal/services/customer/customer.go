package customer

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go_payment_microservice/internal/clients"
	"go_payment_microservice/internal/constants"
	"time"
)

var (
	ErrCustomerAlreadyExists = errors.New("Customer already exists")
	ErrCustomerNotFound      = errors.New("Customer not found")
)

type Service struct {
	customerColl *mongo.Collection
}

func NewService(clients *clients.Clients) *Service {
	return &Service{
		customerColl: clients.Mongo.Db.Collection(constants.CollectionCustomers),
	}
}

// Customer represents a customer in the system.
// swagger:model
type Customer struct {
	ID        string    `json:"telegram_id" bson:"telegram_id" example:"123456789"`
	Name      string    `json:"name" bson:"name" example:"John Doe"`
	IsActive  bool      `json:"is_active" bson:"is_active" example:"true"`
	CreatedAt time.Time `json:"created_at" bson:"created_at" example:"2020-01-01T00:00:00+09:00"`
}

func (s *Service) AddCustomer(ctx context.Context, telegramId string, name string) (*Customer, error) {
	t := time.Now().UTC()

	c := &Customer{
		ID:        telegramId,
		Name:      name,
		IsActive:  true,
		CreatedAt: t,
	}

	_, err := s.customerColl.InsertOne(ctx, c)
	if mongo.IsDuplicateKeyError(err) {
		return nil, ErrCustomerAlreadyExists
	}
	if err != nil {
		return nil, fmt.Errorf("failed to insert Customer: %w", err)
	}

	return c, nil
}

func (s *Service) GetCustomerByID(ctx context.Context, telegramId string) (*Customer, error) {
	var result *Customer
	err := s.customerColl.FindOne(ctx, bson.M{"telegram_id": telegramId}).Decode(&result)
	if err != nil {
		return nil, ErrCustomerNotFound
	}
	return result, nil
}

func (s *Service) UpdateCustomer(ctx context.Context, telegramId string, name string, isActive bool) error {
	c, err := s.GetCustomerByID(ctx, telegramId)
	if err != nil {
		return err
	}

	c.IsActive = isActive
	c.Name = name

	update := bson.M{
		"$set": bson.M{
			"name":      c.Name,
			"is_active": c.IsActive,
		},
	}

	_, err = s.customerColl.UpdateOne(ctx, bson.M{"telegram_id": telegramId}, update)
	return err
}

func (s *Service) isActiveCustomer(ctx context.Context, telegramId string) (bool, error) {
	c, err := s.GetCustomerByID(ctx, telegramId)
	if err != nil {
		return false, err
	}
	return c.IsActive, nil
}

func (s *Service) GetAllCustomers(ctx context.Context) ([]*Customer, error) {
	var customers []*Customer
	cursor, err := s.customerColl.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var customer Customer
		if err := cursor.Decode(&customer); err != nil {
			return nil, err
		}
		customers = append(customers, &customer)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}
