package repository

import (
	"context"
	"errors"
)

// MockSecureDataRepository — мок-репозиторий для тестирования
type MockSecureDataRepository struct {
	data    map[string][]SecureData // userID -> данные
	updates map[string][]SecureDataUpdate
}

// NewMockSecureDataRepository создаёт мок-репозиторий
func NewMockSecureDataRepository() *MockSecureDataRepository {
	return &MockSecureDataRepository{data: make(map[string][]SecureData)}
}

// SaveSecureData сохраняет данные
func (m *MockSecureDataRepository) SaveSecureData(ctx context.Context, userID, dataType, data, metaInfo string) error {
	m.data[userID] = append(m.data[userID], SecureData{ID: "123", Data: data, MetaInfo: metaInfo})
	return nil
}

// GetSecureData возвращает данные пользователя
func (m *MockSecureDataRepository) GetSecureData(ctx context.Context, userID, dataType string) ([]SecureData, error) {
	return m.data[userID], nil
}

// DeleteSecureData удаляет данные пользователя
func (m *MockSecureDataRepository) DeleteSecureData(ctx context.Context, userID, dataID string) error {
	if _, exists := m.data[userID]; !exists {
		return errors.New("данные не найдены")
	}
	m.data[userID] = nil
	return nil
}

// GetLastUpdatedData возвращает последние обновлённые данные пользователя
func (m *MockSecureDataRepository) GetLastUpdatedData(ctx context.Context, userID string) ([]SecureDataUpdate, error) {
	data, exists := m.updates[userID]
	if !exists || len(data) == 0 {
		return nil, errors.New("обновлённые данные не найдены")
	}
	return data, nil
}
