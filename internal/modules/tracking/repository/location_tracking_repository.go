package repository

import (
	"context"
	"service/internal/shared/database"
	"service/internal/shared/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LocationTrackingRepository handles location tracking data operations
type LocationTrackingRepository struct {
	db *gorm.DB
}

// NewLocationTrackingRepository creates a new location tracking repository
func NewLocationTrackingRepository() *LocationTrackingRepository {
	return &LocationTrackingRepository{
		db: database.DB,
	}
}

// CreateLocationHistory creates a new location tracking record
func (r *LocationTrackingRepository) CreateLocationHistory(ctx context.Context, tracking *model.LocationTracking) error {
	return r.db.WithContext(ctx).Create(tracking).Error
}

// GetLocationHistory retrieves location history for an order
func (r *LocationTrackingRepository) GetLocationHistory(ctx context.Context, orderID uuid.UUID, limit int) ([]model.LocationTracking, error) {
	var history []model.LocationTracking
	query := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("timestamp DESC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	err := query.Find(&history).Error
	return history, err
}

// GetLatestLocation retrieves the latest location for an order
func (r *LocationTrackingRepository) GetLatestLocation(ctx context.Context, orderID uuid.UUID) (*model.LocationTracking, error) {
	var tracking model.LocationTracking
	err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("timestamp DESC").
		First(&tracking).Error
	if err != nil {
		return nil, err
	}
	return &tracking, nil
}

// CurrentLocationRepository handles current location data operations
type CurrentLocationRepository struct {
	db *gorm.DB
}

// NewCurrentLocationRepository creates a new current location repository
func NewCurrentLocationRepository() *CurrentLocationRepository {
	return &CurrentLocationRepository{
		db: database.DB,
	}
}

// UpsertCurrentLocation creates or updates current location for an order
func (r *CurrentLocationRepository) UpsertCurrentLocation(ctx context.Context, location *model.CurrentLocation) error {
	// Check if location exists
	var existing model.CurrentLocation
	err := r.db.WithContext(ctx).
		Where("order_id = ?", location.OrderID).
		First(&existing).Error
	
	if err == gorm.ErrRecordNotFound {
		// Create new
		return r.db.WithContext(ctx).Create(location).Error
	} else if err != nil {
		return err
	}
	
	// Update existing
	location.ID = existing.ID
	location.CreatedAt = existing.CreatedAt
	location.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(location).Error
}

// GetCurrentLocation retrieves current location for an order
func (r *CurrentLocationRepository) GetCurrentLocation(ctx context.Context, orderID uuid.UUID) (*model.CurrentLocation, error) {
	var location model.CurrentLocation
	err := r.db.WithContext(ctx).
		Preload("Order").
		Preload("User").
		Where("order_id = ?", orderID).
		First(&location).Error
	if err != nil {
		return nil, err
	}
	return &location, nil
}

// GetCurrentLocationByUser retrieves current location for a user (courier/provider)
func (r *CurrentLocationRepository) GetCurrentLocationByUser(ctx context.Context, userID uuid.UUID) (*model.CurrentLocation, error) {
	var location model.CurrentLocation
	err := r.db.WithContext(ctx).
		Preload("Order").
		Preload("User").
		Where("user_id = ?", userID).
		Order("updated_at DESC").
		First(&location).Error
	if err != nil {
		return nil, err
	}
	return &location, nil
}

// DeleteCurrentLocation deletes current location for an order
func (r *CurrentLocationRepository) DeleteCurrentLocation(ctx context.Context, orderID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Delete(&model.CurrentLocation{}).Error
}

// GetActiveLocations retrieves all active current locations
func (r *CurrentLocationRepository) GetActiveLocations(ctx context.Context, limit int) ([]model.CurrentLocation, error) {
	var locations []model.CurrentLocation
	query := r.db.WithContext(ctx).
		Preload("Order").
		Preload("User").
		Where("updated_at > ?", time.Now().Add(-1*time.Hour)). // Only active in last hour
		Order("updated_at DESC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	err := query.Find(&locations).Error
	return locations, err
}

