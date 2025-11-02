package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SwaggerHandler handles Swagger documentation endpoints
type SwaggerHandler struct{}

// NewSwaggerHandler creates a new Swagger handler
func NewSwaggerHandler() *SwaggerHandler {
	return &SwaggerHandler{}
}

// SetupSwaggerRoutes sets up Swagger documentation routes
func (h *SwaggerHandler) SetupSwaggerRoutes(r *gin.Engine) {
	// Swagger documentation endpoint (handles all /swagger/* paths including index.html)
	// Create a route group so we can set a post-handler middleware to relax CSP for Swagger UI only
	swagGroup := r.Group("/swagger")
	swagGroup.Use(func(c *gin.Context) {
		// Execute handler first
		c.Next()
		// Then override CSP for swagger responses (allows inline script in index.html)
		c.Header("Content-Security-Policy", "default-src 'self' 'unsafe-inline' 'unsafe-eval' data: blob:")
	})
	swagGroup.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API documentation redirect
	r.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	// API documentation info endpoint
	r.GET("/api-docs", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"title":       "iPhone Service POS API",
			"version":     "1.0",
			"description": "Backend API for iPhone Service Point of Sales system supporting 50 branches across Indonesia",
			"swagger_url": "/swagger/index.html",
			"docs_url":    "/docs",
		})
	})
}
