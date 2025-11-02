package dto

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

// BranchRequest represents the request payload for creating/updating a branch
type BranchRequest struct {
	Name      string  `json:"name" validate:"required"`
	Address   string  `json:"address" validate:"required"`
	City      string  `json:"city" validate:"required"`
	Province  string  `json:"province" validate:"required"`
    Phone     string  `json:"phone" validate:"required,phone"`
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

// BranchResponse represents the response payload for branch data
type BranchResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	City      string    `json:"city"`
	Province  string    `json:"province"`
	Phone     string    `json:"phone"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converts Branch to BranchResponse
func (b *Branch) ToResponse() BranchResponse {
	return BranchResponse{
		ID:        b.ID,
		Name:      b.Name,
		Address:   b.Address,
		City:      b.City,
		Province:  b.Province,
		Phone:     b.Phone,
		Latitude:  b.Latitude,
		Longitude: b.Longitude,
		IsActive:  b.IsActive,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}

// BranchDistance represents a branch with distance information
type BranchDistance struct {
	Branch   BranchResponse `json:"branch"`
	Distance float64        `json:"distance"` // in kilometers
}
