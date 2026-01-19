package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"form2mail/internal/email"
)

type ContactForm struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type ContactHandler struct {
	emailSender *email.Sender
	corsOrigin  string
}

func NewContactHandler(emailSender *email.Sender, corsOrigin string) *ContactHandler {
	return &ContactHandler{
		emailSender: emailSender,
		corsOrigin:  corsOrigin,
	}
}

func (h *ContactHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers first (before any method checks)
	w.Header().Set("Access-Control-Allow-Origin", h.corsOrigin)
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight OPTIONS request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow POST requests for actual form submission
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
	if err := h.emailSender.SendContactNotification(form.Name, form.Email, form.Subject, form.Message); err != nil {
		log.Printf("Failed to send email to recipient: %v", err)
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	// Send confirmation email to customer
	if err := h.emailSender.SendConfirmation(form.Name, form.Email, form.Message); err != nil {
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
