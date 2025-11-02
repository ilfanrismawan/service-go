package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// AppConfig holds the application configuration
type AppConfig struct {
	// Server configuration
	Port        string
	Environment string
	BaseURL     string

	// Database configuration
	DatabaseURL string

	// Redis configuration
	RedisURL string

	// JWT configuration
	JWTSecret     string
	JWTExpiry     time.Duration
	RefreshExpiry time.Duration

	// Payment gateway configuration
	MidtransServerKey    string
	MidtransClientKey    string
	MidtransIsProduction bool

	// File storage configuration
	S3Endpoint   string
	S3AccessKey  string
	S3SecretKey  string
	S3BucketName string
	S3Region     string

	// Notification configuration
	TwilioAccountSID  string
	TwilioAuthToken   string
	TwilioPhoneNumber string
	FirebaseServerKey string
	WhatsAppAPIKey    string
	WhatsAppAPIURL    string

	// Email configuration
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string

	// Rate limiting
	RateLimitRequests int
	RateLimitWindow   time.Duration

	// Localization
	DefaultLanguage string
	Timezone       string
	Currency       string
	DateFormat     string

    // Reconciliation
    ReconcileInterval time.Duration

    // Observability
    SentryDSN string
}

var Config *AppConfig

// LoadConfig loads configuration from environment variables
func LoadConfig() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	Config = &AppConfig{
		// Server configuration
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
		BaseURL:     getEnv("BASE_URL", "http://localhost:8080"),

		// Database configuration
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/iphone_service?sslmode=disable"),

		// Redis configuration
		RedisURL: getEnv("REDIS_URL", "redis://localhost:6379/0"),

		// JWT configuration
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key-change-this-in-production"),
		JWTExpiry:     getDurationEnv("JWT_EXPIRY", 24*time.Hour),
		RefreshExpiry: getDurationEnv("REFRESH_EXPIRY", 7*24*time.Hour),

		// Payment gateway configuration
		MidtransServerKey:    getEnv("MIDTRANS_SERVER_KEY", ""),
		MidtransClientKey:    getEnv("MIDTRANS_CLIENT_KEY", ""),
		MidtransIsProduction: getBoolEnv("MIDTRANS_IS_PRODUCTION", false),

		// File storage configuration
		S3Endpoint:   getEnv("S3_ENDPOINT", "http://localhost:9000"),
		S3AccessKey:  getEnv("S3_ACCESS_KEY", "minioadmin"),
		S3SecretKey:  getEnv("S3_SECRET_KEY", "minioadmin"),
		S3BucketName: getEnv("S3_BUCKET_NAME", "iphone-service"),
		S3Region:     getEnv("S3_REGION", "us-east-1"),

		// Notification configuration
		TwilioAccountSID:  getEnv("TWILIO_ACCOUNT_SID", ""),
		TwilioAuthToken:   getEnv("TWILIO_AUTH_TOKEN", ""),
		TwilioPhoneNumber: getEnv("TWILIO_PHONE_NUMBER", ""),
		FirebaseServerKey: getEnv("FIREBASE_SERVER_KEY", ""),
        WhatsAppAPIKey:    getEnv("WHATSAPP_API_KEY", ""),
        WhatsAppAPIKey:    getEnv("WHATSAPP_API_KEY", ""),
        WhatsAppAPIURL:    getEnv("WHATSAPP_API_URL", "https://api.fonnte.com/send"),

		// Email configuration
		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getIntEnv("SMTP_PORT", 587),
		SMTPUsername: getEnv("SMTP_USERNAME", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:     getEnv("SMTP_FROM", ""),

		// Rate limiting
		RateLimitRequests: getIntEnv("RATE_LIMIT_REQUESTS", 100),
		RateLimitWindow:   getDurationEnv("RATE_LIMIT_WINDOW", time.Minute),

		// Localization
		DefaultLanguage: getEnv("DEFAULT_LANGUAGE", "id-ID"),
		Timezone:       getEnv("TIMEZONE", "Asia/Jakarta"),
		Currency:       getEnv("CURRENCY", "IDR"),
		DateFormat:     getEnv("DATE_FORMAT", "02/01/2006"),

        // Reconciliation
        ReconcileInterval: getDurationEnv("RECONCILE_INTERVAL", 5*time.Minute),

        // Observability
        SentryDSN: getEnv("SENTRY_DSN", ""),
	}

	// Fail fast if JWT secret is left as the insecure default
	if Config.JWTSecret == "your-secret-key-change-this-in-production" {
		log.Fatal("JWT_SECRET is using the insecure default value. Set a strong secret in environment variables.")
	}

	// Set global time location based on configured timezone
	if loc, err := time.LoadLocation(Config.Timezone); err == nil {
		time.Local = loc
		log.Printf("Timezone set to %s", Config.Timezone)
	} else {
		log.Printf("Failed to load timezone %s: %v", Config.Timezone, err)
	}

	log.Printf("Configuration loaded successfully")
	log.Printf("Environment: %s", Config.Environment)
	log.Printf("Port: %s", Config.Port)
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getIntEnv gets an integer environment variable with a default value
func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getBoolEnv gets a boolean environment variable with a default value
func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// getDurationEnv gets a duration environment variable with a default value
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
