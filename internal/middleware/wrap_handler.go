package middleware

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

func WrapHandler[Req, Resp any](f func(ctx context.Context, req Req) (Resp, error)) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req Req
		if err := c.Bind(&req); err != nil {
			return err
		}

		res, err := f(c.Request().Context(), req)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, res)
	}
}
