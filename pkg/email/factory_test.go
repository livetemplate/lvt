package email

import (
	"strings"
	"testing"
)

func TestNewEmailSenderFromEnv_DefaultConsole(t *testing.T) {
	t.Setenv("EMAIL_PROVIDER", "")
	sender, err := NewEmailSenderFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := sender.(*ConsoleEmailSender); !ok {
		t.Errorf("expected *ConsoleEmailSender, got %T", sender)
	}
}

func TestNewEmailSenderFromEnv_Console(t *testing.T) {
	t.Setenv("EMAIL_PROVIDER", "console")
	sender, err := NewEmailSenderFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := sender.(*ConsoleEmailSender); !ok {
		t.Errorf("expected *ConsoleEmailSender, got %T", sender)
	}
}

func TestNewEmailSenderFromEnv_Noop(t *testing.T) {
	t.Setenv("EMAIL_PROVIDER", "noop")
	sender, err := NewEmailSenderFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := sender.(*NoopEmailSender); !ok {
		t.Errorf("expected *NoopEmailSender, got %T", sender)
	}
}

func TestNewEmailSenderFromEnv_SMTP(t *testing.T) {
	t.Setenv("EMAIL_PROVIDER", "smtp")
	t.Setenv("SMTP_HOST", "smtp.example.com")
	t.Setenv("SMTP_PORT", "587")
	t.Setenv("EMAIL_FROM", "test@example.com")
	sender, err := NewEmailSenderFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := sender.(*SMTPEmailSender); !ok {
		t.Errorf("expected *SMTPEmailSender, got %T", sender)
	}
}

func TestNewEmailSenderFromEnv_SMTP_DefaultPort(t *testing.T) {
	t.Setenv("EMAIL_PROVIDER", "smtp")
	t.Setenv("SMTP_HOST", "smtp.example.com")
	t.Setenv("SMTP_PORT", "")
	t.Setenv("EMAIL_FROM", "test@example.com")
	sender, err := NewEmailSenderFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	smtp, ok := sender.(*SMTPEmailSender)
	if !ok {
		t.Fatalf("expected *SMTPEmailSender, got %T", sender)
	}
	if smtp.config.Port != 587 {
		t.Errorf("expected default port 587, got %d", smtp.config.Port)
	}
}

func TestNewEmailSenderFromEnv_SMTP_MissingHost(t *testing.T) {
	t.Setenv("EMAIL_PROVIDER", "smtp")
	t.Setenv("SMTP_HOST", "")
	t.Setenv("EMAIL_FROM", "test@example.com")
	_, err := NewEmailSenderFromEnv()
	if err == nil {
		t.Fatal("expected error for missing SMTP_HOST")
	}
	if !strings.Contains(err.Error(), "SMTP_HOST") {
		t.Errorf("error should mention SMTP_HOST, got: %v", err)
	}
}

func TestNewEmailSenderFromEnv_SMTP_MissingFrom(t *testing.T) {
	t.Setenv("EMAIL_PROVIDER", "smtp")
	t.Setenv("SMTP_HOST", "smtp.example.com")
	t.Setenv("EMAIL_FROM", "")
	_, err := NewEmailSenderFromEnv()
	if err == nil {
		t.Fatal("expected error for missing EMAIL_FROM")
	}
	if !strings.Contains(err.Error(), "EMAIL_FROM") {
		t.Errorf("error should mention EMAIL_FROM, got: %v", err)
	}
}

func TestNewEmailSenderFromEnv_SMTP_InvalidPort(t *testing.T) {
	t.Setenv("EMAIL_PROVIDER", "smtp")
	t.Setenv("SMTP_HOST", "smtp.example.com")
	t.Setenv("SMTP_PORT", "abc")
	t.Setenv("EMAIL_FROM", "test@example.com")
	_, err := NewEmailSenderFromEnv()
	if err == nil {
		t.Fatal("expected error for invalid SMTP_PORT")
	}
	if !strings.Contains(err.Error(), "invalid SMTP_PORT") {
		t.Errorf("error should mention invalid SMTP_PORT, got: %v", err)
	}
}

func TestNewEmailSenderFromEnv_UnknownProvider(t *testing.T) {
	t.Setenv("EMAIL_PROVIDER", "sendgrid")
	_, err := NewEmailSenderFromEnv()
	if err == nil {
		t.Fatal("expected error for unknown provider")
	}
	if !strings.Contains(err.Error(), "unknown EMAIL_PROVIDER") {
		t.Errorf("error should mention unknown EMAIL_PROVIDER, got: %v", err)
	}
}
