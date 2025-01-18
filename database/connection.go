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
	// Drop the existing unique constraint if it exists
	var constraintExists bool
	err := db.Raw(`
        SELECT EXISTS (
            SELECT 1
            FROM information_schema.table_constraints
            WHERE constraint_name = 'users_email_key'
            AND table_name = 'users'
        )`).Scan(&constraintExists).Error
	if err != nil {
		return err
	}

	if constraintExists {
		if err := db.Exec(`ALTER TABLE "users" DROP CONSTRAINT "users_email_key"`).Error; err != nil {
			return err
		}
	}

	// Continue with your other migrations...
	return db.AutoMigrate(&models.User{})
}
