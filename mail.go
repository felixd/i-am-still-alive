package main

import (
	"net/smtp"
	"strings"
)

// SSL/TLS Email

func SendEmail(recipients []string, subject, body string) error {
	auth := smtp.PlainAuth("", Config.SmtpUser, Config.SmtpPass, Config.SmtpHost)

	to := strings.Join(recipients, ", ")
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"From: " + Config.SmtpFromEmail + "\r\n" +
		"\r\n" + body + "\r\n")
	return smtp.SendMail(Config.SmtpHost+":"+Config.SmtpPort, auth, Config.SmtpFromEmail, recipients, msg)

	/*

		from := mail.Address{os.Getenv("SMTP_FROM"), os.Getenv("SMTP_SENDER")}
		to := mail.Address{"", "felixd@wp.pl"}
		subj := "FlameIT - Dead Man Switch"
		body := "This is an example body.\n With two lines."

		// Setup headers
		headers := make(map[string]string)
		headers["From"] = from.String()
		headers["To"] = to.String()
		headers["Subject"] = subj

		// Setup message
		message := ""
		for k, v := range headers {
			message += fmt.Sprintf("%s: %s\r\n", k, v)
		}
		message += "\r\n" + body

		// Connect to the SMTP Server
		servername := os.Getenv("SMTP_HOST") + ":" + os.Getenv("SMTP_PORT")

		host, _, _ := net.SplitHostPort(servername)

		auth := smtp.PlainAuth("", os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASSWORD"), host)

		// TLS config
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         host,
		}

		// Here is the key, you need to call tls.Dial instead of smtp.Dial
		// for smtp servers running on 465 that require an ssl connection
		// from the very beginning (no starttls)
		conn, err := tls.Dial("tcp", servername, tlsconfig)
		if err != nil {
			log.Panic(err)
		}

		c, err := smtp.NewClient(conn, host)
		if err != nil {
			log.Panic(err)
		}

		// Auth
		if err = c.Auth(auth); err != nil {
			log.Panic(err)
		}

		// To && From
		if err = c.Mail(from.Address); err != nil {
			log.Panic(err)
		}

		if err = c.Rcpt(to.Address); err != nil {
			log.Panic(err)
		}

		// Data
		w, err := c.Data()
		if err != nil {
			log.Panic(err)
		}

		_, err = w.Write([]byte(message))
		if err != nil {
			log.Panic(err)
		}

		err = w.Close()
		if err != nil {
			log.Panic(err)
		}

		c.Quit()
	*/
}
