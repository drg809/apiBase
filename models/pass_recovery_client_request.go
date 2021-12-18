package models

// RecoveryUserPasswordRequest
//
// Esta es la estructura usada para responder con frases
//
// swagger:model RecoveryUserPasswordRequest
type RecoveryUserPasswordRequest struct {
	Email    string `json:"email" xml:"email" form:"email"`
}
