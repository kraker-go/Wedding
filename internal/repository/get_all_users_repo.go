package repository

import (
	"context"
	"fmt"
	"wedding/internal/domain"
)

func (u UserRepository) GetAllUsers(ctx context.Context) ([]models.Guest, error) {
	rows, err := u.Db.QueryContext(ctx, GetAll)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	defer rows.Close()

	var users []models.Guest

	for rows.Next() {

		var user models.Guest

		err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	return users, nil
}
