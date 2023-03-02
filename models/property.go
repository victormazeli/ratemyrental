package models

import (
	"database/sql"
	"encoding/json"
	"gopkg.in/guregu/null.v4"
	_ "gopkg.in/guregu/null.v4"
	"time"
)

type Property struct {
	ID                   uint             `json:"id" gorm:"primaryKey"`
	PropertyTitle        string           `json:"property_title" gorm:"type:varchar(255)"`
	Description          null.String      `json:"description" gorm:"type:text"`
	Floors               null.String      `json:"floors" gorm:"type:varchar(255)"`
	NumberOfRooms        null.String      `json:"number_of_rooms" gorm:"type:varchar(255)"`
	BedRooms             null.String      `json:"bed_rooms" gorm:"type:varchar(255)"`
	BathRooms            null.String      `json:"bath_rooms" gorm:"type:varchar(255)"`
	CloakRooms           null.String      `json:"cloak_rooms" gorm:"type:varchar(255)"`
	UtilityRooms         null.String      `json:"utility_rooms" gorm:"type:varchar(255)"`
	Conservatory         null.String      `json:"conservatory" gorm:"type:varchar(255)"`
	EntranceHall         null.String      `json:"entrance_hall" gorm:"type:varchar(255)"`
	FrontYard            null.String      `json:"front_yard" gorm:"type:varchar(255)"`
	MudRoom              null.String      `json:"mud_room" gorm:"type:varchar(255)"`
	FurnishedRoom        null.String      `json:"furnished_room" gorm:"type:varchar(255)"`
	Garden               null.String      `json:"garden" gorm:"type:varchar(255)"`
	Garage               null.String      `json:"garage" gorm:"type:varchar(255)"`
	Ensuite              null.String      `json:"ensuite" gorm:"type:varchar(255)"`
	CharacterFeature     null.String      `json:"character_feature" gorm:"type:text"`
	EpcRatings           null.String      `json:"epc_ratings" gorm:"type:varchar(255)"`
	PetsAllowed          null.String      `json:"pets_allowed" gorm:"type:varchar(255)"`
	SmokingAllowed       null.String      `json:"smoking_allowed" gorm:"type:varchar(255)"`
	DssAllowed           null.String      `json:"dss_allowed" gorm:"type:varchar(255)"`
	SharersAllowed       null.String      `json:"sharers_allowed" gorm:"type:varchar(255)"`
	AverageRating        null.Float       `json:"average_rating" gorm:"type:float"`
	Location             null.String      `json:"location" gorm:"type:text"`
	State                string           `json:"state" gorm:"type:varchar(255)"`
	Country              string           `json:"country" gorm:"type:varchar(255)"`
	PostalCode           string           `json:"postal_code" gorm:"type:varchar(255)"`
	PropertyImages       []*PropertyImage `json:"property_images"`
	PropertyType         string           `json:"property_type"`
	PropertyDetachedType string           `json:"property_detached_type"`
	UserID               uint             `json:"user_id"`
	IsVisible            int64            `json:"is_visible" gorm:"default:0"`
	CreatedAt            time.Time        `json:"created_at"`
	UpdatedAt            time.Time        `json:"updated_at"`
	DeletedAt            null.Time        `json:"deleted_at"`
}

type NullString struct {
	sql.NullString
}

// MarshalJSON for NullString
func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}
