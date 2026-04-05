package auth

import "time"

type User struct {
	ID              string    `json:"_id" bson:"_id,omitempty"`
	Email           string    `json:"email" bson:"email"`
	Password        string    `json:"password" bson:"password"` // Şifre doğrulaması için lazım
	FirstName       string    `json:"firstName" bson:"firstName"`
	LastName        string    `json:"lastName" bson:"lastName"`
	IsEmailVerified bool      `json:"isEmailVerified" bson:"isEmailVerified"`
	CreatedOn       time.Time `json:"createdOn" bson:"createdOn"`
}

type RegisterInput struct {
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
