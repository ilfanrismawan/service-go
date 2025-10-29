package notification

import (
	"context"
	"fmt"
	"log"
	"service/internal/core"
	"service/internal/repository"
	"time"

	"github.com/google/uuid"
)

// NotificationService handles notification business logic
type NotificationService struct {
	notificationRepo *repository.NotificationRepository
	userRepo         *repository.UserRepository
	orderRepo        *repository.ServiceOrderRepository
}

// NewNotificationService creates a new notification service
func NewNotificationService() *NotificationService {
	return &NotificationService{
		notificationRepo: repository.NewNotificationRepository(),
		userRepo:         repository.NewUserRepository(),
		orderRepo:        repository.NewServiceOrderRepository(),
	}
}

// SendEmailNotification sends email notification (mock implementation)
func (s *NotificationService) SendEmailNotification(ctx context.Context, email, subject, body string) error {
	// Mock email sending - in production, integrate with email service like SendGrid, AWS SES, etc.
	log.Printf("📧 Mock Email sent to %s: %s - %s", email, subject, body)

	// Simulate email sending delay
	time.Sleep(100 * time.Millisecond)

	return nil
}

// SendSMSNotification sends SMS notification (mock implementation)
func (s *NotificationService) SendSMSNotification(ctx context.Context, phone, message string) error {
	// Mock SMS sending - in production, integrate with SMS service like Twilio, AWS SNS, etc.
	log.Printf("📱 Mock SMS sent to %s: %s", phone, message)

	// Simulate SMS sending delay
	time.Sleep(200 * time.Millisecond)

	return nil
}

// SendWhatsAppNotification sends WhatsApp notification (mock implementation)
func (s *NotificationService) SendWhatsAppNotification(ctx context.Context, phone, message string) error {
	// Mock WhatsApp sending - in production, integrate with WhatsApp Business API
	log.Printf("💬 Mock WhatsApp sent to %s: %s", phone, message)

	// Simulate WhatsApp sending delay
	time.Sleep(150 * time.Millisecond)

	return nil
}

// SendPushNotification sends push notification (mock implementation)
func (s *NotificationService) SendPushNotification(ctx context.Context, userID uuid.UUID, title, body string) error {
	// Mock push notification - in production, integrate with Firebase FCM, AWS SNS, etc.
	log.Printf("🔔 Mock Push notification sent to user %s: %s - %s", userID.String(), title, body)

	// Simulate push notification delay
	time.Sleep(50 * time.Millisecond)

	return nil
}

