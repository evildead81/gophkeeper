package password

import (
	"context"
	"database/sql"
	"testing"

	pb "github.com/evildead81/gophkeeper/api/password"
	"github.com/evildead81/gophkeeper/pkg/db"
	"github.com/evildead81/gophkeeper/pkg/jwt"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

// setupDB подключается к тестовой БД
func setupDB() *sql.DB {
	dsn := "postgres://postgres:password@localhost:5432/gophkeeper_test?sslmode=disable"
	conn, err := db.Connect(dsn)
	if err != nil {
		panic(err)
	}
	return conn
}

// setupTest создает тестового пользователя и возвращает токен
func setupTest(db *sql.DB) string {
	_, _ = db.Exec(`DELETE FROM users; DELETE FROM password_entries;`) // Чистим тестовые данные
	_, _ = db.Exec(`INSERT INTO users (id, username, password_hash) VALUES ('00000000-0000-0000-0000-000000000001', 'testuser', '$2a$10$abcdefg')`)
	token, _ := jwt.GenerateToken("testuser")
	return token
}

func TestAddPassword(t *testing.T) {
	db := setupDB()
	defer db.Close()

	passwordService := NewPasswordService(db)
	token := setupTest(db)

	req := &pb.AddPasswordRequest{
		Token:    token,
		Site:     "example.com",
		Login:    "testuser",
		Password: "mypassword",
		MetaInfo: `{"note": "test entry"}`,
	}

	resp, err := passwordService.AddPassword(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, "Пароль сохранён", resp.Message)
}

func TestGetPassword(t *testing.T) {
	db := setupDB()
	defer db.Close()

	passwordService := NewPasswordService(db)
	token := setupTest(db)

	_, _ = passwordService.AddPassword(context.Background(), &pb.AddPasswordRequest{
		Token:    token,
		Site:     "example.com",
		Login:    "testuser",
		Password: "mypassword",
		MetaInfo: `{"note": "test entry"}`,
	})

	req := &pb.GetPasswordRequest{
		Token: token,
		Site:  "example.com",
	}

	resp, err := passwordService.GetPassword(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, "testuser", resp.Login)
	assert.Equal(t, "mypassword", resp.Password)
}

func TestDeletePassword(t *testing.T) {
	db := setupDB()
	defer db.Close()

	passwordService := NewPasswordService(db)
	token := setupTest(db)

	_, _ = passwordService.AddPassword(context.Background(), &pb.AddPasswordRequest{
		Token:    token,
		Site:     "example.com",
		Login:    "testuser",
		Password: "mypassword",
		MetaInfo: `{"note": "test entry"}`,
	})

	req := &pb.DeletePasswordRequest{
		Token: token,
		Site:  "example.com",
	}

	resp, err := passwordService.DeletePassword(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, "Запись удалена", resp.Message)
}
