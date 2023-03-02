package dtos

import "time"

type PropertyDTO struct {
	Id                   uint                     `json:"id"`
	PropertyTitle        string                   `json:"property_title"`
	Description          string                   `json:"description"`
	Floors               string                   `json:"floors"`
	NumberOfRooms        string                   `json:"number_of_rooms"`
	BedRooms             string                   `json:"bed_rooms"`
	BathRooms            string                   `json:"bath_rooms"`
	CloakRooms           string                   `json:"cloak_rooms"`
	UtilityRooms         string                   `json:"utility_rooms"`
	Conservatory         string                   `json:"conservatory"`
	EntranceHall         string                   `json:"entrance_hall"`
	FrontYard            string                   `json:"front_yard"`
	MudRoom              string                   `json:"mud_room"`
	FurnishedRoom        string                   `json:"furnished_room"`
	Garden               string                   `json:"garden"`
	Garage               string                   `json:"garage"`
	Ensuite              string                   `json:"ensuite"`
	CharacterFeature     string                   `json:"character_feature"`
	EpcRatings           string                   `json:"epc_ratings"`
	PetsAllowed          string                   `json:"pets_allowed"`
	SmokingAllowed       string                   `json:"smoking_allowed"`
	DssAllowed           string                   `json:"dss_allowed"`
	SharersAllowed       string                   `json:"sharers_allowed"`
	AverageRating        float32                  `json:"average_rating"`
	Location             string                   `json:"location"`
	State                string                   `json:"state"`
	Country              string                   `json:"country"`
	PostalCode           string                   `json:"postal_code"`
	PropertyImages       []*PropertyImageDTO      `json:"property_images,omitempty"`
	PropertyType         *PropertyTypeDTO         `json:"property_type,omitempty"`
	PropertyDetachedType *PropertyDetachedTypeDTO `json:"property_detached_type,omitempty"`
	UserID               uint                     `json:"user_id"`
	CreatedAt            time.Time                `json:"created_at"`
	UpdatedAt            time.Time                `json:"updated_at"`
	DeletedAt            time.Time                `json:"deleted_at"`
}

type PropertyImageDTO struct {
	Id         uint      `json:"id,omitempty"`
	ImageUrl   string    `json:"image_url,omitempty"`
	PropertyID uint      `json:"property_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at"`
}

type PropertyTypeDTO struct {
	Id         uint      `json:"id,omitempty"`
	Type       string    `json:"type,omitempty"`
	PropertyID uint      `json:"property_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at"`
}

type PropertyDetachedTypeDTO struct {
	Id         uint      `json:"id,omitempty"`
	Type       string    `json:"type,omitempty"`
	PropertyID uint      `json:"property_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at"`
}
