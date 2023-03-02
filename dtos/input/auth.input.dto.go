package input

type LoginInput struct {
	Email    string `json:"email" mod:"trim" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type RegisterInput struct {
	Email      string `json:"email" mod:"trim" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=8"`
	FirstName  string `json:"first_name" mod:"trim" binding:"required"`
	Lastname   string `json:"lastname" mod:"trim" binding:"required"`
	Address    string `json:"address" mod:"trim" binding:"required"`
	PostalCode string `json:"postal_code" mod:"trim" binding:"required"`
	Location   string `json:"location" mod:"trim" binding:"required"`
}

type ForgotPasswordInput struct {
	Email string `json:"email" mod:"trim"  binding:"required,email"`
}

type ResetPassword struct {
	Otp         string `json:"otp" mod:"trim" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type ValidateOtp struct {
	Otp string `json:"otp" mod:"trim" binding:"required"`
}

type ResendOtpInput struct {
	Email           string `json:"email" mod:"trim" binding:"required,email"`
	IsPasswordReset *bool  `json:"is_password_reset" binding:"required"`
}
