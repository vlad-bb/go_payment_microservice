package liqpay

// TODO fix validation

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

func orderIDValidation(fl validator.FieldLevel) bool {
	orderID := fl.Field().String()
	return strings.Contains(orderID, ":")
}

var validate *validator.Validate

func init() {
	validate = validator.New()

	err := validate.RegisterValidation("orderid", orderIDValidation)
	if err != nil {
		fmt.Println("Error registering validation function:", err)
	}
}

func ValidateCreateSubRequestBody(subRequest *CreateSubRequestBody) error {
	err := validate.Struct(subRequest)
	if err != nil {
		fmt.Println("Validation failed:", err)
		return err
	}
	return nil
}
