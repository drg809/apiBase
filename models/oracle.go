package models

import (
	"github.com/drg809/apiBase/models/base"
)

type Oracle struct {
	base.CustomGormModel
	LastPriceRead float64 `json:"last_price_read"`
	LastTimeRead  int64   `json:"last_time_read"`
}
