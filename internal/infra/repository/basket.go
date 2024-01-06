package basketsql

import (
	"context"
	"time"

	"midterm/internal/domain/model"
	"midterm/internal/domain/repository/basketrepo"

	"gorm.io/gorm"
)

type BasketDTO struct {
	model.Basket
}

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Add(ctx context.Context, model model.Basket, userID uint64) error {
	var result BasketDTO
	r.db.Raw("INSERT INTO basket_dtos VALUES (?,?,?,?,?,?)", model.ID, model.Data, model.State, userID, time.Now(), time.Now()).Scan(&result)
	return nil
}

func (r *Repository) Get(_ context.Context, cmd basketrepo.GetCommand, userID uint64) ([]model.Basket, error) {
	var basketDTOs []BasketDTO
	if err := r.db.Table("basket_dtos").Select("id", "Data", "State", "created_at", "updated_at").Where("user_id = ?", userID).Scan(&basketDTOs).Error; err != nil {
		return nil, err
	}
	baskets := make([]model.Basket, len(basketDTOs))
	for index, dto := range basketDTOs {
		baskets[index] = dto.Basket
	}
	return baskets, nil
}

func (r *Repository) GetByID(_ context.Context, cmd basketrepo.GetCommand, userID uint64) ([]model.Basket, error) {
	var basketDTOs []BasketDTO
	if err := r.db.Table("basket_dtos").Select("id", "Data", "State", "created_at", "updated_at").Where("ID = ? and user_id=?", cmd.ID, userID).Scan(&basketDTOs).Error; err != nil {
		return nil, err
	}
	baskets := make([]model.Basket, len(basketDTOs))
	for index, dto := range basketDTOs {
		baskets[index] = dto.Basket
	}
	return baskets, nil
}

func (r *Repository) Update(_ context.Context, cmd basketrepo.UpdateCommand, userID uint64) error {
	var result BasketDTO
	if err := r.db.Raw("UPDATE basket_dtos SET Data=?, State=?, updated_at=? WHERE ID=? and user_id = ?", cmd.Data, cmd.State, time.Now(), cmd.ID, userID).Scan(&result).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) Delete(_ context.Context, id uint64, userID uint64) error {
	var result BasketDTO
	if err := r.db.Raw("DELETE FROM basket_dtos WHERE ID=? and user_id=?", id, userID).Scan(&result).Error; err != nil {
		return err
	}
	return nil
}
