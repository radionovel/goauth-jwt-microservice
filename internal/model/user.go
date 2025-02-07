package model

import (
	"strconv"
	"time"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type Token struct {
	Value     string        // Сам токен
	Type      TokenType     // Тип токена (например, Access/Refresh)
	ExpiresIn time.Duration // Время жизни токена
}

type ContextKey string

const UserContextKey ContextKey = "user"

type NewUserDTO struct {
	Username     string
	PasswordHash string
	Salt         string
}

type UserID string

func (uid UserID) Int() int {
	id, err := strconv.Atoi(uid.String())
	if err != nil {
		return 0
	}

	return id
}

func (uid UserID) String() string {
	return string(uid)
}

type User struct {
	ID           UserID
	Username     string
	PasswordHash string
	Salt         string
}
