package models

import (
	"errors"
)

type User struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Email    string `json:"email" gorm:"not null;unique"`
	Password string `json:"password" gorm:"not null"`
}

type UserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate TODO validate email with regex
func (u *UserInput) Validate() error {
	if u.Email == "" {
		return errors.New("not valid user email")
	}
	if u.Password == "" {
		return errors.New("not user password")
	}

	return nil
}
