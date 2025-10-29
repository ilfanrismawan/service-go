// @title iPhone Service POS API
// @version 1.0
// @description Backend API for iPhone Service Point of Sales system supporting 50 branches across Indonesia
// @termsOfService http://swagger.io/terms/

// @contact.name iPhone Service API Support
// @contact.email support@iphoneservice.com
// @contact.url https://iphoneservice.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"log"
	_ "service/docs" // Import docs for Swagger
	"service/internal/config"
	"service/internal/database"
	"service/internal/delivery"
	"service/internal/middleware"
	"service/internal/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize database
	database.InitPostgres()

	// Initialize Redis
	database.InitRedis()

	// Initialize validator
	utils.InitValidator()

	// Setup Gin router
	r := setupRouter()

	// Start server
	log.Printf("🚀 iPhone Service API starting on port %s\n", config.Config.Port)
	log.Printf("📊 Environment: %s\n", config.Config.Environment)
	log.Printf("🔗 Health check: http://localhost:%s/health\n", config.Config.Port)
	log.Printf("📚 API Documentation: http://localhost:%s/swagger/index.html\n", config.Config.Port)
	log.Printf("📖 API Docs: http://localhost:%s/docs\n", config.Config.Port)

	if err := r.Run(":" + config.Config.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// setupRouter configures and returns the Gin router
func setupRouter() *gin.Engine {
	// Set Gin mode based on environment
	if config.Config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Enable CORS for Swagger UI
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Add middleware
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.SecurityHeadersMiddleware())
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.ErrorLoggingMiddleware())
	r.Use(middleware.SecurityLoggingMiddleware())
	r.Use(middleware.PerformanceLoggingMiddleware())
	r.Use(gin.Recovery())

	// Setup API routes
	delivery.SetupRoutes(r)

	return r
}
