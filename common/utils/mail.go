package utils

import (
	"context"
	"crypto/tls"
	"github.com/mailgun/mailgun-go/v3"
	"gopkg.in/gomail.v2"
	"log"
	"rateMyRentalBackend/config"
	"time"
)

func SendTestEmail(receiver, body, subject string) error {
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

func SendEmail(env *config.Env, receiver, body, subject string) (string, error) {
	mg := mailgun.NewMailgun(env.MailGunDomain, env.MailGunApiKey)
	mg.SetAPIBase(mailgun.APIBaseEU)
	m := mg.NewMessage(
		"no-reply"+"@"+env.MailGunDomain,
		subject,
		body,
		receiver,
	)
	m.SetHtml(body)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, id, err := mg.Send(ctx, m)
	return id, err
}
