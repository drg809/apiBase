package services

import (
	"time"

	database "github.com/nikola43/fibergormapitemplate/database"
	"github.com/nikola43/fibergormapitemplate/models"
)

func GetPresaleByUserID(id uint64, userLogged *models.UserTokenClaims) (*models.Presale, error) {
	dbPresale := new(models.Presale)

	findResult := database.GormDB.Where("user_id = ?", userLogged.ID).First(dbPresale, id)
	if findResult.Error != nil {
		return nil, findResult.Error
	}

	return dbPresale, nil
}

func InsertPresale(insertPresaleRequest *models.InsertPresaleRequest, userLogged *models.UserTokenClaims) error {
	dbPresale := new(models.Presale)
	dbPresale.Donated = insertPresaleRequest.Donated
	dbPresale.UserID = userLogged.ID
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
	diff := time.Since(dbPresale.CreatedAt)
	if diff.Hours() < 720 {
		return claimResponse, findResult.Error
	} else {
		if !dbPresale.ClaimedFirst && !dbPresale.ClaimedSecond && (dbPresale.ClaimedAmount == 0) {
			claimResponse.ClaimedAmount = uint((float64(35) / float64(100)) * float64(dbPresale.Donated))
			claimResponse.Status = 1
			updateResult := database.GormDB.Model(&dbPresale).Updates(models.Presale{ClaimedFirst: true, ClaimedAmount: (claimResponse.ClaimedAmount + dbPresale.ClaimedAmount)})
			if updateResult.Error != nil {
				return nil, findResult.Error
			}
		}
		if diff.Hours() > 1440 && !dbPresale.ClaimedSecond && (dbPresale.ClaimedAmount < dbPresale.Donated) {
			claimResponse.ClaimedAmount += uint((float64(35) / float64(100)) * float64(dbPresale.Donated))
			totalClaimed := uint(0)
			if claimResponse.Status == 1 {
				totalClaimed = claimResponse.ClaimedAmount + uint((float64(30)/float64(100))*float64(dbPresale.Donated))
			} else {
				totalClaimed = claimResponse.ClaimedAmount + dbPresale.ClaimedAmount
			}
			claimResponse.Status = 1
			updateResult := database.GormDB.Model(&dbPresale).Updates(models.Presale{ClaimedSecond: true, ClaimedAmount: totalClaimed})
			if updateResult.Error != nil {
				return nil, findResult.Error
			}
		}
		return claimResponse, nil
	}

}
