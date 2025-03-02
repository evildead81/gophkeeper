package repository

import (
	"context"
	"database/sql"

	"github.com/evildead81/gophkeeper/pkg/crypto"
)

type SecureData struct {
	ID       string
	Data     string
	MetaInfo string
}

type SecureDataRepository struct {
	db *sql.DB
}

// SecureDataRepositoryInterface — интерфейс репозитория данных
type SecureDataRepositoryInterface interface {
	SaveSecureData(ctx context.Context, userID, dataType, data, metaInfo string) error
	GetSecureData(ctx context.Context, userID, dataType string) ([]SecureData, error)
	DeleteSecureData(ctx context.Context, userID, dataID string) error
	GetLastUpdatedData(ctx context.Context, userID string) ([]SecureDataUpdate, error)
}

func NewSecureDataRepository(db *sql.DB) *SecureDataRepository {
	return &SecureDataRepository{db: db}
}

func (r *SecureDataRepository) SaveSecureData(ctx context.Context, userID, dataType, data, metaInfo string) error {
	encryptedData, err := crypto.Encrypt(data)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO secure_data (user_id, data_type, data_encrypted, meta_info) 
		VALUES ($1, $2, $3, $4)`,
		userID, dataType, encryptedData, metaInfo)

	return err
}

func (r *SecureDataRepository) GetSecureData(ctx context.Context, userID, dataType string) ([]SecureData, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, data_encrypted, meta_info FROM secure_data 
		WHERE user_id = $1 AND data_type = $2`, userID, dataType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []SecureData

	for rows.Next() {
		var id, encryptedData, metaInfo string
		if err := rows.Scan(&id, &encryptedData, &metaInfo); err != nil {
			return nil, err
		}

		data, err := crypto.Decrypt(encryptedData)
		if err != nil {
			return nil, err
		}

		results = append(results, SecureData{
			ID:       id,
			Data:     data,
			MetaInfo: metaInfo,
		})
	}

	return results, nil
}

func (r *SecureDataRepository) DeleteSecureData(ctx context.Context, userID, id string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM secure_data WHERE user_id = $1 AND id = $2`, userID, id)
	return err
}

func (r *SecureDataRepository) GetLastUpdatedData(ctx context.Context, userID string) ([]SecureDataUpdate, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, data_type, data_encrypted, meta_info, updated_at, action FROM secure_data 
		WHERE user_id = $1 AND updated_at > NOW() - INTERVAL '5 minutes'
		ORDER BY updated_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var updates []SecureDataUpdate

	for rows.Next() {
		var update SecureDataUpdate
		err := rows.Scan(&update.ID, &update.DataType, &update.EncryptedData, &update.MetaInfo, &update.UpdatedAt, &update.Action)
		if err != nil {
			return nil, err
		}

		update.Data, err = crypto.Decrypt(update.EncryptedData)
		if err != nil {
			return nil, err
		}

		updates = append(updates, update)
	}

	return updates, nil
}

type SecureDataUpdate struct {
	ID            string
	DataType      string
	EncryptedData string
	Data          string
	MetaInfo      string
	UpdatedAt     string
	Action        string
}
