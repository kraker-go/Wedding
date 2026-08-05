package models

import (
	"errors"
	"time"
)

var (
	ErrorAlreadyExists = errors.New("Вы уже в списке гостей")
	ErrorsNotFound     = errors.New("запись не найдена")
	ErrorNotFound      = errors.New("запись не найдена")

	ErrorEmptyName    = errors.New("имя или фамилия не могут быть пустыми")
	ErrorNameTooLong  = errors.New("имя или фамилия слишком длинные")
	ErrorInvalidChars = errors.New("имя может содержать только буквы, дефис и апостроф (без пробелов)")
)

type Guest struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	CreatedAt time.Time `json:"created_at"`
}
