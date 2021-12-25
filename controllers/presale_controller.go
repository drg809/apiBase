package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/drg809/apiBase/models"
	"github.com/drg809/apiBase/services"
	"github.com/drg809/apiBase/utils"
	"github.com/gofiber/fiber/v2"
)

func GetPresalesByUserID(context *fiber.Ctx) error {
	userLogged, _ := utils.GetUserTokenClaims(context)

	dbUser, err := services.GetPresalesByUserID(userLogged)
	if err != nil {
		return utils.ReturnErrorResponse(fiber.StatusNotFound, err, context)
	}

	return context.Status(fiber.StatusOK).JSON(dbUser)
}

func InsertPresale(context *fiber.Ctx) error {
	insertPresaleRequest := new(models.InsertPresaleRequest)
	userLogged, _ := utils.GetUserTokenClaims(context)
	bnbValue, _ := GetLastOracleRead(context)
	parseError := utils.ParseAndValidateRequestBody(context, insertPresaleRequest)
	if parseError != nil {
		return utils.ReturnErrorResponse(fiber.StatusBadRequest, parseError, context)
	}
	tokenAmount := (float64(insertPresaleRequest.Donated) * float64(bnbValue.LastPriceRead)) / 0.5
	insertPresaleRequest.TokenAmount = tokenAmount
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

func CallBSC(context *fiber.Ctx) (*models.Oracle, error) {
	dbOracle := new(models.Oracle)
	fmt.Println("call")
	parseError := utils.ParseAndValidateRequestBody(context, dbOracle)
	if parseError != nil {
		return nil, parseError
	}

	type Data struct {
		Result struct {
			Ethusd string
		}
	}
	resp, err := http.Get("https://api.bscscan.com/api?module=stats&action=bnbprice&apikey=" + os.Getenv("BSCSCAN_KEY"))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	d := Data{}
	json.Unmarshal([]byte(body), &d)
	dbOracle.LastPriceRead, _ = strconv.ParseFloat(d.Result.Ethusd, 64)

	if err != nil {
		return nil, err
	}
	updatedOracle, err := services.InsertOracleEntrie(dbOracle)
	if err != nil {
		return nil, err
	}

	return updatedOracle, err
}

func GetLastOracleRead(context *fiber.Ctx) (*models.Oracle, error) {
	dbOracle, err := services.GetLastOracleRead()
	if err != nil {
		return nil, err
	}

	if dbOracle == nil || dbOracle.LastTimeRead < (time.Now().Unix()-5) {
		dbOracle, err = CallBSC(context)
		if err != nil {
			return nil, err
		}
	}
	return dbOracle, nil
}

func CalcTokenQuantity(context *fiber.Ctx) error {
	tokenRespose := new(models.CalcTokenQuantityResponse)
	parseError := utils.ParseAndValidateRequestBody(context, tokenRespose)
	if parseError != nil {
		return utils.ReturnErrorResponse(fiber.StatusBadRequest, parseError, context)
	}
	dbOracle, err := GetLastOracleRead(context)
	if err != nil {
		return utils.ReturnErrorResponse(fiber.StatusBadRequest, err, context)
	}

	tokenRespose.TokenAmount = (dbOracle.LastPriceRead * tokenRespose.BnbAmount) * 2
	tokenRespose.LastPrice = dbOracle.LastPriceRead
	tokenRespose.LastRead = dbOracle.LastTimeRead
	return context.Status(fiber.StatusOK).JSON(tokenRespose)
}
