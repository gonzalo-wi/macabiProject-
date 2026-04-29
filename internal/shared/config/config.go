package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	Port              string
	JWTSecret         string
	JWTExpiration     time.Duration
	BrevoAPIKey       string
	BrevoEmailFrom    string
	FrontendPublicURL string
	InvitationTTL     time.Duration
	PasswordResetTTL  time.Duration
}

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}

func (c *Config) Validate() error {
	required := map[string]string{
		"DB_HOST":             c.DBHost,
		"DB_PORT":             c.DBPort,
		"DB_USER":             c.DBUser,
		"DB_PASSWORD":         c.DBPassword,
		"DB_NAME":             c.DBName,
		"JWT_SECRET":          c.JWTSecret,
		"BREVO_API_KEY":       c.BrevoAPIKey,
		"BREVO_EMAIL_FROM":    c.BrevoEmailFrom,
		"FRONTEND_PUBLIC_URL": c.FrontendPublicURL,
	}
	for key, val := range required {
		if val == "" {
			return fmt.Errorf("missing required environment variable: %s", key)
		}
	}
	return nil
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	jwtExpiration := 24 * time.Hour
	if v := os.Getenv("JWT_EXPIRATION_HOURS"); v != "" {
		if d, err := time.ParseDuration(v + "h"); err == nil {
			jwtExpiration = d
		}
	}

	invTTL := 7 * 24 * time.Hour
	if v := strings.TrimSpace(os.Getenv("INVITATION_TTL_HOURS")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			invTTL = time.Duration(n) * time.Hour
		}
	}

	resetTTL := time.Hour
	if v := strings.TrimSpace(os.Getenv("PASSWORD_RESET_TTL_MINUTES")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			resetTTL = time.Duration(n) * time.Minute
		}
	}

	brevoKey := firstNonEmpty(os.Getenv("BREVO_API_KEY"), os.Getenv("brevo_api_key"))
	brevoFrom := firstNonEmpty(os.Getenv("BREVO_EMAIL_FROM"), os.Getenv("brevo_email_from"))

	return &Config{
		DBHost:            os.Getenv("DB_HOST"),
		DBPort:            os.Getenv("DB_PORT"),
		DBUser:            os.Getenv("DB_USER"),
		DBPassword:        os.Getenv("DB_PASSWORD"),
		DBName:            os.Getenv("DB_NAME"),
		Port:              port,
		JWTSecret:         os.Getenv("JWT_SECRET"),
		JWTExpiration:     jwtExpiration,
		BrevoAPIKey:       brevoKey,
		BrevoEmailFrom:    brevoFrom,
		FrontendPublicURL: strings.TrimSpace(os.Getenv("FRONTEND_PUBLIC_URL")),
		InvitationTTL:     invTTL,
		PasswordResetTTL:  resetTTL,
	}
}

func firstNonEmpty(a, b string) string {
	if strings.TrimSpace(a) != "" {
		return strings.TrimSpace(a)
	}
	return strings.TrimSpace(b)
}
