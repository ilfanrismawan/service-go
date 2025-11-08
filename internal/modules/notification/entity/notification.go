package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	userEntity "service-go/internal/modules/users/entity"
	orderEntity "service-go/internal/modules/orders/entity"
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
	User      userEntity.User    `json:"user" gorm:"foreignKey:UserID"`
	OrderID   *uuid.UUID         `json:"order_id,omitempty" gorm:"type:uuid"`
	Order     *orderEntity.ServiceOrder `json:"order,omitempty" gorm:"foreignKey:OrderID"`
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

