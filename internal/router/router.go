package router

import (
	"github.com/gorilla/mux"
	"net/http"
	"wedding/internal/handler"
)

func InitRouter(hand *handler.UserHandler) (*mux.Router, error) {
	router := mux.NewRouter()

	// ✅ Главная страница
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./HTML/index.html")
	}).Methods("GET")

	// ✅ Отдаём CSS, JS из папки HTML
	router.PathPrefix("/").Handler(
		http.StripPrefix("/", http.FileServer(http.Dir("./HTML/"))),
	)

	// ✅ API
	router.HandleFunc("/user", hand.AddUserHandler).Methods("POST")

	return router, nil
}
