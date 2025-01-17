package services

import (
	"auth-server/models"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// HashPassword hashes the given password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a hashed password with its possible plaintext equivalent
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken generates a JWT token valid for 7 days
func GenerateToken(userID uint, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	})

	return token.SignedString([]byte(secret))
}

// RegisterUser creates a new user in the database
func RegisterUser(db *gorm.DB, email, password string) (*models.User, string, error) {
	fmt.Printf("Attempting to register user with email: %s\n", email)

	// Basic validation
	if email == "" || password == "" {
		return nil, "", errors.New("email and password are required")
	}

	// Check if user already exists
	var existing models.User
	result := db.Where("email = ?", email).First(&existing)

	// If it's "record not found", proceed with creation
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println("User not found, proceeding with creation")

		// Hash password
		hashed, err := HashPassword(password)
		if err != nil {
			return nil, "", fmt.Errorf("failed to hash password: %w", err)
		}

		user := &models.User{
			Email:    email,
			Password: hashed,
		}

		// Create new user
		if err := db.Create(user).Error; err != nil {
			return nil, "", fmt.Errorf("failed to create user: %w", err)
		}

		fmt.Printf("Successfully created user with ID: %d, email: %s\n", user.ID, email)
		return user, user.Password, nil
	}

	// If user exists, return error
	if result.Error == nil {
		return nil, "", errors.New("user already exists with that email")
	}

	// If it's any other error, return it
	return nil, "", fmt.Errorf("database error: %w", result.Error)
}

// LoginUser checks user credentials and returns a JWT if successful
func LoginUser(db *gorm.DB, email, password, secret string) (string, error) {
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return "", errors.New("invalid email or password")
	}

	if !CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid email or password")
	}

	token, err := GenerateToken(user.ID, secret)
	if err != nil {
		return "", err
	}

	return token, nil
}
