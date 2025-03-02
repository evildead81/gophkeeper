package service

import (
	"context"
	"testing"

	pb "github.com/evildead81/gophkeeper/api/secure_data"
	"github.com/evildead81/gophkeeper/internal/repository"
	"github.com/evildead81/gophkeeper/pkg/session"

	"github.com/stretchr/testify/assert"
)

// Тестирование SecureDataService
func TestSecureDataService(t *testing.T) {
	dataRepo := repository.NewMockSecureDataRepository()
	userRepo := repository.NewMockUserRepository()
	sessionManager := session.NewSessionManager()

	service := NewSecureDataService(dataRepo, userRepo, sessionManager)

	token := sessionManager.CreateSession("testuser")

	t.Run("AddSecureData", func(t *testing.T) {
		resp, err := service.AddSecureData(context.Background(), &pb.AddSecureDataRequest{
			Token:    token,
			DataType: "password",
			Data:     "mypassword",
			MetaInfo: "example.com",
		})

		assert.NoError(t, err)
		assert.Equal(t, "Данные сохранены", resp.Message)
	})

	t.Run("GetSecureData", func(t *testing.T) {
		resp, err := service.GetSecureData(context.Background(), &pb.GetSecureDataRequest{
			Token:    token,
			DataType: "password",
		})

		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Items)
		assert.Equal(t, "mypassword", resp.Items[0].Data)
	})

	t.Run("DeleteSecureData", func(t *testing.T) {
		resp, err := service.DeleteSecureData(context.Background(), &pb.DeleteSecureDataRequest{
			Token: token,
			Id:    "123",
		})

		assert.NoError(t, err)
		assert.Equal(t, "Данные удалены", resp.Message)
	})

	t.Run("Unauthorized Access", func(t *testing.T) {
		_, err := service.GetSecureData(context.Background(), &pb.GetSecureDataRequest{
			Token:    "INVALID_TOKEN",
			DataType: "password",
		})

		assert.Error(t, err)
		assert.Equal(t, "пользователь не авторизован", err.Error())
	})
}
