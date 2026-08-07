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

func (uh *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	uh.logg.Info("Удаляем гостя")

	ctx := r.Context()

	userID := mux.Vars(r)["id"]
	id, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "ID не указан", http.StatusBadRequest)
		return
	}

	user, err := uh.hand.DeleteUser(ctx, id)
	if err != nil {
		if errors.Is(err, models.ErrorsNotFound) {
			http.Error(w, "Гость не найден", http.StatusNotFound)
		} else {
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			uh.logg.Error("ошибка удаления", zap.Error(err))
		}
		return
	}
	if err == nil && uh.notifier != nil {
		go uh.notifier.Notify(*user, "delete") // ✅ отправляет в уже существующий канал
	}
	uh.logg.Info("гость успешно удален")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(map[string]string{
		"message": "✅ Гость успешно удалён",
	}); err != nil {
		uh.logg.Error("не удалось закодировать ответ", zap.Error(err))
	}
}
