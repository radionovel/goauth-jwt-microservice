package handler

import (
	"github.com/radionovel/goauth-jwt-microservice/internal/model"
)

func FromTokensToResponse(tokens *model.Tokens) *TokenResponse {
	return &TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
}
