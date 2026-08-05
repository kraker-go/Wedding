// internal/app/run.go
package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"wedding/internal/api"
	"wedding/internal/logger"
)

func Run() {
	logg, err := logger.InitLogger()
	if err != nil {
		panic(err)
	}

	// Создаём сервер
	srv, err := api.StartServer(logg)
	if err != nil {
		logg.Fatal("Ошибка создания сервера", zap.Error(err))
	}

	// Запускаем в горутине
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logg.Fatal("Ошибка запуска сервера", zap.Error(err))
		}
	}()

	logg.Info("Сервер запущен", zap.String("addr", srv.Addr))

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logg.Info("Получен сигнал завершения, останавливаем сервер...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logg.Error("Ошибка при остановке сервера", zap.Error(err))
	} else {
		logg.Info("Сервер остановлен")
	}
}
