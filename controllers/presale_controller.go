package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/nikola43/fibergormapitemplate/models"
	"github.com/nikola43/fibergormapitemplate/services"
	"github.com/nikola43/fibergormapitemplate/utils"
)

func GetPresaleByUserID(context *fiber.Ctx) error {
	id, _ := strconv.ParseUint(context.Params("id"), 10, 64)
	userLogged, _ := utils.GetUserTokenClaims(context)

	dbUser, err := services.GetPresaleByUserID(id, userLogged)
	if err != nil {
		return utils.ReturnErrorResponse(fiber.StatusNotFound, err, context)
	}

	return context.Status(fiber.StatusOK).JSON(dbUser)
}

func InsertPresale(context *fiber.Ctx) error {
	insertPresaleRequest := new(models.InsertPresaleRequest)
	userLogged, _ := utils.GetUserTokenClaims(context)
	parseError := utils.ParseAndValidateRequestBody(context, insertPresaleRequest)
	if parseError != nil {
		return utils.ReturnErrorResponse(fiber.StatusBadRequest, parseError, context)
	}
	err := services.InsertPresale(insertPresaleRequest, userLogged)
	if err != nil {
		return utils.ReturnErrorResponse(fiber.StatusBadRequest, err, context)
	}
	return utils.ReturnSuccessResponse(context)
}

func SetUserClaim(context *fiber.Ctx) error {
	dbUser := new(models.SetUserClaimRequest)
	parseError := utils.ParseAndValidateRequestBody(context, dbUser)
	if parseError != nil {
		return utils.ReturnErrorResponse(fiber.StatusBadRequest, parseError, context)
	}
	userLogged, _ := utils.GetUserTokenClaims(context)
	updatedUser, err := services.SetUserClaim(dbUser, userLogged)
	if err != nil {
		return utils.ReturnErrorResponse(fiber.StatusForbidden, err, context)
	}

	return context.Status(fiber.StatusOK).JSON(updatedUser)
}
