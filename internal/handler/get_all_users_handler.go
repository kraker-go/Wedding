package handler

import (
	"encoding/json"
	"net/http"
)

func (uh *UserHandler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	uh.logg.Info("Поиск всех гостей...")

	ctx := r.Context()

	users, err := uh.hand.GetAllUsers(ctx)
	if err != nil {
		http.Error(w, "ошибка поиска гостей", http.StatusBadRequest)
		uh.logg.Error(err.Error())
		return
	}

	uh.logg.Info("список гостей собран успешно")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(users); err != nil {
		uh.logg.Error(err.Error())
		return
	}
}
