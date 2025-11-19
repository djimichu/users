package models

import "github.com/google/uuid"

type UserHttpDTO struct {
	Id           uuid.UUID
	Name         string
	PasswordHash string
}
type UserHashDTO struct {
	Name         string
	PasswordHash string
}

func NewUserHashDto(name, passHash string) UserHashDTO {
	return UserHashDTO{
		Name:         name,
		PasswordHash: passHash,
	}
}
