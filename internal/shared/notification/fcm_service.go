package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"service/internal/shared/config"
	"time"
)

// FCMService handles Firebase Cloud Messaging operations
type FCMService struct {
	serverKey string
	client    *http.Client
}

// NewFCMService creates a new FCM service
func NewFCMService() *FCMService {
	return &FCMService{
		serverKey: config.Config.FirebaseServerKey,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// FCMMessage represents FCM message payload
type FCMMessage struct {
	To               string                 `json:"to,omitempty"`               // Single device token
	RegistrationIDs  []string                `json:"registration_ids,omitempty"` // Multiple device tokens
	Condition        string                  `json:"condition,omitempty"`         // Topic condition
	Notification     *FCMNotification        `json:"notification,omitempty"`
	Data             map[string]interface{}  `json:"data,omitempty"`
	Priority         string                  `json:"priority,omitempty"`         // "normal" or "high"
	ContentAvailable bool                    `json:"content_available,omitempty"`
	MutableContent   bool                    `json:"mutable_content,omitempty"`
}

// FCMNotification represents FCM notification payload
type FCMNotification struct {
	Title        string `json:"title,omitempty"`
	Body         string `json:"body,omitempty"`
	Sound        string `json:"sound,omitempty"`
	Badge        string `json:"badge,omitempty"`
	ClickAction  string `json:"click_action,omitempty"`
	ImageURL     string `json:"image,omitempty"`
	Tag          string `json:"tag,omitempty"`
	Color        string `json:"color,omitempty"`
}

// FCMResponse represents FCM API response
type FCMResponse struct {
	MulticastID  int64                  `json:"multicast_id,omitempty"`
	Success      int                    `json:"success,omitempty"`
	Failure      int                    `json:"failure,omitempty"`
	CanonicalIDs int                    `json:"canonical_ids,omitempty"`
	Results      []FCMResult            `json:"results,omitempty"`
	MessageID    int64                  `json:"message_id,omitempty"`
	Error        map[string]interface{} `json:"error,omitempty"`
}

// FCMResult represents result for each token
type FCMResult struct {
	MessageID      string `json:"message_id,omitempty"`
	RegistrationID string `json:"registration_id,omitempty"`
	Error          string `json:"error,omitempty"`
}

const (
	FCMAPIEndpoint = "https://fcm.googleapis.com/fcm/send"
	PriorityHigh   = "high"
	PriorityNormal = "normal"
	SoundDefault   = "default"
)

// SendToToken sends FCM message to a single device token
func (s *FCMService) SendToToken(ctx context.Context, token string, notification *FCMNotification, data map[string]interface{}) error {
	if s.serverKey == "" {
		log.Printf("🔔 Mock FCM notification (FIREBASE_SERVER_KEY not set): %s - %s", notification.Title, notification.Body)
		time.Sleep(50 * time.Millisecond)
		return nil
	}

	if token == "" {
		return fmt.Errorf("FCM token is empty")
	}

	message := &FCMMessage{
		To:           token,
		Notification: notification,
		Data:         data,
		Priority:     PriorityHigh,
	}

	return s.sendMessage(ctx, message)
}

// SendToTokens sends FCM message to multiple device tokens
func (s *FCMService) SendToTokens(ctx context.Context, tokens []string, notification *FCMNotification, data map[string]interface{}) error {
	if s.serverKey == "" {
		log.Printf("🔔 Mock FCM notification (FIREBASE_SERVER_KEY not set): %s - %s", notification.Title, notification.Body)
		time.Sleep(50 * time.Millisecond)
		return nil
	}

	if len(tokens) == 0 {
		return fmt.Errorf("no FCM tokens provided")
	}

	// If single token, use SendToToken
	if len(tokens) == 1 {
		return s.SendToToken(ctx, tokens[0], notification, data)
	}

	message := &FCMMessage{
		RegistrationIDs: tokens,
		Notification:    notification,
		Data:            data,
		Priority:        PriorityHigh,
	}

	return s.sendMessage(ctx, message)
}

// SendToTopic sends FCM message to a topic
func (s *FCMService) SendToTopic(ctx context.Context, topic string, notification *FCMNotification, data map[string]interface{}) error {
	if s.serverKey == "" {
		log.Printf("🔔 Mock FCM notification (FIREBASE_SERVER_KEY not set): %s - %s", notification.Title, notification.Body)
		time.Sleep(50 * time.Millisecond)
		return nil
	}

	if topic == "" {
		return fmt.Errorf("topic is empty")
	}

	message := &FCMMessage{
		Condition:    fmt.Sprintf("'%s' in topics", topic),
		Notification: notification,
		Data:         data,
		Priority:     PriorityHigh,
	}

	return s.sendMessage(ctx, message)
}

// sendMessage sends message to FCM API
func (s *FCMService) sendMessage(ctx context.Context, message *FCMMessage) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal FCM message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", FCMAPIEndpoint, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create FCM request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("key=%s", s.serverKey))

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send FCM request: %w", err)
	}
	defer resp.Body.Close()

	var fcmResp FCMResponse
	if err := json.NewDecoder(resp.Body).Decode(&fcmResp); err != nil {
		return fmt.Errorf("failed to decode FCM response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("FCM API error: status=%d, error=%v", resp.StatusCode, fcmResp.Error)
	}

	// Check for failures in multicast response
	if fcmResp.Failure > 0 {
		log.Printf("⚠️ FCM: %d failed, %d succeeded out of %d total", fcmResp.Failure, fcmResp.Success, fcmResp.Failure+fcmResp.Success)
		for i, result := range fcmResp.Results {
			if result.Error != "" {
				log.Printf("  Token %d error: %s", i, result.Error)
			}
		}
		return fmt.Errorf("FCM delivery failed for %d tokens", fcmResp.Failure)
	}

	log.Printf("✅ FCM notification sent successfully")
	return nil
}

