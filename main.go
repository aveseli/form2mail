package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"
)

type Config struct {
	SMTPHost       string
	SMTPPort       string
	SMTPUser       string
	SMTPPassword   string
	RecipientEmail string
	FromEmail      string
	ServerPort     string
}

type ContactForm struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func loadConfig() Config {
	return Config{
		SMTPHost:       getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:       getEnv("SMTP_PORT", "587"),
		SMTPUser:       getEnv("SMTP_USER", ""),
		SMTPPassword:   getEnv("SMTP_PASSWORD", ""),
		RecipientEmail: getEnv("RECIPIENT_EMAIL", ""),
		FromEmail:      getEnv("FROM_EMAIL", ""),
		ServerPort:     getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func sendEmail(config Config, to, subject, body string) error {
	auth := smtp.PlainAuth("", config.SMTPUser, config.SMTPPassword, config.SMTPHost)

	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", config.FromEmail, to, subject, body))

	addr := fmt.Sprintf("%s:%s", config.SMTPHost, config.SMTPPort)
	return smtp.SendMail(addr, auth, config.FromEmail, []string{to}, msg)
}

func handleContact(config Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow POST requests
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Enable CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight OPTIONS request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Parse form data
		var form ContactForm
		contentType := r.Header.Get("Content-Type")

		if strings.Contains(contentType, "application/json") {
			// Parse JSON
			if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
				http.Error(w, "Invalid JSON format", http.StatusBadRequest)
				return
			}
		} else {
			// Parse form data
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Failed to parse form", http.StatusBadRequest)
				return
			}
			form.Name = r.FormValue("name")
			form.Email = r.FormValue("email")
			form.Subject = r.FormValue("subject")
			form.Message = r.FormValue("message")
		}

		// Validate required fields
		if form.Name == "" || form.Email == "" || form.Message == "" {
			http.Error(w, "Name, email, and message are required", http.StatusBadRequest)
			return
		}

		// Send email to recipient (site owner)
		recipientSubject := fmt.Sprintf("New Contact Form Submission: %s", form.Subject)
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
		`, form.Name, form.Email, form.Subject, strings.ReplaceAll(form.Message, "\n", "<br>"))

		if err := sendEmail(config, config.RecipientEmail, recipientSubject, recipientBody); err != nil {
			log.Printf("Failed to send email to recipient: %v", err)
			http.Error(w, "Failed to send email", http.StatusInternalServerError)
			return
		}

		// Send confirmation email to customer
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
		`, form.Name, strings.ReplaceAll(form.Message, "\n", "<br>"))

		if err := sendEmail(config, form.Email, confirmationSubject, confirmationBody); err != nil {
			log.Printf("Failed to send confirmation email to customer: %v", err)
			// Don't fail the request if confirmation email fails
		}

		// Send success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "success",
			"message": "Your message has been sent successfully",
		})
	}
}

func main() {
	config := loadConfig()

	// Validate required config
	if config.SMTPUser == "" || config.SMTPPassword == "" || config.RecipientEmail == "" {
		log.Fatal("SMTP_USER, SMTP_PASSWORD, and RECIPIENT_EMAIL must be set")
	}

	http.HandleFunc("/contact", handleContact(config))

	log.Printf("Server starting on port %s...", config.ServerPort)
	if err := http.ListenAndServe(":"+config.ServerPort, nil); err != nil {
		log.Fatal(err)
	}
}
