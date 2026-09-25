package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"sort"
)

func (uh *UserHandler) GetPhotosHandler(w http.ResponseWriter, r *http.Request) {
	uh.logg.Info("Получаем список фотографий")

	files, err := os.ReadDir("./uploads")
	if err != nil {
		uh.logg.Error("Ошибка чтения папки uploads")
		http.Error(w, "Ошибка получения фотографий", http.StatusInternalServerError)
		return
	}

	sort.Slice(files, func(i, j int) bool {
		infoI, errI := files[i].Info()
		infoJ, errJ := files[j].Info()
		if errI != nil || errJ != nil {
			return false
		}

		return infoI.ModTime().After(infoJ.ModTime())
	})

	var photos []string

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		photos = append(photos, file.Name())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(photos); err != nil {
		uh.logg.Error("Ошибка формирования ответа")
	}
}
