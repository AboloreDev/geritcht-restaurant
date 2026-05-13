package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server     ServerConfig
	Cloudinary CloudinaryConfig
	Resend     ResendConfig
	JWT        JWTConfig
	Uploads    UploadConfig
	Database   DatabaseConfig
	Redis      RedisConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	Mode     string
	URL      string
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type CloudinaryConfig struct {
	CloudinaryName   string
	CloudinaryAPIKey string
	CloudinarySecret string
	CloudinaryFolder string
}

type UploadConfig struct {
	UploadPath     string
	UploadSize     int64
	UploadProvider string
}
type ResendConfig struct {
	ResendAPIKey   string
	ResendFromMail string
}

type JWTConfig struct {
	JWTSecret                 string
	JWTTokenExpiration        time.Duration
	JWTRefreshTokenExpiration time.Duration
}

type RedisConfig struct {
	URL        string
	QUEUE_NAME string
}

func LoadEnv() (*Config, error) {
	_ = godotenv.Load()

	jwtExpiresAt, _ := time.ParseDuration(GetEnv("JWT_EXPIRES_IN", "30m"))
	refreshTokenExpiresAt, _ := time.ParseDuration(GetEnv("REFRESH_TOKEN_EXPIRES_IN", "72h"))
	maxUpload, _ := strconv.ParseInt(GetEnv("MAX_UPLOAD_SIZE", "10485760"), 10, 64)
	// smtpPort, _ := strconv.Atoi(GetEnv("SMTP_PORT", "1025"))

	return &Config{
		Server: ServerConfig{
			Port:    GetEnv("PORT", "6000"),
			GinMode: GetEnv("GIN_MODE", "debug"),
		},

		Cloudinary: CloudinaryConfig{
			CloudinaryName:   GetEnv("CLOUDINARY_NAME", "test"),
			CloudinaryAPIKey: GetEnv("CLOUDINARY_API_KEY", "123456789"),
			CloudinarySecret: GetEnv("CLOUDINARY_SECRET", "123456789"),
			CloudinaryFolder: GetEnv("CLOUDINARY_FOLDER", "geritcht"),
		},

		Resend: ResendConfig{
			ResendAPIKey:   GetEnv("RESEND_API_KEY", "your_resend_api_key"),
			ResendFromMail: GetEnv("FROM_MAIL", "onboarding@resend.dev"),
		},

		JWT: JWTConfig{
			JWTSecret:                 GetEnv("JWT_SECRET", "your_jwt_secret"),
			JWTTokenExpiration:        jwtExpiresAt,
			JWTRefreshTokenExpiration: refreshTokenExpiresAt,
		},

		Uploads: UploadConfig{
			UploadPath:     GetEnv("UPLOAD_DIR", "/tmp/uploads"),
			UploadSize:     maxUpload,
			UploadProvider: GetEnv("UPLOAD_PROVIDER", "local"),
		},
		Database: DatabaseConfig{
			Host:     GetEnv("DB_HOST", "localhost"),
			Name:     GetEnv("DB_NAME", "armory_db"),
			Port:     GetEnv("DB_PORT", "5432"),
			Password: GetEnv("DB_PASSWORD", "password"),
			User:     GetEnv("DB_USER", "postgres"),
			Mode:     GetEnv("DB_SSL_MODE", "disable"),
			URL:      GetEnv("DATABASE_URL", ""),
		},

		Redis: RedisConfig{
			URL:        GetEnv("REDIS_URL", "test"),
			QUEUE_NAME: GetEnv("QUEUE_NAME", "queue"),
		},
	}, nil
}

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}
