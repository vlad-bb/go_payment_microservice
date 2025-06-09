package cutomer

import (
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go_payment_microservice/internal/clients"
	"go_payment_microservice/internal/constants"
	"time"
)

var (
	ErrUserAlreadyExists = errors.New("User already exists")
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
type User struct {
	ID        string    `json:"id" bson:"id" example:"0001M2PVBD5Q1DAMYJ0S2HADD6"`
	Email     string    `json:"email" bson:"email" example:"johndoe@example.com"`
	IsActive  bool      `json:"is_active" bson:"is_active" example:"true"`
	PwdHash   string    `json:"pwd_hash" bson:"pwd_hash" example:"$2a$10$pikzoSYzIs1GRRPi0vermeY1mPH4"`
	CreatedAt time.Time `json:"created_at" bson:"created_at" example:"2020-01-01T00:00:00+09:00"`
}
