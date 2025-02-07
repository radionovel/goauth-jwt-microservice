package service

import (
	"context"

	"github.com/radionovel/goauth-jwt-microservice/internal/model"
	"github.com/radionovel/goauth-jwt-microservice/internal/pkg/logger"
	"github.com/radionovel/goauth-jwt-microservice/internal/storage"
	"github.com/radionovel/goauth-jwt-microservice/internal/util"
)

type AuthServiceConfig struct {
	UserStorage  storage.UserStorage
	TokenStorage storage.TokenStorage
	Logger       logger.Logger
}

type AuthService struct {
	userStorage  storage.UserStorage
	tokenService *TokenService
	hasher       util.Hasher
	logger       logger.Logger
}

func NewAuthService(cfg AuthServiceConfig) *AuthService {
	return &AuthService{
		userStorage: cfg.UserStorage,
		logger:      cfg.Logger,
	}
}

func (a *AuthService) Login(ctx context.Context, login, password string) (*model.Tokens, error) {
	// Ищем пользователя по email
	user, err := a.userStorage.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, model.ErrUserNotFound
	}

	// Проверяем пароль
	if !a.hasher.Compare(user.Password, password) {
		return nil, model.ErrInvalidPassword
	}

	// Генерируем токены
	return a.tokenService.GenerateTokens(ctx, user.ID)
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken model.Token) (*model.Tokens, error) {
	// Валидируем refresh_token
	claims, err := s.tokenService.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, model.ErrInvalidToken
	}

	userID := claims.UserID

	// Генерируем новые токены
	newTokens, err := s.tokenService.GenerateTokens(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Аннулируем старый refresh_token
	err = s.tokenService.RevokeToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	return newTokens, nil
}

func (s *AuthService) Logout(ctx context.Context, accessToken model.Token) error {
	return nil
}
