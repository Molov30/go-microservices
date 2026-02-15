package model

import "time"

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type RefreshToken struct {
	ID        uint64
	UserID    uint64
	Token     string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}
