package models

type SetUserClaimResponse struct {
	Status        uint    `json:"status"`
	ClaimedAmount float64 `json:"claimed_amount"`
}
