package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	User      User               `json:"user" gorm:"foreignKey:UserID"`
	OrderID   *uuid.UUID         `json:"order_id,omitempty" gorm:"type:uuid"`
	Order     *ServiceOrder      `json:"order,omitempty" gorm:"foreignKey:OrderID"`
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
	ID        uuid.UUID             `json:"id"`
	UserID    uuid.UUID             `json:"user_id"`
	User      UserResponse          `json:"user"`
	OrderID   *uuid.UUID            `json:"order_id,omitempty"`
	Order     *ServiceOrderResponse `json:"order,omitempty"`
	Type      NotificationType      `json:"type"`
	Title     string                `json:"title"`
	Message   string                `json:"message"`
	Status    NotificationStatus    `json:"status"`
	IsRead    bool                  `json:"is_read"`
	SentAt    *time.Time            `json:"sent_at,omitempty"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
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

// ChatMessage represents a chat message between customer and technician
type ChatMessage struct {
	ID         uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrderID    uuid.UUID      `json:"order_id" gorm:"type:uuid;not null"`
	Order      ServiceOrder   `json:"order" gorm:"foreignKey:OrderID"`
	SenderID   uuid.UUID      `json:"sender_id" gorm:"type:uuid;not null"`
	Sender     User           `json:"sender" gorm:"foreignKey:SenderID"`
	ReceiverID uuid.UUID      `json:"receiver_id" gorm:"type:uuid;not null"`
	Receiver   User           `json:"receiver" gorm:"foreignKey:ReceiverID"`
	Message    string         `json:"message" gorm:"not null"`
	IsRead     bool           `json:"is_read" gorm:"default:false"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName returns the table name for ChatMessage
func (ChatMessage) TableName() string {
	return "chat_messages"
}

// ChatMessageRequest represents the request payload for creating a chat message
type ChatMessageRequest struct {
	OrderID    string `json:"order_id" validate:"required"`
	ReceiverID string `json:"receiver_id" validate:"required"`
	Message    string `json:"message" validate:"required"`
}

// ChatMessageResponse represents the response payload for chat message data
type ChatMessageResponse struct {
	ID         uuid.UUID            `json:"id"`
	OrderID    uuid.UUID            `json:"order_id"`
	Order      ServiceOrderResponse `json:"order"`
	SenderID   uuid.UUID            `json:"sender_id"`
	Sender     UserResponse         `json:"sender"`
	ReceiverID uuid.UUID            `json:"receiver_id"`
	Receiver   UserResponse         `json:"receiver"`
	Message    string               `json:"message"`
	IsRead     bool                 `json:"is_read"`
	CreatedAt  time.Time            `json:"created_at"`
	UpdatedAt  time.Time            `json:"updated_at"`
}

// ToResponse converts ChatMessage to ChatMessageResponse
func (cm *ChatMessage) ToResponse() ChatMessageResponse {
	return ChatMessageResponse{
		ID:         cm.ID,
		OrderID:    cm.OrderID,
		Order:      cm.Order.ToResponse(),
		SenderID:   cm.SenderID,
		Sender:     cm.Sender.ToResponse(),
		ReceiverID: cm.ReceiverID,
		Receiver:   cm.Receiver.ToResponse(),
		Message:    cm.Message,
		IsRead:     cm.IsRead,
		CreatedAt:  cm.CreatedAt,
		UpdatedAt:  cm.UpdatedAt,
	}
}
