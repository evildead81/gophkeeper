package repository

import (
	"database/sql"
	"errors"
)

// User представляет сущность пользователя
type User struct {
	ID           string
	Username     string
	PasswordHash string
}

// UserRepositoryInterface — интерфейс для работы с пользователями
type UserRepositoryInterface interface {
	CreateUser(username, passwordHash string) error
	GetUser(username string) (*User, error)
}

// UserRepository хранит данные пользователей в БД
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository создаёт новый репозиторий пользователей
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser добавляет пользователя в БД
func (r *UserRepository) CreateUser(username, passwordHash string) error {
	_, err := r.db.Exec("INSERT INTO users (username, password_hash) VALUES ($1, $2)", username, passwordHash)
	if err != nil {
		return errors.New("ошибка создания пользователя")
	}
	return nil
}

// GetUser получает пользователя по имени
func (r *UserRepository) GetUser(username string) (*User, error) {
	var user User
	err := r.db.QueryRow("SELECT id, username, password_hash FROM users WHERE username=$1", username).
		Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("пользователь не найден")
		}
		return nil, err
	}
	return &user, nil
}
