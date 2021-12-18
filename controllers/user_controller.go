package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/nikola43/fibergormapitemplate/models"
	"github.com/nikola43/fibergormapitemplate/services"
	"github.com/nikola43/fibergormapitemplate/utils"
)

func GetUserByID(context *fiber.Ctx) error {
	id, _ := strconv.ParseUint(context.Params("id"), 10, 64)

	dbUser, err := services.GetUserByID(id)
	if err != nil {
		return utils.ReturnErrorResponse(fiber.StatusNotFound, err, context)
	}

	return context.Status(fiber.StatusOK).JSON(dbUser)
}

func UpdateUserByID(context *fiber.Ctx) error {
	dbUser := new(models.User)
	parseError := utils.ParseAndValidateRequestBody(context, dbUser)
	if parseError != nil {
		return utils.ReturnErrorResponse(fiber.StatusBadRequest, parseError, context)
	}
	updatedUser, err := services.UpdateUserByID(dbUser)
	if err != nil {
		return utils.ReturnErrorResponse(fiber.StatusNotFound, err, context)
	}

	return context.Status(fiber.StatusOK).JSON(updatedUser)
}

func DeleteUserByID(context *fiber.Ctx) error {
	id, _ := strconv.ParseUint(context.Params("id"), 10, 64)

	err := services.DeleteUserByID(id)
	if err != nil {
		return utils.ReturnErrorResponse(fiber.StatusNotFound, err, context)
	}

	return utils.ReturnSuccessResponse(context)
}

func RefreshUser(context *fiber.Ctx) error {
	id, _ := strconv.ParseUint(context.Params("id"), 10, 64)

	dbUser, err := services.RefreshUser(uint(id))
	if err != nil {
		return utils.ReturnErrorResponse(fiber.StatusBadRequest, err, context)
	}

	return context.Status(fiber.StatusOK).JSON(dbUser)
}
