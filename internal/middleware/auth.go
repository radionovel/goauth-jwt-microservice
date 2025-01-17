package middleware

import (
	"context"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/radionovel/goauth-jwt-microservice/internal/model"
)

type TokenValidator interface {
	ValidateToken(token string) (string, error)
}

func AuthMiddleware(validator TokenValidator) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			h := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(h, "Bearer ") {
				return echo.ErrUnauthorized
			}

			token := strings.TrimPrefix(h, "Bearer ")
			sub, err := validator.ValidateToken(token)
			if err != nil {
				return echo.ErrUnauthorized
			}

			ctxWithUser := context.WithValue(c.Request().Context(), model.UserContextKey, sub)
			reqWithUser := c.Request().WithContext(ctxWithUser)
			c.SetRequest(reqWithUser)

			return next(c)
		}
	}
}
