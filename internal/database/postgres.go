package database

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"net/http"
	"time"
	"wedding/internal/config"
)

func ConnectPostgres(cfg config.Postgres) (*sql.DB, error) {
	str := fmt.Sprintf("port=%s host=%s user=%s password=%s dbname=%s sslmode=%s", cfg.PORT, cfg.HOST, cfg.USER, cfg.PASSWORD, cfg.DBNAME, cfg.SSL_MODE)
	postgres, err := sql.Open("postgres", str)
	if err != nil {
		return nil, fmt.Errorf("Ошибка подключения postgres: %w", err)
	}

	if err = postgres.Ping(); err != nil {
		return nil, fmt.Errorf("Нет соединения с Postgres: %w", err)
	}

	return postgres, nil

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
