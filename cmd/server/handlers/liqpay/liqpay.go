package liqpay

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go_payment_microservice/cmd/server/middlewares/auth"
	"go_payment_microservice/cmd/server/parser"
	"go_payment_microservice/internal/config"
	"go_payment_microservice/internal/services"
	lp "go_payment_microservice/internal/services/liqpay"
	ss "go_payment_microservice/internal/services/subscription"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	svc *services.Services
	cfg *config.Config
}

func NewHandler(svc *services.Services, cfg *config.Config) *Handler {
	return &Handler{svc: svc, cfg: cfg}
}

type CreateSubRequestBody struct {
	Phone                string  `json:"phone"`  // Телефон вказується в міжнародному форматі (Україна +380). Наприклад: +380950000001 (з +) або 380950000001 (без +)
	Amount               float64 `json:"amount"` // число більше 0
	Description          string  `json:"description"`
	OrderId              string  `json:"order_id"` //validate:"orderid"` // має бути символ розподілювача : TODO fix validation
	Card                 string  `json:"card"`
	CardExpMonth         string  `json:"card_exp_month"`
	CardExpYear          string  `json:"card_exp_year"`
	CardCvv              string  `json:"card_cvv"`
	Subscribe            string  `json:"subscribe"`             // ціле число більше 1
	SubscribeDateStart   string  `json:"subscribe_date_start"`  // Час необхідно вказувати в такому форматі 2015-03-31 00:00:00 по UTC
	SubscribePeriodicity string  `json:"subscribe_periodicity"` //Можливі значення: day - щодня, week - щотижня, month - раз на місяць, year - раз на рік
}

type CreateSubResponseBody struct {
	Body *lp.LPCreateSub `json:"body"`
}

