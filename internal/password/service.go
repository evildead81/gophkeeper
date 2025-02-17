package password

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"

	"github.com/evildead81/gophkeeper/pkg/crypto"
	"github.com/evildead81/gophkeeper/pkg/jwt"

	pbPassword "github.com/evildead81/gophkeeper/api/password"
)

type PasswordService struct {
	db *sql.DB
	pbPassword.UnimplementedPasswordServiceServer
}

// NewPasswordService создаёт новый сервис
func NewPasswordService(db *sql.DB) *PasswordService {
	return &PasswordService{db: db}
}

// AddPassword добавляет пароль
func (s *PasswordService) AddPassword(ctx context.Context, req *pbPassword.AddPasswordRequest) (*pbPassword.AddPasswordResponse, error) {
	username, err := jwt.ValidateToken(req.Token)
	if err != nil {
		return nil, errors.New("неверный токен")
	}

	// Шифруем пароль перед сохранением
	encryptedPassword, err := crypto.Encrypt(req.Password)
	if err != nil {
		return nil, err
	}

	metaInfo, _ := json.Marshal(req.MetaInfo)

	_, err = s.db.Exec(`
		INSERT INTO password_entries (user_id, site, login, password_encrypted, meta_info) 
		VALUES ((SELECT id FROM users WHERE username=$1), $2, $3, $4, $5)
	`, username, req.Site, req.Login, encryptedPassword, string(metaInfo))

	if err != nil {
		return nil, err
	}

	log.Printf("Добавлен пароль для %s", req.Site)
	return &pbPassword.AddPasswordResponse{Message: "Пароль сохранён"}, nil
}

// GetPassword получает пароль
func (s *PasswordService) GetPassword(ctx context.Context, req *pbPassword.GetPasswordRequest) (*pbPassword.GetPasswordResponse, error) {
	username, err := jwt.ValidateToken(req.Token)
	if err != nil {
		return nil, errors.New("неверный токен")
	}

	var login, encryptedPassword, metaInfo string
	err = s.db.QueryRow(`
		SELECT login, password_encrypted, meta_info 
		FROM password_entries 
		WHERE user_id = (SELECT id FROM users WHERE username=$1) AND site=$2
	`, username, req.Site).Scan(&login, &encryptedPassword, &metaInfo)

	if err == sql.ErrNoRows {
		return nil, errors.New("данные не найдены")
	}
	if err != nil {
		return nil, err
	}

	// Расшифровываем пароль
	password, err := crypto.Decrypt(encryptedPassword)
	if err != nil {
		return nil, err
	}

	return &pbPassword.GetPasswordResponse{
		Login:    login,
		Password: password,
		MetaInfo: metaInfo,
	}, nil
}

// DeletePassword удаляет пароль
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

	log.Printf("Удалена запись для %s", req.Site)
	return &pbPassword.DeletePasswordResponse{Message: "Запись удалена"}, nil
}
