package database

import (
	"auth-server/config"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	// Find the .env.test file by walking up directories if needed
	envPath := findEnvFile(".env.test")
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Warning: Error loading .env.test file: %v", err)
	}

	code := m.Run()
	os.Exit(code)
}

// findEnvFile walks up the directory tree looking for the env file
func findEnvFile(filename string) string {
	dir, err := os.Getwd()
	if err != nil {
		return filename
	}

	for {
		path := filepath.Join(dir, filename)
		if _, err := os.Stat(path); err == nil {
			return path
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return filename
}

func TestConnectDB(t *testing.T) {
	// Debug logging
	t.Log("Testing with environment:")
	t.Logf("DB_HOST: %s", os.Getenv("DB_HOST"))
	t.Logf("DB_PORT: %s", os.Getenv("DB_PORT"))
	t.Logf("DB_USER: %s", os.Getenv("DB_USER"))
	t.Logf("DB_NAME: %s", os.Getenv("DB_NAME"))

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	db, err := ConnectDB(cfg)
	if err != nil {
		t.Fatalf("failed to connect to the database: %v", err)
	}

	if db == nil {
		t.Fatal("db instance is nil, connection likely failed")
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get underlying sqlDB: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("failed to ping the database: %v", err)
	}

	defer func() {
		if sqlDB != nil {
			sqlDB.Close()
		}
	}()

	t.Log("Successfully connected to the database")
}

func TestRunMigrations(t *testing.T) {
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	db, err := ConnectDB(cfg)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	defer func() {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	if err := RunMigrations(db); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	t.Log("Migrations ran successfully")
}
