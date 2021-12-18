package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/nikola43/fibergormapitemplate/models"
	"github.com/nikola43/fibergormapitemplate/services"
	"github.com/nikola43/fibergormapitemplate/utils"
)

func SignupUser(context *fiber.Ctx) error {
	signupUserRequest := new(models.SignupUserRequest)

	parseError := utils.ParseAndValidateRequestBody(context, signupUserRequest)
	if parseError != nil {
		return utils.ReturnErrorResponse(fiber.StatusBadRequest, parseError, context)
	}

	fmt.Println(signupUserRequest)

	createUserResponse, signupUserError := services.SignupUser(signupUserRequest)
	if signupUserError != nil {
		return utils.ReturnErrorResponse(fiber.StatusBadRequest, signupUserError, context)
	}

	return context.JSON(createUserResponse)
}
