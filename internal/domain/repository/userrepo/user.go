package userrepo

import (
	"context"
	"errors"

	"midterm/internal/domain/model"
)

var ErrUserIDDuplicate = errors.New("user id already exists")

type LoginCommand struct {
	Username *string
	Password *string
}

type Repository interface {
	Login(ctx context.Context, model LoginCommand) ([]model.User1, error)
	Signup(ctx context.Context, model model.User1) error
	CheckUsername(username string) (bool, error)
}
