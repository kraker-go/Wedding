package repository

import (
	"context"
	"fmt"
	models "wedding/internal/domain"
)

func (u UserRepository) RequestUpdateUser(ctx context.Context, guestID int, user *models.Guest) error {
	_, err := u.Db.ExecContext(ctx, RequestUpdate, guestID, user.FirstName, user.LastName)
	if err != nil {
		return fmt.Errorf("UserRepository: не удалось сохранить запрос на изменение: %w", err)
	}

	return nil
}
