package mail

import "testing"

func TestNewEmailSender(t *testing.T) {
	config := SMTPConfig{
		Host:     "smtp.gmail.com",
		Port:     465,
		Username: "example@gmail.com", // Replace with actual email
		Password: "password",          // Replace with actual password
		From:     "example@gmail.com", // Replace with actual email
		UseTLS:   true,
	}
	sender := NewSMTPSender(config)
	if sender == nil {
		t.Fatal("Expected non-nil EmailSender instance")
	}
}

func TestSendEmail(t *testing.T) {
	config := SMTPConfig{
		Host:     "smtp.gmail.com",
		Port:     465,
		Username: "example@gmail.com", // Replace with actual email
		Password: "password",          // Replace with actual password
		From:     "example@gmail.com", // Replace with actual email
		UseTLS:   true,
	}
	sender := NewSMTPSender(config)
	err := sender.Send([]string{"pehkulon@example.com"}, "Test Subject", "Test Body")
	if err != nil {
		t.Errorf("Failed to send email: %v", err)
	}
}
