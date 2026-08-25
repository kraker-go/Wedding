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

	if err = database.InitMigration(db); err != nil {
		return nil, fmt.Errorf("migration: %w", err)
	}

	tg, err := config.InitTelegramm()
	if err != nil {
		return nil, fmt.Errorf("telegram init: %w", err)
	}

	notif := handler.Telegramm(tg.Bot, tg.ChatID, nil)

	repo := repository.NewRepository(db)
	serv := service.NewUserService(repo)
	hand := handler.NewUserHandler(serv, logg, notif)

	notif.SetHandler(hand)

	rout, err := router.InitRouter(hand)
	if err != nil {
		return nil, err
	}

	logg.Info("Starting server on port", zap.String("port", port.Port))
	srv := &http.Server{
		Addr:    port.Port,
		Handler: rout,
	}

	srv, err = server.ConnectServer(port.Port, rout)
	if err != nil {
		return nil, err
	}

	return srv, nil
}
