package model

import (
	"time"
)

type Basket struct {
	ID        uint64  `json:"id,omitempty"`
	Data      JSONMap `json:"data,omitempty"`
	State     string  `json:"state,omitempty"`
	UserID    uint64
	User      User1     `gorm:"foreignKey:UserID"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type User1 struct {
	ID        uint64    `json:"id,omitempty"`
	UserName  string    `json:"username,omitempty"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
