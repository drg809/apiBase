package models

import (
	"github.com/drg809/apiBase/models/base"
)

type User struct {
	base.CustomGormModel
	WalletAddress string `gorm:"index; unique; type:varchar(64) not null" json:"wallet_address" xml:"wallet_address" form:"wallet_address"`
	Username      string `gorm:"index; unique; type:varchar(64) not null" json:"username" xml:"username" form:"username"`
	AvatarUrl     string `gorm:"index; unique; type:varchar(255) not null" json:"avatar_url" xml:"avatar_url" form:"avatar_url"`
}
