package storage

import (
	"context"

	"github.com/radionovel/goauth-jwt-microservice/internal/model"
)

type UserStorage interface {
	Insert(ctx context.Context, dto *model.NewUserDTO) (model.UserID, error)
	Exists(ctx context.Context, username string) (bool, error)
}

type TokenStorage interface {
	Insert(ctx context.Context, dto *model.NewUserDTO) (model.UserID, error)
	Exists(ctx context.Context, username string) (bool, error)
}.

