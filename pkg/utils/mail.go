package utils

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type MailEnvConfig struct {
	SenderEmail string
	SmtpHost    string
	SmtpPort    string
	SmtpUser    string
	SmtpPass    string
}

func LoadMailEnv() (*MailEnvConfig, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	config := &MailEnvConfig{
		SenderEmail: os.Getenv("SENDER_EMAIL"),
		SmtpHost:    os.Getenv("SMTP_HOST"),
		SmtpPort:    os.Getenv("SMTP_PORT"),
		SmtpUser:    os.Getenv("SMTP_USER"),
		SmtpPass:    os.Getenv("SMTP_PASS"),
	}

	if config.SenderEmail == "" {
		return nil, errors.New("SENDER_EMAIL environment variable is not set")
	}

	return config, nil
}
