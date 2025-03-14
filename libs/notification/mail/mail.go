package mail

// MailService is an interface to send email
type MailService interface {
	Send(to []string, subject, body string) error
}
