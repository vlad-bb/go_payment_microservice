package customer

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"go_payment_microservice/cmd/server/parser"
	"go_payment_microservice/internal/services"
	cs "go_payment_microservice/internal/services/customer"
)

type Handler struct {
	svc *services.Services
}

func NewHandler(svc *services.Services) *Handler {
	return &Handler{svc: svc}
}

type CreateCustomerRequestBody struct {
	ID   string `json:"telegram_id" bson:"telegram_id" example:"123456789"`
	Name string `json:"name" bson:"name" example:"John Doe"`
}

type CreateCustomerResponse200Body struct {
	Customer *cs.Customer `json:"customer"`
}

// @Summary Create Customer
// @Description Create new customer
// @Tags customer
// @Param body body CreateCustomerRequestBody true "Create customer request body"
// @Success 200 {object} CreateCustomerResponse200Body "Create customer succeeded"
// @Failure 400 {object} any "Client error"
// @Failure 409 {object} any "Customer already exists"
// @Failure 500 {object} any "Internal server error"
// @Router /api/customer/create [post]
func (h *Handler) CreateCustomer(ctx *fiber.Ctx) error {
	req := &CreateCustomerRequestBody{}
	err := parser.ParseBody(ctx, req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	c, err := h.svc.Customer.AddCustomer(ctx.Context(), req.ID, req.Name)

	if errors.Is(err, cs.ErrCustomerAlreadyExists) {
		return fiber.NewError(fiber.StatusConflict, err.Error())
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(&CreateCustomerResponse200Body{
		Customer: c,
	})
}

type UpdateCustomerRequestBody struct {
	ID       string `json:"telegram_id" bson:"telegram_id" example:"123456789"`
	Name     string `json:"name" bson:"name" example:"John Doe"`
	IsActive bool   `json:"is_active" bson:"is_active" example:"true"`
}

type UpdateCustomerResponse200Body struct {
	Status string `json:"status" bson:"status" example:"OK"`
}

// @Summary Update Customer
// @Description Update existing customer
// @Tags customer
// @Param body body UpdateCustomerRequestBody true "Update customer request body"
// @Success 200 {object} UpdateCustomerResponse200Body "Update customer succeeded"
// @Failure 400 {object} any "Client error"
// @Failure 409 {object} any "Customer already exists"
// @Failure 500 {object} any "Internal server error"
// @Router /api/customer/update [put]
func (h *Handler) UpdateCustomer(ctx *fiber.Ctx) error {
	req := &UpdateCustomerRequestBody{}
	err := parser.ParseBody(ctx, req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err = h.svc.Customer.UpdateCustomer(ctx.Context(), req.ID, req.Name, req.IsActive)

	if errors.Is(err, cs.ErrCustomerNotFound) {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(&UpdateCustomerResponse200Body{
		Status: "Success",
	})
}

type GetCustomerResponse200Body struct {
	Customer *cs.Customer `json:"customer"`
}

// @Summary Get Customer
// @Description Get customer by telegram_id
// @Tags customer
// @Param telegram_id path string true "Telegram ID of the customer"
// @Success 200 {object} GetCustomerResponse200Body "Get customer succeeded"
// @Failure 404 {object} any "Customer not found"
// @Failure 500 {object} any "Internal server error"
// @Router /api/customer/{telegram_id} [get]
func (h *Handler) GetCustomer(ctx *fiber.Ctx) error {
	telegramId := ctx.Params("telegram_id")
	customer, err := h.svc.Customer.GetCustomerByID(ctx.Context(), telegramId)
	fmt.Printf("%+v\n", customer)
	fmt.Printf("%+v\n", err)

	if err != nil {
		if errors.Is(err, cs.ErrCustomerNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(&GetCustomerResponse200Body{
		Customer: customer,
	})
}

type ListCustomerResponse200Body struct {
	Customers []*cs.Customer `json:"customers"`
}

// @Summary Get list of customers
// @Description List customers
// @Tags customer
// @Success 200 {object} ListCustomerResponse200Body "List customers succeeded"
// @Failure 404 {object} any "Customers not found"
// @Failure 500 {object} any "Internal server error"
// @Router /api/customer/list [get]
func (h *Handler) ListCustomer(ctx *fiber.Ctx) error {
	customers, err := h.svc.Customer.GetAllCustomers(ctx.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(&ListCustomerResponse200Body{
		Customers: customers,
	})
}
