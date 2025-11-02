package delivery

import (
	"net/http"
	"service/internal/core"
	"service/internal/orders/service"
	"service/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// MembershipHandler handles membership endpoints
type MembershipHandler struct {
	membershipService *service.MembershipService
}

// NewMembershipHandler creates a new membership handler
func NewMembershipHandler() *MembershipHandler {
	return &MembershipHandler{
		membershipService: service.NewMembershipService(),
	}
}

// CreateMembership godoc
// @Summary Create a new membership
// @Description Create a new membership for a user
// @Tags membership
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body core.MembershipRequest true "Membership data"
// @Success 201 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /membership [post]
func (h *MembershipHandler) CreateMembership(c *gin.Context) {
	var req core.MembershipRequest
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

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, core.CreateErrorResponse(
			"unauthorized",
			"User not authenticated",
			nil,
		))
		return
	}

	membership, err := h.membershipService.CreateMembership(c.Request.Context(), userID.(uuid.UUID), req.Tier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, core.CreateErrorResponse(
			"membership_creation_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusCreated, core.SuccessResponse(membership.ToResponse(), "Membership created successfully"))
}

// GetMembership godoc
// @Summary Get user membership
// @Description Get membership details for the authenticated user
// @Tags membership
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} core.APIResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 404 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /membership [get]
func (h *MembershipHandler) GetMembership(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, core.CreateErrorResponse(
			"unauthorized",
			"User not authenticated",
			nil,
		))
		return
	}

	membership, err := h.membershipService.GetMembership(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, core.CreateErrorResponse(
			"membership_not_found",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, core.SuccessResponse(membership.ToResponse(), "Membership retrieved successfully"))
}

// UpdateMembership godoc
// @Summary Update membership
// @Description Update membership tier or status
// @Tags membership
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body core.MembershipRequest true "Membership update data"
// @Success 200 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 404 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /membership [put]
func (h *MembershipHandler) UpdateMembership(c *gin.Context) {
	var req core.MembershipRequest
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

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, core.CreateErrorResponse(
			"unauthorized",
			"User not authenticated",
			nil,
		))
		return
	}

	membership, err := h.membershipService.UpdateMembership(c.Request.Context(), userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, core.CreateErrorResponse(
			"membership_update_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, core.SuccessResponse(membership.ToResponse(), "Membership updated successfully"))
}

// ListMemberships godoc
// @Summary List memberships
// @Description Get list of memberships with pagination and filters
// @Tags membership
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param tier query string false "Filter by tier"
// @Param status query string false "Filter by status"
// @Success 200 {object} core.PaginatedResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /membership/list [get]
func (h *MembershipHandler) ListMemberships(c *gin.Context) {
	// Parse query parameters
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	tierStr := c.Query("tier")
	statusStr := c.Query("status")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	// Set filter pointers
	var tierPtr *core.MembershipTier
	var statusPtr *core.MembershipStatus
	if tierStr != "" {
		tier := core.MembershipTier(tierStr)
		tierPtr = &tier
	}
	if statusStr != "" {
		status := core.MembershipStatus(statusStr)
		statusPtr = &status
	}

	memberships, total, err := h.membershipService.ListMemberships(c.Request.Context(), page, limit, tierPtr, statusPtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, core.CreateErrorResponse(
			"membership_list_failed",
			err.Error(),
			nil,
		))
		return
	}

	// Convert to response format
	var responses []core.MembershipResponse
	for _, membership := range memberships {
		responses = append(responses, membership.ToResponse())
	}

	// Create paginated response
	paginatedResponse := core.PaginatedResponse{
		Status: "success",
		Data:   responses,
		Pagination: core.PaginationResponse{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: int((total + int64(limit) - 1) / int64(limit)),
		},
		Message:   "Memberships retrieved successfully",
		Timestamp: core.GetCurrentTimestamp(),
	}

	c.JSON(http.StatusOK, paginatedResponse)
}

// GetMembershipStats godoc
// @Summary Get membership statistics
// @Description Get membership statistics and analytics
// @Tags membership
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} core.APIResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /membership/stats [get]
func (h *MembershipHandler) GetMembershipStats(c *gin.Context) {
	stats, err := h.membershipService.GetMembershipStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, core.CreateErrorResponse(
			"membership_stats_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, core.SuccessResponse(stats, "Membership statistics retrieved successfully"))
}

// GetTopSpenders godoc
// @Summary Get top spending members
// @Description Get list of top spending members
// @Tags membership
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Number of top spenders" default(10)
// @Success 200 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /membership/top-spenders [get]
func (h *MembershipHandler) GetTopSpenders(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	memberships, err := h.membershipService.GetTopSpenders(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, core.CreateErrorResponse(
			"top_spenders_failed",
			err.Error(),
			nil,
		))
		return
	}

	// Convert to response format
	var responses []core.MembershipResponse
	for _, membership := range memberships {
		responses = append(responses, membership.ToResponse())
	}

	c.JSON(http.StatusOK, core.SuccessResponse(responses, "Top spenders retrieved successfully"))
}

