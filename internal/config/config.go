package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	Email    EmailConfig
	Storage  StorageConfig
}

type ServerConfig struct {
	Port         string
	Environment  string
	ReadTimeout  int
	WriteTimeout int
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type AuthConfig struct {
	Secret     string
	ExpiryTime int // in hours
}

type EmailConfig struct {
	ResendAPIKey string
	FromEmail    string
	FromName     string
}

type StorageConfig struct {
	Type      string // "local", "s3", "gcs", etc.
	LocalPath string
	Minio     MinioConfig
}

type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
	Bucket    string
	Region    string
}

// isDevelopment checks if we're running in development mode
func isDevelopment() bool {
	env := strings.ToLower(os.Getenv("APP_ENV"))
	return env == "" || env == "development" || env == "dev"
}

// LoadConfig reads configuration from .env file (in development)
// and environment variables (prioritized in production)
func LoadConfig() (*Config, error) {
	// Only in development mode, try to load .env file
	if isDevelopment() {
		// First try config/.env, then fall back to .env in the root directory
		if err := godotenv.Load("config/.env"); err != nil {
			if err := godotenv.Load(); err != nil {
				log.Println("No .env file found. Using environment variables only.")
			} else {
				log.Println("Loaded .env file from root directory")
			}
		} else {
			log.Println("Loaded .env file from config directory")
		}
	}

	// Now that environment variables are loaded (either from OS or .env),
	// we can create our config struct
	config := &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			Environment:  getEnv("SERVER_ENVIRONMENT", "development"),
			ReadTimeout:  getEnvAsInt("SERVER_READ_TIMEOUT", 60),
			WriteTimeout: getEnvAsInt("SERVER_WRITE_TIMEOUT", 60),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "crime_management"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Auth: AuthConfig{
			Secret:     getEnv("AUTH_SECRET", "default-secret-key-change-in-production"),
			ExpiryTime: getEnvAsInt("AUTH_EXPIRY_TIME", 24),
		},
		Email: EmailConfig{
			ResendAPIKey: getEnv("EMAIL_RESEND_API_KEY", ""),
			FromEmail:    getEnv("EMAIL_FROM", "no-reply@districtcore.gov"),
			FromName:     getEnv("EMAIL_FROM_NAME", "District Core Crime Management"),
		},
		Storage: StorageConfig{
			Type:      getEnv("STORAGE_TYPE", "local"),
			LocalPath: getEnv("STORAGE_LOCAL_PATH", "./storage"),
			Minio: MinioConfig{
				Endpoint:  getEnv("MINIO_ENDPOINT", "minio"),
				AccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
				SecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
				UseSSL:    getEnvAsBool("MINIO_USE_SSL", false),
				Bucket:    getEnv("MINIO_BUCKET", "crime-management"),
				Region:    getEnv("MINIO_REGION", "us-east-1"),
			},
		},
	}

	return config, nil
}

// Helper functions to get environment variables with defaults
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// Helper function to get environment variables as boolean
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// GetDSN returns the database connection string
func GetDSN(config *DatabaseConfig) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode)
}
