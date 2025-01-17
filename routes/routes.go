package routes

import (
	"auth-server/config"
	"auth-server/controllers"
	"auth-server/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", controllers.Register(db, cfg))
		authGroup.POST("/login", controllers.Login(db, cfg))
	}

	protected := r.Group("/api")
	protected.Use(middleware.AuthRequired(cfg.JWTSecret))
	{
		protected.GET("/profile", func(c *gin.Context) {
			userID := c.GetUint("user_id")
			c.JSON(200, gin.H{"message": "Welcome!", "user_id": userID})
		})
	}
}
