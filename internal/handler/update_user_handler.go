package handler

import (
	"context"
)

func (uh *UserHandler) UpdateUserHandler(ctx context.Context, id int) error {
	err := uh.hand.UpdateUser(ctx, id)
	if err != nil {
		return err
	}

	if uh.tgNotifier != nil {
		go uh.tgNotifier.NotifyMessage("✅ Гость успешно обновлен! ")
	}

	if uh.maxNotifier != nil {
		go uh.maxNotifier.NotifyMessage("✅ Гость успешно обновлен! ")
	}
	return nil
}
