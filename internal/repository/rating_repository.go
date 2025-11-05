package repository

import (
	"context"
	"service/internal/shared/database"
	"service/internal/shared/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RatingRepository handles rating database operations
type RatingRepository struct {
	db *gorm.DB
}

// NewRatingRepository creates a new rating repository
func NewRatingRepository() *RatingRepository {
	return &RatingRepository{
		db: database.DB,
	}
}

// Create creates a new rating
func (r *RatingRepository) Create(ctx context.Context, rating *model.Rating) error {
	return r.db.WithContext(ctx).Create(rating).Error
}

// GetByID retrieves a rating by ID
func (r *RatingRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Rating, error) {
	var rating model.Rating
	err := r.db.WithContext(ctx).
		Preload("Order").
		Preload("Customer").
		Preload("Branch").
		Preload("Technician").
		First(&rating, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &rating, nil
}

// GetByOrderID retrieves a rating by order ID
func (r *RatingRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) (*model.Rating, error) {
	var rating model.Rating
	err := r.db.WithContext(ctx).
		Preload("Order").
		Preload("Customer").
		Preload("Branch").
		Preload("Technician").
		First(&rating, "order_id = ?", orderID).Error
	if err != nil {
		return nil, err
	}
	return &rating, nil
}

// List retrieves ratings with filters
func (r *RatingRepository) List(ctx context.Context, offset, limit int, filters *RatingFilters) ([]model.Rating, int64, error) {
	var ratings []model.Rating
	var total int64

	query := r.db.WithContext(ctx).
		Preload("Order").
		Preload("Customer").
		Preload("Branch").
		Preload("Technician")

	// Apply filters
	if filters != nil {
		if filters.CustomerID != nil {
			query = query.Where("customer_id = ?", *filters.CustomerID)
		}
		if filters.BranchID != nil {
			query = query.Where("branch_id = ?", *filters.BranchID)
		}
		if filters.TechnicianID != nil {
			query = query.Where("technician_id = ?", *filters.TechnicianID)
		}
		if filters.OrderID != nil {
			query = query.Where("order_id = ?", *filters.OrderID)
		}
		if filters.MinRating != nil {
			query = query.Where("rating >= ?", *filters.MinRating)
		}
		if filters.IsPublic != nil {
			query = query.Where("is_public = ?", *filters.IsPublic)
		}
	}

	// Count total
	if err := query.Model(&model.Rating{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&ratings).Error

	return ratings, total, err
}

// Update updates a rating
func (r *RatingRepository) Update(ctx context.Context, rating *model.Rating) error {
	return r.db.WithContext(ctx).Save(rating).Error
}

// Delete soft deletes a rating
func (r *RatingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Rating{}, "id = ?", id).Error
}

// GetAverageRating calculates average rating
func (r *RatingRepository) GetAverageRating(ctx context.Context, filters *RatingFilters) (*model.AverageRating, error) {
	query := r.db.WithContext(ctx).Model(&model.Rating{})

	// Apply filters
	if filters != nil {
		if filters.BranchID != nil {
			query = query.Where("branch_id = ?", *filters.BranchID)
		}
		if filters.TechnicianID != nil {
			query = query.Where("technician_id = ?", *filters.TechnicianID)
		}
		if filters.IsPublic != nil {
			query = query.Where("is_public = ?", *filters.IsPublic)
		}
	}

	var result struct {
		AverageRating float64
		TotalRatings  int64
		Rating5       int64
		Rating4       int64
		Rating3       int64
		Rating2       int64
		Rating1       int64
	}

	err := query.
		Select("AVG(rating) as average_rating, COUNT(*) as total_ratings").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}

	// Count ratings by value
	baseQuery := query
	_ = baseQuery.Where("rating = ?", 5).Count(&result.Rating5).Error
	_ = baseQuery.Where("rating = ?", 4).Count(&result.Rating4).Error
	_ = baseQuery.Where("rating = ?", 3).Count(&result.Rating3).Error
	_ = baseQuery.Where("rating = ?", 2).Count(&result.Rating2).Error
	_ = baseQuery.Where("rating = ?", 1).Count(&result.Rating1).Error

	return &model.AverageRating{
		AverageRating: result.AverageRating,
		TotalRatings:  result.TotalRatings,
		Rating5:       result.Rating5,
		Rating4:       result.Rating4,
		Rating3:       result.Rating3,
		Rating2:       result.Rating2,
		Rating1:       result.Rating1,
	}, nil
}

// RatingFilters represents filters for rating queries
type RatingFilters struct {
	CustomerID   *uuid.UUID
	BranchID     *uuid.UUID
	TechnicianID *uuid.UUID
	OrderID      *uuid.UUID
	MinRating    *int
	IsPublic     *bool
}
