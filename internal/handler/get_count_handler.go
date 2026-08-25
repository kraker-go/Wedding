package handler

import (
	"encoding/json"
	"net/http"
	"strings"
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

func getClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return xrip
	}
	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		return ip[:idx]
	}
	return ip
}
