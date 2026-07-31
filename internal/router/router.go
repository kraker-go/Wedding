package router

import (
	"github.com/gorilla/mux"
	"wedding/internal/handler"
)

func InitRouter(hand *handler.UserHandler) (*mux.Router, error) {
	router := mux.NewRouter()

	router.HandleFunc("/user", hand.AddUserHandler).Methods("POST")

	return router, nil
}
