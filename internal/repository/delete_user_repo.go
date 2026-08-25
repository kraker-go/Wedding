package repository

import (
	"context"
	"fmt"
	models "wedding/internal/domain"
)

func (ur *UserRepository) DeleteUser(ctx context.Context, id int) error {
	var user models.Guest

	err := ur.Db.QueryRowContext(ctx, GetUserByID, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.CreatedAt)
	if err != nil {
		return fmt.Errorf("UserRepository не найден гость: %w", err)
	}

	rows, err := ur.Db.ExecContext(ctx, Delete, id)
	if err != nil {
		return fmt.Errorf("Repository: не удалось удалить %w", err)
	}

	rowAffected, err := rows.RowsAffected()
	if err != nil {
		return fmt.Errorf("Repository %w", err)
	}
	if rowAffected == 0 {
		return fmt.Errorf("Repository: %w", models.ErrorsNotFound)
	}

	return nil
}
