package core

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuditAction represents the type of action performed
type AuditAction string

const (
	AuditActionCreate   AuditAction = "create"
	AuditActionUpdate   AuditAction = "update"
	AuditActionDelete   AuditAction = "delete"
	AuditActionLogin    AuditAction = "login"
	AuditActionLogout   AuditAction = "logout"
	AuditActionView     AuditAction = "view"
	AuditActionExport   AuditAction = "export"
	AuditActionPayment  AuditAction = "payment"
	AuditActionApproval AuditAction = "approval"
)

// AuditTrail represents an audit trail entry
type AuditTrail struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID      *uuid.UUID     `json:"user_id,omitempty" gorm:"type:uuid"`
	User        *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Action      AuditAction    `json:"action" gorm:"not null"`
	Resource    string         `json:"resource" gorm:"not null"`           // e.g., "order", "payment", "user"
	ResourceID  *uuid.UUID     `json:"resource_id,omitempty" gorm:"type:uuid"` // ID of the affected resource
	IPAddress   string         `json:"ip_address,omitempty" gorm:"type:varchar(45)"`
	UserAgent   string         `json:"user_agent,omitempty" gorm:"type:text"`
	RequestID   string         `json:"request_id,omitempty" gorm:"type:varchar(100)"`
	OldValue    string         `json:"old_value,omitempty" gorm:"type:text"` // JSON string of old values
	NewValue    string         `json:"new_value,omitempty" gorm:"type:text"`   // JSON string of new values
	Changes     string         `json:"changes,omitempty" gorm:"type:text"`     // JSON string of changed fields
	Status      string         `json:"status,omitempty" gorm:"type:varchar(50)"` // "success", "failed", "error"
	ErrorMessage string        `json:"error_message,omitempty" gorm:"type:text"`
	CreatedAt time.Time       `json:"created_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName returns the table name for AuditTrail
func (AuditTrail) TableName() string {
	return "audit_trails"
}

// AuditTrailRequest represents request payload for creating audit trail
type AuditTrailRequest struct {
	Action      AuditAction `json:"action" validate:"required"`
	Resource    string      `json:"resource" validate:"required"`
	ResourceID *string     `json:"resource_id,omitempty"`
	IPAddress   string      `json:"ip_address,omitempty"`
	UserAgent   string      `json:"user_agent,omitempty"`
	RequestID   string      `json:"request_id,omitempty"`
	OldValue    string      `json:"old_value,omitempty"`
	NewValue    string      `json:"new_value,omitempty"`
	Changes     string      `json:"changes,omitempty"`
	Status      string      `json:"status,omitempty"`
	ErrorMessage string     `json:"error_message,omitempty"`
}

