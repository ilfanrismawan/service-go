package handlers

import (
	"net/http"
	"service/internal/shared/model"
	service "service/internal/shared/service"
	"service/internal/shared/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// NotificationHandler handles notification endpoints
type NotificationHandler struct {
	notificationService *service.NotificationService
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{
		notificationService: service.NewNotificationService(),
	}
}

// SendNotification godoc
// @Summary Send notification
// @Description Send a notification to a user
// @Tags notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body core.NotificationRequest true "Notification data"
// @Success 201 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /notifications [post]
func (h *NotificationHandler) SendNotification(c *gin.Context) {
	var req core.NotificationRequest
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

	notification, err := h.notificationService.SendNotification(c.Request.Context(), &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == core.ErrUserNotFound {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, core.CreateErrorResponse(
			"notification_send_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusCreated, core.SuccessResponse(notification, "Notification sent successfully"))
}

// GetNotifications godoc
// @Summary Get user notifications
// @Description Get notifications for the current user
// @Tags notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} core.PaginatedResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /notifications [get]
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
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

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, core.CreateErrorResponse(
			"internal_error",
			"Invalid user ID type",
			nil,
		))
		return
	}

	// Parse query parameters
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

	result, err := h.notificationService.GetNotifications(c.Request.Context(), userUUID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, core.CreateErrorResponse(
			"notifications_fetch_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, result)
}

// MarkAsRead godoc
// @Summary Mark notification as read
// @Description Mark a notification as read
// @Tags notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Notification ID"
// @Success 200 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 404 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /notifications/{id}/read [put]
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"invalid_id",
			"Invalid notification ID format",
			nil,
		))
		return
	}

	err = h.notificationService.MarkAsRead(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, core.CreateErrorResponse(
			"mark_read_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, core.SuccessResponse(nil, "Notification marked as read"))
}

// SendOrderStatusNotification godoc
// @Summary Send order status notification
// @Description Send notification when order status changes
// @Tags notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orderId path string true "Order ID"
// @Param status query string true "Order Status"
// @Success 200 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /notifications/order/{orderId}/status [post]
func (h *NotificationHandler) SendOrderStatusNotification(c *gin.Context) {
	orderIDStr := c.Param("orderId")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"invalid_order_id",
			"Invalid order ID format",
			nil,
		))
		return
	}

	statusStr := c.Query("status")
	status := core.OrderStatus(statusStr)

	// Validate status
	validStatuses := []core.OrderStatus{
		core.StatusPendingPickup,
		core.StatusOnPickup,
		core.StatusInService,
		core.StatusReady,
		core.StatusDelivered,
		core.StatusCompleted,
		core.StatusCancelled,
	}

	valid := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			valid = true
			break
		}
	}

	if !valid {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"invalid_status",
			"Invalid order status",
			nil,
		))
		return
	}

	err = h.notificationService.SendOrderStatusNotification(c.Request.Context(), orderID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, core.CreateErrorResponse(
			"order_status_notification_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, core.SuccessResponse(nil, "Order status notification sent successfully"))
}

// SendPaymentNotification godoc
// @Summary Send payment notification
// @Description Send notification for payment updates
// @Tags notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orderId path string true "Order ID"
// @Param status query string true "Payment Status"
// @Success 200 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /notifications/order/{orderId}/payment [post]
func (h *NotificationHandler) SendPaymentNotification(c *gin.Context) {
	orderIDStr := c.Param("orderId")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"invalid_order_id",
			"Invalid order ID format",
			nil,
		))
		return
	}

	statusStr := c.Query("status")
	status := core.PaymentStatus(statusStr)

	// Validate status
	validStatuses := []core.PaymentStatus{
		core.PaymentStatusPending,
		core.PaymentStatusPaid,
		core.PaymentStatusFailed,
		core.PaymentStatusCancelled,
		core.PaymentStatusRefunded,
	}

	valid := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			valid = true
			break
		}
	}

	if !valid {
		c.JSON(http.StatusBadRequest, core.CreateErrorResponse(
			"invalid_status",
			"Invalid payment status",
			nil,
		))
		return
	}

	err = h.notificationService.SendPaymentNotification(c.Request.Context(), orderID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, core.CreateErrorResponse(
			"payment_notification_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, core.SuccessResponse(nil, "Payment notification sent successfully"))
}
