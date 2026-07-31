package api

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"wedding/internal/config"
	"wedding/internal/database"
	"wedding/internal/handler"
	"wedding/internal/repository"
	"wedding/internal/router"
	"wedding/internal/service"
)

func Start(logg *zap.Logger) error {
	cfg, err := config.CfgInit()
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	logg.Info("конфигурация Postgres загружена")

	port, err := config.ServerInit()
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	logg.Info("конфигурация сервера загружена")

	db, err := database.ConnectPostgres(*cfg)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	defer db.Close()

	logg.Info("Postgres подключен")

	err = database.InitMigration(db)

	repo := repository.NewRepository(db)
	serv := service.NewUserService(repo)
	hand := handler.NewUserHandler(serv, logg)

	rout, err := router.InitRouter(hand)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	startServer, err := database.ConnectServer(port.Port, rout)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	logg.Info("Сервер запущен")

	if err = startServer.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		return fmt.Errorf("ошибка запуска сервер %w", err)
	}

	return nil
}
