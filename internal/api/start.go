// internal/api/start.go
package api

import (
	"fmt"
	"wedding/internal/server"

	"go.uber.org/zap"
	"net/http"
	"wedding/internal/config"
	"wedding/internal/database"
	"wedding/internal/handler"
	"wedding/internal/repository"
	"wedding/internal/router"
	"wedding/internal/service"
)

// StartServer создаёт и возвращает настроенный *http.Server
func StartServer(logg *zap.Logger) (*http.Server, error) {
	cfg, err := config.InitConfig()
	if err != nil {
		return nil, fmt.Errorf("config init: %w", err)
	}

	port, err := config.ServerInit()
	if err != nil {
		return nil, err
	}

	db, err := database.ConnectPostgres(*cfg)
	if err != nil {
		return nil, fmt.Errorf("db connect: %w", err)
	}
	// не закрываем здесь, закроем при остановке приложения

	if err = database.InitMigration(db); err != nil {
		return nil, fmt.Errorf("migration: %w", err)
	}

	repo := repository.NewRepository(db)
	serv := service.NewUserService(repo)
	hand := handler.NewUserHandler(serv, logg)

	rout, err := router.InitRouter(hand)
	if err != nil {
		return nil, err
	}

	srv, err := server.ConnectServer(port.Port, rout) // <- теперь это импорт server (нет цикла)
	if err != nil {
		return nil, err
	}

	// Возвращаем сервер, но не запускаем его
	return srv, nil
}
