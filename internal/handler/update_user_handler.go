package handler

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"net/http"
	models "wedding/internal/domain"
)

func (uh *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	uh.logg.Info("Обновляем данные гостя")

	ctx := r.Context()

	var user models.Guest

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		uh.logg.Error("не корректный json")
		return
	}

	err = uh.hand.UpdateUser(ctx, &user)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrorEmptyName):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, models.ErrorInvalidChars):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, models.ErrorNameTooLong):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, models.ErrorNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			uh.logg.Error("ошибка обновления", zap.Error(err))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(map[string]string{
		"message": "✅ Данные успешно обновлены",
	}); err != nil {
		uh.logg.Error("не удалось закодировать ответ", zap.Error(err))
	}
}
