package service

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/radionovel/goauth-jwt-microservice/internal/model"
	"go.uber.org/zap"
)

type TokenStorage interface {
	Save(ctx context.Context, userID model.UserID, token model.Token) error
	Invalidate(ctx context.Context, token model.Token) error
	Find(ctx context.Context, token model.Token) (model.Token, error)
}

type jwtClaims struct {
	UserID model.UserID `json:"user_id"`
	jwt.RegisteredClaims
}

type TokenServiceConfig struct {
	JWTManager   *JWTManager
	TokenStorage TokenStorage
	Logger       *zap.Logger
}

type TokenService struct {
	jwtManager   *JWTManager
	tokenStorage TokenStorage
	logger       *zap.Logger
}

func NewTokenService(cfg TokenServiceConfig) *TokenService {
	return &TokenService{
		jwtManager:   cfg.JWTManager,
		tokenStorage: cfg.TokenStorage,
		logger:       cfg.Logger,
	}
}

func (s *TokenService) GenerateTokens(ctx context.Context, userID model.UserID) (*model.Tokens, error) {
	accessToken, err := s.jwtManager.GenerateAccessToken(userID)
	if err != nil {
		s.logger.Error("failed to generate access token", zap.Error(err))
		return nil, err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(userID)
	if err != nil {
		s.logger.Error("failed to generate refresh token", zap.Error(err))
		return nil, err
	}

	if err := s.tokenStorage.Save(ctx, userID, *accessToken); err != nil {
		s.logger.Error("failed to save access token", zap.Error(err))
		return nil, err
	}

	if err := s.tokenStorage.Save(ctx, userID, *refreshToken); err != nil {
		s.logger.Error("failed to save refresh token", zap.Error(err))
		return nil, err
	}

	s.logger.Info("Tokens successfully generated for user", zap.String("user_id", userID.String()))

	return &model.Tokens{
		AccessToken:  accessToken.Value,
		RefreshToken: refreshToken.Value,
	}, nil
}

func (s *TokenService) RevokeToken(ctx context.Context, refreshToken model.Token) (*jwtClaims, error) {
	return nil, nil
}

// ParseRefreshToken валидирует refresh_token и возвращает клеймы.
func (s *TokenService) ParseRefreshToken(refreshToken model.Token) (*jwtClaims, error) {
	claims, err := s.parseToken(refreshToken, s.secretKey)
	if err != nil {
		return nil, err
	}
	return claims, nil
}
