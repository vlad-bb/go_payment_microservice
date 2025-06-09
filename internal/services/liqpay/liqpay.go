package liqpay

import (
	"context"
	"fmt"
	liqpay "github.com/liqpay/go-sdk"
	"go.mongodb.org/mongo-driver/mongo"
	"go_payment_microservice/internal/clients"
	"go_payment_microservice/internal/config"
	"go_payment_microservice/internal/constants"
	"strings"
	"time"
)

type Service struct {
	subColl *mongo.Collection
	cfg     *config.Config
	lp      *clients.Clients
}

func NewService(cfg *config.Config, clients *clients.Clients) *Service {
	return &Service{
		subColl: clients.Mongo.Db.Collection(constants.CollectionSubscriptions),
		lp:      clients,
		cfg:     cfg,
	}
}

// CreateSub response model
// swagger:model
type LPCreateSub struct {
	AcqId              int64     `json:"acq_id" bson:"acq_id"`
	Action             string    `json:"action" bson:"action"`
	Amount             float64   `json:"amount" bson:"amount"`
	Currency           string    `json:"currency" bson:"currency"`
	Description        string    `json:"description" bson:"description"`
	LiqpayOrderId      string    `json:"liqpay_order_id" bson:"liqpay_order_id"`
	OrderId            string    `json:"order_id" bson:"order_id"`
	CustomerId         string    `json:"customer_id" bson:"customer_id"`
	PaymentId          int64     `json:"payment_id" bson:"payment_id"`
	ReceiverCommission float64   `json:"receiver_commission" bson:"receiver_commission"`
	Result             string    `json:"result" bson:"result"`
	Status             string    `json:"status" bson:"status"`
	TransactionId      int64     `json:"transaction_id" bson:"transaction_id"`
	Type               string    `json:"type" bson:"type"`
	Version            int       `json:"version" bson:"version"`
	CreatedAt          time.Time `json:"created_at" bson:"created_at"`
}

// DeleteSub response model
// swagger:model
type LPDeleteSub struct {
	AcqId      int64     `json:"acq_id" bson:"acq_id"`
	OrderId    string    `json:"order_id" bson:"order_id"`
	CustomerId string    `json:"customer_id" bson:"customer_id"`
	Result     string    `json:"result" bson:"result"`
	Status     string    `json:"status" bson:"status"`
	Version    int       `json:"version" bson:"version"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
}

// UpdateSub response model
// swagger:model
type LPUpdateSub struct {
	AcqId              int64     `json:"acq_id" bson:"acq_id"`
	Action             string    `json:"action" bson:"action"`
	Amount             float64   `json:"amount" bson:"amount"`
	Currency           string    `json:"currency" bson:"currency"`
	Description        string    `json:"description" bson:"description"`
	LiqpayOrderId      string    `json:"liqpay_order_id" bson:"liqpay_order_id"`
	OrderId            string    `json:"order_id" bson:"order_id"`
	CustomerId         string    `json:"customer_id" bson:"customer_id"`
	PaymentId          int64     `json:"payment_id" bson:"payment_id"`
	ReceiverCommission float64   `json:"receiver_commission" bson:"receiver_commission"`
	Result             string    `json:"result" bson:"result"`
	Status             string    `json:"status" bson:"status"`
	TransactionId      int64     `json:"transaction_id" bson:"transaction_id"`
	Type               string    `json:"type" bson:"type"`
	Version            int       `json:"version" bson:"version"`
	CreatedAt          time.Time `json:"created_at" bson:"created_at"`
}

// LPError response model
// swagger:model
type LPError struct {
	Action         string    `json:"action" bson:"action"`
	Code           string    `json:"code" bson:"code"`
	ErrCode        string    `json:"err_code" bson:"err_code"`
	ErrDescription string    `json:"err_description" bson:"err_description"`
	LiqpayOrderId  string    `json:"liqpay_order_id" bson:"liqpay_order_id"`
	Result         string    `json:"result" bson:"result"`
	Status         string    `json:"status" bson:"status"`
	Type           string    `json:"type" bson:"type"`
	Version        int       `json:"version" bson:"version"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
}

//"action":"unsubscribe", "code":"payment_not_subscribed", "err_code":"payment_not_subscribed",
//"err_description":"Платіж не є регулярним", "is_3ds":false, "liqpay_order_id":"customer7:3",
//"public_key":"sandbox_i79440758032", "result":"error", "status":"failure", "type":"buy", "version":3

