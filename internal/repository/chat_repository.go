package repository

import (
	"context"
	"service/internal/shared/database"
	"service/internal/shared/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ChatRepository handles chat message data operations
type ChatRepository struct {
	db *gorm.DB
}

// NewChatRepository creates a new chat repository
func NewChatRepository() *ChatRepository {
	return &ChatRepository{
		db: database.DB,
	}
}

// Create creates a new chat message
func (r *ChatRepository) Create(ctx context.Context, message *model.ChatMessage) error {
	return r.db.WithContext(ctx).Create(message).Error
}

// GetByID retrieves a chat message by ID
func (r *ChatRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.ChatMessage, error) {
	var message model.ChatMessage
	err := r.db.WithContext(ctx).
		Preload("Order").
		Preload("Sender").
		Preload("Receiver").
		First(&message, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// Update updates a chat message
func (r *ChatRepository) Update(ctx context.Context, message *model.ChatMessage) error {
	return r.db.WithContext(ctx).Save(message).Error
}

// Delete soft deletes a chat message
func (r *ChatRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.ChatMessage{}, "id = ?", id).Error
}

// ListByOrderID retrieves chat messages for a specific order with pagination
func (r *ChatRepository) ListByOrderID(ctx context.Context, orderID uuid.UUID, offset, limit int) ([]*model.ChatMessage, int64, error) {
	var messages []*model.ChatMessage
	var total int64

	query := r.db.WithContext(ctx).Model(&model.ChatMessage{}).Where("order_id = ?", orderID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get messages with pagination
	err := query.
		Preload("Order").
		Preload("Sender").
		Preload("Receiver").
		Offset(offset).
		Limit(limit).
		Order("created_at ASC").
		Find(&messages).Error
	return messages, total, err
}

// GetByUserID retrieves chat messages by user ID
func (r *ChatRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.ChatMessage, error) {
	var messages []*model.ChatMessage
	err := r.db.WithContext(ctx).
		Preload("Order").
		Preload("Sender").
		Preload("Receiver").
		Where("sender_id = ? OR receiver_id = ?", userID, userID).
		Order("created_at DESC").
		Find(&messages).Error
	return messages, err
}

// GetByOrderID retrieves chat messages by order ID
func (r *ChatRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*model.ChatMessage, error) {
	var messages []*model.ChatMessage
	err := r.db.WithContext(ctx).
		Preload("Order").
		Preload("Sender").
		Preload("Receiver").
		Where("order_id = ?", orderID).
		Order("created_at ASC").
		Find(&messages).Error
	return messages, err
}

// MarkOrderMessagesAsRead marks all messages in an order as read for a specific user
func (r *ChatRepository) MarkOrderMessagesAsRead(ctx context.Context, orderID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&model.ChatMessage{}).
		Where("order_id = ? AND receiver_id = ? AND is_read = ?", orderID, userID, false).
		Update("is_read", true).Error
}

// GetUnreadCount gets the count of unread messages for a user
func (r *ChatRepository) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.ChatMessage{}).
		Where("receiver_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

// GetUnreadMessagesByOrder gets unread messages for a user in a specific order
func (r *ChatRepository) GetUnreadMessagesByOrder(ctx context.Context, orderID, userID uuid.UUID) ([]*model.ChatMessage, error) {
	var messages []*model.ChatMessage
	err := r.db.WithContext(ctx).
		Preload("Order").
		Preload("Sender").
		Preload("Receiver").
		Where("order_id = ? AND receiver_id = ? AND is_read = ?", orderID, userID, false).
		Order("created_at ASC").
		Find(&messages).Error
	return messages, err
}
