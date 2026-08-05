package repository

import (
	"context"
	"fmt"
	"wedding/internal/domain"
)

func (ur *UserRepository) UpdateUser(ctx context.Context, user *models.Guest) error {
	rows, err := ur.Db.ExecContext(ctx, Update, user.FirstName, user.LastName, user.ID)
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
