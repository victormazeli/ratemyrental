package request

type UserInput struct {
	FullName    string  `json:"full_name"`
	Avatar      string  `json:"avatar"`
	Address     string  `json:"address"`
	PostalCode  string  `json:"postal_code"`
	PhoneNumber string  `json:"phone_number"`
	UserType    uint8   `json:"user_type"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	City        string  `json:"city"`
	Country     string  `json:"country"`
}

type UserSwitchType struct {
	UserType uint `json:"user_type"`
}
