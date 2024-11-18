package services

import (
	"fmt"
	"os"
	"strconv"

	"github.com/wneessen/go-mail"
)

var client *mail.Client

func InitEmailClient() {

	server := os.Getenv("SMTP_SERVER")
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		panic("Could not convert SMTP_PORT to int")
	}
	userName := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	mc, err := mail.NewClient(server, mail.WithPort(port), mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithUsername(userName), mail.WithPassword(password))
	if err != nil {
		panic(err)
	}
	client = mc
}

func SendVerficationEmail(email string, token string) error {
	m := mail.NewMsg()
	if err := m.From(os.Getenv("SMTP_USERNAME")); err != nil {
		return err
	}
	if err := m.To(email); err != nil {
		return err
	}

	api_port := os.Getenv("API_DOMAIN")

	b := fmt.Sprintf(`
	<h1>Verify your email address</h1>
	<p>Click the link below to verify your email address</p>
	<a href='http://%s/verify/%s/%s'>Verify</a>
	`, api_port, token, email)

	m.Subject("Verify your email address")
	m.SetBodyString(mail.TypeTextHTML, b)
	if err := client.DialAndSend(m); err != nil {
		return err
	}
	return nil

}
