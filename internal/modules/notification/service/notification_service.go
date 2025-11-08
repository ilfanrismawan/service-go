package service

import (
	"context"
	notificationRepo "service/internal/modules/notification/repository"
	orderRepo "service/internal/modules/orders/repository"
	userRepo "service/internal/modules/users/repository"
	"service/internal/shared/model"
	notificationEntity "service-go/internal/modules/notification/entity"
	notificationDto "service-go/internal/modules/notification/dto"
	orderEntity "service-go/internal/modules/orders/entity"
	paymentEntity "service-go/internal/modules/payments/entity"

	"github.com/google/uuid"
)

// NotificationService handles notification business logic
type NotificationService struct {
	notificationRepo *notificationRepo.NotificationRepository
	userRepo         *userRepo.UserRepository
	orderRepo        *orderRepo.ServiceOrderRepository
}

// NewNotificationService creates a new notification service
func NewNotificationService() *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo.NewNotificationRepository(),
		userRepo:         userRepo.NewUserRepository(),
		orderRepo:        orderRepo.NewServiceOrderRepository(),
	}
}

// SendNotification sends a notification to a user
func (s *NotificationService) SendNotification(ctx context.Context, req *notificationDto.NotificationRequest) (*notificationDto.NotificationResponse, error) {
	// Validate user exists
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, err
	}

	_, err = s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, model.ErrUserNotFound
	}

	// Parse order ID if provided
	var orderID *uuid.UUID
	if req.OrderID != nil {
		parsedOrderID, err := uuid.Parse(*req.OrderID)
		if err != nil {
			return nil, err
		}
		orderID = &parsedOrderID
	}

	// Create notification entity
	notification := &notificationEntity.Notification{
		UserID:  userID,
		OrderID: orderID,
		Type:    req.Type,
		Title:   req.Title,
		Message: req.Message,
		Status:  notificationEntity.NotificationStatusPending,
	}

	// Save to database
	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return nil, err
	}

	// TODO: Send actual notification based on type
	// For now, just mark as sent
	notification.Status = notificationEntity.NotificationStatusSent
	s.notificationRepo.Update(ctx, notification)

	response := notificationDto.ToNotificationResponse(notification)
	return &response, nil
}

// SendOrderStatusNotification sends notification when order status changes
func (s *NotificationService) SendOrderStatusNotification(ctx context.Context, orderID uuid.UUID, status orderEntity.OrderStatus) error {
	// Get order details
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	// Create notification for customer
	notification := &notificationEntity.Notification{
		UserID:  order.CustomerID,
		OrderID: &orderID,
		Type:    notificationEntity.NotificationTypeEmail,
		Title:   "Order Status Update",
		Message: s.getOrderStatusMessage(order.OrderNumber, status),
		Status:  notificationEntity.NotificationStatusPending,
	}

	// Save notification
	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return err
	}

	// TODO: Send actual notification
	notification.Status = notificationEntity.NotificationStatusSent
	s.notificationRepo.Update(ctx, notification)

	return nil
}

// SendPaymentNotification sends notification for payment updates
func (s *NotificationService) SendPaymentNotification(ctx context.Context, orderID uuid.UUID, paymentStatus paymentEntity.PaymentStatus) error {
	// Get order details
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	// Create notification for customer
	notification := &notificationEntity.Notification{
		UserID:  order.CustomerID,
		OrderID: &orderID,
		Type:    notificationEntity.NotificationTypeEmail,
		Title:   "Payment Update",
		Message: s.getPaymentStatusMessage(order.OrderNumber, paymentStatus),
		Status:  notificationEntity.NotificationStatusPending,
	}

	// Save notification
	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return err
	}

	// TODO: Send actual notification
	notification.Status = notificationEntity.NotificationStatusSent
	s.notificationRepo.Update(ctx, notification)

	return nil
}

// GetNotifications retrieves notifications for a user
func (s *NotificationService) GetNotifications(ctx context.Context, userID uuid.UUID, page, limit int) (*model.PaginatedResponse, error) {
	offset := (page - 1) * limit

	notifications, total, err := s.notificationRepo.ListByUserID(ctx, userID, offset, limit)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	var responses []notificationDto.NotificationResponse
	for _, notification := range notifications {
		responses = append(responses, notificationDto.ToNotificationResponse(notification))
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
		Message:    "Notifications retrieved successfully",
		Timestamp:  model.GetCurrentTimestamp(),
	}, nil
}

// MarkAsRead marks a notification as read
func (s *NotificationService) MarkAsRead(ctx context.Context, notificationID uuid.UUID) error {
	_, err := s.notificationRepo.GetByID(ctx, notificationID)
	if err != nil {
		return err
	}

	// Mark as read (if notification has read status)
	// For now, just return success
	return nil
}

// getOrderStatusMessage returns appropriate message for order status
func (s *NotificationService) getOrderStatusMessage(orderNumber string, status model.OrderStatus) string {
	switch status {
	case model.StatusPendingPickup:
		return "Order " + orderNumber + " is pending pickup. Our courier will contact you soon."
	case model.StatusOnPickup:
		return "Order " + orderNumber + " is being picked up. Please prepare your device."
	case model.StatusInService:
		return "Order " + orderNumber + " is now in service. We'll keep you updated on the progress."
	case model.StatusReady:
		return "Order " + orderNumber + " is ready! We'll arrange delivery soon."
	case model.StatusDelivered:
		return "Order " + orderNumber + " has been delivered. Thank you for choosing our service!"
	case model.StatusCompleted:
		return "Order " + orderNumber + " has been completed successfully."
	case model.StatusCancelled:
		return "Order " + orderNumber + " has been cancelled."
	default:
		return "Order " + orderNumber + " status has been updated."
	}
}

// getPaymentStatusMessage returns appropriate message for payment status
func (s *NotificationService) getPaymentStatusMessage(orderNumber string, status model.PaymentStatus) string {
	switch status {
	case paymentEntity.PaymentStatusPaid:
		return "Payment for order " + orderNumber + " has been received successfully."
	case paymentEntity.PaymentStatusFailed:
		return "Payment for order " + orderNumber + " has failed. Please try again."
	case paymentEntity.PaymentStatusCancelled:
		return "Payment for order " + orderNumber + " has been cancelled."
	case paymentEntity.PaymentStatusRefunded:
		return "Payment for order " + orderNumber + " has been refunded."
	default:
		return "Payment status for order " + orderNumber + " has been updated."
	}
}
