package main

import (
	"log"

	"midterm/internal/domain/model"
	"midterm/internal/domain/repository/basketrepo"
	"midterm/internal/domain/repository/userrepo"
	"midterm/internal/infra/http/handler"
	basketsql "midterm/internal/infra/repository"
	usersql "midterm/internal/infra/userrepository"

	echojwt "github.com/labstack/echo-jwt/v4"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	userdb := "postgres"
	passworddb := "123456"
	portdb := "5432"
	dbname := "test111"
	dsn := "host=localhost user=" + userdb + " password=" + passworddb + " dbname=" + dbname + " port=" + portdb + " sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database %v", err)
	}

	if err := db.AutoMigrate(&model.User1{}, &basketsql.BasketDTO{}); err != nil {
		log.Fatalf("failed to run migrations %v", err)
	}

	app := echo.New()

	var repo basketrepo.Repository = basketsql.New(db)
	var userrepo userrepo.Repository = usersql.New(db)

	basketHandler := handler.NewBasket(repo)
	userHandler := handler.NewUser(userrepo)
	temp1 := app.Group("baskets/")
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(handler.JwtCustomClaims)
		},
		SigningKey: []byte("secret"),
	}
	temp1.Use(echojwt.WithConfig(config))
	basketHandler.Baskets(temp1)
	userHandler.Users(app.Group("users/"))

	if err := app.Start("0.0.0.0:1373"); err != nil {
		log.Fatalf("server failed to start %v", err)
	}
}
