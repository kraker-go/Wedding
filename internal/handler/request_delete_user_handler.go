package handler

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	models "wedding/internal/domain"
)

func (uh *UserHandler) RequestDeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	uh.logg.Info("Запрос на удаление гостя")

	ctx := r.Context()

	userID := mux.Vars(r)["id"]
	id, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "ID не указан", http.StatusBadRequest)
		return
	}
	
	user, err := uh.hand.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, models.ErrorsNotFound) {
			http.Error(w, "Гость не найден", http.StatusNotFound)
		} else {
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			uh.logg.Error("ошибка поиска гостя", zap.Error(err))
		}
		return
	}

	if uh.notifier != nil {
		go uh.notifier.Notify(*user, models.Guest{}, "delete")
	}

	uh.logg.Info("Запрос на удаление отправлен на подтверждение")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "🗑️ Запрос на удаление отправлен на подтверждение.",
	}); err != nil {
		uh.logg.Error("не удалось закодировать ответ", zap.Error(err))
	}
}
