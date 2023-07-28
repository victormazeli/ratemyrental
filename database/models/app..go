package models

import (
	"gorm.io/gorm"
	"time"
)

type App struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Maintenance  string         `json:"maintenance" gorm:"type:varchar(255)"`
	LoginCounter uint           `json:"login_counter" gorm:"type:int"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
