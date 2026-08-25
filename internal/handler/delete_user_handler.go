package handler

import "context"

func (uh *UserHandler) DeleteUserHandler(ctx context.Context, id int) error {
	err := uh.hand.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	if uh.notifier != nil {
		go uh.notifier.NotifyMessage("✅ Пользователь успешно удалён! ")
	}
	return nil
}