// RedeemPoints godoc
// @Summary Redeem membership points
// @Description Redeem membership points for discount
// @Tags membership
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param points query int true "Points to redeem"
// @Success 200 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /membership/redeem-points [post]
func (h *MembershipHandler) RedeemPoints(c *gin.Context) {
	pointsStr := c.Query("points")
	if pointsStr == "" {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"validation_error",
			"Points parameter is required",
			nil,
		))
		return
	}

	points, err := strconv.ParseInt(pointsStr, 10, 64)
	if err != nil || points <= 0 {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"validation_error",
			"Invalid points value",
			nil,
		))
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, core.CreateErrorResponse(
			"unauthorized",
			"User not authenticated",
			nil,
		))
		return
	}

	result, err := h.membershipService.RedeemPoints(c.Request.Context(), userID.(uuid.UUID), points)
	if err != nil {
		c.JSON(http.StatusInternalServerError, core.CreateErrorResponse(
			"points_redemption_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, core.SuccessResponse(result, "Points redeemed successfully"))
}

// SubscribeToMembership godoc
// @Summary Subscribe to membership
// @Description Subscribe to a membership tier with payment
// @Tags membership
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body core.SubscriptionRequest true "Subscription data"
// @Success 201 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /membership/subscribe [post]
func (h *MembershipHandler) SubscribeToMembership(c *gin.Context) {
	var req core.SubscriptionRequest
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

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, core.CreateErrorResponse(
			"unauthorized",
			"User not authenticated",
			nil,
		))
		return
	}

	result, err := h.membershipService.SubscribeToMembership(c.Request.Context(), userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, core.CreateErrorResponse(
			"subscription_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusCreated, core.SuccessResponse(result, "Subscription created successfully"))
}

// CancelSubscription godoc
// @Summary Cancel subscription
// @Description Cancel user's membership subscription
// @Tags membership
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body core.CancelSubscriptionRequest true "Cancel subscription data"
// @Success 200 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /membership/cancel [post]
func (h *MembershipHandler) CancelSubscription(c *gin.Context) {
	var req core.CancelSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"validation_error",
			"Invalid request data",
			err.Error(),
		))
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, core.CreateErrorResponse(
			"unauthorized",
			"User not authenticated",
			nil,
		))
		return
	}

	err := h.membershipService.CancelSubscription(c.Request.Context(), userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, core.CreateErrorResponse(
			"cancellation_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, core.SuccessResponse(nil, "Subscription cancelled successfully"))
}

// StartTrial godoc
// @Summary Start trial membership
// @Description Start a 7-day trial membership
// @Tags membership
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param tier query string true "Membership tier" Enums(basic, premium, vip, elite)
// @Success 201 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /membership/trial [post]
func (h *MembershipHandler) StartTrial(c *gin.Context) {
	tierStr := c.Query("tier")
	if tierStr == "" {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"validation_error",
			"Tier parameter is required",
			nil,
		))
		return
	}

	tier := core.MembershipTier(tierStr)
	if tier != core.MembershipTierBasic && tier != core.MembershipTierPremium &&
		tier != core.MembershipTierVIP && tier != core.MembershipTierElite {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"validation_error",
			"Invalid tier. Must be one of: basic, premium, vip, elite",
			nil,
		))
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, core.CreateErrorResponse(
			"unauthorized",
			"User not authenticated",
			nil,
		))
		return
	}

	membership, err := h.membershipService.StartTrial(c.Request.Context(), userID.(uuid.UUID), tier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, core.CreateErrorResponse(
			"trial_start_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusCreated, core.SuccessResponse(membership.ToResponse(), "Trial started successfully"))
}

// GetMembershipTiers godoc
// @Summary Get membership tiers
// @Description Get all available membership tiers and their benefits
// @Tags membership
// @Accept json
// @Produce json
// @Success 200 {object} core.APIResponse
// @Router /membership/tiers [get]
func (h *MembershipHandler) GetMembershipTiers(c *gin.Context) {
	tiers := h.membershipService.GetMembershipTiers(c.Request.Context())
	c.JSON(http.StatusOK, core.SuccessResponse(tiers, "Membership tiers retrieved successfully"))
}

// UpgradeMembership godoc
// @Summary Upgrade membership
// @Description Upgrade user's membership to a higher tier
// @Tags membership
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param tier query string true "New membership tier" Enums(basic, premium, vip, elite)
// @Param subscription_type query string true "Subscription type" Enums(monthly, yearly)
// @Success 200 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /membership/upgrade [post]
func (h *MembershipHandler) UpgradeMembership(c *gin.Context) {
	tierStr := c.Query("tier")
	subscriptionTypeStr := c.Query("subscription_type")

	if tierStr == "" || subscriptionTypeStr == "" {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"validation_error",
			"Tier and subscription_type parameters are required",
			nil,
		))
		return
	}

	tier := core.MembershipTier(tierStr)
	subscriptionType := core.SubscriptionType(subscriptionTypeStr)

	if tier != core.MembershipTierBasic && tier != core.MembershipTierPremium &&
		tier != core.MembershipTierVIP && tier != core.MembershipTierElite {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"validation_error",
			"Invalid tier. Must be one of: basic, premium, vip, elite",
			nil,
		))
		return
	}

	if subscriptionType != core.SubscriptionTypeMonthly && subscriptionType != core.SubscriptionTypeYearly {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"validation_error",
			"Invalid subscription type. Must be one of: monthly, yearly",
			nil,
		))
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, core.CreateErrorResponse(
			"unauthorized",
			"User not authenticated",
			nil,
		))
		return
	}

	result, err := h.membershipService.UpgradeMembership(c.Request.Context(), userID.(uuid.UUID), tier, subscriptionType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, core.CreateErrorResponse(
			"upgrade_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, core.SuccessResponse(result, "Membership upgraded successfully"))
}

// GetMembershipUsage godoc
// @Summary Get membership usage
// @Description Get current month usage for membership benefits
// @Tags membership
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} core.APIResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /membership/usage [get]
func (h *MembershipHandler) GetMembershipUsage(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, core.CreateErrorResponse(
			"unauthorized",
			"User not authenticated",
			nil,
		))
		return
	}

	usage, err := h.membershipService.GetMembershipUsage(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, core.CreateErrorResponse(
			"usage_retrieval_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, core.SuccessResponse(usage, "Membership usage retrieved successfully"))
}
