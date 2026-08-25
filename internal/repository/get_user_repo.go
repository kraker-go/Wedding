package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	models "wedding/internal/domain"
)

func (u UserRepository) GetUser(ctx context.Context, userID int) (*models.Guest, error) {
	var user models.Guest
	err := u.Db.QueryRowContext(ctx, GetUserID, userID).Scan(&user.ID, &user.FirstName, &user.LastName, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrorsNotFound
		}
		return nil, fmt.Errorf("GetUser: %w", err)
	}

	return &user, nil
}
