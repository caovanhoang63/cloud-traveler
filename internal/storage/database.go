package storage

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg DBConfig) *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	for {
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			log.Printf("Failed to open database: %v. Retrying in 5s...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if err := db.Ping(); err != nil {
			db.Close()
			log.Printf("Failed to ping database: %v. Retrying in 5s...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(5)
		log.Println("Database connection established")
		return db
	}
}

func RunMigrations(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS uploaded_files (
			id SERIAL PRIMARY KEY,
			filename VARCHAR(255) NOT NULL,
			s3_key VARCHAR(255) NOT NULL UNIQUE,
			content_type VARCHAR(100),
			size_bytes BIGINT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_uploaded_files_s3_key ON uploaded_files(s3_key);
	`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
