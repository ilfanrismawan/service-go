package middleware

import (
	"service/internal/shared/monitoring"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsMiddleware creates a metrics middleware
func MetricsMiddleware(metrics *monitoring.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Record metrics
		metrics.RecordHTTPRequest(
			c.Request.Method,
			c.Request.URL.Path,
			string(rune(c.Writer.Status())),
			duration,
		)
	}
}

// PrometheusHandler returns a Prometheus metrics handler
func PrometheusHandler() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}

// MetricsEndpoint creates a metrics endpoint
func MetricsEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "text/plain")
		c.String(200, "# iPhone Service POS Metrics\n")
		c.String(200, "# This endpoint provides Prometheus metrics\n")
		c.String(200, "# Access: http://localhost:8080/metrics\n")
	}
}
