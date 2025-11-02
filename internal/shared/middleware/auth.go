package middleware

import (
	"net/http"
	"service/internal/shared/model"
	"service/internal/shared/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token and sets user context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, model.CreateErrorResponse(
				"unauthorized",
				"Authorization header is required",
				nil,
			))
			c.Abort()
			return
		}

		// Extract token from header
		token, err := utils.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, model.CreateErrorResponse(
				"unauthorized",
				err.Error(),
				nil,
			))
			c.Abort()
			return
		}

		// Validate token
		claims, err := utils.ValidateAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, model.CreateErrorResponse(
				"unauthorized",
				"Invalid or expired token",
				nil,
			))
			c.Abort()
			return
		}

		// Set user context
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}

// RoleMiddleware checks if user has required role
func RoleMiddleware(allowedRoles ...model.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.CreateErrorResponse(
				"unauthorized",
				"User role not found in context",
				nil,
			))
			c.Abort()
			return
		}

		role, ok := userRole.(model.UserRole)
		if !ok {
			c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
				"internal_error",
				"Invalid user role type",
				nil,
			))
			c.Abort()
			return
		}

		// Check if user role is allowed
		allowed := false
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, model.CreateErrorResponse(
				"forbidden",
				"Insufficient permissions",
				nil,
			))
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminMiddleware checks if user is admin (admin_pusat or admin_cabang)
func AdminMiddleware() gin.HandlerFunc {
	return RoleMiddleware(model.RoleAdminPusat, model.RoleAdminCabang)
}

// StaffMiddleware checks if user is staff (admin, kasir, teknisi, kurir)
func StaffMiddleware() gin.HandlerFunc {
	return RoleMiddleware(
		model.RoleAdminPusat,
		model.RoleAdminCabang,
		model.RoleKasir,
		model.RoleTeknisi,
		model.RoleKurir,
	)
}

// CustomerMiddleware checks if user is customer
func CustomerMiddleware() gin.HandlerFunc {
	return RoleMiddleware(model.RolePelanggan)
}

// OptionalAuthMiddleware validates JWT token if present but doesn't require it
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Extract token from header
		token, err := utils.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.Next()
			return
		}

		// Validate token
		claims, err := utils.ValidateAccessToken(token)
		if err != nil {
			c.Next()
			return
		}

		// Set user context if token is valid
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}
