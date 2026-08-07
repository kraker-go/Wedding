package handler

import (
	"go.uber.org/zap"
	"wedding/internal/service"
)

type UserHandler struct {
	hand     *service.UserService
	logg     *zap.Logger
	notifier *Notifier
}

func NewUserHandler(hand *service.UserService, logg *zap.Logger, not *Notifier) *UserHandler {
	return &UserHandler{hand: hand, logg: logg, notifier: not}
}
