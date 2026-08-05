package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"log"
	"wedding/internal/config"
)

func ConnectPostgres(str config.Config) (*sql.DB, error) {
	// Убрал лишние пробелы между password и sslmode
	connStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		str.DB_HOST, str.DB_PORT, str.DB_DBNAME, str.DB_USER, str.DB_PASSWORD, str.DB_SSL_MODE)

	log.Printf("Connecting to DB: host=%s port=%s dbname=%s user=%s sslmode=%s",
		str.DB_HOST, str.DB_PORT, str.DB_DBNAME, str.DB_USER, str.DB_SSL_MODE)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("SQL Open ERROR: %v", err) // ← ТЕПЕРЬ ВИДНО
		return nil, fmt.Errorf("sql.Open failed: %w", err)
	}

	if err = db.Ping(); err != nil {
		log.Printf("Ping ERROR: %v", err) // ← ТЕПЕРЬ ВИДНО
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL")
	return db, nil
}
func InitMigration(postgres *sql.DB) error {
	return goose.Up(postgres, "migrations")
}
