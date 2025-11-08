package repository

import (
	"context"
	"service/internal/shared/database"
	inventoryEntity "service-go/internal/modules/inventory/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SparePartInventoryRepository handles spare part inventory database operations
type SparePartInventoryRepository struct {
	db *gorm.DB
}

// NewSparePartInventoryRepository creates a new spare part inventory repository
func NewSparePartInventoryRepository() *SparePartInventoryRepository {
	return &SparePartInventoryRepository{
		db: database.DB,
	}
}

// Create creates a new spare part inventory entry
func (r *SparePartInventoryRepository) Create(ctx context.Context, sparePart *inventoryEntity.SparePartInventory) error {
	return r.db.WithContext(ctx).Create(sparePart).Error
}

// GetByID retrieves a spare part by ID
func (r *SparePartInventoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*inventoryEntity.SparePartInventory, error) {
	var sparePart inventoryEntity.SparePartInventory
	err := r.db.WithContext(ctx).
		Preload("Branch").
		First(&sparePart, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &sparePart, nil
}

// GetByPartCode retrieves a spare part by part code
func (r *SparePartInventoryRepository) GetByPartCode(ctx context.Context, partCode string, branchID uuid.UUID) (*inventoryEntity.SparePartInventory, error) {
	var sparePart inventoryEntity.SparePartInventory
	err := r.db.WithContext(ctx).
		Preload("Branch").
		Where("part_code = ? AND branch_id = ?", partCode, branchID).
		First(&sparePart).Error
	if err != nil {
		return nil, err
	}
	return &sparePart, nil
}

// List retrieves spare parts with filters
func (r *SparePartInventoryRepository) List(ctx context.Context, offset, limit int, filters *SparePartInventoryFilters) ([]inventoryEntity.SparePartInventory, int64, error) {
	var spareParts []inventoryEntity.SparePartInventory
	var total int64

	query := r.db.WithContext(ctx).Model(&inventoryEntity.SparePartInventory{}).Preload("Branch")

	// Apply filters
	if filters != nil {
		if filters.BranchID != nil {
			query = query.Where("branch_id = ?", *filters.BranchID)
		}
		if filters.PartCode != "" {
			query = query.Where("part_code = ?", filters.PartCode)
		}
		if filters.Supplier != "" {
			query = query.Where("supplier = ?", filters.Supplier)
		}
		if filters.LowStockOnly {
			query = query.Where("stock < min_stock")
		}
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&spareParts).Error

	return spareParts, total, err
}

// Update updates a spare part
func (r *SparePartInventoryRepository) Update(ctx context.Context, sparePart *inventoryEntity.SparePartInventory) error {
	return r.db.WithContext(ctx).Save(sparePart).Error
}

// Delete soft deletes a spare part
func (r *SparePartInventoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&inventoryEntity.SparePartInventory{}, "id = ?", id).Error
}

// SparePartInventoryFilters represents filters for spare part queries
type SparePartInventoryFilters struct {
	BranchID     *uuid.UUID
	PartCode     string
	Supplier     string
	LowStockOnly bool
}
