package utils

import "github.com/matcornic/hermes/v2"

func GenerateOTPEmailTemplate(otp string) string {
	// Configure hermes by setting a theme and your product info
	h := hermes.Hermes{
		// Optional Theme
		// Theme: new(Default)
		Product: hermes.Product{
			// Appears in header & footer of e-mails
			Name: "AuthCentral",
			Link: "https://example-hermes.com/",
			// Optional product logo
			Logo: "http://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png",
		},
	}

	email := hermes.Email{
		Body: hermes.Body{
			Name: "",
			Intros: []string{
				"Welcome to Auth Central! We're very excited to have you on board.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "To get started with Auth Central, please use the OTP code to verify your email:",
					InviteCode:   otp,
				},
			},
			Outros: []string{
				"Need help, or have questions? Just reply to this email, we'd love to help.",
			},
		},
	}

	// Generate an HTML email with the provided contents (for modern clients)
	emailBody, err := h.GenerateHTML(email)
	if err != nil {
		panic(err) // Tip: Handle error with something else than a panic ;)
	}
	return emailBody
}

func GenerateForgetPasswordEmailTemplate(otp string) string {
	h := hermes.Hermes{
		// Optional Theme
		// Theme: new(Default)
		Product: hermes.Product{
			// Appears in header & footer of e-mails
			Name: "Shop Vending",
			Link: "https://example-hermes.com/",
			// Optional product logo
			Logo: "http://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png",
		},
	}

	email := hermes.Email{
		Body: hermes.Body{
			Name: "",
			Intros: []string{
				"You have received this email because a password reset request for Shop Vending account was received.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Please use the OTP code to reset your password, otp expires in 5 minutes:",
					InviteCode:   otp,
				},
			},
			Outros: []string{
				"Need help, or have questions? Just reply to this email, we'd love to help.",
			},
		},
	}

	// Generate an HTML email with the provided contents (for modern clients)
	emailBody, err := h.GenerateHTML(email)
	if err != nil {
		panic(err) // Tip: Handle error with something else than a panic ;)
	}
	return emailBody
}
