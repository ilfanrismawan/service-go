package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	userEntity "service-go/internal/modules/users/entity"
	orderEntity "service-go/internal/modules/orders/entity"
)

// ChatMessage represents a chat message between customer and technician
type ChatMessage struct {
	ID         uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrderID    uuid.UUID      `json:"order_id" gorm:"type:uuid;not null"`
	Order      orderEntity.ServiceOrder `json:"order" gorm:"foreignKey:OrderID"`
	SenderID   uuid.UUID      `json:"sender_id" gorm:"type:uuid;not null"`
	Sender     userEntity.User `json:"sender" gorm:"foreignKey:SenderID"`
	ReceiverID uuid.UUID      `json:"receiver_id" gorm:"type:uuid;not null"`
	Receiver   userEntity.User `json:"receiver" gorm:"foreignKey:ReceiverID"`
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

