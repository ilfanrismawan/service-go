package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LoggingMiddleware creates a logging middleware
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Get request body
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get response content type
		contentType := ""
		if c.Writer != nil {
			contentType = c.Writer.Header().Get("Content-Type")
		}

		// Create log entry
		logEntry := logrus.WithFields(logrus.Fields{
			"timestamp":    start.Format(time.RFC3339),
			"method":       c.Request.Method,
			"path":         c.Request.URL.Path,
			"query":        c.Request.URL.RawQuery,
			"status":       c.Writer.Status(),
			"latency":      latency.String(),
			"ip":           c.ClientIP(),
			"user_agent":   c.Request.UserAgent(),
			"content_type": contentType,
			"request_id":   c.GetString("request_id"),
			"user_id":      c.GetString("user_id"),
			"user_role":    c.GetString("user_role"),
		})

		// Add request body for non-GET requests
		if c.Request.Method != "GET" && len(requestBody) > 0 {
			var body interface{}
			if err := json.Unmarshal(requestBody, &body); err == nil {
				logEntry = logEntry.WithField("request_body", body)
			}
		}

		// Log based on status code
		switch {
		case c.Writer.Status() >= 500:
			logEntry.Error("Server error")
		case c.Writer.Status() >= 400:
			logEntry.Warn("Client error")
		default:
			logEntry.Info("Request completed")
		}
	}
}

// AuditLogMiddleware creates an audit logging middleware
func AuditLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only log for authenticated users
		userID, exists := c.Get("user_id")
		if !exists {
			c.Next()
			return
		}

		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Create audit log entry
		auditEntry := logrus.WithFields(logrus.Fields{
			"timestamp":  start.Format(time.RFC3339),
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"query":      c.Request.URL.RawQuery,
			"status":     c.Writer.Status(),
			"latency":    latency.String(),
			"ip":         c.ClientIP(),
			"user_id":    userID,
			"user_role":  c.GetString("user_role"),
			"request_id": c.GetString("request_id"),
		})

		// Log audit entry
		auditEntry.Info("Audit log")
	}
}

// ErrorLoggingMiddleware creates an error logging middleware
func ErrorLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Check for errors
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logrus.WithFields(logrus.Fields{
					"timestamp":  time.Now().Format(time.RFC3339),
					"method":     c.Request.Method,
					"path":       c.Request.URL.Path,
					"status":     c.Writer.Status(),
					"ip":         c.ClientIP(),
					"user_id":    c.GetString("user_id"),
					"request_id": c.GetString("request_id"),
					"error":      err.Error(),
				}).Error("Request error")
			}
		}
	}
}

// SecurityLoggingMiddleware creates a security logging middleware
func SecurityLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for suspicious activity
		userAgent := c.Request.UserAgent()
		ip := c.ClientIP()

		// Log suspicious user agents
		if isSuspiciousUserAgent(userAgent) {
			logrus.WithFields(logrus.Fields{
				"timestamp":  time.Now().Format(time.RFC3339),
				"ip":         ip,
				"user_agent": userAgent,
				"path":       c.Request.URL.Path,
				"method":     c.Request.Method,
				"request_id": c.GetString("request_id"),
			}).Warn("Suspicious user agent detected")
		}

		// Log failed authentication attempts
		if c.Writer.Status() == 401 {
			logrus.WithFields(logrus.Fields{
				"timestamp":  time.Now().Format(time.RFC3339),
				"ip":         ip,
				"user_agent": userAgent,
				"path":       c.Request.URL.Path,
				"method":     c.Request.Method,
				"request_id": c.GetString("request_id"),
			}).Warn("Failed authentication attempt")
		}

		c.Next()
	}
}

// isSuspiciousUserAgent checks if user agent is suspicious
func isSuspiciousUserAgent(userAgent string) bool {
	suspiciousPatterns := []string{
		"bot", "crawler", "spider", "scraper", "scanner",
		"sqlmap", "nikto", "nmap", "masscan",
		"curl", "wget", "python-requests",
	}

	for _, pattern := range suspiciousPatterns {
		if strings.Contains(strings.ToLower(userAgent), pattern) {
			return true
		}
	}

	return false
}

// PerformanceLoggingMiddleware creates a performance logging middleware
func PerformanceLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Log slow requests
		if latency > 5*time.Second {
			logrus.WithFields(logrus.Fields{
				"timestamp":  time.Now().Format(time.RFC3339),
				"method":     c.Request.Method,
				"path":       c.Request.URL.Path,
				"latency":    latency.String(),
				"ip":         c.ClientIP(),
				"user_id":    c.GetString("user_id"),
				"request_id": c.GetString("request_id"),
			}).Warn("Slow request detected")
		}
	}
}
