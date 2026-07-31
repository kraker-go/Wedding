package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Postgres struct {
	PORT     string
	HOST     string
	USER     string
	PASSWORD string
	DBNAME   string
	SSL_MODE string
}

type Server struct {
	Port string
}

func CfgInit() (*Postgres, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("Error loading .env file %w", err)
	}
	return &Postgres{
		PORT:     os.Getenv("DB_PORT"),
		HOST:     os.Getenv("DB_HOST"),
		USER:     os.Getenv("DB_USER"),
		PASSWORD: os.Getenv("DB_PASSWORD"),
		DBNAME:   os.Getenv("DB_DBNAME"),
		SSL_MODE: os.Getenv("DB_SSL_MODE"),
	}, nil
}

func ServerInit() (*Server, error) {
	return &Server{os.Getenv("SERVER")}, nil
}
