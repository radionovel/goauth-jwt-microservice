package service

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/radionovel/goauth-jwt-microservice/internal/model"
)

type jwtClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

type AuthServiceConfig struct {
	SecretKey       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type AuthService struct {
	secretKey       string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewAuthService(cfg AuthServiceConfig) *AuthService {
	return &AuthService{
		secretKey:       cfg.SecretKey,
		accessTokenTTL:  cfg.AccessTokenTTL,
		refreshTokenTTL: cfg.RefreshTokenTTL,
	}
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*model.Tokens, error) {
	claims, err := s.parseToken(refreshToken, s.secretKey)
	if err != nil {
		return nil, err
	}

	return s.GenerateTokens(claims.UserID)
}

func (s *AuthService) LoginWithCredentials(ctx context.Context, credentials *model.Credentials) (*model.Tokens, error) {
	user := &model.User{
		ID:       100,
		Username: "Jon Doe",
	}

	return s.GenerateTokens(user.ID)
}

func (s *AuthService) Register(ctx context.Context, credentials *model.Credentials) (*model.Tokens, error) {
	user := &model.User{
		ID:       100,
		Username: credentials.Username,
	}

	return s.GenerateTokens(user.ID)
}

func (s *AuthService) ValidateToken(tokenString string) (string, error) {
	claims, err := s.parseToken(tokenString, s.secretKey)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(claims.UserID), err
}

func (s *AuthService) GenerateTokens(userID int) (*model.Tokens, error) {
	accessToken, err := s.generateToken(userID, s.secretKey, s.accessTokenTTL)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateToken(userID, s.secretKey, s.refreshTokenTTL)
	if err != nil {
		return nil, err
	}

	return &model.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) generateToken(userID int, secretKey string, ttl time.Duration) (string, error) {
	claims := &jwtClaims{
		userID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s *AuthService) parseToken(tokenString, secretKey string) (*jwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, model.ErrUnexpectedSigningMethod
		}

		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, model.ErrInvalidToken
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok {
		slog.Error("invalid token claims")

		return nil, model.ErrInvalidClaims
	}

	slog.Info("parsed token", "claims", claims)

	return claims, nil
}
