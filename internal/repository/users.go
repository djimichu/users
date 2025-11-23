package repository

import (
	"JWTproject/internal/logger"
	"JWTproject/internal/models"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

//TODO: debug логи о выполнении запросов

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (u *UserRepo) CreateUser(user models.UserNamePassForDB) (uuid.UUID, error) {

	start := time.Now()

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

	logger.Logger.Debug("Database query executed",
		zap.String("operation", "CreateUser"),
		zap.String("name", user.Name),
		zap.String("password_hash", user.PasswordHash),
		zap.Duration("duration", time.Since(start)),
		zap.Error(err),
	)

	return id, nil
}
func (u *UserRepo) LoginUser(name string) (uuid.UUID, string, error) {

	start := time.Now()

	var id uuid.UUID
	var pass string

	query := `
		SELECT id, password_hash
		FROM users
		WHERE username = $1
	`
	err := u.db.QueryRow(query, name).Scan(&id, &pass)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, "", ErrNotFoundUser
		}
		return uuid.Nil, "", err
	}

	logger.Logger.Debug("Database query executed",
		zap.String("operation", "LoginUser"),
		zap.String("name", name),
		zap.String("password_hash", pass),
		zap.Duration("duration", time.Since(start)),
		zap.Error(err),
	)

	return id, pass, nil
}

func (u *UserRepo) GetUserByID(id uuid.UUID) (string, error) {

	start := time.Now()

	var name string

	query := `
		SELECT username
		FROM users
		WHERE id = $1
	`
	err := u.db.QueryRow(query, id).Scan(&name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNotFoundUser
		}
		return "", err
	}

	logger.Logger.Debug("Database query executed",
		zap.String("operation", "GetUserByID"),
		zap.String("user_id", id.String()),
		zap.String("name", name),
		zap.Duration("duration", time.Since(start)),
		zap.Error(err),
	)

	return name, nil
}
func (u *UserRepo) ChangeUserNameById(id uuid.UUID, name string) (string, error) {

	start := time.Now()

	var updatedName string

	query := `
		UPDATE users
		SET username = $1
		WHERE id = $2
	`
	err := u.db.QueryRow(query, name, id).Scan(&updatedName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNotFoundUser
		}
		return "", err
	}

	logger.Logger.Debug("Database query executed",
		zap.String("operation", "ChangeUserNameById"),
		zap.String("user_id", id.String()),
		zap.String("new_name", updatedName),
		zap.Duration("duration", time.Since(start)),
		zap.Error(err),
	)

	return updatedName, nil
}
func (u *UserRepo) DeleteUserById(id uuid.UUID) error {

	start := time.Now()

	query := `
		DELETE FROM users
		WHERE id = $1
	`
	res, err := u.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFoundUser
	}

	logger.Logger.Debug("Database query executed",
		zap.String("operation", "DeleteUserById"),
		zap.String("user_id", id.String()),
		zap.Duration("duration", time.Since(start)),
		zap.Error(err),
	)

	return nil
}
