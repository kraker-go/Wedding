package router

import (
	"github.com/gorilla/mux"
	"net/http"
	"wedding/internal/handler"
)

func InitRouter(hand *handler.UserHandler) (*mux.Router, error) {
	router := mux.NewRouter()

	router.HandleFunc("/user", hand.AddUserHandler).Methods("POST")
	router.HandleFunc("/user/count", hand.GetCountUsersHandler).Methods("GET")
	router.HandleFunc("/user/get", hand.GetAllUsersHandler).Methods("GET")

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./HTML/index.html")
	}).Methods("GET")

	router.PathPrefix("/").Handler(
		http.StripPrefix("/", http.FileServer(http.Dir("./HTML/"))),
	)

	return router, nil
}
