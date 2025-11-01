package delivery

import (
	"net/http"
	"service/internal/core"
	"service/internal/repository"
	"service/internal/service"
	"service/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RatingHandler handles rating endpoints
type RatingHandler struct {
	ratingService *service.RatingService
}

// NewRatingHandler creates a new rating handler
func NewRatingHandler() *RatingHandler {
	return &RatingHandler{
		ratingService: service.NewRatingService(),
	}
}

// CreateRating godoc
// @Summary Create a rating
// @Description Create a rating for a completed order
// @Tags ratings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body core.RatingRequest true "Rating request"
// @Success 201 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /ratings [post]
func (h *RatingHandler) CreateRating(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, core.CreateErrorResponse(
			"unauthorized",
			"User ID not found in context",
			nil,
		))
		return
	}

	customerID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, core.CreateErrorResponse(
			"internal_error",
			"Invalid user ID type",
			nil,
		))
		return
	}

	var req core.RatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"invalid_request",
			err.Error(),
			nil,
		))
		return
	}

	// Sanitize strings to prevent XSS
	utils.SanitizeStructStrings(&req)

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"validation_error",
			err.Error(),
			nil,
		))
		return
	}

	rating, err := h.ratingService.CreateRating(c.Request.Context(), customerID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"rating_creation_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusCreated, core.SuccessResponse(rating, "Rating created successfully"))
}

// GetRating godoc
// @Summary Get a rating
// @Description Get a rating by ID
// @Tags ratings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Rating ID"
// @Success 200 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 404 {object} core.ErrorResponse
// @Router /ratings/{id} [get]
func (h *RatingHandler) GetRating(c *gin.Context) {
	ratingIDStr := c.Param("id")
	ratingID, err := uuid.Parse(ratingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"invalid_rating_id",
			"Invalid rating ID format",
			nil,
		))
		return
	}

	rating, err := h.ratingService.GetRating(c.Request.Context(), ratingID)
	if err != nil {
		c.JSON(http.StatusNotFound, core.CreateErrorResponse(
			"rating_not_found",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, core.SuccessResponse(rating, "Rating retrieved successfully"))
}

// ListRatings godoc
// @Summary List ratings
// @Description List ratings with pagination and filters
// @Tags ratings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param branch_id query string false "Branch ID"
// @Param technician_id query string false "Technician ID"
// @Param min_rating query int false "Minimum rating (1-5)"
// @Success 200 {object} core.APIResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /ratings [get]
func (h *RatingHandler) ListRatings(c *gin.Context) {
	// Parse pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Parse filters
	filters := &repository.RatingFilters{}

	if branchIDStr := c.Query("branch_id"); branchIDStr != "" {
		if branchID, err := uuid.Parse(branchIDStr); err == nil {
			filters.BranchID = &branchID
		}
	}

	if technicianIDStr := c.Query("technician_id"); technicianIDStr != "" {
		if technicianID, err := uuid.Parse(technicianIDStr); err == nil {
			filters.TechnicianID = &technicianID
		}
	}

	if minRatingStr := c.Query("min_rating"); minRatingStr != "" {
		if minRating, err := strconv.Atoi(minRatingStr); err == nil && minRating >= 1 && minRating <= 5 {
			filters.MinRating = &minRating
		}
	}

	// Only show public ratings
	isPublic := true
	filters.IsPublic = &isPublic

	result, err := h.ratingService.ListRatings(c.Request.Context(), page, limit, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, core.CreateErrorResponse(
			"list_ratings_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateRating godoc
// @Summary Update a rating
// @Description Update a rating (only by owner)
// @Tags ratings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Rating ID"
// @Param request body core.RatingRequest true "Rating request"
// @Success 200 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 404 {object} core.ErrorResponse
// @Router /ratings/{id} [put]
func (h *RatingHandler) UpdateRating(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, core.CreateErrorResponse(
			"unauthorized",
			"User ID not found in context",
			nil,
		))
		return
	}

	customerID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, core.CreateErrorResponse(
			"internal_error",
			"Invalid user ID type",
			nil,
		))
		return
	}

	ratingIDStr := c.Param("id")
	ratingID, err := uuid.Parse(ratingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"invalid_rating_id",
			"Invalid rating ID format",
			nil,
		))
		return
	}

	var req core.RatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"invalid_request",
			err.Error(),
			nil,
		))
		return
	}

	// Sanitize strings to prevent XSS
	utils.SanitizeStructStrings(&req)

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"validation_error",
			err.Error(),
			nil,
		))
		return
	}

	rating, err := h.ratingService.UpdateRating(c.Request.Context(), ratingID, customerID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"rating_update_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, core.SuccessResponse(rating, "Rating updated successfully"))
}

// DeleteRating godoc
// @Summary Delete a rating
// @Description Delete a rating (only by owner)
// @Tags ratings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Rating ID"
// @Success 200 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 404 {object} core.ErrorResponse
// @Router /ratings/{id} [delete]
func (h *RatingHandler) DeleteRating(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, core.CreateErrorResponse(
			"unauthorized",
			"User ID not found in context",
			nil,
		))
		return
	}

	customerID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, core.CreateErrorResponse(
			"internal_error",
			"Invalid user ID type",
			nil,
		))
		return
	}

	ratingIDStr := c.Param("id")
	ratingID, err := uuid.Parse(ratingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"invalid_rating_id",
			"Invalid rating ID format",
			nil,
		))
		return
	}

	if err := h.ratingService.DeleteRating(c.Request.Context(), ratingID, customerID); err != nil {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"rating_deletion_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, core.SuccessResponse(nil, "Rating deleted successfully"))
}

// GetAverageRating godoc
// @Summary Get average rating
// @Description Get average rating statistics
// @Tags ratings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param branch_id query string false "Branch ID"
// @Param technician_id query string false "Technician ID"
// @Success 200 {object} core.APIResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /ratings/average [get]
func (h *RatingHandler) GetAverageRating(c *gin.Context) {
	var branchID, technicianID *uuid.UUID

	if branchIDStr := c.Query("branch_id"); branchIDStr != "" {
		if parsed, err := uuid.Parse(branchIDStr); err == nil {
			branchID = &parsed
		}
	}

	if technicianIDStr := c.Query("technician_id"); technicianIDStr != "" {
		if parsed, err := uuid.Parse(technicianIDStr); err == nil {
			technicianID = &parsed
		}
	}

	average, err := h.ratingService.GetAverageRating(c.Request.Context(), branchID, technicianID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, core.CreateErrorResponse(
			"average_rating_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, core.SuccessResponse(average, "Average rating retrieved successfully"))
}

