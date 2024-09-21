package main

import (
	"fmt"
	"log"
	"strconv"

	"gopkg.in/gomail.v2"
)

// https://pkg.go.dev/gopkg.in/gomail.v2#section-readme
// SSL/TLS Email
func SendEmail(recipients []string, subject, body string) error {
	// Convert the SMTP port from string to int
	smtpPort, err := strconv.Atoi(Config.SmtpPort)
	if err != nil {
		log.Println("Error converting SMTP port:", err)
		return err
	}

	// Create a new message
	msg := gomail.NewMessage()

	// Set the sender and recipients
	// Set the sender with label and email address (Format: "Label <email@example.com>")
	from := fmt.Sprintf("%s <%s>", Config.SmtpFromLabel, Config.SmtpFromEmail)
	msg.SetHeader("From", from)
	msg.SetHeader("From", Config.SmtpFromEmail)
	msg.SetHeader("To", recipients...)
	msg.SetHeader("Subject", subject)

	// Set the email body
	msg.SetBody("text/plain", body)

	// Create a new dialer
	d := gomail.NewDialer(Config.SmtpHost, smtpPort, Config.SmtpUser, Config.SmtpPass)

	// If you need to skip TLS verification for the SMTP server certificate, uncomment the following line:
	// d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send the email
	if err := d.DialAndSend(msg); err != nil {
		log.Println("Error sending email:", err)
		return err
	}

	if Config.AppEnv == "development" {
		log.Println("Email sent")
	}
	return nil
}
