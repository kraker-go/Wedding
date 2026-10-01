package handler

import (
	"context"
	"go.uber.org/zap"
)

func (uh *UserHandler) DeleteUserHandler(ctx context.Context, id int) error {
	uh.logg.Info("DeleteUserHandler выполняем удаление пользователя")

	user, err := uh.hand.GetUser(ctx, id)
	if err != nil {
		uh.logg.Error(
			"не удалось получить пользователя перед удалением",
			zap.Error(err),
		)

		return err
	}

	err = uh.hand.DeleteUser(ctx, id)
	if err != nil {
		uh.logg.Error("DeleteUserHandler Handler", zap.Error(err))
		return err
	}

	uh.logg.Info("DeleteUserHandler пользователь удален")

	if uh.tgNotifier != nil {
		go uh.tgNotifier.NotifyMessage("✅ Пользователь: " + user.FirstName + " " + user.LastName + ": успешно удалён!")
	}

	if uh.maxNotifier != nil {
		go uh.maxNotifier.NotifyMessage("✅ Пользователь: " + user.FirstName + " " + user.LastName + " - успешно удалён!")
	}

	return nil
}
