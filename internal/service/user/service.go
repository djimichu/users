package user

import (
	"JWTproject/internal/auth"
	"JWTproject/internal/logger"
	"JWTproject/internal/models"
	"JWTproject/internal/repository"
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

//TODO: info логи о бизнес логике

type UserService struct {
	repo      *repository.UserRepo
	jwtManger *auth.JWTManager
}

func NewUserService(repo *repository.UserRepo, jwtManager *auth.JWTManager) *UserService {
	return &UserService{
		repo:      repo,
		jwtManger: jwtManager,
	}
}

func NewUserServiceHash(repo *repository.UserRepo) *UserService {

	return &UserService{
		repo: repo,
	}
}
func hashPassword(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", ErrHashPass
	}

	return string(hashBytes), nil
}

func (s *UserService) CreateUser(user models.UserRequestDto) (uuid.UUID, error) {

	logger.Logger.Info("register attempt",
		zap.String("username", user.Name),
		zap.String("pass", user.Password),
	)

	pass, err := hashPassword(user.Password)
	if err != nil {
		logger.Logger.Error("failed hashing password",
			zap.String("user", user.Name),
			zap.String("pass", user.Password),
			zap.Error(err),
		)
		return uuid.Nil, err
	}

	userHash := models.NewUserHashDto(user.Name, pass)

	id, err := s.repo.CreateUser(userHash)
	if err != nil {
		logger.Logger.Error("failed create user in database",
			zap.String("user", user.Name),
			zap.String("pass", user.Password),
			zap.Error(err),
		)
		return uuid.Nil, err
	}

	logger.Logger.Info("successful create user",
		zap.String("user", user.Name),
		zap.String("pass", userHash.PasswordHash),
		zap.String("id", id.String()),
	)

	return id, nil
}
func (s *UserService) Login(user models.UserRequestDto) (string, error) {

	logger.Logger.Info("login attempt",
		zap.String("username", user.Name),
	)

	id, passDB, err := s.repo.LoginUser(user.Name)
	if err != nil {
		if errors.Is(err, repository.ErrNotFoundUser) {
			logger.Logger.Warn("failed login user in database - user not found",
				zap.String("user", user.Name),
				zap.String("pass", user.Password),
			)
			return "", repository.ErrNotFoundUser
		}
		logger.Logger.Error("failed login user in database",
			zap.String("user", user.Name),
			zap.String("pass", user.Password),
			zap.Error(err),
		)
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passDB), []byte(user.Password))
	if err != nil {
		logger.Logger.Warn("failed login user, incorrect pass",
			zap.String("user", user.Name),
			zap.String("pass", user.Password),
			zap.String("user_id", id.String()),
		)
		return "", ErrIncorrectPass
	}

	token, err := s.jwtManger.Generate(id)
	if err != nil {
		logger.Logger.Error("failed generate JWT-token",
			zap.String("user", user.Name),
			zap.String("pass", user.Password),
			zap.Error(err),
		)
		return "", err
	}

	logger.Logger.Info("successful login user",
		zap.String("user", user.Name),
		zap.String("id", id.String()),
	)

	return token, nil
}

func (s *UserService) GetUserByID(id uuid.UUID) (*models.GetUserDTO, error) {

	logger.Logger.Info("attempting to get user",
		zap.String("user_id", id.String()),
	)

	name, err := s.repo.GetUserByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFoundUser) {
			logger.Logger.Warn("failed get user from DB, this user_id not exist",
				zap.String("id", id.String()),
			)
		}
		logger.Logger.Error("failed get user from DB",
			zap.String("id", id.String()),
			zap.Error(err),
		)
		return nil, err
	}

	user := models.NewGetUserDTO(name, id)

	logger.Logger.Info("successful get user from DB",
		zap.String("user", name),
		zap.String("msg", user.Msg),
		zap.String("id", id.String()),
	)

	return &user, nil

}
func (s *UserService) ChangeUserName(id uuid.UUID, name string) (*models.GetUserDTO, error) {

	logger.Logger.Info("attempting to change username",
		zap.String("user_id", id.String()),
	)

	updatedName, err := s.repo.ChangeUserNameById(id, name)
	if err != nil {
		if errors.Is(err, repository.ErrNotFoundUser) {
			logger.Logger.Warn("failed get user from DB, this user_id not exist",
				zap.String("id", id.String()),
				zap.String("name", name),
			)
		}
		logger.Logger.Error("failed get user from DB",
			zap.String("id", id.String()),
			zap.Error(err),
		)
		return nil, err
	}
	userNew := models.NewGetUserDTOChange(updatedName, id)

	logger.Logger.Info("successful change username from DB",
		zap.String("name", updatedName),
		zap.String("msg", userNew.Msg),
		zap.String("id", id.String()),
	)
	return &userNew, nil
}
func (s *UserService) DeleteUserByID(id uuid.UUID) error {

	logger.Logger.Info("attempting to delete user",
		zap.String("user_id", id.String()),
	)

	if err := s.repo.DeleteUserById(id); err != nil {
		if errors.Is(err, repository.ErrNotFoundUser) {
			logger.Logger.Warn("failed delete user from DB, this user_id not exist",
				zap.String("id", id.String()),
			)
		}
		logger.Logger.Error("failed get user from DB",
			zap.String("id", id.String()),
			zap.Error(err),
		)
		return err
	}

	logger.Logger.Info("successful delete user",
		zap.String("id", id.String()),
	)

	return nil

}
