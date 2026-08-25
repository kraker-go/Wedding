package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func (uh *UserHandler) UploadPhotoHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // максимум 10 МБ
	if err != nil {
		http.Error(w, "Ошибка загрузки файла", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "Фото не найдено", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if err := os.MkdirAll("./uploads", 0755); err != nil {
		http.Error(w, "Ошибка создания папки", http.StatusInternalServerError)
		return
	}

	filename := filepath.Base(header.Filename)
	path := filepath.Join("./uploads", filename)

	dst, err := os.Create(path)
	if err != nil {
		http.Error(w, "Ошибка сохранения фото", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = dst.ReadFrom(file)
	if err != nil {
		http.Error(w, "Ошибка записи фото", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	fmt.Fprintf(w, `{"message":"Фото успешно загружено","filename":%q}`, filename)
}
