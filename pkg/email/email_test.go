package email

import (
	"bytes"
	"strings"
	"testing"
)

func TestConsoleEmailSender(t *testing.T) {
	var buf bytes.Buffer
	sender := &ConsoleEmailSender{Writer: &buf}

	err := sender.Send("test@example.com", "Test Subject", "Test Body")
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "To: test@example.com") {
		t.Error("Send() output missing recipient")
	}

	if !strings.Contains(output, "Subject: Test Subject") {
		t.Error("Send() output missing subject")
	}

	if !strings.Contains(output, "Test Body") {
		t.Error("Send() output missing body")
	}
}

func TestConsoleEmailSenderDefaultWriter(t *testing.T) {
	sender := NewConsoleEmailSender()

	// Should not panic with default writer
	err := sender.Send("test@example.com", "Test", "Body")
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
}

func TestNoopEmailSender(t *testing.T) {
	sender := NewNoopEmailSender()

	err := sender.Send("test@example.com", "Test Subject", "Test Body")
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
}

func TestEmailSenderInterface(t *testing.T) {
	// Verify both types implement the EmailSender interface
	var _ EmailSender = &ConsoleEmailSender{}
	var _ EmailSender = &NoopEmailSender{}
}
