package service

import (
	"github.com/radionovel/goauth-jwt-microservice/internal/pkg/logger"
	"github.com/radionovel/goauth-jwt-microservice/internal/storage"
)

type UserService struct {
	storage storage.UserStorage
	logger  logger.Logger
}

func NewUserService(storage storage.UserStorage, logger logger.Logger) *UserService {
	return &UserService{
		storage: storage,
		logger:  logger,
	}
}

/*
func (s *UserService) userExist(ctx context.Context, username string) (bool, error) {
	filter := model.Filter{
		Username: Username,
	}

	user, err := s.storage.FindOne(ctx, filter)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return false, nil
		}

		return false, err
	}

	return user != nil, nil
}

func (s *UserService) CreateUser(ctx context.Context, dto *model.NewUserDTO) (model.UserID, error) {
	s.logger.Debug("create new user", "username", dto.Username)

	// find user with username
	exists, err := s.userExist(ctx, dto.Username)
	if err != nil {
		return 0, err
	}

	if exists {
		return 0, model.ErrUserAlreadyExists
	}

	userID, err := s.storage.Insert(ctx, dto)
	if err != nil {
		return 0, err
	}

	// publish event to event broker via outbox

	return userID, nil
}
*/
