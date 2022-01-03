package controllers

import (
	"github.com/drg809/apiBase/models"
	"github.com/drg809/apiBase/services"
	"github.com/drg809/apiBase/utils"
	"github.com/gofiber/fiber/v2"
)

func SignupUser(context *fiber.Ctx) error {
	signupUserRequest := new(models.SignupUserRequest)

	parseError := utils.ParseAndValidateRequestBody(context, signupUserRequest)
	if parseError != nil {
		return utils.ReturnErrorResponse(fiber.StatusBadRequest, parseError, context)
	}

	createUserResponse, signupUserError := services.SignupUser(signupUserRequest)
	if signupUserError != nil {
		return utils.ReturnErrorResponse(fiber.StatusBadRequest, signupUserError, context)
	}

	return context.JSON(createUserResponse)
}