// @Summary Create Subscription
// @Description Create Subscription for customer
// @Tags subscription
// @Param body body CreateSubRequestBody true "Create subscription request body"
// @Success 200 {object} CreateSubResponseBody "Create subscription succeeded"
// @Failure 400 {object} any "Client error"
// @Failure 500 {object} any "Internal server error"
// @Router /api/subscription/create-sub [post]
func (h *Handler) CreateSub(ctx *fiber.Ctx) error {
	request := &CreateSubRequestBody{}
	err := parser.ParseBody(ctx, request)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	//err = ValidateCreateSubRequestBody(msg) // TODO not work, raise panic: Undefined validation function 'orderid' on field 'OrderId'
	valid := strings.Contains(request.OrderId, ":")
	if valid != true {
		errStr := fmt.Sprintf("Invalid order_id, not contain ':' symbol%v", request.OrderId)
		return fiber.NewError(fiber.StatusBadRequest, errStr)
	}
	_ = auth.MustGetUser(ctx)
	successResponse, failureResponse, err := h.svc.Liqpay.CreateSub(ctx.Context(), request.Phone, request.Amount, request.Description,
		request.OrderId, request.SubscribePeriodicity, request.SubscribeDateStart, request.SubscribePeriodicity,
		request.Card, request.CardExpMonth, request.CardExpYear, request.CardCvv)
	if err != nil {
		//telemetry.AddMessageSend(h.cfg.InstanceName, telemetry.MessageStatusFailure)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	//telemetry.AddMessageSend(h.cfg.InstanceName, telemetry.MessageStatusSuccess)

	if successResponse != nil {
		subscribeInt, _ := strconv.Atoi(request.Subscribe)
		endDate, err := h.svc.Subscription.CalculateSubscribeDateEnd(request.SubscribeDateStart, subscribeInt, request.SubscribePeriodicity)
		data := ss.Subscription{
			Amount:               request.Amount,
			Description:          request.Description,
			OrderId:              request.OrderId,
			CustomerId:           successResponse.CustomerId,
			Status:               successResponse.Status,
			Subscribe:            request.Subscribe,            // ціле число більше 1 - кількість списань
			SubscribeDateStart:   request.SubscribeDateStart,   // "2025-06-08 00:00:00",
			SubscribeDateEnd:     endDate,                      // приклад, // треба розрахувати
			SubscribePeriodicity: request.SubscribePeriodicity, //Можливі значення: day - щодня, week - щотижня, month - раз на місяць, year - раз на рік
			CreatedAt:            time.Now().UTC(),
		}
		err = h.svc.Subscription.SaveNewSub(ctx.Context(), data)
		if err != nil {
			fmt.Println(err)
		}
		return ctx.JSON(&CreateSubResponseBody{
			Body: successResponse,
		})
	}
	if failureResponse != nil {
		return ctx.JSON(&LPErrorResponseBody{
			Body: failureResponse,
		})
	}
	response := fiber.NewError(fiber.StatusInternalServerError, "unexpected error")
	return response
}

type DeleteSubRequestBody struct {
	OrderId string `json:"order_id"` //validate:"orderid"` // має бути символ розподілювача : TODO fix validation
}

type DeleteSubResponseBody struct {
	Body *lp.LPDeleteSub `json:"body"`
}

type LPErrorResponseBody struct {
	Body *lp.LPError `json:"body"`
}

// @Summary Delete Subscription
// @Description Delete Subscription for customer
// @Tags subscription
// @Param body body DeleteSubRequestBody true "Delete subscription request body"
// @Success 200 {object} DeleteSubResponseBody "Delete subscription succeeded"
// @Failure 400 {object} any "Client error"
// @Failure 500 {object} any "Internal server error"
// @Router /api/subscription/delete-sub [delete]
func (h *Handler) DeleteSub(ctx *fiber.Ctx) error {
	request := &DeleteSubRequestBody{}
	err := parser.ParseBody(ctx, request)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	//err = ValidateCreateSubRequestBody(msg) // TODO not work, raise panic: Undefined validation function 'orderid' on field 'OrderId'
	valid := strings.Contains(request.OrderId, ":")
	if valid != true {
		errStr := fmt.Sprintf("Invalid order_id, not contain ':' symbol%v", request.OrderId)
		return fiber.NewError(fiber.StatusBadRequest, errStr)
	}
	_ = auth.MustGetUser(ctx)
	successResponse, failureResponse, err := h.svc.Liqpay.DeleteSub(ctx.Context(), request.OrderId)
	if err != nil {
		//telemetry.AddMessageSend(h.cfg.InstanceName, telemetry.MessageStatusFailure)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	//telemetry.AddMessageSend(h.cfg.InstanceName, telemetry.MessageStatusSuccess)
	if successResponse != nil {
		err = h.svc.Subscription.UpdateStatusSub(ctx.Context(), successResponse.OrderId, successResponse.Status)
		if err != nil {
			fmt.Println(err)
		}
		return ctx.JSON(&DeleteSubResponseBody{
			Body: successResponse,
		})
	}
	if failureResponse != nil {
		return ctx.JSON(&LPErrorResponseBody{
			Body: failureResponse,
		})
	}
	response := fiber.NewError(fiber.StatusInternalServerError, "unexpected error")
	return response
}

type UpdateSubRequestBody struct {
	Phone        string  `json:"phone"`  // Телефон вказується в міжнародному форматі (Україна +380). Наприклад: +380950000001 (з +) або 380950000001 (без +)
	Amount       float64 `json:"amount"` // число більше 0
	Description  string  `json:"description"`
	OrderId      string  `json:"order_id"` //validate:"orderid"` // має бути символ розподілювача : TODO fix validation
	Card         string  `json:"card"`
	CardExpMonth string  `json:"card_exp_month"`
	CardExpYear  string  `json:"card_exp_year"`
	CardCvv      string  `json:"card_cvv"`
}

type UpdateSubResponseBody struct {
	Body *lp.LPUpdateSub `json:"body"`
}

// @Summary Update Subscription
// @Description Update Subscription for customer
// @Tags subscription
// @Param body body UpdateSubRequestBody true "Update subscription request body"
// @Success 200 {object} UpdateSubResponseBody "Update subscription succeeded"
// @Failure 400 {object} any "Client error"
// @Failure 500 {object} any "Internal server error"
// @Router /api/subscription/update-sub [put]
func (h *Handler) UpdateSub(ctx *fiber.Ctx) error {
	request := &UpdateSubRequestBody{}
	err := parser.ParseBody(ctx, request)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	//err = ValidateCreateSubRequestBody(msg) // TODO not work, raise panic: Undefined validation function 'orderid' on field 'OrderId'
	valid := strings.Contains(request.OrderId, ":")
	if valid != true {
		errStr := fmt.Sprintf("Invalid order_id, not contain ':' symbol%v", request.OrderId)
		return fiber.NewError(fiber.StatusBadRequest, errStr)
	}
	_ = auth.MustGetUser(ctx)
	successResponse, failureResponse, err := h.svc.Liqpay.UpdateSub(ctx.Context(), request.Phone, request.Amount,
		request.Description, request.OrderId,
		request.Card, request.CardExpMonth, request.CardExpYear, request.CardCvv)
	if err != nil {
		//telemetry.AddMessageSend(h.cfg.InstanceName, telemetry.MessageStatusFailure)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	//telemetry.AddMessageSend(h.cfg.InstanceName, telemetry.MessageStatusSuccess)

	if successResponse != nil {
		err = h.svc.Subscription.UpdateSub(ctx.Context(), successResponse.OrderId, successResponse.Description, successResponse.Amount)
		if err != nil {
			fmt.Println(err)
		}
		return ctx.JSON(&UpdateSubResponseBody{
			Body: successResponse,
		})
	}
	if failureResponse != nil {
		return ctx.JSON(&LPErrorResponseBody{
			Body: failureResponse,
		})
	}
	response := fiber.NewError(fiber.StatusInternalServerError, "unexpected error")
	return response
}

type GetSubResponse200Body struct {
	Sub *ss.Subscription `json:"subscription"`
}

// @Summary Get Subscription
// @Description Get subscription by telegram_id
// @Tags subscription
// @Param telegram_id path string true "Telegram ID of the customer"
// @Success 200 {object} GetSubResponse200Body "Get subscription succeeded"
// @Failure 404 {object} any "Subscription not found"
// @Failure 500 {object} any "Internal server error"
// @Router /api/subscription/{telegram_id} [get]
func (h *Handler) GetSub(ctx *fiber.Ctx) error {
	telegramId := ctx.Params("telegram_id")
	sub, err := h.svc.Subscription.GetSubscriptionByID(ctx.Context(), telegramId)
	fmt.Printf("%+v\n", sub)
	fmt.Printf("%+v\n", err)

	if err != nil {
		if errors.Is(err, ss.ErrorSubNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(&GetSubResponse200Body{
		Sub: sub,
	})
}
