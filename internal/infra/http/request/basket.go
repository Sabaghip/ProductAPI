package request

import (
	"fmt"

	"midterm/internal/domain/model"

	"github.com/labstack/echo/v4"

	"github.com/go-playground/validator/v10"
)

type BasketCreate struct {
	Data  model.JSONMap `json:"data,omitempty"   validate:"required"`
	State string        `json:"state,omitempty" validate:"required"`
}

func (bc BasketCreate) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(bc); err != nil {
		return fmt.Errorf("create request validation failed %w", err)
	}
	if bc.State == "COMPLETED" || bc.State == "PENDING" {
		return nil
	}

	return fmt.Errorf("State should be COMPLETED or PENDING")
}

type UserCreate struct {
	UserName string `json:"username,omitempty"   validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}

func (uc UserCreate) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(uc); err != nil {
		return echo.NewHTTPError(400, "create request validation failed %w", err)
	}
	if len(uc.Password) < 8 {
		return echo.NewHTTPError(400, "password should be more than 7 chars")
	}

	return nil
}
