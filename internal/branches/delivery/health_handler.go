package delivery

import (
	"context"
	"net/http"

	"service/internal/shared/database"
	"service/internal/shared/model"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check endpoints
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck godoc
// @Summary Health check
// @Description Check the health of the application and its dependencies
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} model.APIResponse
// @Failure 503 {object} model.ErrorResponse
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	// Check database connection
	dbStatus := h.checkDatabase()

	// Check Redis connection
	redisStatus := h.checkRedis()

	// Determine overall health
	overallHealth := "healthy"
	if !dbStatus["healthy"].(bool) || !redisStatus["healthy"].(bool) {
		overallHealth = "unhealthy"
	}

	// Create health response
	healthData := gin.H{
		"status":    overallHealth,
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "1.0.0",
		"services": gin.H{
			"database": dbStatus,
			"redis":    redisStatus,
		},
	}

	// Set appropriate status code
	statusCode := http.StatusOK
	if overallHealth == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, model.SuccessResponse(healthData, "Health check completed"))
}

// LivenessCheck godoc
// @Summary Liveness check
// @Description Check if the application is alive
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} model.APIResponse
// @Router /health/live [get]
func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, model.SuccessResponse(gin.H{
		"status":    "alive",
		"timestamp": time.Now().Format(time.RFC3339),
	}, "Application is alive"))
}

// ReadinessCheck godoc
// @Summary Readiness check
// @Description Check if the application is ready to serve requests
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} model.APIResponse
// @Failure 503 {object} model.ErrorResponse
// @Router /health/ready [get]
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	// Check database connection
	dbStatus := h.checkDatabase()

	// Check Redis connection
	redisStatus := h.checkRedis()

	// Determine readiness
	ready := dbStatus["healthy"].(bool) && redisStatus["healthy"].(bool)

	// Create readiness response
	readinessData := gin.H{
		"ready":     ready,
		"timestamp": time.Now().Format(time.RFC3339),
		"services": gin.H{
			"database": dbStatus,
			"redis":    redisStatus,
		},
	}

	// Set appropriate status code
	statusCode := http.StatusOK
	if !ready {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, model.SuccessResponse(readinessData, "Readiness check completed"))
}

// checkDatabase checks database connection health
func (h *HealthHandler) checkDatabase() gin.H {
	status := gin.H{
		"service": "postgresql",
		"healthy": false,
		"error":   nil,
	}

	// Get database connection
	db := database.DB
	if db == nil {
		status["error"] = "Database connection not initialized"
		return status
	}

	// Test database connection
	sqlDB, err := db.DB()
	if err != nil {
		status["error"] = err.Error()
		return status
	}

	// Ping database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		status["error"] = err.Error()
		return status
	}

	// Get database stats
	stats := sqlDB.Stats()
	status["healthy"] = true
	status["stats"] = gin.H{
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration.String(),
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_idle_time_closed": stats.MaxIdleTimeClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}

	return status
}

// checkRedis checks Redis connection health
func (h *HealthHandler) checkRedis() gin.H {
	status := gin.H{
		"service": "redis",
		"healthy": false,
		"error":   nil,
	}

	// Get Redis connection
	redis := database.Redis
	if redis == nil {
		status["error"] = "Redis connection not initialized"
		return status
	}

	// Test Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redis.Ping(ctx).Err(); err != nil {
		status["error"] = err.Error()
		return status
	}

	// Get Redis info
	info, err := redis.Info(ctx).Result()
	if err != nil {
		status["error"] = err.Error()
		return status
	}

	status["healthy"] = true
	status["info"] = info

	return status
}

// Metrics godoc
// @Summary Application metrics
// @Description Get application metrics and statistics
// @Tags health
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.APIResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /health/metrics [get]
func (h *HealthHandler) Metrics(c *gin.Context) {
	// Get database stats
	db := database.DB
	var dbStats gin.H
	if db != nil {
		if sqlDB, err := db.DB(); err == nil {
			stats := sqlDB.Stats()
			dbStats = gin.H{
				"open_connections":     stats.OpenConnections,
				"in_use":               stats.InUse,
				"idle":                 stats.Idle,
				"wait_count":           stats.WaitCount,
				"wait_duration":        stats.WaitDuration.String(),
				"max_idle_closed":      stats.MaxIdleClosed,
				"max_idle_time_closed": stats.MaxIdleTimeClosed,
				"max_lifetime_closed":  stats.MaxLifetimeClosed,
			}
		}
	}

	// Get Redis stats
	redis := database.Redis
	var redisStats gin.H
	if redis != nil {
		ctx := context.Background()
		if info, err := redis.Info(ctx).Result(); err == nil {
			redisStats = gin.H{
				"info": info,
			}
		}
	}

	// Create metrics response
	metricsData := gin.H{
		"timestamp": time.Now().Format(time.RFC3339),
		"database":  dbStats,
		"redis":     redisStats,
		"uptime":    time.Since(startTime).String(),
	}

	c.JSON(http.StatusOK, model.SuccessResponse(metricsData, "Metrics retrieved successfully"))
}

// startTime tracks application start time
var startTime = time.Now()
