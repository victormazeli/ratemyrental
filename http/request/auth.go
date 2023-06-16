package request

type LoginInput struct {
	Email    string `json:"email" mod:"trim" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type RegisterInput struct {
	Email    string `json:"email" mod:"trim" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"full_name" mod:"trim" binding:"required"`
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

type GoogleAuth struct {
	AccessToken string `json:"access_token" binding:"required"`
}

type GoogleAuthModel struct {
	Sub           *string `json:"sub"`
	Name          *string `json:"name"`
	GivenName     *string `json:"given_name"`
	FamilyName    *string `json:"family_name"`
	Email         *string `json:"email"`
	Picture       *string `json:"picture"`
	EmailVerified *bool   `json:"email_verified"`
	Locale        *string `json:"locale"`
}
