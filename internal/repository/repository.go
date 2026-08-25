package repository

import (
	"database/sql"
)

type UserRepository struct {
	Db *sql.DB
}

func NewRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

// константы для базы данных

const GetUserID = "SELECT id, firstname, lastname, created_at FROM guests WHERE id = $1"

const GetUserName = "SELECT id FROM guests WHERE firstname = $1 and lastname = $2"

const AddUser = "INSERT INTO guests(firstname, lastname) VALUES ($1, $2) returning id, created_at"

const GetCount = "SELECT count(*) FROM guests"

const GetAll = "SELECT id, firstname, lastname, created_at FROM guests ORDER BY created_at DESC"

const Delete = "DELETE FROM guests WHERE id = $1"

const Update = "UPDATE guests SET firstname = $1, lastname = $2 WHERE id = $3"

const GetUserByID = "SELECT id, firstname, lastname, created_at FROM guests WHERE id = $1"

// константы для временной таблицы

const RequestUpdate = "INSERT INTO changeuser (guest_id, firstname, lastname) VALUES ($1, $2, $3)"

const TakeUserID = "SELECT id, guest_id, firstname, lastname, created_at FROM changeuser WHERE guest_id = $1"
