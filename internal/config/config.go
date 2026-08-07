package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Server struct {
	Port string
}
type Config struct {
	DB_HOST     string
	DB_PORT     string
	DB_DBNAME   string
	DB_USER     string
	DB_PASSWORD string
	DB_SSL_MODE string
}

func InitConfig() (*Config, error) {
	_ = godotenv.Load(".env")

	return &Config{
		DB_HOST:     os.Getenv("DB_HOST"),
		DB_PORT:     os.Getenv("DB_PORT"),
		DB_DBNAME:   os.Getenv("DB_DBNAME"),
		DB_USER:     os.Getenv("DB_USER"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		DB_SSL_MODE: os.Getenv("DB_SSL_MODE"),
	}, nil

}

type Telegramm struct {
	Bot    string
	ChatID string
}

func InitTelegramm() (*Telegramm, error) {
	_ = godotenv.Load(".env")
	return &Telegramm{
		Bot:    os.Getenv("TG_BOT"),
		ChatID: os.Getenv("CHAT_ID"),
	}, nil
}

func ServerInit() (*Server, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("PORT")
	}

	if port == "" {
		port = "8080"
	}

	if port[0] != ':' {
		port = ":" + port
	}

	return &Server{Port: port}, nil
}
