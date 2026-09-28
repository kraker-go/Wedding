package handler

import (
	"context"
	"go.uber.org/zap"
)

func (uh *UserHandler) DeleteUserHandler(ctx context.Context, id int) error {
	uh.logg.Info("DeleteUserHandler выполняем удаление пользователя")
	err := uh.hand.DeleteUser(ctx, id)
	if err != nil {
		uh.logg.Error("DeleteUserHandler Handler", zap.Error(err))
		return err
	}

	uh.logg.Info("DeleteUserHandler пользователь удален")

	if uh.notifier != nil {
		go uh.notifier.NotifyMessage("✅ Пользователь успешно удалён! ")
	}

	return nil
}