func (s *Service) CreateSub(ctx context.Context, phone string, amount float64, description string,
	orderId string, periodicity string, dateStart string, subscribeCount string, card string, month string,
	year string, cvv string) (*LPCreateSub, *LPError, error) {
	r := liqpay.Request{
		"action":                "subscribe",
		"version":               3,
		"public_key":            s.cfg.LiqPayPublicKey,
		"phone":                 phone,
		"amount":                amount,
		"currency":              "UAH",
		"description":           description,
		"order_id":              orderId,
		"card":                  card,
		"card_exp_month":        month,
		"card_exp_year":         year,
		"card_cvv":              cvv,
		"subscribe":             subscribeCount,
		"subscribe_date_start":  dateStart,
		"subscribe_periodicity": periodicity,
	}
	resp, err := s.lp.LiqPay.LP.Send("request", r)
	if err != nil {
		fmt.Printf("error %v\n", err.Error())
		return nil, nil, err
	}
	fmt.Printf("response %#v\n", resp)
	status, _ := resp["status"].(string)
	if status == "failure" || status == "error" {
		Action, _ := resp["action"].(string)
		Code, _ := resp["code"].(string)
		ErrCode, _ := resp["err_code"].(string)
		ErrDescription, _ := resp["err_description"].(string)
		LiqpayOrderId, _ := resp["liqpay_order_id"].(string)
		Result, _ := resp["result"].(string)
		Status, _ := resp["status"].(string)
		Type, _ := resp["type"].(string)
		Version, _ := resp["version"].(float64)
		return nil,
			&LPError{
				Action:         Action,
				Code:           Code,
				ErrCode:        ErrCode,
				ErrDescription: ErrDescription,
				LiqpayOrderId:  LiqpayOrderId,
				Result:         Result,
				Status:         Status,
				Type:           Type,
				Version:        int(Version),
				CreatedAt:      time.Now().UTC(),
			}, nil
	}
	AcqId, _ := resp["acq_id"].(float64)
	Action, _ := resp["action"].(string)
	Amount, _ := resp["amount"].(float64)
	Currency, _ := resp["currency"].(string)
	Description, _ := resp["description"].(string)
	LiqpayOrderId, _ := resp["liqpay_order_id"].(string)
	OrderId, _ := resp["order_id"].(string)
	parts := strings.Split(OrderId, ":")
	CustomerId := parts[0]
	PaymentId, _ := resp["payment_id"].(float64)
	ReceiverCommission, _ := resp["receiver_commission"].(float64)
	Result, _ := resp["result"].(string)
	Status, _ := resp["status"].(string)
	TransactionId, _ := resp["transaction_id"].(float64)
	Type, _ := resp["type"].(string)
	Version, _ := resp["version"].(float64)
	return &LPCreateSub{
		AcqId:              int64(AcqId),
		Action:             Action,
		Amount:             Amount,
		Currency:           Currency,
		Description:        Description,
		LiqpayOrderId:      LiqpayOrderId,
		OrderId:            OrderId,
		CustomerId:         CustomerId,
		PaymentId:          int64(PaymentId),
		ReceiverCommission: ReceiverCommission,
		Result:             Result,
		Status:             Status,
		TransactionId:      int64(TransactionId),
		Type:               Type,
		Version:            int(Version),
		CreatedAt:          time.Now().UTC(),
	}, nil, nil
}

