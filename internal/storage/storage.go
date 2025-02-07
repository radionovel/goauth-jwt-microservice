package storage

import (
	"context"
	"time"

	"github.com/radionovel/goauth-jwt-microservice/internal/model"
)

type UserStorage interface {
	Insert(ctx context.Context, dto *model.NewUserDTO) (model.UserID, error)
	Exists(ctx context.Context, username string) (bool, error)
}

type TokenStorage interface {
	Save(ctx context.Context, userID model.UserID, token model.Token, tokenType model.TokenType, exp time.Duration) error
	Invalidate(ctx context.Context, token model.Token, tokenType model.TokenType) error
	Find(ctx context.Context, token model.Token, tokenType model.TokenType) (model.Token, error)
}
