package dtos

type AuthResponse struct {
	User  UserDTO
	Token string `json:"token"`
}
