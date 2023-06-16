package response

type AuthResponse struct {
	//User  UserDTO
	Token string `json:"token"`
}
