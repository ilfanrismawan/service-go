package service

import (
	"context"
	"fmt"
	"service/internal/core"
	"service/internal/repository"

	"github.com/google/uuid"
)

// InventoryService handles inventory business logic
type InventoryService struct {
	inventoryRepo *repository.SparePartInventoryRepository
}

// NewInventoryService creates a new inventory service
func NewInventoryService() *InventoryService {
	return &InventoryService{
		inventoryRepo: repository.NewSparePartInventoryRepository(),
	}
}

// CreateSparePart creates a new spare part inventory entry
func (s *InventoryService) CreateSparePart(ctx context.Context, req *core.SparePartInventoryRequest) (*core.SparePartInventoryResponse, error) {
	// Parse branch ID
	branchID, err := uuid.Parse(req.BranchID)
	if err != nil {
		return nil, fmt.Errorf("invalid branch ID: %w", err)
	}

	sparePart := &core.SparePartInventory{
		BranchID: branchID,
		PartName: req.PartName,
		PartCode: req.PartCode,
		Stock:    req.Stock,
		MinStock: req.MinStock,
		Price:    req.Price,
		Supplier: req.Supplier,
	}

	if err := s.inventoryRepo.Create(ctx, sparePart); err != nil {
		return nil, fmt.Errorf("failed to create spare part: %w", err)
	}

	response := sparePart.ToResponse()
	return &response, nil
}

// GetSparePart retrieves a spare part by ID
func (s *InventoryService) GetSparePart(ctx context.Context, id uuid.UUID) (*core.SparePartInventoryResponse, error) {
	sparePart, err := s.inventoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("spare part not found: %w", err)
	}

	response := sparePart.ToResponse()
	return &response, nil
}

// ListSpareParts retrieves spare parts with filters
func (s *InventoryService) ListSpareParts(ctx context.Context, page, limit int, branchID *uuid.UUID, lowStockOnly bool) (*core.PaginatedResponse, error) {
	offset := (page - 1) * limit

	filters := &repository.SparePartInventoryFilters{
		BranchID:    branchID,
		LowStockOnly: lowStockOnly,
	}

	spareParts, total, err := s.inventoryRepo.List(ctx, offset, limit, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list spare parts: %w", err)
	}

	var responses []core.SparePartInventoryResponse
	for _, part := range spareParts {
		responses = append(responses, part.ToResponse())
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	pagination := core.PaginationResponse{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return &core.PaginatedResponse{
		Status:     "success",
		Data:       responses,
		Pagination: pagination,
		Message:    "Spare parts retrieved successfully",
		Timestamp:  core.GetCurrentTimestamp(),
	}, nil
}

// UpdateSparePart updates a spare part
func (s *InventoryService) UpdateSparePart(ctx context.Context, id uuid.UUID, req *core.SparePartInventoryRequest) (*core.SparePartInventoryResponse, error) {
	sparePart, err := s.inventoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("spare part not found: %w", err)
	}

	// Update fields
	sparePart.PartName = req.PartName
	sparePart.PartCode = req.PartCode
	sparePart.Stock = req.Stock
	sparePart.MinStock = req.MinStock
	sparePart.Price = req.Price
	sparePart.Supplier = req.Supplier

	if err := s.inventoryRepo.Update(ctx, sparePart); err != nil {
		return nil, fmt.Errorf("failed to update spare part: %w", err)
	}

	response := sparePart.ToResponse()
	return &response, nil
}

// UpdateStock updates stock for a spare part
func (s *InventoryService) UpdateStock(ctx context.Context, id uuid.UUID, quantity int, operation string) (*core.SparePartInventoryResponse, error) {
	sparePart, err := s.inventoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("spare part not found: %w", err)
	}

	// Update stock based on operation
	switch operation {
	case "add", "increase":
		sparePart.Stock += quantity
	case "subtract", "decrease":
		sparePart.Stock -= quantity
		if sparePart.Stock < 0 {
			sparePart.Stock = 0
		}
	case "set":
		sparePart.Stock = quantity
		if sparePart.Stock < 0 {
			sparePart.Stock = 0
		}
	default:
		return nil, fmt.Errorf("invalid operation: %s", operation)
	}

	if err := s.inventoryRepo.Update(ctx, sparePart); err != nil {
		return nil, fmt.Errorf("failed to update stock: %w", err)
	}

	response := sparePart.ToResponse()
	return &response, nil
}

// GetLowStockItems retrieves all items with stock below minimum threshold
func (s *InventoryService) GetLowStockItems(ctx context.Context, branchID *uuid.UUID) ([]core.SparePartInventoryResponse, error) {
	filters := &repository.SparePartInventoryFilters{
		BranchID:     branchID,
		LowStockOnly: true,
	}

	spareParts, _, err := s.inventoryRepo.List(ctx, 0, 1000, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get low stock items: %w", err)
	}

	var responses []core.SparePartInventoryResponse
	for _, part := range spareParts {
		responses = append(responses, part.ToResponse())
	}

	return responses, nil
}

// DeleteSparePart soft deletes a spare part
func (s *InventoryService) DeleteSparePart(ctx context.Context, id uuid.UUID) error {
	return s.inventoryRepo.Delete(ctx, id)
}

