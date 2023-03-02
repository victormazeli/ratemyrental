package models

import (
	"gopkg.in/guregu/null.v4"
	"time"
)

type FavoriteProperty struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id"`
	PropertyID uint      `json:"property_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  null.Time `json:"deleted_at"`
}
