package handlers

import (
	"net/http"
	"service/internal/modules/admin/service"
	"service/internal/shared/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DashboardHandler handles dashboard endpoints
type DashboardHandler struct {
	dashboardService *service.DashboardService
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{
		dashboardService: service.NewDashboardService(),
	}
}

// GetDashboardStats godoc
// @Summary Get dashboard statistics
// @Description Get dashboard statistics for the current user
// @Tags dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param branch_id query string false "Branch ID (for admin users)"
// @Success 200 {object} core.APIResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /dashboard/stats [get]
func (h *DashboardHandler) GetDashboardStats(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
			"internal_error",
			"Invalid user ID type",
			nil,
		))
		return
	}

	// Parse branch ID if provided
	var branchID *uuid.UUID
	if branchIDStr := c.Query("branch_id"); branchIDStr != "" {
		if parsedBranchID, err := uuid.Parse(branchIDStr); err == nil {
			branchID = &parsedBranchID
		}
	}

	stats, err := h.dashboardService.GetDashboardStats(c.Request.Context(), &userUUID, branchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
			"dashboard_stats_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(stats, "Dashboard statistics retrieved successfully"))
}

// GetServiceStats godoc
// @Summary Get service statistics
// @Description Get service type statistics
// @Tags dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param branch_id query string false "Branch ID (for admin users)"
// @Success 200 {object} core.APIResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /dashboard/service-stats [get]
func (h *DashboardHandler) GetServiceStats(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
			"internal_error",
			"Invalid user ID type",
			nil,
		))
		return
	}

	// Parse branch ID if provided
	var branchID *uuid.UUID
	if branchIDStr := c.Query("branch_id"); branchIDStr != "" {
		if parsedBranchID, err := uuid.Parse(branchIDStr); err == nil {
			branchID = &parsedBranchID
		}
	}

	stats, err := h.dashboardService.GetServiceStats(c.Request.Context(), &userUUID, branchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
			"service_stats_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(stats, "Service statistics retrieved successfully"))
}

// GetRevenueReport godoc
// @Summary Get revenue report
// @Description Get revenue report for a date range
// @Tags dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param branch_id query string false "Branch ID (for admin users)"
// @Param date_from query string true "Start date (YYYY-MM-DD)"
// @Param date_to query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /dashboard/revenue-report [get]
func (h *DashboardHandler) GetRevenueReport(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
			"internal_error",
			"Invalid user ID type",
			nil,
		))
		return
	}

	// Parse query parameters
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	if dateFrom == "" || dateTo == "" {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"missing_parameters",
			"date_from and date_to are required",
			nil,
		))
		return
	}

	// Parse branch ID if provided
	var branchID *uuid.UUID
	if branchIDStr := c.Query("branch_id"); branchIDStr != "" {
		if parsedBranchID, err := uuid.Parse(branchIDStr); err == nil {
			branchID = &parsedBranchID
		}
	}

	report, err := h.dashboardService.GetRevenueReport(c.Request.Context(), &userUUID, branchID, dateFrom, dateTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
			"revenue_report_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(report, "Revenue report retrieved successfully"))
}

// GetBranchStats godoc
// @Summary Get branch statistics
// @Description Get statistics for a specific branch
// @Tags dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param branchId path string true "Branch ID"
// @Success 200 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /dashboard/branches/{branchId}/stats [get]
func (h *DashboardHandler) GetBranchStats(c *gin.Context) {
	branchIDStr := c.Param("branchId")
	branchID, err := uuid.Parse(branchIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"invalid_branch_id",
			"Invalid branch ID format",
			nil,
		))
		return
	}

	stats, err := h.dashboardService.GetBranchStats(c.Request.Context(), branchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
			"branch_stats_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(stats, "Branch statistics retrieved successfully"))
}

// GetUserStats godoc
// @Summary Get user statistics
// @Description Get statistics for a specific user
// @Tags dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User ID"
// @Success 200 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /dashboard/users/{userId}/stats [get]
func (h *DashboardHandler) GetUserStats(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"invalid_user_id",
			"Invalid user ID format",
			nil,
		))
		return
	}

	stats, err := h.dashboardService.GetUserStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
			"user_stats_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(stats, "User statistics retrieved successfully"))
}

// GetPopularServices godoc
// @Summary Get popular services
// @Description Get most popular service types
// @Tags dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Number of services to return" default(5)
// @Param branch_id query string false "Branch ID (for admin users)"
// @Success 200 {object} core.APIResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /dashboard/popular-services [get]
func (h *DashboardHandler) GetPopularServices(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
			"internal_error",
			"Invalid user ID type",
			nil,
		))
		return
	}

	// Parse limit
	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 20 {
		limit = 5
	}

	// Parse branch ID if provided
	var branchID *uuid.UUID
	if branchIDStr := c.Query("branch_id"); branchIDStr != "" {
		if parsedBranchID, err := uuid.Parse(branchIDStr); err == nil {
			branchID = &parsedBranchID
		}
	}

	stats, err := h.dashboardService.GetServiceStats(c.Request.Context(), &userUUID, branchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
			"popular_services_failed",
			err.Error(),
			nil,
		))
		return
	}

	// Limit results
	if len(stats) > limit {
		stats = stats[:limit]
	}

	c.JSON(http.StatusOK, model.SuccessResponse(stats, "Popular services retrieved successfully"))
}

// GetOverview godoc
// @Summary Dashboard overview
// @Description Get overview statistics for dashboard
// @Tags dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} core.APIResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /dashboard/overview [get]
func (h *DashboardHandler) GetOverview(c *gin.Context) {
	h.GetDashboardStats(c)
}

// GetOrderStats godoc
// @Summary Order statistics
// @Description Get order related statistics
// @Tags dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} core.APIResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /dashboard/orders [get]
func (h *DashboardHandler) GetOrderStats(c *gin.Context) {
	// Reuse service stats for now
	h.GetServiceStats(c)
}

// GetRevenueStats godoc
// @Summary Revenue statistics
// @Description Get revenue related statistics
// @Tags dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} core.APIResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /dashboard/revenue [get]
func (h *DashboardHandler) GetRevenueStats(c *gin.Context) {
	// If date range provided, use revenue report; otherwise return overview
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")
	if dateFrom != "" && dateTo != "" {
		h.GetRevenueReport(c)
		return
	}
	h.GetDashboardStats(c)
}

// GetAdminDashboard returns admin dashboard (admin)
func (h *DashboardHandler) GetAdminDashboard(c *gin.Context) {
	// Call GetDashboardStats with nil user to get overall stats
	stats, err := h.dashboardService.GetDashboardStats(c.Request.Context(), nil, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
			"admin_dashboard_failed",
			err.Error(),
			nil,
		))
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(stats, "Admin dashboard retrieved successfully"))
}
