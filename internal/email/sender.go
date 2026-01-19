package email

import (
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
	auth := smtp.PlainAuth("", s.config.SMTPUser, s.config.SMTPPassword, s.config.SMTPHost)

	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", s.config.FromEmail, to, subject, body))

	addr := fmt.Sprintf("%s:%s", s.config.SMTPHost, s.config.SMTPPort)
	return smtp.SendMail(addr, auth, s.config.FromEmail, []string{to}, msg)
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
