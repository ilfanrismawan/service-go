package core

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SparePartInventory represents spare parts inventory
type SparePartInventory struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	BranchID    uuid.UUID `json:"branch_id" gorm:"type:uuid;not null"`
	Branch      Branch    `json:"branch" gorm:"foreignKey:BranchID"`
	PartName    string    `json:"part_name" gorm:"not null"`
	PartCode    string    `json:"part_code" gorm:"not null;uniqueIndex"`
	Stock       int       `json:"stock" gorm:"default:0"`
	MinStock    int       `json:"min_stock" gorm:"default:5"` // Threshold untuk reorder
	Price       int64     `json:"price" gorm:"not null"`      // Price in Rupiah
	Supplier    string    `json:"supplier" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName returns the table name for SparePartInventory
func (SparePartInventory) TableName() string {
	return "spare_part_inventory"
}

// IsLowStock checks if stock is below minimum threshold
func (s *SparePartInventory) IsLowStock() bool {
	return s.Stock < s.MinStock
}

// NeedsReorder checks if reorder is needed
func (s *SparePartInventory) NeedsReorder() bool {
	return s.IsLowStock()
}

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
	ID          uuid.UUID          `json:"id"`
	BranchID    uuid.UUID          `json:"branch_id"`
	Branch      BranchResponse     `json:"branch"`
	PartName    string             `json:"part_name"`
	PartCode    string             `json:"part_code"`
	Stock       int                `json:"stock"`
	MinStock    int                `json:"min_stock"`
	Price       int64              `json:"price"`
	Supplier    string             `json:"supplier"`
	IsLowStock  bool               `json:"is_low_stock"`
	NeedsReorder bool              `json:"needs_reorder"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// ToResponse converts SparePartInventory to SparePartInventoryResponse
func (s *SparePartInventory) ToResponse() SparePartInventoryResponse {
	response := SparePartInventoryResponse{
		ID:          s.ID,
		BranchID:    s.BranchID,
		PartName:    s.PartName,
		PartCode:    s.PartCode,
		Stock:       s.Stock,
		MinStock:    s.MinStock,
		Price:       s.Price,
		Supplier:    s.Supplier,
		IsLowStock:  s.IsLowStock(),
		NeedsReorder: s.NeedsReorder(),
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}

	if s.Branch.ID != (uuid.UUID{}) {
		response.Branch = s.Branch.ToResponse()
	}

	return response
}

