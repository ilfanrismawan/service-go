package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	branchEntity "service-go/internal/modules/branches/entity"
)

// SparePartInventory represents spare parts inventory
type SparePartInventory struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	BranchID  uuid.UUID      `json:"branch_id" gorm:"type:uuid;not null"`
	Branch    branchEntity.Branch `json:"branch" gorm:"foreignKey:BranchID"`
	PartName  string         `json:"part_name" gorm:"not null"`
	PartCode  string         `json:"part_code" gorm:"not null;uniqueIndex"`
	Stock     int            `json:"stock" gorm:"default:0"`
	MinStock  int            `json:"min_stock" gorm:"default:5"` // Threshold untuk reorder
	Price     int64          `json:"price" gorm:"not null"`      // Price in Rupiah
	Supplier  string         `json:"supplier" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
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

