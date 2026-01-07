package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID               bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId           string        `json:"user_id" bson:"user_id"`
	Email            string        `json:"email" bson:"email" validate:"required,email"`
	Password         string        `json:"password" bson:"password" validate:"required,min=6"`
	FirstName        string        `json:"first_name" bson:"first_name" validate:"required,min=2,max=100"`
	LastName         string        `json:"last_name" bson:"last_name" `
	Role             string        `json:"role" bson:"role" validate:"oneof USER ADMIN"`
	Token            string        `json:"token" bson:"token"`
	RefreshToken     string        `json:"refresh_token" bson:"refresh_token"`
	VerficationToken string        `json:"verification_token" bson:"verification_token"`
	EmailVerified    bool          `json:"email_verified" bson:"email_verified"`
	Avatar           string        `json:"avatar" bson:"avatar"`
	Karma            string        `json:"karma" bson:"karma"`
	ResetToken       string        `json:"reset_token" bson:"reset_token"`
	CreatedAt        time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at" bson:"updated_at"`
}

type UserLogin struct {
	Email    string `json:"email" bson:"email" validate:"required,email"`
	Password string `json:"password" bson:"password" validate:"required,min=6"`
}

type VefiryEmail struct {
	EmailVerified     bool   `json:"email_verified" bson:"email_verified"`
	VerificationToken string `json:"verification_token" bson:"verification_token"`
}

type UserDTO struct {
	UserId       string    `json:"user_id"`
	Email        string    `json:"email"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name" `
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" `
}

type UserLogout struct {
	UserId string `json:"user_id"`
}

type ForgetPasswordRequestDTO struct {
	Email      string `json:"email" bson:"email" validate:"required,email"`
	ResetToken string `json:"reset_token" bson:"reset_token"`
}

type ForgetPasswordDTO struct {
	Password string `json:"password" bson:"password" validate:"required,min=6"`
}

type UpdateAvatar struct {
	Avatar string `json:"avatar" bson:"avatar"`
}
