package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/radionovel/goauth-jwt-microservice/internal/model"
)

type JWTManager struct {
	SecretKey       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewJWTManager(secretKey string, accessTTL, refreshTTL time.Duration) *JWTManager {
	return &JWTManager{
		SecretKey:       secretKey,
		AccessTokenTTL:  accessTTL,
		RefreshTokenTTL: refreshTTL,
	}
}

func (j *JWTManager) GenerateAccessToken(userID model.UserID) (*model.Token, error) {
	return j.generateToken(userID, model.AccessToken, j.AccessTokenTTL)
}

func (j *JWTManager) GenerateRefreshToken(userID model.UserID) (*model.Token, error) {
	return j.generateToken(userID, model.RefreshToken, j.RefreshTokenTTL)
}

// Общая логика генерации токена
func (j *JWTManager) generateToken(userID model.UserID, tokenType model.TokenType, ttl time.Duration) (*model.Token, error) {
	claims := &jwtClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return nil, err
	}

	return &model.Token{
		Value:     signedToken,
		Type:      tokenType,
		ExpiresIn: ttl,
	}, nil
}
