package handler

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/radionovel/goauth-jwt-microservice/internal/middleware"
	"github.com/radionovel/goauth-jwt-microservice/internal/service"
)

type GetUserRequest struct{}

type GetUserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserHandler struct {
	authSvc *service.AuthService
}

func NewUserHandler(authSvc *service.AuthService) *UserHandler {
	return &UserHandler{authSvc: authSvc}
}

func (h *UserHandler) GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error) {
	return nil, nil
}

func (h *UserHandler) RegisterRoutes(router *echo.Echo) {
	protected := router.Group("", middleware.AuthMiddleware(h.authSvc))
	protected.GET("/api/v1/users", middleware.WrapHandler(h.GetUser))
}
