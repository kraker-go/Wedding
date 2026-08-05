package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"wedding/internal/domain"
)

func (u UserRepository) GetCountUsers(ctx context.Context) (int, error) {
	var count int
	err := u.Db.QueryRowContext(ctx, GetCount).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrorsNotFound
		}

		return 0, fmt.Errorf(err.Error())
	}

	return count, nil
}
