package config

import "os"

type Config struct {
	SMTPHost       string
	SMTPPort       string
	SMTPUser       string
	SMTPPassword   string
	RecipientEmail string
	FromEmail      string
	ServerPort     string
	CORSOrigin     string
}

func Load() Config {
	return Config{
		SMTPHost:       getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:       getEnv("SMTP_PORT", "587"),
		SMTPUser:       getEnv("SMTP_USER", ""),
		SMTPPassword:   getEnv("SMTP_PASSWORD", ""),
		RecipientEmail: getEnv("RECIPIENT_EMAIL", ""),
		FromEmail:      getEnv("FROM_EMAIL", ""),
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		CORSOrigin:     getEnv("CORS_ORIGIN", "*"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
