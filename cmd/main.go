package main

import (
	"fmt"
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
		fmt.Println(err)
	}
	select {}
}
