package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	userModel "service-go/internal/modules/users/model"
	orderModel "service-go/internal/modules/orders/model"
)

// ChatMessage represents a chat message between customer and technician
type ChatMessage struct {
	ID         uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrderID    uuid.UUID      `json:"order_id" gorm:"type:uuid;not null"`
	Order      orderModel.ServiceOrder `json:"order" gorm:"foreignKey:OrderID"`
	SenderID   uuid.UUID      `json:"sender_id" gorm:"type:uuid;not null"`
	Sender     userModel.User `json:"sender" gorm:"foreignKey:SenderID"`
	ReceiverID uuid.UUID      `json:"receiver_id" gorm:"type:uuid;not null"`
	Receiver   userModel.User `json:"receiver" gorm:"foreignKey:ReceiverID"`
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
	ID         uuid.UUID                      `json:"id"`
	OrderID    uuid.UUID                      `json:"order_id"`
	Order      orderModel.ServiceOrderResponse `json:"order"`
	SenderID   uuid.UUID                      `json:"sender_id"`
	Sender     userModel.UserResponse          `json:"sender"`
	ReceiverID uuid.UUID                      `json:"receiver_id"`
	Receiver   userModel.UserResponse          `json:"receiver"`
	Message    string                         `json:"message"`
	IsRead     bool                           `json:"is_read"`
	CreatedAt  time.Time                      `json:"created_at"`
	UpdatedAt  time.Time                      `json:"updated_at"`
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

