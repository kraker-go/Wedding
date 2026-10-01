package handler

import (
	"go.uber.org/zap"
	models "wedding/internal/domain"
	"wedding/internal/service"
)

type UserHandler struct {
	hand        *service.UserService
	logg        *zap.Logger
	tgNotifier  *Notifier
	maxNotifier interface {
		NotifyDelete(models.Guest)
		NotifyAdd(models.Guest)
		NotifyUpdate(models.Guest, models.Guest)
		NotifyMessage(string)
	}
}

func NewUserHandler(hand *service.UserService, logg *zap.Logger, not *Notifier) *UserHandler {
	return &UserHandler{hand: hand, logg: logg, tgNotifier: not}
}

func (uh *UserHandler) SetMaxNotifier(max interface {
	NotifyDelete(models.Guest)
	NotifyAdd(models.Guest)
	NotifyUpdate(models.Guest, models.Guest)
	NotifyMessage(string)
},
) {
	uh.maxNotifier = max
}
