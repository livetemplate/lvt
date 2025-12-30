// Package email provides an interface and implementations for sending emails.
package email

import (
	"fmt"
	"io"
	"log"
	"os"
)

// EmailSender is the interface for sending emails.
type EmailSender interface {
	Send(to, subject, body string) error
}

// Sender is an alias for EmailSender for convenience.
type Sender = EmailSender

// ConsoleEmailSender logs emails to a writer (default: stdout) for development.
type ConsoleEmailSender struct {
	Writer io.Writer
}

// NewConsoleEmailSender creates a new ConsoleEmailSender that writes to stdout.
// Useful for development to see emails in the console.
//
// Example:
//
//	sender := email.NewConsoleEmailSender()
//	sender.Send("user@example.com", "Welcome", "Hello!")
func NewConsoleEmailSender() *ConsoleEmailSender {
	return &ConsoleEmailSender{Writer: os.Stdout}
}

// Send logs the email to the configured writer.
func (s *ConsoleEmailSender) Send(to, subject, body string) error {
	writer := s.Writer
	if writer == nil {
		writer = os.Stdout
	}
	log.Printf("📧 EMAIL (Console Mode)\n")
	fmt.Fprintf(writer, "To: %s\n", to)
	fmt.Fprintf(writer, "Subject: %s\n", subject)
	fmt.Fprintf(writer, "Body:\n%s\n", body)
	fmt.Fprintf(writer, "---\n")
	return nil
}

// NoopEmailSender discards all emails silently.
// Useful for testing when you don't want email output.
type NoopEmailSender struct{}

// NewNoopEmailSender creates a new NoopEmailSender.
//
// Example:
//
//	sender := email.NewNoopEmailSender()
//	sender.Send("user@example.com", "Welcome", "Hello!") // Does nothing
func NewNoopEmailSender() *NoopEmailSender {
	return &NoopEmailSender{}
}

// Send does nothing and returns nil.
func (s *NoopEmailSender) Send(to, subject, body string) error {
	return nil
}
