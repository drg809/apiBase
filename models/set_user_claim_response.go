package models

type SetUserClaimResponse struct {
	PresaleID     uint    `json:"presale_id"`
	Status        uint    `json:"status"`
	ClaimedAmount float64 `json:"claimed_amount"`
}
