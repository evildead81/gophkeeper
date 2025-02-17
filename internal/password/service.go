package password

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"sync"

	pbPassword "github.com/evildead81/gophkeeper/api/password"
	"github.com/evildead81/gophkeeper/pkg/crypto"
	"github.com/evildead81/gophkeeper/pkg/jwt"
)

type PasswordService struct {
	db      *sql.DB
	clients map[string]chan *pbPassword.SyncResponse
	mu      sync.Mutex
	pbPassword.UnimplementedPasswordServiceServer
}

// NewPasswordService создаёт новый сервис
func NewPasswordService(db *sql.DB) *PasswordService {
	return &PasswordService{
		db:      db,
		clients: make(map[string]chan *pbPassword.SyncResponse),
	}
}

// AddPassword добавляет пароль и отправляет обновление клиентам
func (s *PasswordService) AddPassword(ctx context.Context, req *pbPassword.AddPasswordRequest) (*pbPassword.AddPasswordResponse, error) {
	username, err := jwt.ValidateToken(req.Token)
	if err != nil {
		return nil, errors.New("неверный токен")
	}

	encryptedPassword, err := crypto.Encrypt(req.Password)
	if err != nil {
		return nil, err
	}

	_, err = s.db.Exec(`
		INSERT INTO password_entries (user_id, site, login, password_encrypted, meta_info) 
		VALUES ((SELECT id FROM users WHERE username=$1), $2, $3, $4, $5)
	`, username, req.Site, req.Login, encryptedPassword, req.MetaInfo)
	if err != nil {
		return nil, err
	}

	s.notifyClients(username, &pbPassword.SyncResponse{
		Site:     req.Site,
		Login:    req.Login,
		Password: req.Password,
		MetaInfo: req.MetaInfo,
		Action:   "added",
	})

	log.Printf("Добавлен пароль для %s", req.Site)
	return &pbPassword.AddPasswordResponse{Message: "Пароль сохранён"}, nil
}

// DeletePassword удаляет пароль и отправляет обновление клиентам
func (s *PasswordService) DeletePassword(ctx context.Context, req *pbPassword.DeletePasswordRequest) (*pbPassword.DeletePasswordResponse, error) {
	username, err := jwt.ValidateToken(req.Token)
	if err != nil {
		return nil, errors.New("неверный токен")
	}

	_, err = s.db.Exec(`
		DELETE FROM password_entries 
		WHERE user_id = (SELECT id FROM users WHERE username=$1) AND site=$2
	`, username, req.Site)
	if err != nil {
		return nil, err
	}

	s.notifyClients(username, &pbPassword.SyncResponse{
		Site:   req.Site,
		Action: "deleted",
	})

	log.Printf("Удалена запись для %s", req.Site)
	return &pbPassword.DeletePasswordResponse{Message: "Запись удалена"}, nil
}

// SyncPasswords отправляет клиенту обновления в реальном времени
func (s *PasswordService) SyncPasswords(req *pbPassword.SyncRequest, stream pbPassword.PasswordService_SyncPasswordsServer) error {
	username, err := jwt.ValidateToken(req.Token)
	if err != nil {
		return errors.New("неверный токен")
	}

	s.mu.Lock()
	ch := make(chan *pbPassword.SyncResponse, 10)
	s.clients[username] = ch
	s.mu.Unlock()

	log.Printf("Клиент %s подключился к синхронизации", username)

	for update := range ch {
		if err := stream.Send(update); err != nil {
			break
		}
	}

	s.mu.Lock()
	delete(s.clients, username)
	s.mu.Unlock()

	log.Printf("Клиент %s отключился от синхронизации", username)
	return nil
}

func (s *PasswordService) notifyClients(username string, update *pbPassword.SyncResponse) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if ch, ok := s.clients[username]; ok {
		ch <- update
	}
}
