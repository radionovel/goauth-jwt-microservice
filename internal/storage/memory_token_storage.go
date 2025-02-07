package storage

import (
	"context"
	"sync"

	"github.com/radionovel/goauth-jwt-microservice/internal/model"
)

type MemoryTokenStorage struct {
	store map[model.Token]model.UserID
	mu    sync.Mutex
}

func NewInMemoryTokenStorage() *MemoryTokenStorage {
	return &MemoryTokenStorage{
		store: make(map[model.Token]model.UserID),
	}
}

func (s *MemoryTokenStorage) Save(ctx context.Context, userID model.UserID, token model.Token) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[token] = userID

	return nil
}

func (s *MemoryTokenStorage) Invalidate(ctx context.Context, token model.Token) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.store, token)

	return nil
}

func (s *MemoryTokenStorage) Find(ctx context.Context, token model.Token) (model.Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.store[token]
	if !exists {
		return model.Token{}, model.ErrTokenNotFound
	}

	return token, nil
}
