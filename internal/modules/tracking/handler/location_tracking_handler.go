package handler

import (
	"net/http"
	"service/internal/modules/tracking/service"
	"service/internal/shared/middleware"
	"service/internal/shared/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LocationTrackingHandler handles location tracking HTTP requests
type LocationTrackingHandler struct {
	trackingService *service.LocationTrackingService
}

// NewLocationTrackingHandler creates a new location tracking handler
func NewLocationTrackingHandler() *LocationTrackingHandler {
	return &LocationTrackingHandler{
		trackingService: service.NewLocationTrackingService(),
	}
}

// RegisterRoutes registers location tracking routes
func (h *LocationTrackingHandler) RegisterRoutes(r *gin.RouterGroup) {
	tracking := r.Group("/orders/:id")
	tracking.Use(middleware.AuthMiddleware())
	{
		// Update location (courier/provider)
		tracking.POST("/location", h.UpdateLocation)
		
		// Get current location (customer)
		tracking.GET("/location", h.GetCurrentLocation)
		
		// Get location history
		tracking.GET("/location/history", h.GetLocationHistory)
		
		// Calculate ETA
		tracking.POST("/eta", h.CalculateETA)
		tracking.GET("/eta", h.GetETA)
	}
}

// UpdateLocation updates location for an order (courier/provider)
// @Summary Update location
// @Description Update location for an order (Courier/Provider only)
// @Tags Tracking
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param request body model.LocationUpdateRequest true "Location update request"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /orders/{id}/location [post]
func (h *LocationTrackingHandler) UpdateLocation(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"invalid_id",
			"Invalid ID format",
			nil,
		))
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.CreateErrorResponse(
			"unauthorized",
			"User ID not found in context",
			nil,
		))
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, model.CreateErrorResponse(
			"unauthorized",
			"Invalid user ID type",
			nil,
		))
		return
	}

	var req model.LocationUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"validation_error",
			"Invalid request data",
			err.Error(),
		))
		return
	}

	location, err := h.trackingService.UpdateLocation(c.Request.Context(), orderID, userUUID, &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == model.ErrOrderNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, model.CreateErrorResponse(
			"operation_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(location, "Location updated successfully"))
}

// GetCurrentLocation retrieves current location for an order (customer)
// @Summary Get current location
// @Description Get current location for an order (Customer can view)
// @Tags Tracking
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} model.APIResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /orders/{id}/location [get]
func (h *LocationTrackingHandler) GetCurrentLocation(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"invalid_id",
			"Invalid ID format",
			nil,
		))
		return
	}

	location, err := h.trackingService.GetCurrentLocation(c.Request.Context(), orderID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == model.ErrOrderNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, model.CreateErrorResponse(
			"operation_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(location, "Location retrieved successfully"))
}

// GetLocationHistory retrieves location history for an order
// @Summary Get location history
// @Description Get location history for an order
// @Tags Tracking
// @Produce json
// @Param id path string true "Order ID"
// @Param limit query int false "Limit number of records" default(50)
// @Success 200 {object} model.APIResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /orders/{id}/location/history [get]
func (h *LocationTrackingHandler) GetLocationHistory(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"invalid_id",
			"Invalid ID format",
			nil,
		))
		return
	}

	limit := 50 // default
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	history, err := h.trackingService.GetLocationHistory(c.Request.Context(), orderID, limit)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == model.ErrOrderNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, model.CreateErrorResponse(
			"operation_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(history, "Location history retrieved successfully"))
}

// CalculateETA calculates ETA based on current location
// @Summary Calculate ETA
// @Description Calculate ETA based on current location and destination
// @Tags Tracking
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param request body model.ETACalculationRequest true "ETA calculation request"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.ErrorResponse
// @Router /orders/{id}/eta [post]
func (h *LocationTrackingHandler) CalculateETA(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"invalid_id",
			"Invalid ID format",
			nil,
		))
		return
	}

	var req model.ETACalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"validation_error",
			"Invalid request data",
			err.Error(),
		))
		return
	}

	eta, err := h.trackingService.CalculateETA(c.Request.Context(), orderID, &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == model.ErrOrderNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, model.CreateErrorResponse(
			"operation_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(eta, "ETA calculated successfully"))
}

// GetETA retrieves ETA from current location stored in order
// @Summary Get ETA
// @Description Get ETA from current location stored in order
// @Tags Tracking
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} model.APIResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /orders/{id}/eta [get]
func (h *LocationTrackingHandler) GetETA(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"invalid_id",
			"Invalid ID format",
			nil,
		))
		return
	}

	location, err := h.trackingService.GetCurrentLocation(c.Request.Context(), orderID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == model.ErrOrderNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, model.CreateErrorResponse(
			"operation_failed",
			err.Error(),
			nil,
		))
		return
	}

	etaResponse := map[string]interface{}{
		"eta":      location.ETA,
		"distance": location.Distance,
	}

	c.JSON(http.StatusOK, model.SuccessResponse(etaResponse, "ETA retrieved successfully"))
}

