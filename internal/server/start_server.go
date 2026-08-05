package server

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

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
