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

func (uh *UserHandler) RequestUpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	uh.logg.Info("Запрос гостя")

	ctx := r.Context()

	userID := mux.Vars(r)["id"]
	id, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Неверный id", http.StatusBadRequest)
		uh.logg.Error("не корректный id")
		return
	}

	var requestUser models.Guest
	var user *models.Guest

	err = json.NewDecoder(r.Body).Decode(&requestUser)
	if err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		uh.logg.Error("не корректный json")
		return
	}

	user, err = uh.hand.GetUser(ctx, id)
	if err != nil {
		http.Error(w, "пользователь не найден", http.StatusBadRequest)
		uh.logg.Error("пользователь не найден")
		return
	}

	err = uh.hand.RequestUpdateUser(ctx, id, &requestUser)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrorEmptyName):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, models.ErrorInvalidChars):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, models.ErrorNameTooLong):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, models.ErrorsNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			uh.logg.Error("ошибка обновления", zap.Error(err))
		}
		return
	}

	if err == nil && uh.notifier != nil {
		go uh.notifier.Notify(*user, requestUser, "update") // ✅ отправляет в уже существующий канал
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(map[string]interface{}{
		"Данные успешно обновлены": requestUser.FirstName + " " + requestUser.LastName,
	}); err != nil {
		uh.logg.Error("не удалось закодировать ответ", zap.Error(err))
	}
}
