package repository

import (
	"context"
	"fmt"
	models "wedding/internal/domain"
)

func (u UserRepository) RequestDeleteUser(ctx context.Context, guestID int) error {

	result, err := u.Db.ExecContext(ctx, Delete, guestID)
	if err != nil {
		return fmt.Errorf("UserRepository: ошибка удаления гостя: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("UserRepository: ошибка проверки удаления гостя: %w", err)
	}

	if rows == 0 {
		return models.ErrorsNotFound
	}

	return nil
}
