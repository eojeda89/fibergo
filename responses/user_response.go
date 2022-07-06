package responses

import (
	"fibergo/models"
)

type UserResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type SigninResponse struct {
	User  models.User `json:"user"`
	Token string      `json:"token"`
}
