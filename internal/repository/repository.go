package repository

import (
	"context"
	"database/sql"
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
