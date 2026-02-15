package model

import "time"

type User struct {
	ID           uint64
	Login        string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserRegister struct {
	Login    string
	Email    string
	Password string
}

type UserLogin struct {
	LoginOrEmail string
	Password     string
}
