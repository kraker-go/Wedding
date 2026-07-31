package service

import (
	"context"
	models "wedding/internal/domain"
)

type CRUD interface {
	AddUser(context.Context, string, string) (*models.Guest, error)
}

type UserService struct {
	serv CRUD
}

func NewUserService(serv CRUD) *UserService {
	return &UserService{serv}
}

func (us *UserService) AddUser(ctx context.Context, firstname, lastname string) (*models.Guest, error) {
	if firstname == "" || lastname == "" {
		return &models.Guest{}, models.ErrorNotString
	}

	return us.serv.AddUser(ctx, firstname, lastname)
}
