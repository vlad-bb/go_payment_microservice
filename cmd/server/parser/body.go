package parser

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func ParseBody(ctx *fiber.Ctx, out any) error {
	err := ctx.BodyParser(out)
	if err != nil {
		return fmt.Errorf("failed to parse body: %w", err)
	}
	validateErr := validate.Struct(out)
	if validateErr != nil {
		return fmt.Errorf("validation failed: %w", validateErr)
	}
	return nil
}
