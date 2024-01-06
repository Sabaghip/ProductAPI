package basketrepo

import (
	"context"
	"errors"

	"midterm/internal/domain/model"
)

var ErrBasketIDDuplicate = errors.New("basket id already exists")

type GetCommand struct {
	ID *uint64
}

type UpdateCommand struct {
	ID    *uint64
	Data  *model.JSONMap
	State *string
}

type Repository interface {
	Add(ctx context.Context, model model.Basket, userID uint64) error
	Get(ctx context.Context, cmd GetCommand, userID uint64) ([]model.Basket, error)
	GetByID(ctx context.Context, cmd GetCommand, userID uint64) ([]model.Basket, error)
	Update(ctx context.Context, cmd UpdateCommand, userID uint64) error
	Delete(ctx context.Context, id uint64, userID uint64) error
}
