package database

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"net/http"
	"os"
	"time"
)

func ConnectPostgres() (*sql.DB, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return nil, fmt.Errorf("переменная DATABASE_URL не установлена")
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
func InitMigration(postgres *sql.DB) error {
	return goose.Up(postgres, "migrations")
}

func ConnectServer(port string, rout *mux.Router) (*http.Server, error) {
	srv := &http.Server{
		Addr:         port,
		Handler:      rout,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	return srv, nil
}
