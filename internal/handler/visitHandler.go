package handler

import (
	"net/http"
	"time"
	"wedding/internal/handler/geo"
)

func (uh *UserHandler) VisitHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	uh.logg.Info("Посещение сайта")

	go func() {
		ip := getClientIP(r)
		userAgent := r.UserAgent()
		city, _ := geo.GetCityByIP(ip)

		msg := "👀 **Новое посещение сайта**\n" +
			"IP: " + ip + "\n" +
			"Город: " + city + "\n" +
			"Устройство: " + userAgent + "\n" +
			"Время: " + time.Now().Format("2006-01-02 15:04")

		// Telegram
		if uh.tgNotifier != nil {
			uh.tgNotifier.NotifyMessage(msg)
		}
	}()
	w.WriteHeader(http.StatusOK)
}
