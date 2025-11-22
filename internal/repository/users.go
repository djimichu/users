package repository

import (
	"JWTproject/internal/models"
	"database/sql"
	"github.com/google/uuid"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (u *UserRepo) CreateUser(user models.UserNamePassForDB) (uuid.UUID, error) {

	var id uuid.UUID
	query := `
        INSERT INTO users (username, password_hash)
        VALUES ($1, $2) 
        RETURNING id
    `
	err := u.db.QueryRow(query, user.Name, user.PasswordHash).Scan(&id)

	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
func (u *UserRepo) LoginUser(name string) (uuid.UUID, string, error) {

	var id uuid.UUID
	var pass string

	query := `
		SELECT id, password_hash
		FROM users
		WHERE username = $1
	`
	err := u.db.QueryRow(query, name).Scan(&id, &pass)
	if err != nil {
		return uuid.Nil, "", err
	}
	return id, pass, nil
}

func (u *UserRepo) GetUserByID(id uuid.UUID) (string, error) {

	var name string

	query := `
		SELECT username
		FROM users
		WHERE id = $1
	`
	err := u.db.QueryRow(query, id).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}
func (u *UserRepo) ChangeUserNameById(id uuid.UUID, name string) error {

	query := `
		UPDATE users
		SET username = $1
		WHERE id = $2
	`
	_, err := u.db.Exec(query, name, id)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepo) DeleteUserById(id uuid.UUID) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	_, err := u.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