func (s *Service) DeleteSub(ctx context.Context, orderId string) (*LPDeleteSub, *LPError, error) {
	r := liqpay.Request{
		"action":     "unsubscribe",
		"version":    3,
		"public_key": s.cfg.LiqPayPublicKey,
		"order_id":   orderId,
	}
	resp, err := s.lp.LiqPay.LP.Send("request", r)
	fmt.Printf("response %#v\n", resp)
	if err != nil {
		fmt.Printf("error %v\n", err.Error())
		return nil, nil, err
	}
	status, _ := resp["status"].(string)
	if status == "failure" || status == "error" {
		Action, _ := resp["action"].(string)
		Code, _ := resp["code"].(string)
		ErrCode, _ := resp["err_code"].(string)
		ErrDescription, _ := resp["err_description"].(string)
		LiqpayOrderId, _ := resp["liqpay_order_id"].(string)
		Result, _ := resp["result"].(string)
		Status, _ := resp["status"].(string)
		Type, _ := resp["type"].(string)
		Version, _ := resp["version"].(float64)
		return nil,
			&LPError{
				Action:         Action,
				Code:           Code,
				ErrCode:        ErrCode,
				ErrDescription: ErrDescription,
				LiqpayOrderId:  LiqpayOrderId,
				Result:         Result,
				Status:         Status,
				Type:           Type,
				Version:        int(Version),
				CreatedAt:      time.Now().UTC(),
			}, nil
	}

	AcqId, _ := resp["acq_id"].(float64)
	OrderId, _ := resp["order_id"].(string)
	parts := strings.Split(OrderId, ":")
	CustomerId := parts[0]
	Result, _ := resp["result"].(string)
	Status, _ := resp["status"].(string)
	Version, _ := resp["version"].(float64)
	return &LPDeleteSub{
			AcqId:      int64(AcqId),
			OrderId:    OrderId,
			CustomerId: CustomerId,
			Result:     Result,
			Status:     Status,
			Version:    int(Version),
			CreatedAt:  time.Now().UTC(),
		},
		nil,
		nil
}

func (s *Service) UpdateSub(ctx context.Context, phone string, amount float64, description string,
	orderId string, card string, month string,
	year string, cvv string) (*LPUpdateSub, *LPError, error) {
	r := liqpay.Request{
		"action":         "subscribe_update",
		"version":        3,
		"public_key":     s.cfg.LiqPayPublicKey,
		"phone":          phone,
		"amount":         amount,
		"currency":       "UAH",
		"description":    description,
		"order_id":       orderId,
		"card":           card,
		"card_exp_month": month,
		"card_exp_year":  year,
		"card_cvv":       cvv,
	}
	resp, err := s.lp.LiqPay.LP.Send("request", r)
	if err != nil {
		fmt.Printf("error %v\n", err.Error())
		return nil, nil, err
	}
	fmt.Printf("response %#v\n", resp)
	status, _ := resp["status"].(string)
	if status == "failure" || status == "error" {
		Action, _ := resp["action"].(string)
		Code, _ := resp["code"].(string)
		ErrCode, _ := resp["err_code"].(string)
		ErrDescription, _ := resp["err_description"].(string)
		LiqpayOrderId, _ := resp["liqpay_order_id"].(string)
		Result, _ := resp["result"].(string)
		Status, _ := resp["status"].(string)
		Type, _ := resp["type"].(string)
		Version, _ := resp["version"].(float64)
		return nil,
			&LPError{
				Action:         Action,
				Code:           Code,
				ErrCode:        ErrCode,
				ErrDescription: ErrDescription,
				LiqpayOrderId:  LiqpayOrderId,
				Result:         Result,
				Status:         Status,
				Type:           Type,
				Version:        int(Version),
				CreatedAt:      time.Now().UTC(),
			}, nil
	}
	AcqId, _ := resp["acq_id"].(float64)
	Action, _ := resp["action"].(string)
	Amount, _ := resp["amount"].(float64)
	Currency, _ := resp["currency"].(string)
	Description, _ := resp["description"].(string)
	LiqpayOrderId, _ := resp["liqpay_order_id"].(string)
	OrderId, _ := resp["order_id"].(string)
	parts := strings.Split(OrderId, ":")
	CustomerId := parts[0]
	PaymentId, _ := resp["payment_id"].(float64)
	ReceiverCommission, _ := resp["receiver_commission"].(float64)
	Result, _ := resp["result"].(string)
	Status, _ := resp["status"].(string)
	TransactionId, _ := resp["transaction_id"].(float64)
	Type, _ := resp["type"].(string)
	Version, _ := resp["version"].(float64)
	return &LPUpdateSub{
		AcqId:              int64(AcqId),
		Action:             Action,
		Amount:             Amount,
		Currency:           Currency,
		Description:        Description,
		LiqpayOrderId:      LiqpayOrderId,
		OrderId:            OrderId,
		CustomerId:         CustomerId,
		PaymentId:          int64(PaymentId),
		ReceiverCommission: ReceiverCommission,
		Result:             Result,
		Status:             Status,
		TransactionId:      int64(TransactionId),
		Type:               Type,
		Version:            int(Version),
		CreatedAt:          time.Now().UTC(),
	}, nil, nil
}
