package services

import (
	"fmt"
	"time"

	database "github.com/drg809/apiBase/database"
	"github.com/drg809/apiBase/models"
	"github.com/drg809/apiBase/utils"
	"github.com/gofiber/fiber/v2"
)

func GetPresalesByUserID(userLogged *models.UserTokenClaims) (*models.UserPresales, error) {
	var dbPresales *models.UserPresales

	findResult := database.GormDB.Where("user_id = ?", userLogged.ID).Find(&dbPresales)
	if findResult.Error != nil {
		return nil, findResult.Error
	}

	return dbPresales, nil
}

func GetPresaleByID(id uint64, userLogged *models.UserTokenClaims) (*models.Presale, error) {
	dbPresale := new(models.Presale)

	findResult := database.GormDB.Where("user_id = ?", userLogged.ID).First(dbPresale, id)
	if findResult.Error != nil {
		return nil, findResult.Error
	}

	return dbPresale, nil
}

func InsertPresale(insertPresaleRequest *models.InsertPresaleRequest, userLogged *models.UserTokenClaims) error {
	dbPresale := new(models.Presale)
	fmt.Println(insertPresaleRequest)
	dbPresale.DonatedAmount = insertPresaleRequest.Donated
	dbPresale.TokenAmount = insertPresaleRequest.TokenAmount
	dbPresale.TxHash = insertPresaleRequest.TxHash
	dbPresale.UserID = userLogged.ID
	fmt.Println(dbPresale)
	insertResult := database.GormDB.Create(&dbPresale)
	if insertResult.Error != nil {
		return insertResult.Error
	}
	return nil
}

func SetUserClaim(userClaim *models.SetUserClaimRequest, userLogged *models.UserTokenClaims) (*models.SetUserClaimResponse, error) {
	dbPresale := new(models.Presale)
	findResult := database.GormDB.Where("user_id = ?", userLogged.ID).First(dbPresale, userClaim.PresaleID)
	if findResult.Error != nil {
		return nil, findResult.Error
	}
	claimResponse := new(models.SetUserClaimResponse)
	if dbPresale.ClaimedFirst {
		return claimResponse, findResult.Error
	} else {
		if !dbPresale.ClaimedFirst && (dbPresale.ClaimedAmount == 0) {
			claimResponse.ClaimedAmount = (float64(30) / float64(100)) * float64(dbPresale.DonatedAmount)
			claimResponse.Status = 1
			updateResult := database.GormDB.Model(&dbPresale).Updates(models.Presale{ClaimedFirst: true, ClaimedAmount: (claimResponse.ClaimedAmount)})
			if updateResult.Error != nil {
				return nil, findResult.Error
			}
		}
		return claimResponse, nil
	}

}

func SetUserVesting(userClaim *models.SetUserClaimRequest, userLogged *models.UserTokenClaims) (*models.SetUserClaimResponse, error) {
	dbPresale := new(models.Presale)
	findResult := database.GormDB.Where("user_id = ?", userLogged.ID).First(dbPresale, userClaim.PresaleID)
	if findResult.Error != nil {
		return nil, findResult.Error
	}
	claimResponse := new(models.SetUserClaimResponse)
	diff := time.Since(dbPresale.CreatedAt)
	if diff.Hours() < 720 {
		return claimResponse, findResult.Error
	} else {
		if dbPresale.ClaimedFirst && !dbPresale.ClaimedSecond && !dbPresale.ClaimedThird {
			claimResponse.ClaimedAmount = (float64(35) / float64(100)) * float64(dbPresale.DonatedAmount)
			claimResponse.Status = 1
			totalClaimed := claimResponse.ClaimedAmount + dbPresale.ClaimedAmount
			if totalClaimed > dbPresale.DonatedAmount {
				return nil, findResult.Error
			}
			updateResult := database.GormDB.Model(&dbPresale).Updates(models.Presale{ClaimedSecond: true, ClaimedAmount: totalClaimed})
			if updateResult.Error != nil {
				return nil, findResult.Error
			}
		}
		if diff.Hours() > 1440 && !dbPresale.ClaimedThird && (dbPresale.ClaimedAmount < dbPresale.DonatedAmount) {
			claimResponse.ClaimedAmount += (float64(35) / float64(100)) * float64(dbPresale.DonatedAmount)
			totalClaimed := float64(0)
			if claimResponse.Status == 1 {
				totalClaimed = claimResponse.ClaimedAmount + (float64(30)/float64(100))*float64(dbPresale.DonatedAmount)
			} else {
				totalClaimed = claimResponse.ClaimedAmount + dbPresale.ClaimedAmount
			}
			if totalClaimed > dbPresale.DonatedAmount {
				return nil, findResult.Error
			}
			claimResponse.Status = 1
			updateResult := database.GormDB.Model(&dbPresale).Updates(models.Presale{ClaimedThird: true, ClaimedAmount: totalClaimed})
			if updateResult.Error != nil {
				return nil, findResult.Error
			}
		}
		return claimResponse, nil
	}

}

func CallBSC(context *fiber.Ctx) error {
	fmt.Println("call")

	return utils.ReturnSuccessResponse(context)
}

func InsertOracleEntrie(insertOracleEntrie *models.Oracle) (*models.Oracle, error) {
	dbOracle := new(models.Oracle)
	fmt.Println(insertOracleEntrie)
	dbOracle.LastPriceRead = insertOracleEntrie.LastPriceRead
	dbOracle.LastTimeRead = time.Now().Unix()
	fmt.Println(dbOracle)
	insertResult := database.GormDB.Create(&dbOracle)
	if insertResult.Error != nil {
		return nil, insertResult.Error
	}
	return dbOracle, nil
}

func GetLastOracleRead() (*models.Oracle, error) {
	dbPresale := new(models.Oracle)

	findResult := database.GormDB.Last(dbPresale)
	if findResult.Error != nil {
		return nil, findResult.Error
	}

	return dbPresale, nil
}
