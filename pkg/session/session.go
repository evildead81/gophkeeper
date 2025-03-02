package session

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
)

// SessionManager управляет сессиями пользователей
type SessionManager struct {
	store *cache.Cache
}

// NewSessionManager создаёт новый менеджер сессий
func NewSessionManager() *SessionManager {
	return &SessionManager{
		store: cache.New(24*time.Hour, 10*time.Minute), // Сессии хранятся 24 часа
	}
}

// CreateSession создаёт новую сессию
func (s *SessionManager) CreateSession(userID string) string {
	token := uuid.NewString() // Генерируем случайный UUID как токен
	s.store.Set(token, userID, cache.DefaultExpiration)
	return token
}

// GetUserID по токену возвращает userID
func (s *SessionManager) GetUserID(token string) (string, error) {
	userID, found := s.store.Get(token)
	if !found {
		return "", errors.New("сессия недействительна")
	}
	return userID.(string), nil
}

// DeleteSession удаляет токен из хранилища
func (s *SessionManager) DeleteSession(token string) {
	s.store.Delete(token)
}
