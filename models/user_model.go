package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id        primitive.ObjectID  `json:"id,omitempty"`
	Email     string              `json:"email,omitempty" validate:"required" bson:"email"`
	Password  string              `json:"password,omitempty" validate:"required"`
	Status    int                 `json:"status,omitempty"`
	CreatedAt primitive.Timestamp `json:"created_at,omitempty"`
	UpdatedAt primitive.Timestamp `json:"updated_at,omitempty"`
	Roles     [6]string           `json:"roles"`
	Profile   Profile             `json:"profile"`
}

type Profile struct {
	Age      int    `json:"age"`
	Country  string `json:"country"`
	Image    string `json:"image"`
	FullName string `json:"full_name"`
}

type UserPost struct {
	Email     string              `json:"email,omitempty" validate:"required"`
	Password  string              `json:"password,omitempty" validate:"required"`
	Status    int                 `json:"status,omitempty"`
	CreatedAt primitive.Timestamp `json:"created_at,omitempty"`
	UpdatedAt primitive.Timestamp `json:"updated_at,omitempty"`
	Roles     [6]string           `json:"roles"`
	Age       int                 `json:"age"`
	Country   string              `json:"country"`
	Image     string              `json:"image"`
	FullName  string              `json:"full_name"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
