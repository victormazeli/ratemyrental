package input

type PropertyUpdateInput struct {
	PropertyTitle        string `json:"property_title"`
	Description          string `json:"description"`
	Floors               string `json:"floors"`
	NumberOfRooms        string `json:"number_of_rooms"`
	BedRooms             string `json:"bed_rooms"`
	BathRooms            string `json:"bath_rooms"`
	CloakRooms           string `json:"cloak_rooms"`
	UtilityRooms         string `json:"utility_rooms"`
	Conservatory         string `json:"conservatory"`
	EntranceHall         string `json:"entrance_hall"`
	FrontYard            string `json:"front_yard"`
	MudRoom              string `json:"mud_room"`
	FurnishedRoom        string `json:"furnished_room"`
	Garden               string `json:"garden"`
	Garage               string `json:"garage"`
	Ensuite              string `json:"ensuite"`
	CharacterFeature     string `json:"character_feature"`
	EpcRatings           string `json:"epc_ratings"`
	PetsAllowed          string `json:"pets_allowed"`
	SmokingAllowed       string `json:"smoking_allowed"`
	DssAllowed           string `json:"dss_allowed"`
	SharersAllowed       string `json:"sharers_allowed"`
	Location             string `json:"location"`
	State                string `json:"state"`
	Country              string `json:"country"`
	PostalCode           string `json:"postal_code"`
	PropertyType         string `json:"property_type"`
	PropertyDetachedType string `json:"property_detached_type"`
}

type PropertyInput struct {
	PropertyTitle string `json:"property_title" binding:"required"`
	Description   string `json:"description"`
	Location      string `json:"location" binding:"required"`
	State         string `json:"state" binding:"required"`
	Country       string `json:"country" binding:"required"`
	PostalCode    string `json:"postal_code" binding:"required"`
}

type FavoritePropertyInput struct {
	PropertyID uint `json:"property_id" binding:"required"`
}
