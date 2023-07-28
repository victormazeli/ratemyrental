package models

import (
	"gopkg.in/guregu/null.v4"
	"time"
)

type PropertyImage struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	ImageUrl   string    `json:"image_url" gorm:"type:varchar(255)"`
	PropertyID uint      `json:"property_id" gorm:"foreignkey:PropertyID"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  null.Time `json:"deleted_at"`
}
