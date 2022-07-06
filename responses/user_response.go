package responses

import (
	"fibergo/models"
)

// User model info
// @Description User account information
type UserResponse struct {
	Status  int         `json:"status"`  //this is the status of the response
	Message string      `json:"message"` //this is the message of the response
	Data    interface{} `json:"data"`    //this is the status of the response
}

// Signin model info
// @Description Signin model response information
type SigninResponse struct {
	User  models.User `json:"user"`  //this is the user that  has been signed in
	Token string      `json:"token"` //this is the generated token
}
