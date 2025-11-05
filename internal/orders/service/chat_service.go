package service

import (
	"context"
	"service/internal/orders/repository"
	"service/internal/shared/model"

	"github.com/google/uuid"
)

// ChatService handles chat business logic
type ChatService struct {
	chatRepo  *repository.ChatRepository
	userRepo  *repository.UserRepository
	orderRepo *repository.ServiceOrderRepository
}

// NewChatService creates a new chat service
func NewChatService() *ChatService {
	return &ChatService{
		chatRepo:  repository.NewChatRepository(),
		userRepo:  repository.NewUserRepository(),
		orderRepo: repository.NewServiceOrderRepository(),
	}
}

// SendMessage sends a chat message
func (s *ChatService) SendMessage(ctx context.Context, senderID uuid.UUID, req *model.ChatMessageRequest) (*model.ChatMessageResponse, error) {
	// Validate order exists
	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		return nil, err
	}

	_, err = s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, model.ErrOrderNotFound
	}

	// Validate receiver exists
	receiverID, err := uuid.Parse(req.ReceiverID)
	if err != nil {
		return nil, err
	}

	_, err = s.userRepo.GetByID(ctx, receiverID)
	if err != nil {
		return nil, model.ErrUserNotFound
	}

	// Create chat message entity
	message := &model.ChatMessage{
		OrderID:    orderID,
		SenderID:   senderID,
		ReceiverID: receiverID,
		Message:    req.Message,
		IsRead:     false,
	}

	// Save to database
	if err := s.chatRepo.Create(ctx, message); err != nil {
		return nil, err
	}

	response := message.ToResponse()
	return &response, nil
}

// GetChatMessages retrieves chat messages for an order
func (s *ChatService) GetChatMessages(ctx context.Context, orderID uuid.UUID, page, limit int) (*model.PaginatedResponse, error) {
	offset := (page - 1) * limit

	messages, total, err := s.chatRepo.ListByOrderID(ctx, orderID, offset, limit)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	var responses []model.ChatMessageResponse
	for _, message := range messages {
		responses = append(responses, message.ToResponse())
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
		Message:    "Chat messages retrieved successfully",
		Timestamp:  model.GetCurrentTimestamp(),
	}, nil
}

// GetUserChats retrieves all chat conversations for a user
func (s *ChatService) GetUserChats(ctx context.Context, userID uuid.UUID) ([]model.ChatMessageResponse, error) {
	messages, err := s.chatRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var responses []model.ChatMessageResponse
	for _, message := range messages {
		responses = append(responses, message.ToResponse())
	}

	return responses, nil
}

// MarkAsRead marks a chat message as read
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
