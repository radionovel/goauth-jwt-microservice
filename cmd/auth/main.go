package main

import (
	"context"
	"time"

	"github.com/radionovel/goauth-jwt-microservice/internal/logger"
	"github.com/radionovel/goauth-jwt-microservice/internal/model"
	"github.com/radionovel/goauth-jwt-microservice/internal/service"
	"github.com/radionovel/goauth-jwt-microservice/internal/storage"
	"go.uber.org/zap"
)

func main() {
	jwtManager := service.NewJWTManager("my-secret-key", 15*time.Minute, 7*24*time.Hour)
	tokenStorage := storage.NewInMemoryTokenStorage() // Предположим, есть реализация InMemory
	logger := logger.NewConsoleLogger()               // Примерный логгер

	cfgTokenService := service.TokenServiceConfig{
		JWTManager:   jwtManager,
		TokenStorage: tokenStorage,
		Logger:       logger,
	}
	tokenService := service.NewTokenService(cfgTokenService)

	ctx := context.Background()
	userID := model.UserID("12345")

	tokens, err := tokenService.GenerateTokens(ctx, userID)
	if err != nil {
		logger.Error("Failed to generate tokens", zap.Error(err))
		return
	}

	logger.Info("Access token", zap.Any("token", tokens.AccessToken))
	logger.Info("Refresh Token", zap.Any("token", tokens.RefreshToken))
}
