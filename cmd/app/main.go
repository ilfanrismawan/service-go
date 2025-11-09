// @title iPhone Service POS API
// @version 1.0.0
// @description Backend API untuk sistem Point of Sales jasa service iPhone dengan dukungan 50 cabang di seluruh Indonesia. Sistem ini mendukung multi-service, real-time tracking, pembayaran online/tunai, membership system, dan berbagai fitur lainnya.
//
// ## Fitur Utama
// - Authentication & Authorization dengan JWT
// - Multi-service support (Service Catalog)
// - Order Management dengan real-time tracking
// - Payment Processing (Midtrans, Cash, Bank Transfer, GoPay, QRIS)
// - Membership System (4-tier: Bronze, Silver, Gold, Platinum)
// - Location Tracking dengan WebSocket
// - Real-time Chat dengan WebSocket
// - Notification System (Email, WhatsApp, FCM)
// - File Management (S3-compatible storage)
// - Reporting & Analytics
// - Rating System
// - Dashboard untuk Admin
//
// ## Role-based Access Control
// - **admin_pusat**: Admin pusat dengan akses penuh
// - **admin_cabang**: Admin cabang dengan akses terbatas pada cabangnya
// - **kasir**: Kasir untuk proses pembayaran dan order
// - **teknisi**: Teknisi untuk service order management
// - **kurir**: Kurir untuk pickup dan delivery
// - **pelanggan**: Pelanggan untuk order dan tracking
//
// ## Authentication
// API menggunakan JWT (JSON Web Token) untuk authentication. Setiap request ke protected endpoint harus menyertakan header:
// ```
// Authorization: Bearer <access_token>
// ```
//
// Access token memiliki expiry 24 jam, dan refresh token memiliki expiry 7 hari.
//
// ## Base URL
// - Development: http://localhost:8080
// - Production: https://api.iphoneservice.com
//
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
// @description Type "Bearer" followed by a space and JWT token. Example: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

package main

import (
	"context"
	"log"
	"time"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "service/docs" // Swagger docs
	svc "service/internal/domain/payments/service"
	"service/internal/router"
	"service/internal/shared/config"
	"service/internal/shared/database"
	"service/internal/shared/middleware"
	"service/internal/shared/monitoring"
	"service/internal/shared/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize Sentry (if DSN is provided)
	if err := middleware.InitSentry(config.Config.SentryDSN, config.Config.Environment); err != nil {
		log.Printf("Failed to initialize Sentry: %v", err)
	}

	// Initialize database
	database.InitPostgres()

	// Initialize Redis
	database.InitRedis()

	// Initialize validator
	utils.InitValidator()

	// Initialize Swagger
	docs.SwaggerInfo.BasePath = "/"

	// Setup Gin router
	r := setupRouter()

	// Register Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start background reconciliation job
	go func() {
		ticker := time.NewTicker(config.Config.ReconcileInterval)
		defer ticker.Stop()
		ps := svc.NewPaymentService()
		for {
			<-ticker.C
			_ = ps.ReconcilePendingPayments(context.Background())
		}
	}()

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

	// Initialize metrics and expose Prometheus endpoint
	metrics := monitoring.NewMetrics()

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
	// Sentry capture middleware (only if DSN is set)
	if config.Config.SentryDSN != "" {
		r.Use(middleware.SentryMiddleware())
	}
	// Enforce HTTPS only in production
	if config.Config.Environment == "production" {
		r.Use(middleware.HTTPSRedirectMiddleware())
	}
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.SecurityHeadersMiddleware())
	r.Use(middleware.MetricsMiddleware(metrics))
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.ErrorLoggingMiddleware())
	r.Use(middleware.SecurityLoggingMiddleware())
	r.Use(middleware.PerformanceLoggingMiddleware())
	r.Use(gin.Recovery())

	// Expose Prometheus metrics at /metrics (no auth)
	r.GET("/metrics", middleware.PrometheusHandler())

	// Setup API routes
	router.SetupRoutes(r)

	return r
}
