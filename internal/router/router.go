package router

import (
	"github.com/gorilla/mux"
	"net/http"
	"wedding/internal/handler"
)

func InitRouter(hand *handler.UserHandler) (*mux.Router, error) {
	router := mux.NewRouter()

	router.PathPrefix("/uploads/").Handler(
		http.StripPrefix(
			"/uploads/",
			http.FileServer(http.Dir("./uploads")),
		),
	)
	router.HandleFunc("/upload", hand.UploadPhotoHandler).Methods("POST")
	router.HandleFunc("/upload", hand.GetPhotosHandler).Methods("GET")

	router.HandleFunc("/visit", hand.VisitHandler).Methods("GET")
	router.HandleFunc("/user", hand.AddUserHandler).Methods("POST")
	router.HandleFunc("/user/count", hand.GetCountUsersHandler).Methods("GET")
	router.HandleFunc("/user/get", hand.GetAllUsersHandler).Methods("GET")
	router.HandleFunc("/user/update/{id}", hand.RequestUpdateUserHandler).Methods("PUT")
	router.HandleFunc("/user/{id}", hand.RequestDeleteUserHandler).Methods("DELETE")

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./HTML/index.html")
	}).Methods("GET")

	router.PathPrefix("/").Handler(
		http.StripPrefix("/", http.FileServer(http.Dir("./HTML/"))),
	)

	return router, nil
}
