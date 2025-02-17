package auth

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"

	pbAuth "github.com/evildead81/gophkeeper/api/auth"
	"github.com/evildead81/gophkeeper/pkg/jwt"
)

type AuthService struct {
	db *sql.DB
	pbAuth.UnimplementedAuthServiceServer
}

// NewAuthService создаёт новый экземпляр AuthService
func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{db: db}
}

// Register регистрирует пользователя
func (s *AuthService) Register(ctx context.Context, req *pbAuth.RegisterRequest) (*pbAuth.RegisterResponse, error) {
	if s.userExists(req.Username) {
		return nil, errors.New("пользователь уже существует")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	_, err = s.db.Exec(`
		INSERT INTO users (username, password_hash) 
		VALUES ($1, $2)
	`, req.Username, string(hashedPassword))
	if err != nil {
		return nil, err
	}

	log.Printf("Пользователь %s зарегистрирован", req.Username)
	return &pbAuth.RegisterResponse{Message: "Регистрация успешна"}, nil
}

// Login выполняет аутентификацию пользователя
func (s *AuthService) Login(ctx context.Context, req *pbAuth.LoginRequest) (*pbAuth.LoginResponse, error) {
	var hashedPassword string
	err := s.db.QueryRow(`SELECT password_hash FROM users WHERE username=$1`, req.Username).Scan(&hashedPassword)
	if err == sql.ErrNoRows {
		return nil, errors.New("пользователь не найден")
	}
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		return nil, errors.New("неверный пароль")
	}

	token, err := jwt.GenerateToken(req.Username)
	if err != nil {
		return nil, err
	}

	return &pbAuth.LoginResponse{Token: token}, nil
}

// userExists проверяет существование пользователя
func (s *AuthService) userExists(username string) bool {
	var exists bool
	err := s.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)`, username).Scan(&exists)
	return err == nil && exists
}

// Logout — заглушка для завершения сессии
func (s *AuthService) Logout(ctx context.Context, req *pbAuth.LogoutRequest) (*pbAuth.LogoutResponse, error) {
	return &pbAuth.LogoutResponse{Message: "Выход выполнен"}, nil
}
