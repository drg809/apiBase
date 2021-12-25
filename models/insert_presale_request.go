package models

type InsertPresaleRequest struct {
	TxHash      string  `json:"tx_hash"`
	Donated     float64 `json:"donated_amount"`
	TokenAmount float64 `json:"token_amount,omitempty"`
}
