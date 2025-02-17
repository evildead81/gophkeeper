package auth

import (
	"context"
	"database/sql"
	"testing"

	pb "github.com/evildead81/gophkeeper/api/auth"
	"github.com/evildead81/gophkeeper/pkg/db"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupDB() *sql.DB {
	dsn := "postgres://postgres:password@localhost:5432/gophkeeper_test?sslmode=disable"
	conn, err := db.Connect(dsn)
	if err != nil {
		panic(err)
	}
	return conn
}

func TestRegister(t *testing.T) {
	db := setupDB()
	defer db.Close()

	authService := NewAuthService(db)

	req := &pb.RegisterRequest{
		Username: "testuser",
		Password: "testpassword",
	}
	resp, err := authService.Register(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, "Регистрация успешна", resp.Message)
}

func TestLoginSuccess(t *testing.T) {
	db := setupDB()
	defer db.Close()

	authService := NewAuthService(db)

	authService.Register(context.Background(), &pb.RegisterRequest{
		Username: "testuser",
		Password: "testpassword",
	})

	req := &pb.LoginRequest{
		Username: "testuser",
		Password: "testpassword",
	}
	resp, err := authService.Login(context.Background(), req)

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
}

func TestLoginWrongPassword(t *testing.T) {
	db := setupDB()
	defer db.Close()

	authService := NewAuthService(db)

	authService.Register(context.Background(), &pb.RegisterRequest{
		Username: "testuser",
		Password: "testpassword",
	})

	req := &pb.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}
	resp, err := authService.Login(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}
