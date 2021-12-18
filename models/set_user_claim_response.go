package models

type SetUserClaimResponse struct {
	Status        uint `json:"status"`
	ClaimedAmount uint `json:"claimed_amount"`
}
