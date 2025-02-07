package util

import (
	"golang.org/x/crypto/bcrypt"
)

// Интерфейс для хеширования паролей
type Hasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) bool
}

type BcryptHasher struct {
	cost int // параметр сложности
}

func NewBcryptHasher(cost int) *BcryptHasher {
	return &BcryptHasher{cost: cost}
}

func (h *BcryptHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (h *BcryptHasher) Compare(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
