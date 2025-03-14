package mail

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	UseTLS   bool
}

type emailSender struct {
	config SMTPConfig
}

// NewSMTPSender creates a new email sender with the given SMTP configuration.
func NewSMTPSender(config SMTPConfig) MailService {
	return &emailSender{config: config}
}

// Send sends an email to the specified recipients.
func (e *emailSender) Send(to []string, subject, body string) error {
	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", e.config.From, to, subject, body)
	auth := smtp.PlainAuth("", e.config.Username, e.config.Password, e.config.Host)

	address := fmt.Sprintf("%s:%d", e.config.Host, e.config.Port)

	if e.config.UseTLS {
		return e.sendWithTLS(address, auth, e.config.From, to, msg)
	}
	return smtp.SendMail(address, auth, e.config.From, to, []byte(msg))
}

// sendWithTLS sends an email using TLS.
func (e *emailSender) sendWithTLS(addr string, auth smtp.Auth, from string, to []string, msg string) error {
	conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: e.config.Host})
	if err != nil {
		return err
	}
	client, err := smtp.NewClient(conn, e.config.Host)
	if err != nil {
		return err
	}
	defer client.Quit()

	if err = client.Auth(auth); err != nil {
		return err
	}
	if err = client.Mail(from); err != nil {
		return err
	}
	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			return err
		}
	}

	writer, err := client.Data()
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(msg))
	if err != nil {
		return err
	}
	return writer.Close()
}
