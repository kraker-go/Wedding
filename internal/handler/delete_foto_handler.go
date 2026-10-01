package handler

import (
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"path/filepath"
)

func (uh *UserHandler) DeleteFotoHandler(w http.ResponseWriter, r *http.Request) {
	uh.logg.Info("Запрос на удаление foto")

	nameFoto := mux.Vars(r)
	filename := nameFoto["filename"]

	if filename == "" {
		http.Error(w, "filename is required", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join("./uploads", filepath.Base(filename))

	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	uh.tgNotifier.NotifyMessage("❌ удалена фотография из галереи")
	uh.maxNotifier.NotifyMessage("❌ удалена фотография из галереи")

	w.WriteHeader(http.StatusNoContent)
}
