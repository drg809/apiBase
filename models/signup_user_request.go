package models

type SignupUserRequest struct {
	WalletAddress string `json:"wallet_address" xml:"wallet_address" form:"wallet_address" validate:"required"`
}
