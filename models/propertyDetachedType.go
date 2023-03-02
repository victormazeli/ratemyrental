package models

import (
	"gopkg.in/guregu/null.v4"
	"time"
)

type PropertyDetachedType struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Type      string    `json:"type" gorm:"type:varchar(255)"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at"`
}
