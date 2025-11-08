package repository

import (
	"context"
	"service/internal/shared/database"
	notificationEntity "service-go/internal/modules/notification/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// NotificationRepository handles notification data operations
type NotificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository creates a new notification repository
func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{
		db: database.DB,
	}
}

// Create creates a new notification
func (r *NotificationRepository) Create(ctx context.Context, notification *notificationEntity.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

// GetByID retrieves a notification by ID
func (r *NotificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*notificationEntity.Notification, error) {
	var notification notificationEntity.Notification
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Order").
		First(&notification, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

// Update updates a notification
func (r *NotificationRepository) Update(ctx context.Context, notification *notificationEntity.Notification) error {
	return r.db.WithContext(ctx).Save(notification).Error
}

// Delete soft deletes a notification
func (r *NotificationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&notificationEntity.Notification{}, "id = ?", id).Error
}

// ListByUserID retrieves notifications for a specific user with pagination
func (r *NotificationRepository) ListByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*notificationEntity.Notification, int64, error) {
	var notifications []*notificationEntity.Notification
	var total int64

	query := r.db.WithContext(ctx).Model(&notificationEntity.Notification{}).Where("user_id = ?", userID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get notifications with pagination
	err := query.
		Preload("User").
		Preload("Order").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, total, err
}

// List retrieves notifications with pagination and filters
func (r *NotificationRepository) List(ctx context.Context, offset, limit int, filters *NotificationFilters) ([]*notificationEntity.Notification, int64, error) {
	var notifications []*notificationEntity.Notification
	var total int64

	query := r.db.WithContext(ctx).Model(&notificationEntity.Notification{})

	if filters != nil {
		if filters.UserID != nil {
			query = query.Where("user_id = ?", *filters.UserID)
		}
		if filters.OrderID != nil {
			query = query.Where("order_id = ?", *filters.OrderID)
		}
		if filters.Type != nil {
			query = query.Where("type = ?", *filters.Type)
		}
		if filters.Status != nil {
			query = query.Where("status = ?", *filters.Status)
		}
		if filters.DateFrom != nil {
			query = query.Where("created_at >= ?", *filters.DateFrom)
		}
		if filters.DateTo != nil {
			query = query.Where("created_at <= ?", *filters.DateTo)
		}
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get notifications with pagination
	err := query.
		Preload("User").
		Preload("Order").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, total, err
}

// GetByUserID retrieves notifications by user ID
func (r *NotificationRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*notificationEntity.Notification, error) {
	var notifications []*notificationEntity.Notification
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Order").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, err
}

// GetByOrderID retrieves notifications by order ID
func (r *NotificationRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*notificationEntity.Notification, error) {
	var notifications []*notificationEntity.Notification
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Order").
		Where("order_id = ?", orderID).
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, err
}

// GetByStatus retrieves notifications by status
func (r *NotificationRepository) GetByStatus(ctx context.Context, status notificationEntity.NotificationStatus) ([]*notificationEntity.Notification, error) {
	var notifications []*notificationEntity.Notification
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Order").
		Where("status = ?", status).
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, err
}

// GetByType retrieves notifications by type
func (r *NotificationRepository) GetByType(ctx context.Context, notificationType notificationEntity.NotificationType) ([]*notificationEntity.Notification, error) {
	var notifications []*notificationEntity.Notification
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Order").
		Where("type = ?", notificationType).
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, err
}

// MarkAsSent marks a notification as sent
func (r *NotificationRepository) MarkAsSent(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&notificationEntity.Notification{}).
		Where("id = ?", id).
		Update("status", notificationEntity.NotificationStatusSent).Error
}

// MarkAsFailed marks a notification as failed
func (r *NotificationRepository) MarkAsFailed(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&notificationEntity.Notification{}).
		Where("id = ?", id).
		Update("status", notificationEntity.NotificationStatusFailed).Error
}

// NotificationFilters represents filters for notification queries
type NotificationFilters struct {
	UserID   *uuid.UUID
	OrderID  *uuid.UUID
	Type     *notificationEntity.NotificationType
	Status   *notificationEntity.NotificationStatus
	DateFrom *string
	DateTo   *string
}
