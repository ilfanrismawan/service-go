package repository

import (
	"context"
	"service/internal/shared/database"
	"service/internal/shared/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuditRepository handles audit trail database operations
type AuditRepository struct {
	db *gorm.DB
}

// NewAuditRepository creates a new audit repository
func NewAuditRepository() *AuditRepository {
	return &AuditRepository{
		db: database.DB,
	}
}

// Create creates a new audit trail entry
func (r *AuditRepository) Create(ctx context.Context, audit *model.AuditTrail) error {
	return r.db.WithContext(ctx).Create(audit).Error
}

// GetByID retrieves an audit trail entry by ID
func (r *AuditRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.AuditTrail, error) {
	var audit model.AuditTrail
	err := r.db.WithContext(ctx).
		Preload("User").
		First(&audit, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &audit, nil
}

// List retrieves audit trails with filters
func (r *AuditRepository) List(ctx context.Context, offset, limit int, filters *AuditFilters) ([]model.AuditTrail, int64, error) {
	var audits []model.AuditTrail
	var total int64

	query := r.db.WithContext(ctx).Model(&model.AuditTrail{}).Preload("User")

	// Apply filters
	if filters != nil {
		if filters.UserID != nil {
			query = query.Where("user_id = ?", *filters.UserID)
		}
		if filters.Action != "" {
			query = query.Where("action = ?", filters.Action)
		}
		if filters.Resource != "" {
			query = query.Where("resource = ?", filters.Resource)
		}
		if filters.ResourceID != nil {
			query = query.Where("resource_id = ?", *filters.ResourceID)
		}
		if filters.Status != "" {
			query = query.Where("status = ?", filters.Status)
		}
		if filters.DateFrom != nil {
			query = query.Where("created_at >= ?", *filters.DateFrom)
		}
		if filters.DateTo != nil {
			query = query.Where("created_at <= ?", *filters.DateTo)
		}
		if filters.IPAddress != "" {
			query = query.Where("ip_address = ?", filters.IPAddress)
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
		Find(&audits).Error

	return audits, total, err
}

// GetByResource retrieves audit trails for a specific resource
func (r *AuditRepository) GetByResource(ctx context.Context, resource string, resourceID uuid.UUID, limit int) ([]model.AuditTrail, error) {
	var audits []model.AuditTrail
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("resource = ? AND resource_id = ?", resource, resourceID).
		Order("created_at DESC").
		Limit(limit).
		Find(&audits).Error
	return audits, err
}

// GetByUser retrieves audit trails for a specific user
func (r *AuditRepository) GetByUser(ctx context.Context, userID uuid.UUID, limit int) ([]model.AuditTrail, error) {
	var audits []model.AuditTrail
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&audits).Error
	return audits, err
}

// AuditFilters represents filters for audit trail queries
type AuditFilters struct {
	UserID     *uuid.UUID
	Action     model.AuditAction
	Resource   string
	ResourceID *uuid.UUID
	Status     string
	DateFrom   *time.Time
	DateTo     *time.Time
	IPAddress  string
}
