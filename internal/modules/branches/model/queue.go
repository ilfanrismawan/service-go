package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	userModel "service-go/internal/modules/users/model"
)

// QueueStatus represents the status of a queue
type QueueStatus string

const (
	QueueStatusWaiting   QueueStatus = "waiting"
	QueueStatusServed    QueueStatus = "served"
	QueueStatusCancelled QueueStatus = "cancelled"
)

// Queue represents a queue for walk-in customers
type Queue struct {
	ID         uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	BranchID   uuid.UUID      `json:"branch_id" gorm:"type:uuid;not null"`
	Branch     Branch          `json:"branch" gorm:"foreignKey:BranchID"`
	QueueNo    string          `json:"queue_no" gorm:"not null"` // "A001"
	CustomerID *uuid.UUID      `json:"customer_id,omitempty" gorm:"type:uuid"`
	Customer   *userModel.User `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	Status     QueueStatus     `json:"status" gorm:"not null;default:'waiting'"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	DeletedAt  gorm.DeletedAt  `json:"-" gorm:"index"`
}

// TableName returns the table name for Queue
func (Queue) TableName() string {
	return "queues"
}

