package models

// LoginUserResponse
//
// Used for return login result to front
//
// swagger:model LoginUserResponse
type LoginUserResponse struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}
