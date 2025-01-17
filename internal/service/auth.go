package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/radionovel/goauth-jwt-microservice/internal/model"
	"github.com/radionovel/goauth-jwt-microservice/internal/pkg/logger"
	"github.com/radionovel/goauth-jwt-microservice/internal/pkg/password"
	"github.com/radionovel/goauth-jwt-microservice/internal/storage"
)

type jwtClaims struct {
	UserID model.UserID `json:"user_id"`
	jwt.RegisteredClaims
}

type AuthServiceConfig struct {
	UserStorage  storage.UserStorage
	TokenStorage storage.TokenStorage
	Logger       logger.Logger

	SecretKey       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type AuthService struct {
	userStorage storage.UserStorage
	logger      logger.Logger

	secretKey       string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewAuthService(cfg AuthServiceConfig) *AuthService {
	return &AuthService{
		userStorage:     cfg.UserStorage,
		tokenStorage:    cfg.TokenStorage,
		logger:          cfg.Logger,
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

	//password.ComparePasswordHash()

	return s.GenerateTokens(user.ID)
}

func (s *AuthService) Register(ctx context.Context, credentials *model.Credentials) (*model.Tokens, error) {
	//s.logger.Debug("create new user", "username", credentials.Username)

	// find user with username
	exists, err := s.userStorage.Exists(ctx, credentials.Username)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, model.ErrUserAlreadyExists
	}

	salt := "asdfqwerzxcvtgby"
	passwordHash := password.Generate(credentials.Password, salt)

	userID, err := s.userStorage.Insert(ctx, &model.NewUserDTO{
		Username:     credentials.Username,
		PasswordHash: passwordHash,
		Salt:         salt,
	})
	if err != nil {
		return nil, err
	}

	// publish event to event broker via outbox

	return s.GenerateTokens(userID)
}

func (s *AuthService) ValidateToken(tokenString string) (string, error) {
	claims, err := s.parseToken(tokenString, s.secretKey)
	if err != nil {
		return "", err
	}

	return claims.UserID.String(), err
}

func (s *AuthService) GenerateTokens(userID model.UserID) (*model.Tokens, error) {
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

func (s *AuthService) generateToken(userID model.UserID, secretKey string, ttl time.Duration) (string, error) {
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

func (s *AuthService) RevokeToken(ctx context.Context, token string) error {

}
