package handler

import (
	"errors"
	"math/rand"
	"net/http"
	"strconv"

	"midterm/internal/domain/model"
	"midterm/internal/domain/repository/basketrepo"
	"midterm/internal/domain/repository/userrepo"
	"midterm/internal/infra/http/request"

	"github.com/labstack/echo/v4"

	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Basket struct {
	repo basketrepo.Repository
}
type User struct {
	repo userrepo.Repository
}

func NewBasket(repo basketrepo.Repository) *Basket {
	return &Basket{
		repo: repo,
	}
}

func NewUser(repo userrepo.Repository) *User {
	return &User{
		repo: repo,
	}
}

func (b *Basket) GetByID(c echo.Context) error {
	if temp := c.Get("user"); temp != nil {

		user := temp.(*jwt.Token)
		claims := user.Claims.(*JwtCustomClaims)
		userID := claims.ID
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return echo.ErrBadRequest
		}
		baskets, err := b.repo.GetByID(c.Request().Context(), basketrepo.GetCommand{
			ID: &id,
		}, userID)
		if err != nil {
			return echo.ErrInternalServerError
		}
		if len(baskets) == 0 {
			return echo.ErrNotFound
		}

		if len(baskets) > 1 {
			return echo.ErrInternalServerError
		}

		return c.JSON(http.StatusOK, baskets[0])
	}
	return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid jwt token")

}

func (b *Basket) Get(c echo.Context) error {
	if temp := c.Get("user"); temp != nil {

		user := temp.(*jwt.Token)
		claims := user.Claims.(*JwtCustomClaims)
		userID := claims.ID
		baskets, err := b.repo.Get(c.Request().Context(), basketrepo.GetCommand{}, userID)
		if err != nil {
			return echo.ErrInternalServerError
		}
		if len(baskets) == 0 {
			return echo.ErrNotFound
		}
		return c.JSON(http.StatusOK, baskets)
	}
	return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid jwt token")
}

func (b *Basket) Create(c echo.Context) error {
	if temp := c.Get("user"); temp != nil {

		user := temp.(*jwt.Token)
		claims := user.Claims.(*JwtCustomClaims)
		userID := claims.ID
		var req request.BasketCreate

		if err := c.Bind(&req); err != nil {
			return echo.ErrBadRequest
		}
		// we have the filled request
		if err := req.Validate(); err != nil {
			return echo.ErrBadRequest
		}
		// nolint: gosec, gomnd
		id := rand.Uint64() % 1_000_000
		if err := b.repo.Add(c.Request().Context(), model.Basket{
			ID:    id,
			Data:  req.Data,
			State: req.State,
		}, userID); err != nil {

			if errors.Is(err, basketrepo.ErrBasketIDDuplicate) {
				return echo.ErrBadRequest
			}

			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusCreated, id)
	}
	return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid jwt token")
}

func (b *Basket) Update(c echo.Context) error {
	if temp := c.Get("user"); temp != nil {

		user := temp.(*jwt.Token)
		claims := user.Claims.(*JwtCustomClaims)
		userID := claims.ID
		var req request.BasketCreate

		if err := c.Bind(&req); err != nil {
			return echo.ErrBadRequest
		}
		// we have the filled request
		if err := req.Validate(); err != nil {
			return echo.ErrBadRequest
		}

		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return echo.ErrBadRequest
		}

		baskets, err := b.repo.GetByID(c.Request().Context(), basketrepo.GetCommand{
			ID: &id,
		}, userID)
		if err != nil {
			return echo.ErrInternalServerError
		}
		if len(baskets) == 0 {
			return echo.ErrNotFound
		}

		if len(baskets) > 1 {
			return echo.ErrInternalServerError
		}
		if baskets[0].State == "COMPLETED" {
			return echo.NewHTTPError(400, "Completed baskets cannot change!")
		}

		err = b.repo.Update(c.Request().Context(), basketrepo.UpdateCommand{
			ID:    &id,
			Data:  &req.Data,
			State: &req.State,
		}, userID)
		if err != nil {
			return echo.ErrInternalServerError
		}
		baskets, err = b.repo.GetByID(c.Request().Context(), basketrepo.GetCommand{
			ID: &id,
		}, userID)
		if err != nil {
			return echo.ErrInternalServerError
		}
		if len(baskets) == 0 {
			return echo.ErrNotFound
		}

		if len(baskets) > 1 {
			return echo.ErrInternalServerError
		}

		return c.JSON(http.StatusOK, baskets[0])
	}
	return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid jwt token")

}

func (b *Basket) Delete(c echo.Context) error {
	if temp := c.Get("user"); temp != nil {

		user := temp.(*jwt.Token)
		claims := user.Claims.(*JwtCustomClaims)
		userID := claims.ID
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return echo.ErrBadRequest
		}

		baskets, err := b.repo.GetByID(c.Request().Context(), basketrepo.GetCommand{
			ID: &id,
		}, userID)
		if err != nil {
			return echo.ErrInternalServerError
		}
		if len(baskets) == 0 {
			return echo.ErrNotFound
		}

		if len(baskets) > 1 {
			return echo.ErrInternalServerError
		}

		err = b.repo.Delete(c.Request().Context(), id, userID)
		if err != nil {
			return echo.ErrInternalServerError
		}
		return nil
	}
	return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid jwt token")
}

type JwtCustomClaims struct {
	ID uint64 `json:"id"`
	jwt.RegisteredClaims
}

func (u *User) login(c echo.Context) error {
	var req request.UserCreate

	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}
	// we have the filled request
	if err := req.Validate(); err != nil {
		return err
	}
	// nolint: gosec, gomnd
	users, err := u.repo.Login(c.Request().Context(), userrepo.LoginCommand{
		Username: &req.UserName,
		Password: &req.Password,
	})
	if err != nil {

		if errors.Is(err, userrepo.ErrUserIDDuplicate) {
			return echo.ErrBadRequest
		}

		return echo.ErrInternalServerError
	}
	if len(users) < 1 {
		return echo.NewHTTPError(400, "invlid creditionals")
	}

	claims := &JwtCustomClaims{
		users[0].ID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 10)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}
func (u *User) signup(c echo.Context) error {
	var req request.UserCreate

	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}
	// we have the filled request
	if err := req.Validate(); err != nil {
		return echo.ErrBadRequest
	}
	flag, err := u.repo.CheckUsername(req.UserName)
	if err != nil {
		return echo.ErrInternalServerError
	}
	if flag {
		return echo.NewHTTPError(400, "Username is unavailable")
	}
	// nolint: gosec, gomnd
	id := rand.Uint64() % 1_000_000
	if err := u.repo.Signup(c.Request().Context(), model.User1{
		ID:       id,
		UserName: req.UserName,
		Password: req.Password,
	}); err != nil {

		if errors.Is(err, basketrepo.ErrBasketIDDuplicate) {
			return echo.ErrBadRequest
		}

		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusCreated, id)
}

func (b *Basket) Baskets(g *echo.Group) {
	g.GET("", b.Get)
	g.POST("", b.Create)
	g.GET(":id", b.GetByID)
	g.PATCH(":id", b.Update)
	g.DELETE(":id", b.Delete)
}

func (u *User) Users(g *echo.Group) {
	g.POST("login", u.login)
	g.POST("signup", u.signup)
}
