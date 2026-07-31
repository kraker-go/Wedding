package models

import (
	"errors"
	"time"
)

var (
	ErrorAlreadyExists = errors.New("Вы уже в списке гостей")
	ErrorNotString     = errors.New("поле не может быть пустым")
)

type Guest struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	CreatedAt time.Time `json:"created_at"`
}
