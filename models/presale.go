package models

import (
	"github.com/nikola43/fibergormapitemplate/models/base"
)

type Presale struct {
	base.CustomGormModel
	UserID        uint `gorm:"type:INTEGER not null" json:"users_id"`
	Donated       uint `gorm:"type:INTEGER not null" json:"donated_amount,omitempty"`
	ClaimedFirst  bool `gorm:"type:bool" json:"claimed_first,omitempty"`
	ClaimedSecond bool `gorm:"type:bool" json:"claimed_second,omitempty"`
}
