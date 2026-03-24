package model

import "time"

type User struct {
	ID        int64     `db:"id" json:"id"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string      `json:"access_token"`
	User        UserSummary `json:"user"`
}

type UserSummary struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}
