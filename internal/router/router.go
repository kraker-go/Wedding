package router

import (
	"github.com/gorilla/mux"
	"net/http"
	"wedding/internal/handler"
)

func InitRouter(hand *handler.UserHandler) (*mux.Router, error) {
	router := mux.NewRouter()

	// ✅ Отдаём HTML-страницу
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./HTML/index.html")
	}).Methods("GET")

	// ✅ Отдаём CSS, JS, картинки (если есть)
	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("./HTML/static/"))),
	)

	router.HandleFunc("/user", hand.AddUserHandler).Methods("POST")

	return router, nil
}
