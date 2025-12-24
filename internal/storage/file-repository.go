package storage

import (
	"context"
	"database/sql"
	"time"
)

type UploadedFile struct {
	ID          int64     `json:"id"`
	Filename    string    `json:"filename"`
	S3Key       string    `json:"s3_key"`
	ContentType string    `json:"content_type"`
	SizeBytes   int64     `json:"size_bytes"`
	CreatedAt   time.Time `json:"created_at"`
}

type FileRepository struct {
	db *sql.DB
}

func NewFileRepository(db *sql.DB) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) Create(ctx context.Context, file *UploadedFile) error {
	query := `
		INSERT INTO uploaded_files (filename, s3_key, content_type, size_bytes)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	return r.db.QueryRowContext(ctx, query,
		file.Filename, file.S3Key, file.ContentType, file.SizeBytes,
	).Scan(&file.ID, &file.CreatedAt)
}

func (r *FileRepository) GetByID(ctx context.Context, id int64) (*UploadedFile, error) {
	query := `
		SELECT id, filename, s3_key, content_type, size_bytes, created_at
		FROM uploaded_files WHERE id = $1
	`

	file := &UploadedFile{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&file.ID, &file.Filename, &file.S3Key,
		&file.ContentType, &file.SizeBytes, &file.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (r *FileRepository) List(ctx context.Context, limit, offset int) ([]UploadedFile, error) {
	query := `
		SELECT id, filename, s3_key, content_type, size_bytes, created_at
		FROM uploaded_files
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []UploadedFile
	for rows.Next() {
		var f UploadedFile
		if err := rows.Scan(&f.ID, &f.Filename, &f.S3Key, &f.ContentType, &f.SizeBytes, &f.CreatedAt); err != nil {
			return nil, err
		}
		files = append(files, f)
	}

	return files, rows.Err()
}

func (r *FileRepository) Delete(ctx context.Context, id int64) (string, error) {
	query := `DELETE FROM uploaded_files WHERE id = $1 RETURNING s3_key`

	var s3Key string
	err := r.db.QueryRowContext(ctx, query, id).Scan(&s3Key)
	if err != nil {
		return "", err
	}

	return s3Key, nil
}

func (r *FileRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM uploaded_files").Scan(&count)
	return count, err
}
