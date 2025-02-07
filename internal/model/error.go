package model

import (
	"errors"
)

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidToken            = errors.New("invalid token")
	ErrTokenNotFound           = errors.New("token not found")
	ErrInvalidClaims           = errors.New("claims has invalid format")
)

var (
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)
