package handler

import (
	"context"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/radionovel/goauth-jwt-microservice/internal/middleware"
	"github.com/radionovel/goauth-jwt-microservice/internal/model"
	"github.com/radionovel/goauth-jwt-microservice/internal/service"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*TokenResponse, error) {
	slog.Info("received refresh token request", "request", req)

	tokens, err := h.svc.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *LoginRequest) (*TokenResponse, error) {
	slog.Info("received login request", "request", req)

	foo := ctx.Value("user")
	slog.Info("received login request", "foo", foo)

	creds := &model.Credentials{
		Username: req.Username,
		Password: req.Password,
	}

	tokens, err := h.svc.LoginWithCredentials(ctx, creds)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (h *AuthHandler) RegisterRoutes(router *echo.Echo) {
	router.POST("/api/v1/login", middleware.WrapHandler(h.Login))
	router.POST("/api/v1/refresh", middleware.WrapHandler(h.RefreshToken))
}
