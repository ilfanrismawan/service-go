package service

import (
	"context"
<<<<<<< HEAD
	"service/internal/core"
	"service/internal/orders/repository"
=======
	"service/internal/branches/dto"
	"service/internal/branches/repository"
	"service/internal/cache"
<<<<<<<< HEAD:internal/service/branch_service.go
	"service/internal/core"
	"service/internal/orders/repository"
========
	"service/internal/shared/model"
>>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b:internal/branches/service/branch_service.go
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b

	"github.com/google/uuid"
)

// BranchService handles branch business logic
type BranchService struct {
	branchRepo *repository.BranchRepository
<<<<<<< HEAD
=======
	cache      *cache.CacheService
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
}

// NewBranchService creates a new branch service
func NewBranchService() *BranchService {
	return &BranchService{
		branchRepo: repository.NewBranchRepository(),
<<<<<<< HEAD
=======
		cache:      cache.NewCacheService(),
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	}
}

// CreateBranch creates a new branch
<<<<<<< HEAD
func (s *BranchService) CreateBranch(ctx context.Context, req *core.BranchRequest) (*core.BranchResponse, error) {
	// Create branch entity
	branch := &core.Branch{
=======
func (s *BranchService) CreateBranch(ctx context.Context, req *dto.BranchRequest) (*dto.BranchResponse, error) {
	// Create branch entity
	branch := &dto.Branch{
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
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

<<<<<<< HEAD
// GetBranch retrieves a branch by ID
func (s *BranchService) GetBranch(ctx context.Context, id uuid.UUID) (*core.BranchResponse, error) {
	branch, err := s.branchRepo.GetByID(ctx, id)
	if err != nil {
		return nil, core.ErrBranchNotFound
	}

=======
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

>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	response := branch.ToResponse()
	return &response, nil
}

// UpdateBranch updates a branch
<<<<<<< HEAD
func (s *BranchService) UpdateBranch(ctx context.Context, id uuid.UUID, req *core.BranchRequest) (*core.BranchResponse, error) {
	// Get existing branch
	branch, err := s.branchRepo.GetByID(ctx, id)
	if err != nil {
		return nil, core.ErrBranchNotFound
=======
func (s *BranchService) UpdateBranch(ctx context.Context, id uuid.UUID, req *dto.BranchRequest) (*dto.BranchResponse, error) {
	// Get existing branch
	branch, err := s.branchRepo.GetByID(ctx, id)
	if err != nil {
		return nil, model.ErrBranchNotFound
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
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

<<<<<<< HEAD
=======
	// Invalidate cache
	_ = s.cache.InvalidateBranch(ctx, id)
	_ = s.cache.InvalidateBranchList(ctx)

>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	response := branch.ToResponse()
	return &response, nil
}

// DeleteBranch soft deletes a branch
func (s *BranchService) DeleteBranch(ctx context.Context, id uuid.UUID) error {
	// Check if branch exists
	_, err := s.branchRepo.GetByID(ctx, id)
	if err != nil {
<<<<<<< HEAD
		return core.ErrBranchNotFound
=======
		return model.ErrBranchNotFound
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	}

	// Soft delete
	return s.branchRepo.Delete(ctx, id)
}

// ListBranches retrieves branches with pagination and filters
<<<<<<< HEAD
func (s *BranchService) ListBranches(ctx context.Context, page, limit int, city, province *string) (*core.PaginatedResponse, error) {
=======
func (s *BranchService) ListBranches(ctx context.Context, page, limit int, city, province *string) (*model.PaginatedResponse, error) {
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	offset := (page - 1) * limit

	branches, total, err := s.branchRepo.List(ctx, offset, limit, city, province)
	if err != nil {
		return nil, err
	}

	// Convert to response format
<<<<<<< HEAD
	var responses []core.BranchResponse
=======
	var responses []dto.BranchResponse
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	for _, branch := range branches {
		responses = append(responses, branch.ToResponse())
	}

	// Calculate pagination
	totalPages := int((total + int64(limit) - 1) / int64(limit))
<<<<<<< HEAD
	pagination := core.PaginationResponse{
=======
	pagination := model.PaginationResponse{
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

<<<<<<< HEAD
	return &core.PaginatedResponse{
=======
	return &model.PaginatedResponse{
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
		Status:     "success",
		Data:       responses,
		Pagination: pagination,
		Message:    "Branches retrieved successfully",
<<<<<<< HEAD
		Timestamp:  core.GetCurrentTimestamp(),
=======
		Timestamp:  model.GetCurrentTimestamp(),
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	}, nil
}

// GetNearbyBranches retrieves branches within a certain radius
<<<<<<< HEAD
func (s *BranchService) GetNearbyBranches(ctx context.Context, latitude, longitude, radiusKm float64) ([]core.BranchDistance, error) {
=======
func (s *BranchService) GetNearbyBranches(ctx context.Context, latitude, longitude, radiusKm float64) ([]dto.BranchDistance, error) {
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	branches, err := s.branchRepo.GetNearbyBranches(ctx, latitude, longitude, radiusKm)
	if err != nil {
		return nil, err
	}

<<<<<<< HEAD
	var results []core.BranchDistance
	for _, branch := range branches {
		distance := core.CalculateDistance(latitude, longitude, branch.Latitude, branch.Longitude)
		results = append(results, core.BranchDistance{
=======
	var results []dto.BranchDistance
	for _, branch := range branches {
		distance := model.CalculateDistance(latitude, longitude, branch.Latitude, branch.Longitude)
		results = append(results, dto.BranchDistance{
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
			Branch:   branch.ToResponse(),
			Distance: distance,
		})
	}

	return results, nil
}

// GetActiveBranches retrieves all active branches
<<<<<<< HEAD
func (s *BranchService) GetActiveBranches(ctx context.Context) ([]core.BranchResponse, error) {
=======
func (s *BranchService) GetActiveBranches(ctx context.Context) ([]dto.BranchResponse, error) {
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	branches, err := s.branchRepo.GetActiveBranches(ctx)
	if err != nil {
		return nil, err
	}

<<<<<<< HEAD
	var responses []core.BranchResponse
=======
	var responses []dto.BranchResponse
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	for _, branch := range branches {
		responses = append(responses, branch.ToResponse())
	}

	return responses, nil
}

// GetBranchesByCity retrieves branches by city
<<<<<<< HEAD
func (s *BranchService) GetBranchesByCity(ctx context.Context, city string) ([]core.BranchResponse, error) {
=======
func (s *BranchService) GetBranchesByCity(ctx context.Context, city string) ([]dto.BranchResponse, error) {
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	branches, err := s.branchRepo.GetByCity(ctx, city)
	if err != nil {
		return nil, err
	}

<<<<<<< HEAD
	var responses []core.BranchResponse
=======
	var responses []dto.BranchResponse
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	for _, branch := range branches {
		responses = append(responses, branch.ToResponse())
	}

	return responses, nil
}

// GetBranchesByProvince retrieves branches by province
<<<<<<< HEAD
func (s *BranchService) GetBranchesByProvince(ctx context.Context, province string) ([]core.BranchResponse, error) {
=======
func (s *BranchService) GetBranchesByProvince(ctx context.Context, province string) ([]dto.BranchResponse, error) {
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	branches, err := s.branchRepo.GetByProvince(ctx, province)
	if err != nil {
		return nil, err
	}

<<<<<<< HEAD
	var responses []core.BranchResponse
=======
	var responses []dto.BranchResponse
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	for _, branch := range branches {
		responses = append(responses, branch.ToResponse())
	}

	return responses, nil
}

// GetBranches retrieves branches with pagination
<<<<<<< HEAD
func (s *BranchService) GetBranches(ctx context.Context, page, limit int) ([]core.BranchResponse, int64, error) {
=======
func (s *BranchService) GetBranches(ctx context.Context, page, limit int) ([]dto.BranchResponse, int64, error) {
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	branches, total, err := s.branchRepo.GetBranches(ctx, page, limit)
	if err != nil {
		return nil, 0, err
	}

<<<<<<< HEAD
	var responses []core.BranchResponse
=======
	var responses []dto.BranchResponse
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	for _, branch := range branches {
		responses = append(responses, branch.ToResponse())
	}

	return responses, total, nil
}
