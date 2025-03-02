package repository

import (
	"errors"
)

// MockUserRepository — фейковый репозиторий пользователей для тестов
type MockUserRepository struct {
	users map[string]string
}

// NewMockUserRepository создаёт новый мок-репозиторий
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{users: make(map[string]string)}
}

// CreateUser добавляет пользователя (используется в тестах)
func (m *MockUserRepository) CreateUser(username, passwordHash string) error {
	if _, exists := m.users[username]; exists {
		return errors.New("пользователь уже существует")
	}
	m.users[username] = passwordHash
	return nil
}

// GetUser получает пользователя по имени (используется в тестах)
func (m *MockUserRepository) GetUser(username string) (*User, error) {
	passwordHash, exists := m.users[username]
	if !exists {
		return nil, errors.New("пользователь не найден")
	}
	return &User{Username: username, PasswordHash: passwordHash}, nil
}
