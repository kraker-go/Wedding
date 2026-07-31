package config

import (
	"os"
)

type Server struct {
	Port string
}

func ServerInit() (*Server, error) {
	port := os.Getenv("SERVER")
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
