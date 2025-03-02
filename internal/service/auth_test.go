package service

import (
	"context"
	"testing"

	pb "github.com/evildead81/gophkeeper/api/auth"
	"github.com/evildead81/gophkeeper/internal/repository"
	"github.com/evildead81/gophkeeper/pkg/session"

	"github.com/stretchr/testify/assert"
)

func TestAuthService(t *testing.T) {
	userRepo := repository.NewMockUserRepository()
	sessionManager := session.NewSessionManager()
	authService := NewAuthService(userRepo, sessionManager)

	t.Run("Register", func(t *testing.T) {
		resp, err := authService.Register(context.Background(), &pb.RegisterRequest{
			Username: "testuser",
			Password: "testpassword",
		})

		assert.NoError(t, err)
		assert.Equal(t, "Регистрация успешна", resp.Message)
	})

	t.Run("Login", func(t *testing.T) {
		resp, err := authService.Login(context.Background(), &pb.LoginRequest{
			Username: "testuser",
			Password: "testpassword",
		})

		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Token)
	})

	t.Run("Logout", func(t *testing.T) {
		token := sessionManager.CreateSession("testuser")

		resp, err := authService.Logout(context.Background(), &pb.LogoutRequest{Token: token})

		assert.NoError(t, err)
		assert.Equal(t, "Вы успешно вышли", resp.Message)
	})
}
