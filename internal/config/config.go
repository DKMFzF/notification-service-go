package config

import (
	"os"

	env "github.com/joho/godotenv"
)

type Config struct {
	Port        string
	SMTPUser    string
	SMTPPass    string
	SMTPHost    string
	SMTPPort    string
	KafkaBroker string
}

func Load() *Config {
	_ = env.Load()

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	return &Config{
		Port:        port,
		SMTPUser:    os.Getenv("SMTP_USER"),
		SMTPPass:    os.Getenv("SMTP_PASS"),
		SMTPHost:    os.Getenv("SMTP_HOST"),
		SMTPPort:    os.Getenv("SMTP_PORT"),
		KafkaBroker: os.Getenv("KAFKA_BROKER"),
	}
}
