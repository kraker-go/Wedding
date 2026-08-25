package service

import (
	"context"
	"fmt"
	"regexp"
	models "wedding/internal/domain"
)

type CRUD interface {
	AddUser(context.Context, string, string) (*models.Guest, error)
	GetUser(context.Context, int) (*models.Guest, error)
	GetAllUsers(context.Context) ([]models.Guest, error)
	GetCountUsers(context.Context) (int, error)
	RequestUpdateUser(context.Context, int, *models.Guest) error
	UpdateUser(context.Context, int) error
	DeleteUser(context.Context, int) error
	RequestDeleteUser(context.Context, int) error
}

type UserService struct {
	serv CRUD
}

func NewUserService(serv CRUD) *UserService {
	return &UserService{serv}
}

var validNameRegex = regexp.MustCompile(`^[\p{L}\-']+$`)

func validateName(name string, maxLen int) error {
	if name == "" {
		return models.ErrorEmptyName
	}
	if !validNameRegex.MatchString(name) {
		return models.ErrorInvalidChars
	}
	if len(name) > maxLen {
		return models.ErrorNameTooLong
	}
	return nil
}

func (us *UserService) AddUser(ctx context.Context, firstname, lastname string) (*models.Guest, error) {
	if err := validateName(firstname, 24); err != nil {
		return nil, err
	}

	if err := validateName(lastname, 30); err != nil {
		return nil, err
	}

	user, err := us.serv.AddUser(ctx, firstname, lastname)
	if err != nil {
		return nil, fmt.Errorf("Service: %w", err)
	}

	return user, nil
}
func (us *UserService) RequestUpdateUser(ctx context.Context, guestID int, user *models.Guest) error {

	if err := validateName(user.FirstName, 24); err != nil {
		return err
	}

	if err := validateName(user.LastName, 30); err != nil {
		return err
	}

	err := us.serv.RequestUpdateUser(ctx, guestID, user)
	if err != nil {
		return fmt.Errorf("Service: %w", err)
	}

	return nil
}

func (us *UserService) UpdateUser(ctx context.Context, id int) error {

	err := us.serv.UpdateUser(ctx, id)
	if err != nil {
		return fmt.Errorf("Service %w", err)
	}

	return nil
}

func (us *UserService) GetUser(ctx context.Context, id int) (*models.Guest, error) {
	user, err := us.serv.GetUser(ctx, id)
	return user, err
}
func (us *UserService) GetAllUsers(ctx context.Context) ([]models.Guest, error) {
	return us.serv.GetAllUsers(ctx)
}

func (us *UserService) GetCountUsers(ctx context.Context) (int, error) {
	return us.serv.GetCountUsers(ctx)
}

func (us *UserService) RequestDeleteUser(ctx context.Context, id int) error {
	return us.serv.RequestDeleteUser(ctx, id)
}

func (us *UserService) DeleteUser(ctx context.Context, id int) error {
	return us.serv.DeleteUser(ctx, id)
}
