package handler

import (
	"net/http"
	"time"
	"wedding/internal/handler/geo"
)

func (uh *UserHandler) VisitHandler(w http.ResponseWriter, r *http.Request) {
	uh.logg.Info("Посещение сайта")

	if uh.notifier != nil {
		go func() {
			ip := getClientIP(r)
			userAgent := r.UserAgent()
			city, _ := geo.GetCityByIP(ip)

			msg := "👀 **Посещение сайта**\n" +
				"IP: " + ip + "\n" +
				"Город: " + city + "\n" +
				"Устройство: " + userAgent + "\n" +
				"Время: " + time.Now().Format("2006-01-02 15:04")

			uh.notifier.NotifyMessage(msg)
		}()
	}

	w.WriteHeader(http.StatusOK)
}
