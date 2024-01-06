package usersql

import (
	"context"
	"time"

	"midterm/internal/domain/model"
	"midterm/internal/domain/repository/userrepo"

	"gorm.io/gorm"
)

type UserDTO struct {
	model.User1
}

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Signup(ctx context.Context, model model.User1) error {
	var result UserDTO
	if err := r.db.Raw("INSERT INTO user1 VALUES (?,?,?,?,?)", model.ID, model.UserName, model.Password, time.Now(), time.Now()).Scan(&result).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) Login(ctx context.Context, cmd userrepo.LoginCommand) ([]model.User1, error) {
	var userDTOs []UserDTO
	if err := r.db.Table("user1").Select("id", "User_name", "created_at", "updated_at").Where("User_Name = ? and Password = ?", cmd.Username, cmd.Password).Scan(&userDTOs).Error; err != nil {
		return nil, err
	}
	users := make([]model.User1, len(userDTOs))
	for index, dto := range userDTOs {
		users[index] = dto.User1
	}
	return users, nil
}

func (r *Repository) CheckUsername(username string) (bool, error) {
	var result []UserDTO
	if err := r.db.Raw("Select * from user1 WHERE user_name=?", username).Scan(&result).Error; err != nil {
		return true, err
	}
	users := make([]model.User1, len(result))
	for index, dto := range result {
		users[index] = dto.User1
	}
	if len(users) > 0 {
		return true, nil
	}
	return false, nil
}
