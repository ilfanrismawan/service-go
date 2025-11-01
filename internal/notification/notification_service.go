package notification

import (
    "bytes"
	"context"
    "encoding/json"
	"fmt"
	"log"
    "net/http"
    "service/internal/config"
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

// SendWhatsAppNotification sends WhatsApp notification with template support
func (s *NotificationService) SendWhatsAppNotification(ctx context.Context, phone, message string) error {
	return s.SendWhatsAppNotificationWithTemplate(ctx, phone, "", nil, message)
}

// SendWhatsAppNotificationWithTemplate sends WhatsApp notification using template
func (s *NotificationService) SendWhatsAppNotificationWithTemplate(ctx context.Context, phone string, templateType WhatsAppTemplateType, templateData map[string]interface{}, fallbackMessage string) error {
	// If no API key configured, fallback to mock to keep dev UX smooth
	if config.Config == nil || (config.Config.TwilioAuthToken == "" && config.Config.FirebaseServerKey == "") { /* noop */ }
	apiKey := config.Config.WhatsAppAPIKey

	// Determine message content
	var message string
	if templateType != "" && templateData != nil {
		message = GetWhatsAppTemplate(templateType, templateData)
	} else {
		message = fallbackMessage
	}

	if apiKey == "" {
		log.Printf("💬 Mock WhatsApp sent to %s: %s", phone, message)
		time.Sleep(150 * time.Millisecond)
		return nil
	}

	// Fonnte simple integration
	// Docs: POST https://api.fonnte.com/send with headers: Authorization: <TOKEN>
	// Body (JSON): { "target": "08xxxx" or "+62...", "message": "..." }
	apiURL := config.Config.WhatsAppAPIURL
	if apiURL == "" {
		apiURL = "https://api.fonnte.com/send"
	}

	payload := map[string]string{
		"target":  phone,
		"message": message,
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("WhatsApp provider error: status=%d", resp.StatusCode)
		return fmt.Errorf("whatsapp provider error: %d", resp.StatusCode)
	}
	return nil
}

// SendPushNotification sends push notification using FCM
func (s *NotificationService) SendPushNotification(ctx context.Context, userID uuid.UUID, title, body string) error {
	// Get user to retrieve FCM token
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// If no FCM token, log and return (don't fail)
	if user.FCMToken == "" {
		log.Printf("⚠️ No FCM token for user %s, skipping push notification", userID.String())
		return nil
	}

	// Initialize FCM service
	fcmService := NewFCMService()

	// Create notification payload
	notification := &FCMNotification{
		Title: title,
		Body:  body,
		Sound: SoundDefault,
	}

	// Send notification
	if err := fcmService.SendToToken(ctx, user.FCMToken, notification, nil); err != nil {
		log.Printf("Failed to send FCM notification to user %s: %v", userID.String(), err)
		return fmt.Errorf("failed to send push notification: %w", err)
	}

	log.Printf("✅ Push notification sent to user %s: %s - %s", userID.String(), title, body)
	return nil
}

// SendPushNotificationWithData sends push notification with additional data
func (s *NotificationService) SendPushNotificationWithData(ctx context.Context, userID uuid.UUID, title, body string, data map[string]interface{}) error {
	// Get user to retrieve FCM token
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// If no FCM token, log and return (don't fail)
	if user.FCMToken == "" {
		log.Printf("⚠️ No FCM token for user %s, skipping push notification", userID.String())
		return nil
	}

	// Initialize FCM service
	fcmService := NewFCMService()

	// Create notification payload
	notification := &FCMNotification{
		Title: title,
		Body:  body,
		Sound: SoundDefault,
	}

	// Send notification with data
	if err := fcmService.SendToToken(ctx, user.FCMToken, notification, data); err != nil {
		log.Printf("Failed to send FCM notification to user %s: %v", userID.String(), err)
		return fmt.Errorf("failed to send push notification: %w", err)
	}

	log.Printf("✅ Push notification with data sent to user %s: %s - %s", userID.String(), title, body)
	return nil
}

// SendPushNotificationToMultiple sends push notification to multiple users
func (s *NotificationService) SendPushNotificationToMultiple(ctx context.Context, userIDs []uuid.UUID, title, body string) error {
	// Get all users and collect FCM tokens
	var tokens []string
	for _, userID := range userIDs {
		user, err := s.userRepo.GetByID(ctx, userID)
		if err != nil {
			continue // Skip users that can't be found
		}
		if user.FCMToken != "" {
			tokens = append(tokens, user.FCMToken)
		}
	}

	if len(tokens) == 0 {
		log.Printf("⚠️ No FCM tokens found for %d users, skipping push notification", len(userIDs))
		return nil
	}

	// Initialize FCM service
	fcmService := NewFCMService()

	// Create notification payload
	notification := &FCMNotification{
		Title: title,
		Body:  body,
		Sound: SoundDefault,
	}

	// Send to all tokens
	if err := fcmService.SendToTokens(ctx, tokens, notification, nil); err != nil {
		log.Printf("Failed to send FCM notification to multiple users: %v", err)
		return fmt.Errorf("failed to send push notification: %w", err)
	}

	log.Printf("✅ Push notification sent to %d users", len(tokens))
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

	// Send notifications through multiple channels with WhatsApp template
	templateData := map[string]interface{}{
		"order_number":  order.OrderNumber,
		"status":        string(status),
		"customer_name": user.FullName,
		"branch_name":   order.Branch.Name,
	}
	go s.sendMultiChannelNotificationWithTemplate(ctx, user, TemplateStatusUpdate, templateData, message, fmt.Sprintf("Order #%s Update", order.OrderNumber))

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

	// Send notifications through multiple channels with WhatsApp template
	templateData := map[string]interface{}{
		"order_number": order.OrderNumber,
		"amount":       order.ActualCost,
		"due_date":     time.Now().Add(24 * time.Hour),
		"customer_name": user.FullName,
	}
	templateType := TemplatePaymentReminder
	if status == core.PaymentStatusPaid {
		templateType = TemplateStatusUpdate
	}
	go s.sendMultiChannelNotificationWithTemplate(ctx, user, templateType, templateData, message, fmt.Sprintf("Payment Update - Order #%s", order.OrderNumber))

	return nil
}

// sendMultiChannelNotification sends notification through multiple channels
func (s *NotificationService) sendMultiChannelNotification(ctx context.Context, user *core.User, message, title string) {
	s.sendMultiChannelNotificationWithTemplate(ctx, user, "", nil, message, title)
}

// sendMultiChannelNotificationWithTemplate sends notification through multiple channels with WhatsApp template support
func (s *NotificationService) sendMultiChannelNotificationWithTemplate(ctx context.Context, user *core.User, templateType WhatsAppTemplateType, templateData map[string]interface{}, fallbackMessage, title string) {
	// Send email notification
	if user.Email != "" {
		if err := s.SendEmailNotification(ctx, user.Email, title, fallbackMessage); err != nil {
			log.Printf("Failed to send email notification: %v", err)
		}
	}

	// Send SMS notification
	if user.Phone != "" {
		if err := s.SendSMSNotification(ctx, user.Phone, fallbackMessage); err != nil {
			log.Printf("Failed to send SMS notification: %v", err)
		}
	}

	// Send WhatsApp notification with template
	if user.Phone != "" {
		if err := s.SendWhatsAppNotificationWithTemplate(ctx, user.Phone, templateType, templateData, fallbackMessage); err != nil {
			log.Printf("Failed to send WhatsApp notification: %v", err)
		}
	}

	// Send push notification
	if err := s.SendPushNotification(ctx, user.ID, title, fallbackMessage); err != nil {
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
