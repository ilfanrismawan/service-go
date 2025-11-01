package delivery

import (
	"net/http"
	"service/internal/users/auth"
	"service/internal/core"
	"service/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *auth.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: auth.NewAuthService(),
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body core.UserRequest true "User registration data"
// @Success 201 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 409 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req core.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"validation_error",
			"Invalid request data",
			err.Error(),
		))
		return
	}

    // Sanitize strings to prevent XSS
    utils.SanitizeStructStrings(&req)

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"validation_error",
			"Validation failed",
			err.Error(),
		))
		return
	}

	user, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == core.ErrEmailExists || err == core.ErrPhoneExists {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, core.CreateErrorResponse(
			"registration_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusCreated, core.SuccessResponse(user, "User registered successfully"))
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body core.LoginRequest true "Login credentials"
// @Success 200 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req core.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"validation_error",
			"Invalid request data",
			err.Error(),
		))
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"validation_error",
			"Validation failed",
			err.Error(),
		))
		return
	}

	response, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == core.ErrUserNotFound || err == core.ErrInvalidPassword {
			statusCode = http.StatusUnauthorized
		}
		c.JSON(statusCode, core.CreateErrorResponse(
			"login_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, core.SuccessResponse(response, "Login successful"))
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body core.RefreshTokenRequest true "Refresh token data"
// @Success 200 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req core.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"validation_error",
			"Invalid request data",
			err.Error(),
		))
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"validation_error",
			"Validation failed",
			err.Error(),
		))
		return
	}

	response, err := h.authService.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == core.ErrInvalidToken {
			statusCode = http.StatusUnauthorized
		}
		c.JSON(statusCode, core.CreateErrorResponse(
			"refresh_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, core.SuccessResponse(response, "Token refreshed successfully"))
}

// Logout godoc
// @Summary Logout
// @Description Revoke refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body core.RefreshTokenRequest true "Refresh token data"
// @Success 200 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
    var req core.RefreshTokenRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
            "validation_error",
            "Invalid request data",
            err.Error(),
        ))
        return
    }

    if err := h.authService.Logout(c.Request.Context(), req.RefreshToken); err != nil {
        status := http.StatusUnauthorized
        c.JSON(status, core.CreateErrorResponse(
            "logout_failed",
            err.Error(),
            nil,
        ))
        return
    }

    c.JSON(http.StatusOK, core.SuccessResponse(nil, "Logged out"))
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current user profile information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} core.APIResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 404 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, core.CreateErrorResponse(
			"unauthorized",
			"User ID not found in context",
			nil,
		))
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, core.CreateErrorResponse(
			"internal_error",
			"Invalid user ID type",
			nil,
		))
		return
	}

	profile, err := h.authService.GetProfile(c.Request.Context(), userUUID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == core.ErrUserNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, core.CreateErrorResponse(
			"profile_fetch_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, core.SuccessResponse(profile, "Profile retrieved successfully"))
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update current user profile information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body core.UserRequest true "Profile update data"
// @Success 200 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 404 {object} core.ErrorResponse
// @Failure 409 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, core.CreateErrorResponse(
			"unauthorized",
			"User ID not found in context",
			nil,
		))
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, core.CreateErrorResponse(
			"internal_error",
			"Invalid user ID type",
			nil,
		))
		return
	}

	var req core.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"validation_error",
			"Invalid request data",
			err.Error(),
		))
		return
	}

    // Sanitize strings to prevent XSS
    utils.SanitizeStructStrings(&req)

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"validation_error",
			"Validation failed",
			err.Error(),
		))
		return
	}

	profile, err := h.authService.UpdateProfile(c.Request.Context(), userUUID, &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == core.ErrUserNotFound {
			statusCode = http.StatusNotFound
		} else if err == core.ErrEmailExists || err == core.ErrPhoneExists {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, core.CreateErrorResponse(
			"profile_update_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, core.SuccessResponse(profile, "Profile updated successfully"))
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change current user password
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body core.ChangePasswordRequest true "Password change data"
// @Success 200 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 404 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, core.CreateErrorResponse(
			"unauthorized",
			"User ID not found in context",
			nil,
		))
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, core.CreateErrorResponse(
			"internal_error",
			"Invalid user ID type",
			nil,
		))
		return
	}

	var req core.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"validation_error",
			"Invalid request data",
			err.Error(),
		))
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"validation_error",
			"Validation failed",
			err.Error(),
		))
		return
	}

	err := h.authService.ChangePassword(c.Request.Context(), userUUID, &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == core.ErrUserNotFound {
			statusCode = http.StatusNotFound
		} else if err == core.ErrInvalidPassword {
			statusCode = http.StatusUnauthorized
		}
		c.JSON(statusCode, core.CreateErrorResponse(
			"password_change_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, core.SuccessResponse(nil, "Password changed successfully"))
}

// ForgotPassword godoc
// @Summary Forgot password
// @Description Initiate password reset process
// @Tags auth
// @Accept json
// @Produce json
// @Param request body core.ForgotPasswordRequest true "Forgot password data"
// @Success 200 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req core.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"validation_error",
			"Invalid request data",
			err.Error(),
		))
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"validation_error",
			"Validation failed",
			err.Error(),
		))
		return
	}

	err := h.authService.ForgotPassword(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, core.CreateErrorResponse(
			"forgot_password_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, core.SuccessResponse(nil, "Password reset instructions sent to your email"))
}

// ResetPassword godoc
// @Summary Reset password
// @Description Reset password using reset token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body core.ResetPasswordRequest true "Reset password data"
// @Success 200 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req core.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"validation_error",
			"Invalid request data",
			err.Error(),
		))
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"validation_error",
			"Validation failed",
			err.Error(),
		))
		return
	}

	err := h.authService.ResetPassword(c.Request.Context(), &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == core.ErrInvalidToken {
			statusCode = http.StatusUnauthorized
		}
		c.JSON(statusCode, core.CreateErrorResponse(
			"reset_password_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, core.SuccessResponse(nil, "Password reset successfully"))
}

// GetUsers godoc
// @Summary Get users (admin)
// @Description List users with pagination and filters (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param role query string false "Filter by role"
// @Param branch_id query string false "Filter by branch ID"
// @Success 200 {object} core.PaginatedResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /admin/users [get]
func (h *AuthHandler) GetUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	var role *core.UserRole
	if roleStr := c.Query("role"); roleStr != "" {
		r := core.UserRole(roleStr)
		role = &r
	}

	var branchID *uuid.UUID
	if branchStr := c.Query("branch_id"); branchStr != "" {
		if b, err := uuid.Parse(branchStr); err == nil {
			branchID = &b
		}
	}

	users, total, err := h.authService.ListUsers(c.Request.Context(), page, limit, role, branchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, core.CreateErrorResponse(
			"users_fetch_failed",
			err.Error(),
			nil,
		))
		return
	}

	var responses []core.UserResponse
	for _, u := range users {
		responses = append(responses, u.ToResponse())
	}

	pagination := core.PaginationResponse{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
	}

	c.JSON(http.StatusOK, core.PaginatedSuccessResponse(responses, pagination, "Users retrieved successfully"))
}

