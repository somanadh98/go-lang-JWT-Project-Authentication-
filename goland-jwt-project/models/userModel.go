package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `json:"_id"`
	FirstName     *string            `json:"first_name" validate:"required,min=2,max=100"`
	LastName      *string            `json:"last_name" validate:"required,min=2,max=100"`
	Email         *string            `json:"email" validate:"required,email"`
	Password      *string            `json:"password" validate:"required,min=6"`
	Token         *string            `json:"token"`
	UserType      *string            `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	Created_At    time.Time          `json:"created_at"`
	Updated_At    time.Time          `json:"updated_at"`
	User_id       *string            `json:"user_id"`
	Phone         *string            `json:"phone" validate:"required"`
	Refresh_Token *string            `json:"refresh_token"`
}
