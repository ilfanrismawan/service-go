package service

import (
	"context"
	"service/internal/cache"
	"service/internal/branches/dto"
	"service/internal/branches/repository"
	"service/internal/shared/model"

	"github.com/google/uuid"
)

// BranchService handles branch business logic
type BranchService struct {
	branchRepo *repository.BranchRepository
	cache      *cache.CacheService
}

// NewBranchService creates a new branch service
func NewBranchService() *BranchService {
	return &BranchService{
		branchRepo: repository.NewBranchRepository(),
		cache:      cache.NewCacheService(),
	}
}

// CreateBranch creates a new branch
func (s *BranchService) CreateBranch(ctx context.Context, req *dto.BranchRequest) (*dto.BranchResponse, error) {
	// Create branch entity
	branch := &dto.Branch{
		Name:      req.Name,
		Address:   req.Address,
		City:      req.City,
		Province:  req.Province,
		Phone:     req.Phone,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		IsActive:  true,
	}

	// Save to database
	if err := s.branchRepo.Create(ctx, branch); err != nil {
		return nil, err
	}

	// Return response
	response := branch.ToResponse()
	return &response, nil
}

// GetBranch retrieves a branch by ID (with cache)
func (s *BranchService) GetBranch(ctx context.Context, id uuid.UUID) (*dto.BranchResponse, error) {
	// Try cache first
	if cachedBranch, err := s.cache.GetBranch(ctx, id); err == nil {
		response := cachedBranch.ToResponse()
		return &response, nil
	}

	// Cache miss, get from database
	branch, err := s.branchRepo.GetByID(ctx, id)
	if err != nil {
		return nil, model.ErrBranchNotFound
	}

	// Cache the result
	_ = s.cache.SetBranch(ctx, branch)

	response := branch.ToResponse()
	return &response, nil
}

// UpdateBranch updates a branch
func (s *BranchService) UpdateBranch(ctx context.Context, id uuid.UUID, req *dto.BranchRequest) (*dto.BranchResponse, error) {
	// Get existing branch
	branch, err := s.branchRepo.GetByID(ctx, id)
	if err != nil {
		return nil, model.ErrBranchNotFound
	}

	// Update fields
	branch.Name = req.Name
	branch.Address = req.Address
	branch.City = req.City
	branch.Province = req.Province
	branch.Phone = req.Phone
	branch.Latitude = req.Latitude
	branch.Longitude = req.Longitude

	// Save changes
	if err := s.branchRepo.Update(ctx, branch); err != nil {
		return nil, err
	}

	// Invalidate cache
	_ = s.cache.InvalidateBranch(ctx, id)
	_ = s.cache.InvalidateBranchList(ctx)

	response := branch.ToResponse()
	return &response, nil
}

// DeleteBranch soft deletes a branch
func (s *BranchService) DeleteBranch(ctx context.Context, id uuid.UUID) error {
	// Check if branch exists
	_, err := s.branchRepo.GetByID(ctx, id)
	if err != nil {
		return model.ErrBranchNotFound
	}

	// Soft delete
	return s.branchRepo.Delete(ctx, id)
}

// ListBranches retrieves branches with pagination and filters
func (s *BranchService) ListBranches(ctx context.Context, page, limit int, city, province *string) (*model.PaginatedResponse, error) {
	offset := (page - 1) * limit

	branches, total, err := s.branchRepo.List(ctx, offset, limit, city, province)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	var responses []dto.BranchResponse
	for _, branch := range branches {
		responses = append(responses, branch.ToResponse())
	}

	// Calculate pagination
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	pagination := model.PaginationResponse{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return &model.PaginatedResponse{
		Status:     "success",
		Data:       responses,
		Pagination: pagination,
		Message:    "Branches retrieved successfully",
		Timestamp:  model.GetCurrentTimestamp(),
	}, nil
}

// GetNearbyBranches retrieves branches within a certain radius
func (s *BranchService) GetNearbyBranches(ctx context.Context, latitude, longitude, radiusKm float64) ([]dto.BranchDistance, error) {
	branches, err := s.branchRepo.GetNearbyBranches(ctx, latitude, longitude, radiusKm)
	if err != nil {
		return nil, err
	}

	var results []dto.BranchDistance
	for _, branch := range branches {
		distance := model.CalculateDistance(latitude, longitude, branch.Latitude, branch.Longitude)
		results = append(results, dto.BranchDistance{
			Branch:   branch.ToResponse(),
			Distance: distance,
		})
	}

	return results, nil
}

// GetActiveBranches retrieves all active branches
func (s *BranchService) GetActiveBranches(ctx context.Context) ([]dto.BranchResponse, error) {
	branches, err := s.branchRepo.GetActiveBranches(ctx)
	if err != nil {
		return nil, err
	}

	var responses []dto.BranchResponse
	for _, branch := range branches {
		responses = append(responses, branch.ToResponse())
	}

	return responses, nil
}

// GetBranchesByCity retrieves branches by city
func (s *BranchService) GetBranchesByCity(ctx context.Context, city string) ([]dto.BranchResponse, error) {
	branches, err := s.branchRepo.GetByCity(ctx, city)
	if err != nil {
		return nil, err
	}

	var responses []dto.BranchResponse
	for _, branch := range branches {
		responses = append(responses, branch.ToResponse())
	}

	return responses, nil
}

// GetBranchesByProvince retrieves branches by province
func (s *BranchService) GetBranchesByProvince(ctx context.Context, province string) ([]dto.BranchResponse, error) {
	branches, err := s.branchRepo.GetByProvince(ctx, province)
	if err != nil {
		return nil, err
	}

	var responses []dto.BranchResponse
	for _, branch := range branches {
		responses = append(responses, branch.ToResponse())
	}

	return responses, nil
}

// GetBranches retrieves branches with pagination
func (s *BranchService) GetBranches(ctx context.Context, page, limit int) ([]dto.BranchResponse, int64, error) {
	branches, total, err := s.branchRepo.GetBranches(ctx, page, limit)
	if err != nil {
		return nil, 0, err
	}

	var responses []dto.BranchResponse
	for _, branch := range branches {
		responses = append(responses, branch.ToResponse())
	}

	return responses, total, nil
}
