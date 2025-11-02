package service

import (
	"context"
	"fmt"
	"service/internal/core"
	"service/internal/repository"

	"github.com/google/uuid"
)

// RatingService handles rating business logic
type RatingService struct {
	ratingRepo *repository.RatingRepository
	orderRepo  *repository.ServiceOrderRepository
	userRepo   *repository.UserRepository
}

// NewRatingService creates a new rating service
func NewRatingService() *RatingService {
	return &RatingService{
		ratingRepo: repository.NewRatingRepository(),
		orderRepo:  repository.NewServiceOrderRepository(),
		userRepo:   repository.NewUserRepository(),
	}
}

// CreateRating creates a new rating for an order
func (s *RatingService) CreateRating(ctx context.Context, customerID uuid.UUID, req *core.RatingRequest) (*core.RatingResponse, error) {
	// Validate rating
	if req.Rating < 1 || req.Rating > 5 {
		return nil, fmt.Errorf("rating must be between 1 and 5")
	}

	// Parse order ID
	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID: %w", err)
	}

	// Check if order exists and belongs to customer
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	if order.UserID != customerID {
		return nil, fmt.Errorf("order does not belong to customer")
	}

	// Check if order is completed
	if order.Status != core.StatusCompleted {
		return nil, fmt.Errorf("can only rate completed orders")
	}

	// Check if rating already exists
	existingRating, _ := s.ratingRepo.GetByOrderID(ctx, orderID)
	if existingRating != nil {
		return nil, fmt.Errorf("rating already exists for this order")
	}

	// Create rating
	rating := &core.Rating{
		OrderID:     orderID,
		CustomerID:  customerID,
		Rating:      req.Rating,
		Review:      req.Review,
		IsPublic:    req.IsPublic,
	}

	// Set branch and technician if available
	if order.BranchID != nil {
		rating.BranchID = order.BranchID
	}
	if order.TechnicianID != nil {
		rating.TechnicianID = order.TechnicianID
	}

	// Save to database
	if err := s.ratingRepo.Create(ctx, rating); err != nil {
		return nil, fmt.Errorf("failed to create rating: %w", err)
	}

	// Get full rating with relations
	fullRating, err := s.ratingRepo.GetByID(ctx, rating.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve rating: %w", err)
	}

	response := fullRating.ToResponse()
	return &response, nil
}

// GetRating retrieves a rating by ID
func (s *RatingService) GetRating(ctx context.Context, ratingID uuid.UUID) (*core.RatingResponse, error) {
	rating, err := s.ratingRepo.GetByID(ctx, ratingID)
	if err != nil {
		return nil, fmt.Errorf("rating not found: %w", err)
	}

	response := rating.ToResponse()
	return &response, nil
}

// ListRatings retrieves ratings with filters
func (s *RatingService) ListRatings(ctx context.Context, page, limit int, filters *repository.RatingFilters) (*core.PaginatedResponse, error) {
	offset := (page - 1) * limit

	ratings, total, err := s.ratingRepo.List(ctx, offset, limit, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list ratings: %w", err)
	}

	// Convert to response format
	var responses []core.RatingResponse
	for _, rating := range ratings {
		responses = append(responses, rating.ToResponse())
	}

	// Calculate pagination
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
		Message:    "Ratings retrieved successfully",
		Timestamp:  core.GetCurrentTimestamp(),
	}, nil
}

// UpdateRating updates a rating
func (s *RatingService) UpdateRating(ctx context.Context, ratingID, customerID uuid.UUID, req *core.RatingRequest) (*core.RatingResponse, error) {
	// Get existing rating
	rating, err := s.ratingRepo.GetByID(ctx, ratingID)
	if err != nil {
		return nil, fmt.Errorf("rating not found: %w", err)
	}

	// Check ownership
	if rating.CustomerID != customerID {
		return nil, fmt.Errorf("unauthorized: rating does not belong to customer")
	}

	// Validate rating
	if req.Rating < 1 || req.Rating > 5 {
		return nil, fmt.Errorf("rating must be between 1 and 5")
	}

	// Update fields
	rating.Rating = req.Rating
	rating.Review = req.Review
	if req.IsPublic {
		rating.IsPublic = req.IsPublic
	}

	// Save changes
	if err := s.ratingRepo.Update(ctx, rating); err != nil {
		return nil, fmt.Errorf("failed to update rating: %w", err)
	}

	// Get updated rating
	updatedRating, err := s.ratingRepo.GetByID(ctx, ratingID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated rating: %w", err)
	}

	response := updatedRating.ToResponse()
	return &response, nil
}

// DeleteRating deletes a rating
func (s *RatingService) DeleteRating(ctx context.Context, ratingID, customerID uuid.UUID) error {
	// Get existing rating
	rating, err := s.ratingRepo.GetByID(ctx, ratingID)
	if err != nil {
		return fmt.Errorf("rating not found: %w", err)
	}

	// Check ownership
	if rating.CustomerID != customerID {
		return fmt.Errorf("unauthorized: rating does not belong to customer")
	}

	// Soft delete
	return s.ratingRepo.Delete(ctx, ratingID)
}

// GetAverageRating retrieves average rating statistics
func (s *RatingService) GetAverageRating(ctx context.Context, branchID, technicianID *uuid.UUID) (*core.AverageRating, error) {
	filters := &repository.RatingFilters{
		IsPublic: boolPtr(true),
	}

	if branchID != nil {
		filters.BranchID = branchID
	}
	if technicianID != nil {
		filters.TechnicianID = technicianID
	}

	return s.ratingRepo.GetAverageRating(ctx, filters)
}

// GetBranchRatings retrieves ratings for a specific branch
func (s *RatingService) GetBranchRatings(ctx context.Context, branchID uuid.UUID, page, limit int) (*core.PaginatedResponse, error) {
	filters := &repository.RatingFilters{
		BranchID: &branchID,
		IsPublic: boolPtr(true),
	}

	return s.ListRatings(ctx, page, limit, filters)
}

// GetTechnicianRatings retrieves ratings for a specific technician
func (s *RatingService) GetTechnicianRatings(ctx context.Context, technicianID uuid.UUID, page, limit int) (*core.PaginatedResponse, error) {
	filters := &repository.RatingFilters{
		TechnicianID: &technicianID,
		IsPublic:     boolPtr(true),
	}

	return s.ListRatings(ctx, page, limit, filters)
}

// Helper function
func boolPtr(b bool) *bool {
	return &b
}

