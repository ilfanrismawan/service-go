package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	userModel "service-go/internal/modules/users/model"
	orderModel "service-go/internal/modules/orders/model"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeEmail            NotificationType = "email"
	NotificationTypeWhatsApp         NotificationType = "whatsapp"
	NotificationTypePush             NotificationType = "push"
	NotificationTypeSMS              NotificationType = "sms"
	NotificationTypeOrderUpdate      NotificationType = "order_update"
	NotificationTypeOrderReady       NotificationType = "order_ready"
	NotificationTypeOrderDelivered   NotificationType = "order_delivered"
	NotificationTypeOrderCompleted   NotificationType = "order_completed"
	NotificationTypeOrderCancelled   NotificationType = "order_cancelled"
	NotificationTypePaymentPending   NotificationType = "payment_pending"
	NotificationTypePaymentReceived  NotificationType = "payment_received"
	NotificationTypePaymentFailed    NotificationType = "payment_failed"
	NotificationTypePaymentCancelled NotificationType = "payment_cancelled"
	NotificationTypePaymentRefunded  NotificationType = "payment_refunded"
	NotificationTypePaymentUpdate    NotificationType = "payment_update"
	NotificationTypeWelcome          NotificationType = "welcome"
	NotificationTypePromotion        NotificationType = "promotion"
	NotificationTypeSystem           NotificationType = "system"
)

// NotificationStatus represents the status of a notification
type NotificationStatus string

const (
	NotificationStatusPending NotificationStatus = "pending"
	NotificationStatusSent    NotificationStatus = "sent"
	NotificationStatusFailed  NotificationStatus = "failed"
)

// Notification represents a notification in the system
type Notification struct {
	ID        uuid.UUID          `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID    uuid.UUID          `json:"user_id" gorm:"type:uuid;not null"`
	User      userModel.User     `json:"user" gorm:"foreignKey:UserID"`
	OrderID   *uuid.UUID         `json:"order_id,omitempty" gorm:"type:uuid"`
	Order     *orderModel.ServiceOrder `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	Type      NotificationType   `json:"type" gorm:"not null"`
	Title     string             `json:"title" gorm:"not null"`
	Message   string             `json:"message" gorm:"not null"`
	Status    NotificationStatus `json:"status" gorm:"not null;default:'pending'"`
	IsRead    bool               `json:"is_read" gorm:"default:false"`
	SentAt    *time.Time         `json:"sent_at,omitempty"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	DeletedAt gorm.DeletedAt     `json:"-" gorm:"index"`
}

// TableName returns the table name for Notification
func (Notification) TableName() string {
	return "notifications"
}

// NotificationRequest represents the request payload for creating a notification
type NotificationRequest struct {
	UserID  string           `json:"user_id" validate:"required"`
	OrderID *string          `json:"order_id,omitempty"`
	Type    NotificationType `json:"type" validate:"required"`
	Title   string           `json:"title" validate:"required"`
	Message string           `json:"message" validate:"required"`
}

// NotificationResponse represents the response payload for notification data
type NotificationResponse struct {
	ID        uuid.UUID                      `json:"id"`
	UserID    uuid.UUID                      `json:"user_id"`
	User      userModel.UserResponse         `json:"user"`
	OrderID   *uuid.UUID                     `json:"order_id,omitempty"`
	Order     *orderModel.ServiceOrderResponse `json:"order,omitempty"`
	Type      NotificationType               `json:"type"`
	Title     string                         `json:"title"`
	Message   string                         `json:"message"`
	Status    NotificationStatus             `json:"status"`
	IsRead    bool                           `json:"is_read"`
	SentAt    *time.Time                     `json:"sent_at,omitempty"`
	CreatedAt time.Time                      `json:"created_at"`
	UpdatedAt time.Time                      `json:"updated_at"`
}

// ToResponse converts Notification to NotificationResponse
func (n *Notification) ToResponse() NotificationResponse {
	response := NotificationResponse{
		ID:        n.ID,
		UserID:    n.UserID,
		User:      n.User.ToResponse(),
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
		orderResponse := n.Order.ToResponse()
		response.Order = &orderResponse
	}

	return response
}

