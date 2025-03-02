package repository

import (
	"context"
	"database/sql"
	"errors"
)

type FileRepository struct {
	db *sql.DB
}

func NewFileRepository(db *sql.DB) *FileRepository {
	return &FileRepository{db: db}
}

// SaveFileInfo сохраняет метаинформацию о файле в БД
func (r *FileRepository) SaveFileInfo(ctx context.Context, userID, fileID, filename string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO secure_data (user_id, data_type, data_encrypted, meta_info) 
		VALUES ($1, 'binary', $2, $3)`, userID, fileID, filename)
	return err
}

// GetFileInfo получает метаинформацию о файле
func (r *FileRepository) GetFileInfo(ctx context.Context, fileID string) (string, string, error) {
	var userID, filename string
	err := r.db.QueryRowContext(ctx, `
		SELECT user_id, meta_info FROM secure_data WHERE id = $1`, fileID).
		Scan(&userID, &filename)

	if err == sql.ErrNoRows {
		return "", "", errors.New("файл не найден")
	}
	return userID, filename, err
}
