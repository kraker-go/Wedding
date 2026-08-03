package main

import (
	"go.uber.org/zap"
	"wedding/internal/api"
	"wedding/internal/logger"
)

func main() {
	logg, err := logger.InitLogger()
	if err != nil {
		panic(err)
	}
	err = api.Start(logg)
	if err != nil {
		logg.Fatal("Ошибка запуска", zap.Error(err))
	}

	select {}
}