// SendOrderStatusNotification sends notification when order status changes
func (s *NotificationService) SendOrderStatusNotification(ctx context.Context, orderID uuid.UUID, status core.OrderStatus) error {
	// Get order details
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	// Get user details
	user, err := s.userRepo.GetByID(ctx, order.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Create notification message based on status
	var message string
	var notificationType core.NotificationType

	switch status {
	case core.StatusPendingPickup:
		message = fmt.Sprintf("Order #%s is pending pickup. Our courier will contact you soon.", order.OrderNumber)
		notificationType = core.NotificationTypeOrderUpdate
	case core.StatusOnPickup:
		message = fmt.Sprintf("Order #%s is being picked up. Our courier is on the way.", order.OrderNumber)
		notificationType = core.NotificationTypeOrderUpdate
	case core.StatusInService:
		message = fmt.Sprintf("Order #%s is now in service. We'll update you on the progress.", order.OrderNumber)
		notificationType = core.NotificationTypeOrderUpdate
	case core.StatusReady:
		message = fmt.Sprintf("Order #%s is ready! You can pick it up or schedule delivery.", order.OrderNumber)
		notificationType = core.NotificationTypeOrderReady
	case core.StatusDelivered:
		message = fmt.Sprintf("Order #%s has been delivered successfully. Thank you for choosing our service!", order.OrderNumber)
		notificationType = core.NotificationTypeOrderDelivered
	case core.StatusCompleted:
		message = fmt.Sprintf("Order #%s has been completed. We hope you're satisfied with our service!", order.OrderNumber)
		notificationType = core.NotificationTypeOrderCompleted
	case core.StatusCancelled:
		message = fmt.Sprintf("Order #%s has been cancelled. Please contact us if you have any questions.", order.OrderNumber)
		notificationType = core.NotificationTypeOrderCancelled
	default:
		message = fmt.Sprintf("Order #%s status has been updated to %s.", order.OrderNumber, string(status))
		notificationType = core.NotificationTypeOrderUpdate
	}

	// Create notification record
	notification := &core.Notification{
		UserID:    user.ID,
		OrderID:   &orderID,
		Type:      notificationType,
		Message:   message,
		IsRead:    false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save notification to database
	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return fmt.Errorf("failed to save notification: %w", err)
	}

	// Send notifications through multiple channels
	go s.sendMultiChannelNotification(ctx, user, message, fmt.Sprintf("Order #%s Update", order.OrderNumber))

	return nil
}

// SendPaymentNotification sends notification for payment updates
func (s *NotificationService) SendPaymentNotification(ctx context.Context, orderID uuid.UUID, status core.PaymentStatus) error {
	// Get order details
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	// Get user details
	user, err := s.userRepo.GetByID(ctx, order.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Create notification message based on payment status
	var message string
	var notificationType core.NotificationType

	switch status {
	case core.PaymentStatusPending:
		message = fmt.Sprintf("Payment for Order #%s is pending. Please complete your payment.", order.OrderNumber)
		notificationType = core.NotificationTypePaymentPending
	case core.PaymentStatusPaid:
		message = fmt.Sprintf("Payment for Order #%s has been received successfully. Thank you!", order.OrderNumber)
		notificationType = core.NotificationTypePaymentReceived
	case core.PaymentStatusFailed:
		message = fmt.Sprintf("Payment for Order #%s failed. Please try again or contact support.", order.OrderNumber)
		notificationType = core.NotificationTypePaymentFailed
	case core.PaymentStatusCancelled:
		message = fmt.Sprintf("Payment for Order #%s has been cancelled.", order.OrderNumber)
		notificationType = core.NotificationTypePaymentCancelled
	case core.PaymentStatusRefunded:
		message = fmt.Sprintf("Payment for Order #%s has been refunded. The amount will be processed within 3-5 business days.", order.OrderNumber)
		notificationType = core.NotificationTypePaymentRefunded
	default:
		message = fmt.Sprintf("Payment status for Order #%s has been updated to %s.", order.OrderNumber, string(status))
		notificationType = core.NotificationTypePaymentUpdate
	}

	// Create notification record
	notification := &core.Notification{
		UserID:    user.ID,
		OrderID:   &orderID,
		Type:      notificationType,
		Message:   message,
		IsRead:    false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save notification to database
	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return fmt.Errorf("failed to save notification: %w", err)
	}

	// Send notifications through multiple channels
	go s.sendMultiChannelNotification(ctx, user, message, fmt.Sprintf("Payment Update - Order #%s", order.OrderNumber))

	return nil
}

// sendMultiChannelNotification sends notification through multiple channels
func (s *NotificationService) sendMultiChannelNotification(ctx context.Context, user *core.User, message, title string) {
	// Send email notification
	if user.Email != "" {
		if err := s.SendEmailNotification(ctx, user.Email, title, message); err != nil {
			log.Printf("Failed to send email notification: %v", err)
		}
	}

	// Send SMS notification
	if user.Phone != "" {
		if err := s.SendSMSNotification(ctx, user.Phone, message); err != nil {
			log.Printf("Failed to send SMS notification: %v", err)
		}
	}

	// Send WhatsApp notification
	if user.Phone != "" {
		if err := s.SendWhatsAppNotification(ctx, user.Phone, message); err != nil {
			log.Printf("Failed to send WhatsApp notification: %v", err)
		}
	}

	// Send push notification
	if err := s.SendPushNotification(ctx, user.ID, title, message); err != nil {
		log.Printf("Failed to send push notification: %v", err)
	}
}

// SendWelcomeNotification sends welcome notification to new users
func (s *NotificationService) SendWelcomeNotification(ctx context.Context, userID uuid.UUID) error {
	// Get user details
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	message := fmt.Sprintf("Welcome to iPhone Service! We're excited to help you with your iPhone repair needs.")
	title := "Welcome to iPhone Service"

	// Create notification record
	notification := &core.Notification{
		UserID:    user.ID,
		OrderID:   nil,
		Type:      core.NotificationTypeWelcome,
		Message:   message,
		IsRead:    false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save notification to database
	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return fmt.Errorf("failed to save notification: %w", err)
	}

	// Send notifications through multiple channels
	go s.sendMultiChannelNotification(ctx, user, message, title)

	return nil
}

// SendPromotionalNotification sends promotional notification
func (s *NotificationService) SendPromotionalNotification(ctx context.Context, userID uuid.UUID, message string) error {
	// Get user details
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	title := "Special Offer"

	// Create notification record
	notification := &core.Notification{
		UserID:    user.ID,
		OrderID:   nil,
		Type:      core.NotificationTypePromotion,
		Message:   message,
		IsRead:    false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save notification to database
	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return fmt.Errorf("failed to save notification: %w", err)
	}

	// Send notifications through multiple channels
	go s.sendMultiChannelNotification(ctx, user, message, title)

	return nil
}

// SendSystemNotification sends system-wide notification
func (s *NotificationService) SendSystemNotification(ctx context.Context, message string) error {
	users, _, err := s.userRepo.List(ctx, 0, 10000, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to get active users: %w", err)
	}

	title := "System Notification"

	// Send notification to all users
	for _, user := range users {
		// Create notification record
		notification := &core.Notification{
			UserID:    user.ID,
			OrderID:   nil,
			Type:      core.NotificationTypeSystem,
			Message:   message,
			IsRead:    false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Save notification to database
		if err := s.notificationRepo.Create(ctx, notification); err != nil {
			log.Printf("Failed to save system notification for user %s: %v", user.ID.String(), err)
			continue
		}

		// Send notifications through multiple channels
		go s.sendMultiChannelNotification(ctx, user, message, title)
	}

	return nil
}
