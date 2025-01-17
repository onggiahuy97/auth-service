package main

import (
	"auth-server/config"
	"auth-server/database"
	"auth-server/routes"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading: %v", err)
	}

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config to database: %v", err)
	}

	// Connect to the db
	db, err := database.ConnectDB(cfg)
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("error running migrations: %v", err)
	}

	// Set up Gin router
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r, db, cfg)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server is running at port: %s\n", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Could not run server: %v", err)
	}
}
