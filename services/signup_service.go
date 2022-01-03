package services

import (
	database "github.com/drg809/apiBase/database"
	"github.com/drg809/apiBase/models"
	"github.com/drg809/apiBase/utils"
)

func SignupUser(signupUserRequest *models.SignupUserRequest) (*models.SignupUserResponse, error) {

	bdUser := new(models.User)
	findResult := database.GormDB.Where("wallet_address = ?", signupUserRequest.WalletAddress).Find(&bdUser)
	if findResult.Error != nil {
		return nil, findResult.Error
	}
	token, err := utils.GenerateUserToken(bdUser.WalletAddress, bdUser.ID)

	if err != nil {
		return nil, findResult.Error
	}
	if bdUser.ID > 0 {
		return &models.SignupUserResponse{
			Token: token,
		}, nil
	}

	bdUser.Username = "USER_" + signupUserRequest.WalletAddress
	bdUser.WalletAddress = signupUserRequest.WalletAddress
	bdUser.AvatarUrl = "url" + signupUserRequest.WalletAddress
	insertUserResult := database.GormDB.Create(&bdUser)
	if insertUserResult.Error != nil {
		return nil, insertUserResult.Error
	}

	token2, err2 := utils.GenerateUserToken(bdUser.WalletAddress, bdUser.ID)
	createUserResponse := &models.SignupUserResponse{
		Token: token2,
	}
	if err2 != nil {
		return nil, err2
	}

	return createUserResponse, nil
}
