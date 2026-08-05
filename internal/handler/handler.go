package handler

import (
	"go.uber.org/zap"
	"wedding/internal/service"
)

type UserHandler struct {
	hand *service.UserService
	logg *zap.Logger
}

func NewUserHandler(hand *service.UserService, logg *zap.Logger) *UserHandler {
	return &UserHandler{hand: hand, logg: logg}
}
