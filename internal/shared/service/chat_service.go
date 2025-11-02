package service

import (
	"context"
	orderDTO "service/internal/orders/dto"
	"service/internal/repository"
	"service/internal/shared/model"
	"time"

	"github.com/google/uuid"
)

// ChatService handles chat business logic
type ChatService struct {
	chatRepo *repository.ChatRepository
}

// NewChatService creates a new chat service
func NewChatService() *ChatService {
	return &ChatService{
		chatRepo: repository.NewChatRepository(),
	}
}

// SendMessage sends a chat message
func (s *ChatService) SendMessage(ctx context.Context, senderID uuid.UUID, req *model.ChatMessageRequest) (*model.ChatMessageResponse, error) {
	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		return nil, model.ErrOrderNotFound
	}

	receiverID, err := uuid.Parse(req.ReceiverID)
	if err != nil {
		return nil, model.ErrUserNotFound
	}

	message := &model.ChatMessage{
		OrderID:    orderID,
		SenderID:   senderID,
		ReceiverID: receiverID,
		Message:    req.Message,
		CreatedAt:  time.Now(),
	}

	if err := s.chatRepo.Create(ctx, message); err != nil {
		return nil, err
	}

	// Load relations for response - GetByID already loads relations
	loadedMsg, err := s.chatRepo.GetByID(ctx, message.ID)
	if err != nil {
		return nil, err
	}

	resp := loadedMsg.ToResponse()
	return &resp, nil
}

// GetChatMessages retrieves chat messages for an order
func (s *ChatService) GetChatMessages(ctx context.Context, orderID uuid.UUID, page, limit int) (*model.PaginatedResponse, error) {
	offset := (page - 1) * limit
	messages, total, err := s.chatRepo.ListByOrderID(ctx, orderID, offset, limit)
	if err != nil {
		return nil, err
	}

	var responses []model.ChatMessageResponse
	for _, msg := range messages {
		responses = append(responses, msg.ToResponse())
	}

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
		Message:    "Chat messages retrieved successfully",
		Timestamp:  model.GetCurrentTimestamp(),
	}, nil
}

// GetUserChats retrieves all chats for a user
func (s *ChatService) GetUserChats(ctx context.Context, userID uuid.UUID) ([]*orderDTO.ServiceOrderResponse, error) {
	// This would typically join orders and chat messages
	// For now, return empty slice
	return []*orderDTO.ServiceOrderResponse{}, nil
}

// MarkAsRead marks a single message as read
func (s *ChatService) MarkAsRead(ctx context.Context, messageID uuid.UUID) error {
	message, err := s.chatRepo.GetByID(ctx, messageID)
	if err != nil {
		return err
	}
	
	message.IsRead = true
	return s.chatRepo.Update(ctx, message)
}

// MarkOrderMessagesAsRead marks all messages in an order as read for a specific user
func (s *ChatService) MarkOrderMessagesAsRead(ctx context.Context, orderID, userID uuid.UUID) error {
	return s.chatRepo.MarkOrderMessagesAsRead(ctx, orderID, userID)
}

// GetUnreadCount gets the count of unread messages for a user
func (s *ChatService) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	return s.chatRepo.GetUnreadCount(ctx, userID)
}
