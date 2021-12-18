package services

import (
	"errors"

	database "github.com/nikola43/fibergormapitemplate/database"
	"github.com/nikola43/fibergormapitemplate/models"
	"github.com/nikola43/fibergormapitemplate/utils"
)

func GetUserByID(id uint64) (*models.User, error) {
	dbUser := new(models.User)

	findResult := database.GormDB.First(dbUser, id)
	if findResult.Error != nil {
		return nil, findResult.Error
	}

	return dbUser, nil
}

func UpdateUserByID(updatedUser *models.User) (*models.User, error) {
	dbUser := new(models.User)
	findResult := database.GormDB.First(dbUser, updatedUser.ID)
	if findResult.Error != nil {
		return nil, findResult.Error
	}
	updateResult := database.GormDB.Model(&dbUser).Updates(models.User{Username: updatedUser.Username})
	if updateResult.Error != nil {
		return nil, findResult.Error
	}

	return dbUser, nil
}

func DeleteUserByID(id uint64) error {

	dbUser := new(models.User)

	findResult := database.GormDB.First(dbUser, id)
	if findResult.Error != nil {
		return findResult.Error
	}

	deleteResult := database.GormDB.Model(&dbUser).Delete(dbUser)
	if deleteResult.Error != nil {
		return deleteResult.Error
	}

	return nil
}

func RefreshUser(id uint) (*models.LoginUserResponse, error) {
	dbUser := new(models.User)

	findResult := database.GormDB.First(dbUser, id)
	if findResult.Error != nil {
		return nil, findResult.Error
	}

	if dbUser.ID < 1 {
		return nil, errors.New("user not found")
	}

	token, err := utils.GenerateUserToken(dbUser.WalletAddress, dbUser.ID)
	if err != nil {
		return nil, err
	}

	// Generate and send response
	loginUserResponse := &models.LoginUserResponse{
		ID:    dbUser.ID,
		Email: dbUser.WalletAddress,
		Token: token,
	}

	return loginUserResponse, nil
}
