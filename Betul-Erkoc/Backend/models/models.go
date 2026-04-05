package models

import "time"

type User struct {
	ID              string    `json:"_id" bson:"_id,omitempty"`
	Email           string    `json:"email" bson:"email"`
	Password        string    `json:"password" bson:"password"`
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

type ProfileUpdateInput struct {
	FirstName                string                   `json:"firstName"`
	LastName                 string                   `json:"lastName"`
	CommunicationPreferences CommunicationPreferences `json:"communicationPreferences"`
}

type CommunicationPreferences struct {
	Newsletter       bool `json:"newsletter"`
	SmsNotifications bool `json:"smsNotifications"`
}

type PasswordUpdateInput struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}
