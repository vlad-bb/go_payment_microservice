package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"go_payment_microservice/cmd/server/parser"
	"go_payment_microservice/internal/services"
	us "go_payment_microservice/internal/services/user"
)

type Handler struct {
	svc *services.Services
}

func NewHandler(svc *services.Services) *Handler {
	return &Handler{svc: svc}
}

type SignUpRequestBody struct {
	Email    string `json:"email" validate:"required,email" example:"johndoe@example.com"`
	Password string `json:"password" validate:"required,min=8,max=20" example:"12345678"`
}

type SignUpResponse200Body struct {
	User *us.User `json:"user"`
}

// @Summary SignUp
// @Description Signing up new User
// @Tags auth
// @Param body body SignUpRequestBody true "Sign up request body"
// @Success 200 {object} SignUpResponse200Body "Sign up succeeded"
// @Failure 400 {object} any "Client error"
// @Failure 409 {object} any "User already exists"
// @Failure 500 {object} any "Internal server error"
// @Router /api/auth/signup [post]
func (h *Handler) SignUp(ctx *fiber.Ctx) error {
	req := &SignUpRequestBody{}
	err := parser.ParseBody(ctx, req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	pwdHash, err := h.svc.Auth.GeneratePasswordHash(req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	u, err := h.svc.User.AddUser(ctx.Context(), req.Email, pwdHash)

	if errors.Is(err, us.ErrUserAlreadyExists) {
		return fiber.NewError(fiber.StatusConflict, err.Error())
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(&SignUpResponse200Body{
		User: u,
	})
}

type SignInRequestBody struct {
	Email    string `json:"email" validate:"required,email" example:"johndoe@example.com"`
	Password string `json:"password" validate:"required,min=8,max=20" example:"12345678"`
}

type SignInResponse200Body struct {
	Token string `json:"token"`
}

// @Summary SignIn
// @Description Signing up new User
// @Tags auth
// @Param body body SignInRequestBody true "Sign up request body"
// @Success 200 {object} SignInResponse200Body "Sign up succeeded"
// @Failure 400 {object} any "Client error"
// @Failure 409 {object} any "User already exists"
// @Failure 500 {object} any "Internal server error"
// @Router /api/auth/signin [post]
func (h *Handler) SignIn(ctx *fiber.Ctx) error {
	req := &SignInRequestBody{}

	err := parser.ParseBody(ctx, req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	u, err := h.svc.User.GetUserByEmail(ctx.Context(), req.Email)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if u.IsActive != true {
		return fiber.NewError(fiber.StatusForbidden, "Permission denied, request to administrator about permission")
	}

	res, err := h.svc.Auth.CompareHashAndPassword(req.Password, u.PwdHash)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if !res {
		return fiber.NewError(fiber.StatusUnauthorized)
	}
	token, err := h.svc.Auth.CreateAuthToken(u.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized)
	}
	return ctx.JSON(&SignInResponse200Body{
		Token: token,
	})
}
