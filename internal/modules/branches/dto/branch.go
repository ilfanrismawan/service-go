package dto

import (
	"time"

	"github.com/google/uuid"

	branchEntity "service-go/internal/modules/branches/entity"
)

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

// BranchDistance represents a branch with distance information
type BranchDistance struct {
	Branch   BranchResponse `json:"branch"`
	Distance float64        `json:"distance"` // in kilometers
}

// ToBranchResponse converts Branch entity to BranchResponse DTO
func ToBranchResponse(b *branchEntity.Branch) BranchResponse {
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

