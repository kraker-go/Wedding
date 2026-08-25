package repository

import (
	"context"
	"database/sql"
	"fmt"
	"wedding/internal/domain"
)

func (u UserRepository) AddUser(ctx context.Context, firstname, lastname string) (*models.Guest, error) {
	var user models.Guest

	err := u.Db.QueryRowContext(ctx, GetUserName, firstname, lastname).Scan(&user.ID)
	if err == nil {
		return &models.Guest{}, models.ErrorAlreadyExists
	}

	if err != sql.ErrNoRows {
		return &models.Guest{}, fmt.Errorf("ошибка запроса в базу данных: %w", err)
	}

	if err == sql.ErrNoRows {
		err = u.Db.QueryRowContext(ctx, AddUser, firstname, lastname).Scan(&user.ID, &user.CreatedAt)
		if err != nil {
			return &models.Guest{}, fmt.Errorf("Ошибка записи гостя: %w", err)
		}
	}
	return &user, err
}
