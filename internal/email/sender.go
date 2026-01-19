package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"form2mail/internal/config"
)

type Sender struct {
	config config.Config
}

func NewSender(cfg config.Config) *Sender {
	return &Sender{config: cfg}
}

func (s *Sender) Send(to, subject, body string) error {
	// Connect to the SMTP server
	addr := fmt.Sprintf("%s:%s", s.config.SMTPHost, s.config.SMTPPort)

	// Connect to server
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer client.Close()

	// Send EHLO/HELO
	if err = client.Hello(s.config.SMTPHost); err != nil {
		return fmt.Errorf("failed to send HELLO: %w", err)
	}

	// Check if STARTTLS is supported and use it
	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{
			ServerName: s.config.SMTPHost,
		}
		if err = client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("failed to start TLS: %w", err)
		}
	}

	// Authenticate
	auth := smtp.PlainAuth("", s.config.SMTPUser, s.config.SMTPPassword, s.config.SMTPHost)
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Set sender
	if err = client.Mail(s.config.FromEmail); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipient
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Send message body
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to open data writer: %w", err)
	}

	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", s.config.FromEmail, to, subject, body))

	if _, err = w.Write(msg); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	// Quit
	return client.Quit()
}

func (s *Sender) SendContactNotification(name, email, subject, message string) error {
	recipientSubject := fmt.Sprintf("New Contact Form Submission: %s", subject)
	recipientBody := fmt.Sprintf(`
		<html>
		<body>
			<h2>New Contact Form Submission</h2>
			<p><strong>Name:</strong> %s</p>
			<p><strong>Email:</strong> %s</p>
			<p><strong>Subject:</strong> %s</p>
			<p><strong>Message:</strong></p>
			<p>%s</p>
		</body>
		</html>
	`, name, email, subject, strings.ReplaceAll(message, "\n", "<br>"))

	return s.Send(s.config.RecipientEmail, recipientSubject, recipientBody)
}

func (s *Sender) SendConfirmation(name, email, message string) error {
	confirmationSubject := "Thank you for contacting us"
	confirmationBody := fmt.Sprintf(`
		<html>
		<body>
			<h2>Thank you for your message, %s!</h2>
			<p>We have received your contact form submission and will get back to you as soon as possible.</p>
			<hr>
			<p><strong>Your message:</strong></p>
			<p>%s</p>
			<hr>
			<p>Best regards</p>
		</body>
		</html>
	`, name, strings.ReplaceAll(message, "\n", "<br>"))

	return s.Send(email, confirmationSubject, confirmationBody)
}
