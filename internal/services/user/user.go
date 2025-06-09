package user

import (
	"context"
	"fmt"
	"github.com/oklog/ulid/v2"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go_payment_microservice/internal/clients"
	"go_payment_microservice/internal/constants"
	"time"
)

var (
	ErrUserAlreadyExists = errors.New("User already exists")
)

type Service struct {
	userColl *mongo.Collection
}

func NewService(clients *clients.Clients) *Service {
	return &Service{
		userColl: clients.Mongo.Db.Collection(constants.CollectionUsers),
	}
}

// User represents a user in the system.
// swagger:model
type User struct {
	ID        string    `json:"id" bson:"id" example:"0001M2PVBD5Q1DAMYJ0S2HADD6"`
	Email     string    `json:"email" bson:"email" example:"johndoe@example.com"`
	IsActive  bool      `json:"is_active" bson:"is_active" example:"true"`
	PwdHash   string    `json:"pwd_hash" bson:"pwd_hash" example:"$2a$10$pikzoSYzIs1GRRPi0vermeY1mPH4"`
	CreatedAt time.Time `json:"created_at" bson:"created_at" example:"2020-01-01T00:00:00+09:00"`
}

func (s *Service) AddUser(ctx context.Context, email string, pwdHash string) (*User, error) {
	t := time.Now().UTC()

	u := &User{
		ID:        ulid.MustNew(uint64(t.Unix()), ulid.DefaultEntropy()).String(),
		Email:     email,
		PwdHash:   pwdHash,
		IsActive:  false,
		CreatedAt: t,
	}

	_, err := s.userColl.InsertOne(ctx, u)
	if mongo.IsDuplicateKeyError(err) {
		return nil, ErrUserAlreadyExists
	}
	if err != nil {
		return nil, fmt.Errorf("failed to insert User: %w", err)
	}

	return u, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var result *User
	err := s.userColl.FindOne(ctx, bson.M{"email": email}).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to find User by Name: %w", err)
	}
	return result, nil
}

func (s *Service) GetUserByID(ctx context.Context, id string) (*User, error) {
	var result *User
	err := s.userColl.FindOne(ctx, bson.M{"id": id}).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to find User by ID: %w", err)
	}
	return result, nil
}
