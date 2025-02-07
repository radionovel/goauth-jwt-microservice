package storage

import (
	"context"
	"sync"

	"github.com/radionovel/goauth-jwt-microservice/internal/model"
	"github.com/radionovel/goauth-jwt-microservice/internal/pkg/logger"
)

type inMemoryUserStorage struct {
	logger     logger.Logger
	users      map[string]model.User
	mu         sync.RWMutex
	lastUserID model.UserID
}

func NewInMemoryUserStorage(logger logger.Logger) UserStorage {
	return &inMemoryUserStorage{
		logger: logger,
		users:  make(map[string]model.User),
	}
}

func (s *inMemoryUserStorage) Insert(ctx context.Context, dto *model.NewUserDTO) (model.UserID, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// s.lastUserID++

	// if _, exists := s.users[dto.Username]; exists {
	// 	return 0, model.ErrUserAlreadyExists
	// }

	// s.users[dto.Username] = model.User{
	// 	ID:           s.lastUserID,
	// 	Username:     dto.Username,
	// 	PasswordHash: dto.PasswordHash,
	// 	Salt:         dto.Salt,
	// }

	return s.lastUserID, nil
}

func (s *inMemoryUserStorage) Exists(ctx context.Context, username string) (bool, error) {
	_, exists := s.users[username]
	return exists, nil
}
