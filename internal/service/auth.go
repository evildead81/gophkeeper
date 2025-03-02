package service

import (
	"context"
	"errors"
	"log"

	pb "github.com/evildead81/gophkeeper/api/auth"
	"github.com/evildead81/gophkeeper/internal/repository"
	"github.com/evildead81/gophkeeper/pkg/crypto"
	"github.com/evildead81/gophkeeper/pkg/session"
)

type AuthService struct {
	repo    repository.UserRepositoryInterface
	session *session.SessionManager
	pb.UnimplementedAuthServiceServer
}

// NewAuthService создаёт новый сервис аутентификации
func NewAuthService(repo repository.UserRepositoryInterface, session *session.SessionManager) *AuthService {
	return &AuthService{repo: repo, session: session}
}

// Register регистрирует нового пользователя
func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// Хешируем пароль
	hashedPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("ошибка хеширования пароля")
	}

	// Добавляем пользователя в БД
	err = s.repo.CreateUser(req.Username, hashedPassword)
	if err != nil {
		return nil, errors.New("ошибка регистрации пользователя")
	}

	log.Printf("Пользователь %s зарегистрирован", req.Username)
	return &pb.RegisterResponse{Message: "Регистрация успешна"}, nil
}

// Login выполняет вход пользователя и создаёт сессию
func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Получаем пользователя по имени
	user, err := s.repo.GetUser(req.Username)
	if err != nil {
		return nil, errors.New("пользователь не найден")
	}

	// Проверяем пароль
	if !crypto.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("неверный пароль")
	}

	// Создаём новую сессию
	token := s.session.CreateSession(user.ID)

	log.Printf("Пользователь %s вошёл в систему. Токен: %s", req.Username, token)
	return &pb.LoginResponse{Token: token}, nil
}

// Logout завершает сессию пользователя
func (s *AuthService) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	// Удаляем сессию
	s.session.DeleteSession(req.Token)

	log.Printf("Сессия %s удалена", req.Token)
	return &pb.LogoutResponse{Message: "Вы успешно вышли"}, nil
}

// ValidateToken проверяет, авторизован ли пользователь
func (s *AuthService) ValidateToken(token string) (string, error) {
	return s.session.GetUserID(token)
}
