package input

type UserInput struct {
	FullName   string `json:"full_name"`
	Avatar     string `json:"avatar"`
	Address    string `json:"address"`
	PostalCode string `json:"postal_code"`
	Location   string `json:"location"`
}
