package dto

import (
	"time"

	"github.com/google/uuid"

	notificationEntity "service-go/internal/modules/notification/entity"
	userDto "service-go/internal/modules/users/dto"
	orderDto "service-go/internal/modules/orders/dto"
)

// NotificationRequest represents the request payload for creating a notification
type NotificationRequest struct {
	UserID  string                      `json:"user_id" validate:"required"`
	OrderID *string                     `json:"order_id,omitempty"`
	Type    notificationEntity.NotificationType `json:"type" validate:"required"`
	Title   string                      `json:"title" validate:"required"`
	Message string                      `json:"message" validate:"required"`
}

// NotificationResponse represents the response payload for notification data
type NotificationResponse struct {
	ID        uuid.UUID                      `json:"id"`
	UserID    uuid.UUID                      `json:"user_id"`
	User      userDto.UserResponse           `json:"user"`
	OrderID   *uuid.UUID                     `json:"order_id,omitempty"`
	Order     *orderDto.ServiceOrderResponse `json:"order,omitempty"`
	Type      notificationEntity.NotificationType `json:"type"`
	Title     string                         `json:"title"`
	Message   string                         `json:"message"`
	Status    notificationEntity.NotificationStatus `json:"status"`
	IsRead    bool                           `json:"is_read"`
	SentAt    *time.Time                     `json:"sent_at,omitempty"`
	CreatedAt time.Time                      `json:"created_at"`
	UpdatedAt time.Time                      `json:"updated_at"`
}

// ToNotificationResponse converts Notification entity to NotificationResponse DTO
func ToNotificationResponse(n *notificationEntity.Notification) NotificationResponse {
	response := NotificationResponse{
		ID:        n.ID,
		UserID:    n.UserID,
		User:      userDto.ToUserResponse(&n.User),
		Type:      n.Type,
		Title:     n.Title,
		Message:   n.Message,
		Status:    n.Status,
		IsRead:    n.IsRead,
		SentAt:    n.SentAt,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}

	if n.OrderID != nil {
		response.OrderID = n.OrderID
		if n.Order != nil {
			orderResp := orderDto.ToServiceOrderResponse(n.Order)
			response.Order = &orderResp
		}
	}

	return response
}

