package middleware

import (
	"context"
	"net/http"
	"service/internal/shared/database"
	"service/internal/shared/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// RateLimiter handles rate limiting
type RateLimiter struct {
	redis *redis.Client
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		redis: database.Redis,
	}
}

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware() gin.HandlerFunc {
	limiter := NewRateLimiter()
	return func(c *gin.Context) {
		// Get client IP
		clientIP := c.ClientIP()

		// Get rate limit configuration
		requests := 100       // Default requests per window
		window := time.Minute // Default window

		// Check rate limit
		allowed, err := limiter.IsAllowed(clientIP, requests, window)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
				"rate_limit_error",
				"Rate limit check failed",
				nil,
			))
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(http.StatusTooManyRequests, model.CreateErrorResponse(
				"rate_limit_exceeded",
				"Too many requests",
				nil,
			))
			c.Abort()
			return
		}

		c.Next()
	}
}

// IsAllowed checks if a request is allowed based on rate limit
func (r *RateLimiter) IsAllowed(key string, requests int, window time.Duration) (bool, error) {
	ctx := context.Background()

	// Create rate limit key
	rateLimitKey := "rate_limit:" + key

	// Get current count
	count, err := r.redis.Get(ctx, rateLimitKey).Int()
	if err != nil && err != redis.Nil {
		return false, err
	}

	// If count is 0, this is the first request in the window
	if count == 0 {
		// Set count to 1 and set expiration
		err = r.redis.Set(ctx, rateLimitKey, 1, window).Err()
		if err != nil {
			return false, err
		}
		return true, nil
	}

	// Check if count exceeds limit
	if count >= requests {
		return false, nil
	}

	// Increment count
	err = r.redis.Incr(ctx, rateLimitKey).Err()
	if err != nil {
		return false, err
	}

	return true, nil
}

// GetRateLimitInfo gets rate limit information for a client
func (r *RateLimiter) GetRateLimitInfo(key string, requests int, window time.Duration) (map[string]interface{}, error) {
	ctx := context.Background()
	rateLimitKey := "rate_limit:" + key

	count, err := r.redis.Get(ctx, rateLimitKey).Int()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	ttl, err := r.redis.TTL(ctx, rateLimitKey).Result()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"limit":     requests,
		"remaining": requests - count,
		"reset":     time.Now().Add(ttl).Unix(),
		"window":    window.Seconds(),
	}, nil
}

// UserRateLimitMiddleware creates a user-specific rate limiting middleware
func UserRateLimitMiddleware() gin.HandlerFunc {
	limiter := NewRateLimiter()
	return func(c *gin.Context) {
		// Get user ID from context
		userID, exists := c.Get("user_id")
		if !exists {
			// If no user ID, fall back to IP-based rate limiting
			clientIP := c.ClientIP()
			allowed, err := limiter.IsAllowed(clientIP, 100, time.Minute)
			if err != nil || !allowed {
				c.JSON(http.StatusTooManyRequests, model.CreateErrorResponse(
					"rate_limit_exceeded",
					"Too many requests",
					nil,
				))
				c.Abort()
				return
			}
			c.Next()
			return
		}

		// Use user ID for rate limiting
		userKey := "user_rate_limit:" + userID.(uuid.UUID).String()
		allowed, err := limiter.IsAllowed(userKey, 200, time.Minute)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
				"rate_limit_error",
				"Rate limit check failed",
				nil,
			))
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(http.StatusTooManyRequests, model.CreateErrorResponse(
				"rate_limit_exceeded",
				"Too many requests",
				nil,
			))
			c.Abort()
			return
		}

		c.Next()
	}
}

// APIKeyRateLimitMiddleware creates an API key-based rate limiting middleware
func APIKeyRateLimitMiddleware() gin.HandlerFunc {
	limiter := NewRateLimiter()
	return func(c *gin.Context) {
		// Get API key from header
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, model.CreateErrorResponse(
				"api_key_required",
				"API key is required",
				nil,
			))
			c.Abort()
			return
		}

		// Use API key for rate limiting
		key := "api_rate_limit:" + apiKey
		allowed, err := limiter.IsAllowed(key, 1000, time.Hour)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
				"rate_limit_error",
				"Rate limit check failed",
				nil,
			))
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(http.StatusTooManyRequests, model.CreateErrorResponse(
				"rate_limit_exceeded",
				"Too many requests",
				nil,
			))
			c.Abort()
			return
		}

		c.Next()
	}
}
