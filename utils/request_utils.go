package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func ReturnErrorResponse(status int, err error, context *fiber.Ctx) error {
	return context.Status(status).JSON(&fiber.Map{
		"error": err.Error(),
	})
}

func ReturnSuccessResponse(context *fiber.Ctx) error {
	return context.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
	})
}

func ParseAndValidateRequestBody(context *fiber.Ctx, object interface{}) error {
	err := context.BodyParser(object)
	if err != nil {
		return err
	}

	v := validator.New()
	err = v.Struct(object)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			if e != nil {
				return e
			}
		}
	}

	return nil
}
