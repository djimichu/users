package models

import "github.com/google/uuid"

type UserHttpDTO struct {
	Id           uuid.UUID
	Name         string
	PasswordHash string
}
type UserNamePassForDB struct {
	Name         string
	PasswordHash string
}

func NewUserHashDto(name, passHash string) UserNamePassForDB {
	return UserNamePassForDB{
		Name:         name,
		PasswordHash: passHash,
	}
}
