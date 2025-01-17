package database

import (
	"auth-server/config"
	"auth-server/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return db, err
}

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{})
}
