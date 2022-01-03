package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
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
	updatedUser, err := services.SetUserClaim(dbUser)
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

	tokenRespose.TokenAmount = (dbOracle.LastPriceRead * tokenRespose.BnbAmount) / 0.05
	tokenRespose.LastPrice = dbOracle.LastPriceRead
	tokenRespose.LastRead = dbOracle.LastTimeRead
	return context.Status(fiber.StatusOK).JSON(tokenRespose)
}

func SendTokens(context *fiber.Ctx) error {
	userLogged, _ := utils.GetUserTokenClaims(context)
	if userLogged.Address == "0xb6e76628beb7872d2ade6ae9641bb390401c18ef" {
		dbPresales, err := services.GetPresales()
		if err != nil {
			return utils.ReturnErrorResponse(fiber.StatusNotFound, err, context)
		}
		for i := 0; i < len(dbPresales); i++ {
			dbUser, err := services.GetUserByID(uint64(dbPresales[i].UserID))
			if err == nil {
				presale := new(models.SetUserClaimRequest)
				presale.PresaleID = dbPresales[i].ID
				dbPresale, err := services.SetUserClaim(presale)
				if err == nil && dbPresale.ClaimedAmount > 0 {
					value := utils.EtherToWei(big.NewFloat(dbPresale.ClaimedAmount))
					res, err := SendTransaction(dbUser.WalletAddress, value.String())
					if err != nil {
						fmt.Println("err send")
						return utils.ReturnErrorResponse(fiber.StatusNotFound, err, context)
					}
					if !res {
						fmt.Println("fail send")
						return utils.ReturnErrorResponse(fiber.StatusNotFound, err, context)
					} else {
						fmt.Println("user claim")
						services.UpdateUserClaim(dbPresale)
					}
				}

			}

		}

		return context.Status(fiber.StatusOK).JSON(dbPresales)
	}
	return context.Status(fiber.StatusUnauthorized).JSON(userLogged)
}

func SendTransaction1(address string, tokenAmount string) {}

func SendTransaction(address string, tokenAmount string) (bool, error) {
	return true, nil
	// client, err := ethclient.Dial("https://infura.io")
	// if err != nil {
	// 	log.Fatal(err)
	// 	return err
	// }
	// chainID, err := client.NetworkID(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// 	return err
	// }
	// fmt.Println("===========")
	// fmt.Println(chainID)
	// fmt.Println(address)
	// fmt.Println(tokenAmount)
	// fmt.Println("===========")
	// privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	// if err != nil {
	// 	log.Fatal(err)
	// 	return err
	// }

	// publicKey := privateKey.Public()
	// _, ok := publicKey.(*ecdsa.PublicKey)
	// if !ok {
	// 	log.Fatal("error casting public key to ECDSA")
	// 	return err
	// }

	// return nil
	// fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	// nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// value := big.NewInt(0) // in wei (0 eth)
	// gasPrice, err := client.SuggestGasPrice(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// toAddress := common.HexToAddress(address)
	// tokenAddress := common.HexToAddress("0x28b149020d2152179873ec60bed6bf7cd705775d")

	// transferFnSignature := []byte("transfer(address,uint256)")
	// hash := sha3.NewLegacyKeccak256()
	// hash.Write(transferFnSignature)
	// methodID := hash.Sum(nil)[:4]
	// fmt.Println(hexutil.Encode(methodID)) // 0xa9059cbb

	// paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	// fmt.Println(hexutil.Encode(paddedAddress)) // 0x0000000000000000000000004592d8f8d7b001e72cb26a73e4fa1806a51ac79d

	// amount := new(big.Int)
	// amount.SetString(tokenAmount, 10) // 1000 tokens
	// paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	// fmt.Println(hexutil.Encode(paddedAmount)) // 0x00000000000000000000000000000000000000000000003635c9adc5dea00000

	// var data []byte
	// data = append(data, methodID...)
	// data = append(data, paddedAddress...)
	// data = append(data, paddedAmount...)

	// gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
	// 	To:   &toAddress,
	// 	Data: data,
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(gasLimit) // 23256

	// tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// _, err = types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = client.SendTransaction(context.Background(), signedTx)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("tx sent: %s", signedTx.Hash().Hex()) // tx sent: 0xa56316b637a94c4cc0331c73ef26389d6c097506d581073f927275e7a6ece0bc
}
