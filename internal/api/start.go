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

	logg.Info("конфигурация Postgres загружена")
	cfg, err := config.InitConfig()
	if err != nil {
		logg.Error("конфиг", zap.Error(err))
		return fmt.Errorf("api %w", err)
	}

	port, err := config.ServerInit()
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	logg.Info("конфигурация сервера загружена")

	db, err := database.ConnectPostgres(*cfg)
	if err != nil {
		logg.Error("Ошибка подключения к Postgres", zap.Error(err))
		return fmt.Errorf("не удалось подключиться к БД: %w", err)
	}

	defer db.Close()

	logg.Info("Postgres подключен")

	if err = database.InitMigration(db); err != nil {
		logg.Error("Ошибка миграции", zap.Error(err))
		return fmt.Errorf("не удалось применить миграции: %w", err)
	}

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
