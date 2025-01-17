package model

import (
	"strconv"
)

type ContextKey string

const UserContextKey ContextKey = "user"

type NewUserDTO struct {
	Username     string
	PasswordHash string
	Salt         string
}

type UserID int

func (v UserID) Int() int {
	return int(v)
}

func (v UserID) String() string {
	return strconv.Itoa(int(v))
}

type User struct {
	ID           UserID
	Username     string
	PasswordHash string
	Salt         string
}
