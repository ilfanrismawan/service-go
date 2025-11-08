package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Branch represents a branch/outlet of the iPhone service company
type Branch struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name      string         `json:"name" gorm:"not null"`
	Address   string         `json:"address" gorm:"not null"`
	City      string         `json:"city" gorm:"not null"`
	Province  string         `json:"province" gorm:"not null"`
	Phone     string         `json:"phone" gorm:"not null"`
	Latitude  float64        `json:"latitude" gorm:"type:decimal(10,6);not null"`
	Longitude float64        `json:"longitude" gorm:"type:decimal(10,6);not null"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName returns the table name for Branch
func (Branch) TableName() string {
	return "branches"
}

