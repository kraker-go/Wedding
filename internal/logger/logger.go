package logger

import (
	"fmt"
	"go.uber.org/zap"
)

func InitLogger() (*zap.Logger, error) {
	Logger, err := zap.NewProduction()
	if err != nil {
		fmt.Errorf("Error initializing logger: %w", err)
	}
	return Logger, nil
}
