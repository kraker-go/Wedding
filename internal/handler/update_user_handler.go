package handler

import (
	"context"
)

func (uh *UserHandler) UpdateUserHandler(ctx context.Context, id int) error {
	err := uh.hand.UpdateUser(ctx, id)
	if err != nil {
		return err
	}

	if uh.notifier != nil {
		go uh.notifier.NotifyMessage("✅ Пользователь успешно обновлен! ")
	}
	return nil
}
