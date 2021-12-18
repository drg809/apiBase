package models

// swagger:parameters LoginUserRequest
type LoginUserRequest struct {
	// User registered email
	// in: body
	// required: true
	// example: pauloxti@gmail.com
	Email    string `json:"email" xml:"email" form:"email"`

	// User plain password
	// in: body
	// required: true
	// example: paulo
	Password string `json:"password" xml:"password" form:"password"`
}
