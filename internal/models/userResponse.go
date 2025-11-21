package models

import (
	"fmt"
	"github.com/google/uuid"
)

type UserRegRespDTO struct {
	Name string    `json:"name"`
	ID   uuid.UUID `json:"id"`
	Msg  string    `json:"msg"`
}

type LoginTokenDTO struct {
	Token string `yaml:"token"`
}

type GetUserDTO struct {
	UserID uuid.UUID `json:"user_id"`
	Msg    string    `json:"msg"`
}

func NewGetUserDTOChange(name string, id uuid.UUID) GetUserDTO {
	return GetUserDTO{
		UserID: id,
		Msg:    fmt.Sprintf("Вы успешно изменили имя профиля: %s", name),
	}
}

func NewGetUserDTO(name string, id uuid.UUID) GetUserDTO {
	return GetUserDTO{
		UserID: id,
		Msg:    fmt.Sprintf("%s вы авторизованы!", name),
	}
}

func NewUserRegRespDTO(name string, id uuid.UUID) UserRegRespDTO {
	return UserRegRespDTO{
		Name: name,
		ID:   id,
		Msg:  "successful register!",
	}
}
