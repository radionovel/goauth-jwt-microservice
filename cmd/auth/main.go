package main

import (
	"github.com/labstack/echo/v4"
	"github.com/radionovel/goauth-jwt-microservice/internal/handler"
	"github.com/radionovel/goauth-jwt-microservice/internal/service"
	"log"
	"time"
)

func main() {
	router := echo.New()

	svcConfig := service.AuthServiceConfig{
		SecretKey:       "secret",
		AccessTokenTTL:  time.Hour,
		RefreshTokenTTL: time.Hour * 10,
	}
	svc := service.NewAuthService(svcConfig)
	h := handler.NewAuthHandler(svc)
	h.RegisterRoutes(router)

	userHandler := handler.NewUserHandler(svc)
	userHandler.RegisterRoutes(router)

	log.Fatal(router.Start(":8080"))
}
