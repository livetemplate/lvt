package email

import (
	"fmt"
	"os"
	"strconv"
)

// NewSMTPEmailSenderFromEnv creates an SMTPEmailSender from environment variables.
//
// Required env vars:
//
//	SMTP_HOST  - SMTP server hostname
//	EMAIL_FROM - Sender email address
//
// Optional env vars:
//
//	SMTP_PORT      - SMTP server port (default: 587)
//	SMTP_USER      - SMTP username
//	SMTP_PASS      - SMTP password
//	EMAIL_FROM_NAME - Sender display name
func NewSMTPEmailSenderFromEnv() (*SMTPEmailSender, error) {
	host := os.Getenv("SMTP_HOST")
	if host == "" {
		return nil, fmt.Errorf("email: SMTP_HOST environment variable is required")
	}

	portStr := os.Getenv("SMTP_PORT")
	if portStr == "" {
		portStr = "587"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("email: invalid SMTP_PORT %q: %w", portStr, err)
	}

	from := os.Getenv("EMAIL_FROM")
	if from == "" {
		return nil, fmt.Errorf("email: EMAIL_FROM environment variable is required")
	}

	return NewSMTPEmailSender(SMTPConfig{
		Host:     host,
		Port:     port,
		Username: os.Getenv("SMTP_USER"),
		Password: os.Getenv("SMTP_PASS"),
		From:     from,
		FromName: os.Getenv("EMAIL_FROM_NAME"),
	}), nil
}

// NewEmailSenderFromEnv creates an EmailSender based on the EMAIL_PROVIDER
// environment variable.
//
// Supported providers:
//
//	"console" (default) - logs emails to stdout (development)
//	"smtp"              - sends via SMTP server (production)
//	"noop"              - discards silently (testing)
func NewEmailSenderFromEnv() (EmailSender, error) {
	provider := os.Getenv("EMAIL_PROVIDER")
	if provider == "" {
		provider = "console"
	}

	switch provider {
	case "console":
		return NewConsoleEmailSender(), nil
	case "noop":
		return NewNoopEmailSender(), nil
	case "smtp":
		return NewSMTPEmailSenderFromEnv()
	default:
		return nil, fmt.Errorf("email: unknown EMAIL_PROVIDER %q (supported: console, smtp, noop)", provider)
	}
}
