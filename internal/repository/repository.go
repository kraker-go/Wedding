package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	models "wedding/internal/domain"
)

type UserRepository struct {
	Db *sql.DB
}

func NewRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

const GetUser = "SELECT id FROM guests WHERE firstname = $1 and lastname = $2"
const AddUser = "INSERT INTO guests(firstname, lastname) VALUES ($1, $2) returning id, created_at"

func (u UserRepository) AddUser(ctx context.Context, firstname, lastname string) (*models.Guest, error) {
	var user models.Guest
	err := u.Db.QueryRowContext(ctx, GetUser, firstname, lastname).Scan(&user.ID)
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

const GetAll = "SELECT firstname, lastname, guests.created_at FROM guests ORDER BY guests.created_at DESC"

func (u UserRepository) GetAllUsers(ctx context.Context) ([]models.Guest, error) {
	rows, err := u.Db.QueryContext(ctx, GetAll)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	defer rows.Close()
	var users []models.Guest
	for rows.Next() {
		var user models.Guest
		err = rows.Scan(&user.FirstName, &user.LastName)
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

const GetCount = "SELECT count(*) FROM guests"

func (u UserRepository) GetCount(ctx context.Context) (int, error) {
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
