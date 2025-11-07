package handlers

import (
	"net/http"
	"service/internal/modules/chat/service"
	"service/internal/shared/model"
	"service/internal/shared/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ChatHandler handles chat endpoints
type ChatHandler struct {
	chatService *service.ChatService
}

// NewChatHandler creates a new chat handler
func NewChatHandler() *ChatHandler {
	return &ChatHandler{
		chatService: service.NewChatService(),
	}
}

// SendMessage godoc
// @Summary Send chat message
// @Description Send a chat message in an order conversation
// @Tags chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body core.ChatMessageRequest true "Chat message data"
// @Success 201 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /chat/messages [post]
func (h *ChatHandler) SendMessage(c *gin.Context) {
	// Get sender ID from context
	senderID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.CreateErrorResponse(
			"unauthorized",
			"User ID not found in context",
			nil,
		))
		return
	}

	senderUUID, ok := senderID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
			"internal_error",
			"Invalid user ID type",
			nil,
		))
		return
	}

	var req model.ChatMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"validation_error",
			"Invalid request data",
			err.Error(),
		))
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"validation_error",
			"Validation failed",
			err.Error(),
		))
		return
	}

	message, err := h.chatService.SendMessage(c.Request.Context(), senderUUID, &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == model.ErrOrderNotFound || err == model.ErrUserNotFound {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, model.CreateErrorResponse(
			"message_send_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusCreated, model.SuccessResponse(message, "Message sent successfully"))
}

// GetChatMessages godoc
// @Summary Get chat messages
// @Description Get chat messages for an order
// @Tags chat
// @Accept json
// @Produce json
// @Param orderId path string true "Order ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} core.PaginatedResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /chat/orders/{orderId}/messages [get]
func (h *ChatHandler) GetChatMessages(c *gin.Context) {
	orderIDStr := c.Param("orderId")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"invalid_order_id",
			"Invalid order ID format",
			nil,
		))
		return
	}

	// Parse query parameters
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	result, err := h.chatService.GetChatMessages(c.Request.Context(), orderID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
			"chat_messages_fetch_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetUserChats godoc
// @Summary Get user chats
// @Description Get all chat conversations for the current user
// @Tags chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} core.APIResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /chat/conversations [get]
func (h *ChatHandler) GetUserChats(c *gin.Context) {
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

	chats, err := h.chatService.GetUserChats(c.Request.Context(), userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
			"user_chats_fetch_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(chats, "User chats retrieved successfully"))
}

// MarkAsRead godoc
// @Summary Mark message as read
// @Description Mark a chat message as read
// @Tags chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Message ID"
// @Success 200 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 404 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /chat/messages/{id}/read [put]
func (h *ChatHandler) MarkAsRead(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"invalid_id",
			"Invalid message ID format",
			nil,
		))
		return
	}

	err = h.chatService.MarkAsRead(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
			"mark_read_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil, "Message marked as read"))
}

// MarkOrderMessagesAsRead godoc
// @Summary Mark order messages as read
// @Description Mark all messages in an order as read for the current user
// @Tags chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orderId path string true "Order ID"
// @Success 200 {object} core.APIResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /chat/orders/{orderId}/read [put]
func (h *ChatHandler) MarkOrderMessagesAsRead(c *gin.Context) {
	orderIDStr := c.Param("orderId")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"invalid_order_id",
			"Invalid order ID format",
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
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
			"internal_error",
			"Invalid user ID type",
			nil,
		))
		return
	}

	err = h.chatService.MarkOrderMessagesAsRead(c.Request.Context(), orderID, userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
			"mark_order_read_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil, "Order messages marked as read"))
}

// GetUnreadCount godoc
// @Summary Get unread message count
// @Description Get the count of unread messages for the current user
// @Tags chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} core.APIResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /chat/unread-count [get]
func (h *ChatHandler) GetUnreadCount(c *gin.Context) {
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

	count, err := h.chatService.GetUnreadCount(c.Request.Context(), userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
			"unread_count_fetch_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(gin.H{"unread_count": count}, "Unread count retrieved successfully"))
}
