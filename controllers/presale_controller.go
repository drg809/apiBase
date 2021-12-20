package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/drg809/apiBase/models"
	"github.com/drg809/apiBase/services"
	"github.com/drg809/apiBase/utils"
	"github.com/gofiber/fiber/v2"
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
	bnbValue, _ := utils.CallBSC(context)
	parseError := utils.ParseAndValidateRequestBody(context, insertPresaleRequest)
	if parseError != nil {
		return utils.ReturnErrorResponse(fiber.StatusBadRequest, parseError, context)
	}
	bnbValueF, _ := strconv.ParseFloat(bnbValue, 64)
	tokenAmount := (float64(insertPresaleRequest.Donated) * bnbValueF) / 0.5
	insertPresaleRequest.Donated = tokenAmount
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

func SetUserVesting(context *fiber.Ctx) error {
	dbUser := new(models.SetUserClaimRequest)
	parseError := utils.ParseAndValidateRequestBody(context, dbUser)
	if parseError != nil {
		return utils.ReturnErrorResponse(fiber.StatusBadRequest, parseError, context)
	}
	userLogged, _ := utils.GetUserTokenClaims(context)
	updatedUser, err := services.SetUserVesting(dbUser, userLogged)
	if err != nil {
		return utils.ReturnErrorResponse(fiber.StatusForbidden, err, context)
	}

	return context.Status(fiber.StatusOK).JSON(updatedUser)
}

// 0.1 BNB compra minima
// 5 BNB compra máxima

func CallBSC(context *fiber.Ctx) error {
	fmt.Println("call")

	type Data struct {
		Result struct {
			Ethusd string
		}
	}
	resp, err := http.Get("https://api.bscscan.com/api?module=stats&action=bnbprice&apikey=UMKZDMNWZE1PTPD4JVUUUXN7WGNR1FWZJW")
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	d := Data{}
	json.Unmarshal([]byte(body), &d)
	fmt.Println("BNB Value:", d.Result.Ethusd)
	fmt.Println(string(body))
	return utils.ReturnSuccessResponse(context)
}
