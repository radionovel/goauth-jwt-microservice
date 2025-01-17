package handler

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/radionovel/goauth-jwt-microservice/internal/middleware"
	"github.com/radionovel/goauth-jwt-microservice/internal/model"
	"github.com/radionovel/goauth-jwt-microservice/internal/pkg/logger"
	"github.com/radionovel/goauth-jwt-microservice/internal/service"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
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

type ValidateTokenRequest struct{}

type ValidateResponse struct {
	Status string `json:"status"`
}

type AuthHandler struct {
	authSvc *service.AuthService
	userSvc *service.UserService

	logger logger.Logger
}

func NewAuthHandler(svc *service.AuthService, logger logger.Logger) *AuthHandler {
	return &AuthHandler{
		authSvc: svc,
		logger:  logger,
	}
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*TokenResponse, error) {
	h.logger.Debug("received refresh token request", "request", req)

	tokens, err := h.authSvc.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		h.logger.Error(err.Error())
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *LoginRequest) (*TokenResponse, error) {
	h.logger.Debug("received login request", "request", req)

	creds := &model.Credentials{
		Username: req.Username,
		Password: req.Password,
	}

	tokens, err := h.authSvc.LoginWithCredentials(ctx, creds)
	if err != nil {
		return nil, err
	}

	return FromTokensToResponse(tokens), nil
}

func (h *AuthHandler) Register(ctx context.Context, req *RegisterRequest) (*TokenResponse, error) {
	tokens, err := h.authSvc.Register(ctx, &model.Credentials{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		h.logger.Error(err.Error())
		return nil, err
	}

	return FromTokensToResponse(tokens), nil
}

func (h *AuthHandler) ValidateToken(ctx context.Context, _ *ValidateTokenRequest) (*ValidateResponse, error) {
	return &ValidateResponse{}, nil
}

func (h *AuthHandler) RegisterRoutes(router *echo.Echo) {
	router.POST("/api/v1/login", middleware.WrapHandler(h.Login))
	router.POST("/api/v1/register", middleware.WrapHandler(h.Register))
	router.POST("/api/v1/refresh", middleware.WrapHandler(h.RefreshToken))

	router.POST("/api/v1/validate", middleware.WrapHandler(h.ValidateToken), middleware.AuthMiddleware(h.authSvc))
}
