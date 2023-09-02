package models

import (
	"gopkg.in/guregu/null.v4"
	"time"
)

type Rating struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Score      uint      `json:"Score" gorm:"type:uint"`
	Comment    string    `json:"comment" gorm:"type:varchar(255)"`
	PropertyID uint      `json:"property_id"`
	UserID     uint      `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  null.Time `json:"deleted_at"`
}
