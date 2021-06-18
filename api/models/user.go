package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email string `json:"email"`
	Fullname string `json:"full_name"`
	Password string `json:"password"`
}

func (user *User) CreateUser() (*User, error) {
	return user, nil
}


