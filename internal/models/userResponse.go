package models

import (
	"fmt"
	"github.com/google/uuid"
)

type UserRegRespDTO struct {
	Name string
	ID   uuid.UUID
	msg  string
}

type LoginTokenDTO struct {
	Token string `yaml:"token"`
}

type GetUserDTO struct {
	UserID uuid.UUID
	msg    string
}

func NewGetUserDTOChange(name string, id uuid.UUID) GetUserDTO {
	return GetUserDTO{
		UserID: id,
		msg:    fmt.Sprintf("Вы успешно изменили имя профиля: %s", name),
	}
}

func NewGetUserDTO(name string, id uuid.UUID) GetUserDTO {
	return GetUserDTO{
		UserID: id,
		msg:    fmt.Sprintf("%s вы авторизованы!", name),
	}
}

func NewUserRegRespDTO(name string, id uuid.UUID) UserRegRespDTO {
	return UserRegRespDTO{
		Name: name,
		ID:   id,
		msg:  "successful register!",
	}
}
