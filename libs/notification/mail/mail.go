package mail

import (
	"fmt"
	"log"
	"sync"

	"github.com/zainulbr/simple-loan-engine/settings"
)

// MailService is an interface to send email
type MailService interface {
	Send(to []string, subject, body string) error
}

const (
	defaultKey = "default"
)

var store sync.Map

// Create instance
func Create(option *settings.SMTPOption, key string) error {
	if _, ok := store.Load(key); ok {
		return fmt.Errorf("Mail instance '%s' already exists", key)
	}

	log.Printf("Initializing smtps: %s", key)

	svc := NewSMTPSender(SMTPConfig{
		Host:     option.Host,
		Port:     option.Port,
		Username: option.Username,
		Password: option.Password,
		From:     option.From,
		UseTLS:   true,
	})

	store.Store(key, svc)

	return nil
}

// Close closes instances
func Close() error {
	store.Range(func(key, value interface{}) bool {
		log.Printf("Closing SQL DB: %s", key)

		store.Delete(key)

		return true
	})

	return nil
}

// Open connection
func Open(settings *settings.Settings) error {
	return Create(&settings.Conn.SMTP, defaultKey)
}

// Mail returns the mail connection of the specified key,
// if none specified, return default connection
func Mail(mailKey ...string) MailService {
	key := defaultKey

	if len(mailKey) > 0 && len(mailKey[0]) > 0 {
		key = mailKey[0]
	}

	instance, ok := store.Load(key)
	if !ok {
		log.Fatalf("Mail Instance '%s' not found, please call Create() or Open() first.", key)
	}

	return instance.(MailService)
}
