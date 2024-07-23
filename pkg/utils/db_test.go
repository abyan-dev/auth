package utils

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Set up the environment variables for testing
	os.Setenv("POSTGRES_HOST", "localhost")
	os.Setenv("POSTGRES_PORT", "5432")
	os.Setenv("POSTGRES_USER", "testuser")
	os.Setenv("POSTGRES_PASSWORD", "testpass")
	os.Setenv("POSTGRES_DB", "testdb")
	defer func() {
		os.Unsetenv("POSTGRES_HOST")
		os.Unsetenv("POSTGRES_PORT")
		os.Unsetenv("POSTGRES_USER")
		os.Unsetenv("POSTGRES_PASSWORD")
		os.Unsetenv("POSTGRES_DB")
	}()

	expectedConfig := &Config{
		DBHost:     "localhost",
		DBPort:     5432,
		DBUser:     "testuser",
		DBPassword: "testpass",
		DBName:     "testdb",
	}

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if *config != *expectedConfig {
		t.Errorf("Config mismatch: got %v, want %v", config, expectedConfig)
	}
}

func TestInitDB(t *testing.T) {
	config := &Config{
		DBHost:     "localhost",
		DBPort:     5432,
		DBUser:     "postgres",
		DBPassword: "postgres",
		DBName:     "postgres",
	}

	db, err := InitDB(config)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get database instance: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
}
