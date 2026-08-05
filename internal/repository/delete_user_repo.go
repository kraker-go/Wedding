package repository

import (
	"context"
	"fmt"
	models "wedding/internal/domain"
)

func (ur *UserRepository) DeleteUser(ctx context.Context, id int) error {
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