// GetUser godoc
// @Summary Get user (admin)
// @Description Get user by ID (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 404 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /admin/users/{id} [get]
func (h *AuthHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"invalid_id",
			"Invalid user ID format",
			nil,
		))
		return
	}

	user, err := h.authService.GetUser(c.Request.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		if err == core.ErrUserNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, core.CreateErrorResponse("user_fetch_failed", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, core.SuccessResponse(user, "User retrieved successfully"))
}

// UpdateUser godoc
// @Summary Update user (admin)
// @Description Update user by ID (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body core.UserRequest true "User data"
// @Success 200 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 404 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /admin/users/{id} [put]
func (h *AuthHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"invalid_id",
			"Invalid user ID format",
			nil,
		))
		return
	}

	var req core.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"validation_error",
			"Invalid request data",
			err.Error(),
		))
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"validation_error",
			"Validation failed",
			err.Error(),
		))
		return
	}

	user, err := h.authService.UpdateUser(c.Request.Context(), id, &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err == core.ErrUserNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, core.CreateErrorResponse("user_update_failed", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, core.SuccessResponse(user, "User updated successfully"))
}

// DeleteUser godoc
// @Summary Delete user (admin)
// @Description Delete a user by ID (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 404 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /admin/users/{id} [delete]
func (h *AuthHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"invalid_id",
			"Invalid user ID format",
			nil,
		))
		return
	}

	if err := h.authService.DeleteUser(c.Request.Context(), id); err != nil {
		status := http.StatusInternalServerError
		if err == core.ErrUserNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, core.CreateErrorResponse("user_delete_failed", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, core.SuccessResponse(nil, "User deleted successfully"))
}
