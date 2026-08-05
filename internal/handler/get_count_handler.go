package handler

import (
	"encoding/json"
	"net/http"
)

func (uh *UserHandler) GetCountUsersHandler(w http.ResponseWriter, r *http.Request) {

	uh.logg.Info("загружаем счетчик количества гостей")

	ctx := r.Context()

	count, err := uh.hand.GetCountUsers(ctx)
	if err != nil {
		uh.logg.Error(err.Error())
		return
	}

	uh.logg.Info("количество гостей определено")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(count); err != nil {
		uh.logg.Error(err.Error())
		return
	}
}
