package models

import (
	"github.com/drg809/apiBase/models/base"
)

type Presale struct {
	base.CustomGormModel
	UserID        uint `gorm:"type:INTEGER not null" json:"users_id"`
	Donated       uint `gorm:"type:INTEGER not null" json:"donated_amount,omitempty"`
	ClaimedAmount uint `gorm:"type:INTEGER not null" json:"claimed_amount,omitempty"`
	ClaimedFirst  bool `gorm:"type:bool" json:"claimed_first,omitempty"`
	ClaimedSecond bool `gorm:"type:bool" json:"claimed_second,omitempty"`
	ClaimedThird  bool `gorm:"type:bool" json:"claimed_third,omitempty"`
}
