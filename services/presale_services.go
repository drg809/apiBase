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
	diff := time.Since(dbPresale.CreatedAt)
	if diff.Hours() < 720 {
		return nil, findResult.Error
	} else {
		claimResponse := new(models.SetUserClaimResponse)
		if !dbPresale.ClaimedFirst && !dbPresale.ClaimedSecond {
			claimResponse.ClaimedAmount = uint((float64(35) / float64(100)) * float64(dbPresale.Donated))
			claimResponse.Status = 1
			updateResult := database.GormDB.Model(&dbPresale).Updates(models.Presale{ClaimedFirst: true})
			if updateResult.Error != nil {
				return nil, findResult.Error
			}
		}
		if diff.Hours() < 1440 && !dbPresale.ClaimedSecond {
			claimResponse.ClaimedAmount = claimResponse.ClaimedAmount + uint((float64(35)/float64(100))*float64(dbPresale.Donated))
			claimResponse.Status = 1
			updateResult := database.GormDB.Model(&dbPresale).Updates(models.Presale{ClaimedSecond: true})
			if updateResult.Error != nil {
				return nil, findResult.Error
			}
		}
		return claimResponse, nil
	}

}
