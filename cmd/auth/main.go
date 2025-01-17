package main

import (
	"fmt"
	"log"
	"log/slog"
	"time"

	elog "github.com/labstack/gommon/log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/radionovel/goauth-jwt-microservice/internal/handler"
	"github.com/radionovel/goauth-jwt-microservice/internal/service"
	"github.com/radionovel/goauth-jwt-microservice/internal/storage"
)

func main() {
	router := echo.New()
	router.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
		LogLevel:  elog.ERROR,
	}))

	router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		Skipper: func(c echo.Context) bool {
			return c.Request().URL.Path == "/health"
		},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			fmt.Printf("REQUEST: uri: %v, status: %v\n", v.URI, v.Status)
			return nil
		},
	}))

	logger := slog.Default()

	// @todo config
	svcConfig := service.AuthServiceConfig{
		Storage:         storage.NewInMemoryUserStorage(logger),
		Logger:          logger,
		SecretKey:       "secret",
		AccessTokenTTL:  time.Hour,
		RefreshTokenTTL: time.Hour * 10,
	}

	svc := service.NewAuthService(svcConfig)
	h := handler.NewAuthHandler(svc, logger)
	h.RegisterRoutes(router)

	//userHandler := handler.NewUserHandler(svc)
	//userHandler.RegisterRoutes(router)

	log.Fatal(router.Start(":8080"))
}
