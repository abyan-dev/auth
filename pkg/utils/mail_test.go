package utils

import (
	"os"
	"testing"
)

func TestLoadMailEnv(t *testing.T) {
	os.Setenv("SENDER_EMAIL", "test@example.com")
	os.Setenv("SMTP_HOST", "smtp.example.com")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_USER", "testuser")
	os.Setenv("SMTP_PASS", "testpass")
	defer func() {
		os.Unsetenv("SENDER_EMAIL")
		os.Unsetenv("SMTP_HOST")
		os.Unsetenv("SMTP_PORT")
		os.Unsetenv("SMTP_USER")
		os.Unsetenv("SMTP_PASS")
	}()

	expectedConfig := &MailEnvConfig{
		SenderEmail: "test@example.com",
		SmtpHost:    "smtp.example.com",
		SmtpPort:    "587",
		SmtpUser:    "testuser",
		SmtpPass:    "testpass",
	}

	config, err := LoadMailEnv()
	if err != nil {
		t.Fatalf("Failed to load mail environment config: %v", err)
	}

	if *config != *expectedConfig {
		t.Errorf("Config mismatch: got %v, want %v", config, expectedConfig)
	}
}

func TestLoadMailEnvMissingSenderEmail(t *testing.T) {
	os.Unsetenv("SENDER_EMAIL")
	os.Setenv("SMTP_HOST", "smtp.example.com")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_USER", "testuser")
	os.Setenv("SMTP_PASS", "testpass")
	defer func() {
		os.Unsetenv("SMTP_HOST")
		os.Unsetenv("SMTP_PORT")
		os.Unsetenv("SMTP_USER")
		os.Unsetenv("SMTP_PASS")
	}()

	_, err := LoadMailEnv()
	if err == nil {
		t.Fatal("Expected error due to missing SENDER_EMAIL, but got none")
	}

	expectedErr := "SENDER_EMAIL environment variable is not set"
	if err.Error() != expectedErr {
		t.Errorf("Error message mismatch: got %v, want %v", err.Error(), expectedErr)
	}
}
