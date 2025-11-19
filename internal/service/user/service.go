package user

import (
	"JWTproject/internal/auth"
	"JWTproject/internal/models"
	"JWTproject/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo      *repository.UserRepo
	jwtManger *auth.JWTManager
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
	pass, err := hashPassword(user.Password)
	if err != nil {
		return uuid.Nil, err
	}

	userHash := models.NewUserHashDto(user.Name, pass)

	id, err := s.repo.CreateUser(userHash)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
func (s *UserService) Login(user models.UserRequestDto) (string, error) {

	id, passDB, err := s.repo.LoginUser(user.Name)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(passDB), []byte(user.Password))
	if err != nil {
		return "", ErrIncorrectPass
	}
	token, err := s.jwtManger.Generate(id)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *UserService) GetUserByID(id uuid.UUID) (*models.GetUserDTO, error) {

	name, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	user := models.NewGetUserDTO(name, id)

	return &user, nil

}
func (s *UserService) ChangeUserName(id uuid.UUID, name string) (*models.GetUserDTO, error) {

	if err := s.repo.ChangeUserNameById(id, name); err != nil {
		return nil, err
	}
	userNew := models.NewGetUserDTOChange(name, id)
	return &userNew, nil
}
func (s *UserService) DeleteUserByID(id uuid.UUID) error {

	if err := s.repo.DeleteUserById(id); err != nil {
		return err
	}
	return nil

}
