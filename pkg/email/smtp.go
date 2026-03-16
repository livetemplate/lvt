package email

import (
	"fmt"
	"time"

	mail "github.com/wneessen/go-mail"
)

// SMTPConfig holds SMTP connection configuration.
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	FromName string
	Timeout  time.Duration
}

// SMTPEmailSender sends emails via SMTP using the go-mail library.
type SMTPEmailSender struct {
	config SMTPConfig
}

// NewSMTPEmailSender creates a new SMTP email sender with the given configuration.
func NewSMTPEmailSender(config SMTPConfig) *SMTPEmailSender {
	if config.Timeout == 0 {
		config.Timeout = 10 * time.Second
	}
	return &SMTPEmailSender{config: config}
}

// Send sends a plain-text email via SMTP.
func (s *SMTPEmailSender) Send(to, subject, body string) error {
	msg := mail.NewMsg()

	if s.config.FromName != "" {
		if err := msg.FromFormat(s.config.FromName, s.config.From); err != nil {
			return fmt.Errorf("email: invalid from address: %w", err)
		}
	} else {
		if err := msg.From(s.config.From); err != nil {
			return fmt.Errorf("email: invalid from address: %w", err)
		}
	}

	if err := msg.To(to); err != nil {
		return fmt.Errorf("email: invalid to address: %w", err)
	}

	msg.Subject(subject)
	msg.SetBodyString(mail.TypeTextPlain, body)

	client, err := mail.NewClient(s.config.Host,
		mail.WithPort(s.config.Port),
		mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
		mail.WithTLSPortPolicy(mail.TLSOpportunistic),
		mail.WithUsername(s.config.Username),
		mail.WithPassword(s.config.Password),
		mail.WithTimeout(s.config.Timeout),
	)
	if err != nil {
		return fmt.Errorf("email: failed to create SMTP client: %w", err)
	}

	if err := client.DialAndSend(msg); err != nil {
		return fmt.Errorf("email: failed to send: %w", err)
	}

	return nil
}
