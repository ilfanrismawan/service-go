package dto

import (
	"time"

	"github.com/google/uuid"

	"service-go/internal/modules/chat/entity"
	orderDto "service-go/internal/modules/orders/dto"
	userDto "service-go/internal/modules/users/dto"
)

// ChatMessageRequest represents the request payload for creating a chat message
type ChatMessageRequest struct {
	OrderID    string `json:"order_id" validate:"required"`
	ReceiverID string `json:"receiver_id" validate:"required"`
	Message    string `json:"message" validate:"required"`
}

// ChatMessageResponse represents the response payload for chat message data
type ChatMessageResponse struct {
	ID         uuid.UUID                     `json:"id"`
	OrderID    uuid.UUID                     `json:"order_id"`
	Order      orderDto.ServiceOrderResponse `json:"order"`
	SenderID   uuid.UUID                     `json:"sender_id"`
	Sender     userDto.UserResponse          `json:"sender"`
	ReceiverID uuid.UUID                     `json:"receiver_id"`
	Receiver   userDto.UserResponse          `json:"receiver"`
	Message    string                        `json:"message"`
	IsRead     bool                          `json:"is_read"`
	CreatedAt  time.Time                     `json:"created_at"`
	UpdatedAt  time.Time                     `json:"updated_at"`
}

// ToChatMessageResponse converts ChatMessage entity to ChatMessageResponse DTO
func ToChatMessageResponse(cm *chatEntity.ChatMessage) ChatMessageResponse {
	return ChatMessageResponse{
		ID:         cm.ID,
		OrderID:    cm.OrderID,
		SenderID:   cm.SenderID,
		ReceiverID: cm.ReceiverID,
		Message:    cm.Message,
		IsRead:     cm.IsRead,
		CreatedAt:  cm.CreatedAt,
		UpdatedAt:  cm.UpdatedAt,
	}
}
