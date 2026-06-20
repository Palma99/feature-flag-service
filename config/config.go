package config

import (
	"os"

	"github.com/joho/godotenv"
)

type PostgresConfig struct {
	Host           string
	Port           string
	User           string
	Password       string
	Database       string
	SSLMode        string
	ChannelBinding string
}

type Config struct {
	DB                PostgresConfig
	Port              string
	JwtSecret         string
	CorsAllowedOrigin string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		DB: PostgresConfig{
			Host:           os.Getenv("DB_HOST"),
			Port:           os.Getenv("DB_PORT"),
			User:           os.Getenv("DB_USER"),
			Password:       os.Getenv("DB_PASSWORD"),
			Database:       os.Getenv("DB_NAME"),
			SSLMode:        os.Getenv("DB_SSLMODE"),
			ChannelBinding: os.Getenv("DB_CHANNEL_BINDING"),
		},
		Port:              os.Getenv("PORT"),
		JwtSecret:         os.Getenv("JWT_SECRET"),
		CorsAllowedOrigin: os.Getenv("CORS_ALLOWED_ORIGIN"),
	}, nil
}
