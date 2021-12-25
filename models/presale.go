package models

import (
	"github.com/drg809/apiBase/models/base"
)

type Presale struct {
	base.CustomGormModel
	TxHash        string  `gorm:"unique" json:"tx_hash"`
	UserID        uint    `gorm:"type:INTEGER not null" json:"users_id"`
	DonatedAmount float64 `json:"donated_amount,omitempty"`
	TokenAmount   float64 `json:"token_amount,omitempty"`
	ClaimedAmount float64 `json:"claimed_amount,omitempty"`
	ClaimedFirst  bool    `gorm:"type:bool" json:"claimed_first,omitempty"`
	ClaimedSecond bool    `gorm:"type:bool" json:"claimed_second,omitempty"`
	ClaimedThird  bool    `gorm:"type:bool" json:"claimed_third,omitempty"`
}
