package main

import (
	"log"
	"net/http"

	"form2mail/internal/config"
	"form2mail/internal/email"
	"form2mail/internal/handler"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Validate required config
	if cfg.SMTPUser == "" || cfg.SMTPPassword == "" || cfg.RecipientEmail == "" {
		log.Fatal("SMTP_USER, SMTP_PASSWORD, and RECIPIENT_EMAIL must be set")
	}

	// Initialize email sender
	emailSender := email.NewSender(cfg)

	// Initialize handler
	contactHandler := handler.NewContactHandler(emailSender, cfg.CORSOrigin)

	// Register routes
	http.Handle("/contact", contactHandler)

	// Start server
	log.Printf("Server starting on port %s...", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, nil); err != nil {
		log.Fatal(err)
	}
}
