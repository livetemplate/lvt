package email

import (
	"testing"
	"time"
)

func TestSMTPEmailSender_ImplementsInterface(t *testing.T) {
	var _ EmailSender = &SMTPEmailSender{}
}

func TestNewSMTPEmailSender_DefaultTimeout(t *testing.T) {
	sender := NewSMTPEmailSender(SMTPConfig{
		Host: "localhost",
		Port: 587,
		From: "test@example.com",
	})
	if sender.config.Timeout != 10*time.Second {
		t.Errorf("expected default timeout 10s, got %v", sender.config.Timeout)
	}
}

func TestNewSMTPEmailSender_CustomTimeout(t *testing.T) {
	sender := NewSMTPEmailSender(SMTPConfig{
		Host:    "localhost",
		Port:    587,
		From:    "test@example.com",
		Timeout: 30 * time.Second,
	})
	if sender.config.Timeout != 30*time.Second {
		t.Errorf("expected custom timeout 30s, got %v", sender.config.Timeout)
	}
}

func TestNewSMTPEmailSender_ConfigFields(t *testing.T) {
	config := SMTPConfig{
		Host:     "smtp.example.com",
		Port:     465,
		Username: "user",
		Password: "pass",
		From:     "noreply@example.com",
		FromName: "My App",
	}
	sender := NewSMTPEmailSender(config)

	if sender.config.Host != "smtp.example.com" {
		t.Errorf("expected host smtp.example.com, got %s", sender.config.Host)
	}
	if sender.config.Port != 465 {
		t.Errorf("expected port 465, got %d", sender.config.Port)
	}
	if sender.config.Username != "user" {
		t.Errorf("expected username user, got %s", sender.config.Username)
	}
	if sender.config.From != "noreply@example.com" {
		t.Errorf("expected from noreply@example.com, got %s", sender.config.From)
	}
	if sender.config.FromName != "My App" {
		t.Errorf("expected from name My App, got %s", sender.config.FromName)
	}
}
