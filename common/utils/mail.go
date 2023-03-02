package utils

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
	"log"
)

func SendEmail(receiver, body, subject string) error {
	var e error
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "test@email.com")

	// Set E-Mail receivers
	m.SetHeader("To", receiver)

	// Set E-Mail subject
	m.SetHeader("Subject", subject)

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/html", body)

	d := gomail.NewDialer("", 1025, "", "")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		log.Println(err)
		e = err
	}
	return e
}
