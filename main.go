package main

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"wedding/internal/api"
	"wedding/internal/logger"
)

func main() {
	logg, err := logger.InitLogger()
	if err != nil {
		panic(err)
	}
	err = api.Start(logg)
	// Принудительно выводим ошибку в stderr (гарантированно попадёт в логи)
	fmt.Fprintf(os.Stderr, "FATAL: %v\n", err)
	logg.Fatal("Ошибка запуска", zap.Error(err))
	select {}
}

// перенес
