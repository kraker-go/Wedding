package repository

import (
	"context"
	"fmt"
	"wedding/internal/domain"
)

func (ur *UserRepository) UpdateUser(ctx context.Context, id int) error {
	var user models.Change_user

	err := ur.Db.QueryRowContext(ctx, TakeUserID, id).Scan(&user.ID, &user.Guest_id, &user.FirstName, &user.LastName, &user.CreatedAt)
	if err != nil {
		return fmt.Errorf("ошибка поиска пользователя %w", err)
	}

	rows, err := ur.Db.ExecContext(ctx, Update, user.FirstName, user.LastName, id)
	if err != nil {
		return fmt.Errorf("Repository %w", err)
	}

	rowAffected, err := rows.RowsAffected()
	if err != nil {
		return fmt.Errorf("Repository %w", err)
	}

	if rowAffected == 0 {
		return models.ErrorsNotFound
	}

	return nil
}
