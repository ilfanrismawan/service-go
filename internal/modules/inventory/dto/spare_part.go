package dto

import (
	"time"

	"github.com/google/uuid"

	branchDto "service-go/internal/modules/branches/dto"
)

// SparePartInventoryRequest represents request payload for creating/updating spare part
type SparePartInventoryRequest struct {
	BranchID string `json:"branch_id" validate:"required"`
	PartName string `json:"part_name" validate:"required"`
	PartCode string `json:"part_code" validate:"required"`
	Stock    int    `json:"stock" validate:"min=0"`
	MinStock int    `json:"min_stock" validate:"min=0"`
	Price    int64  `json:"price" validate:"min=0"`
	Supplier string `json:"supplier" validate:"required"`
}

// SparePartInventoryResponse represents response payload for spare part data
type SparePartInventoryResponse struct {
	ID           uuid.UUID                `json:"id"`
	BranchID     uuid.UUID                 `json:"branch_id"`
	Branch       branchDto.BranchResponse  `json:"branch"`
	PartName     string                    `json:"part_name"`
	PartCode     string                    `json:"part_code"`
	Stock        int                       `json:"stock"`
	MinStock     int                       `json:"min_stock"`
	Price        int64                     `json:"price"`
	Supplier     string                    `json:"supplier"`
	IsLowStock   bool                      `json:"is_low_stock"`
	NeedsReorder bool                      `json:"needs_reorder"`
	CreatedAt    time.Time                 `json:"created_at"`
	UpdatedAt    time.Time                 `json:"updated_at"`
}

// ToSparePartInventoryResponse converts SparePartInventory entity to SparePartInventoryResponse DTO
func ToSparePartInventoryResponse(s *inventoryEntity.SparePartInventory) SparePartInventoryResponse {
	response := SparePartInventoryResponse{
		ID:           s.ID,
		BranchID:     s.BranchID,
		PartName:     s.PartName,
		PartCode:     s.PartCode,
		Stock:        s.Stock,
		MinStock:     s.MinStock,
		Price:        s.Price,
		Supplier:     s.Supplier,
		IsLowStock:   s.IsLowStock(),
		NeedsReorder: s.NeedsReorder(),
		CreatedAt:    s.CreatedAt,
		UpdatedAt:    s.UpdatedAt,
	}

	if s.Branch.ID != (uuid.UUID{}) {
		response.Branch = branchDto.ToBranchResponse(&s.Branch)
	}

	return response
}

