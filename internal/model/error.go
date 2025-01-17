package model

import (
	"errors"
)

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidToken            = errors.New("invalid token")
	ErrInvalidClaims           = errors.New("claims has invalid format")
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)
