package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"wedding/internal/handler/geo"
)

func (uh *UserHandler) GetCountUsersHandler(w http.ResponseWriter, r *http.Request) {

	uh.logg.Info("загружаем счетчик количества гостей")

	ctx := r.Context()

	count, err := uh.hand.GetCountUsers(ctx)
	if err != nil {
		uh.logg.Error(err.Error())
		return
	}
	if uh.notifier != nil {
		go func() {
			ip := getClientIP(r)
			userAgent := r.UserAgent()
			// Получаем город по IP (асинхронно, с таймаутом)
			city, _ := geo.GetCityByIP(ip) // ошибку игнорируем, просто не покажем город

			msg := "👀 **Посещение сайта**\n" +
				"IP: " + ip + "\n" +
				"Город: " + city + "\n" +
				"Устройство: " + userAgent + "\n" +
				"Время: " + time.Now().Format("2006-01-02 15:04")

			uh.notifier.NotifyMessage(msg)
		}()
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
