package models

type CalcTokenQuantityResponse struct {
	TokenAmount float64 `json:"token_amount,omitempty"`
	BnbAmount   float64 `json:"donated_amount"`
	LastPrice   float64 `json:"last_price,omitempty"`
	LastRead    int64   `json:"last_read,omitempty"`
}
