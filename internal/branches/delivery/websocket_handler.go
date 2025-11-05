package delivery

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"service/internal/orders/service"
	"service/internal/shared/model"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	chatService *service.ChatService
	upgrader    websocket.Upgrader
	clients     map[*websocket.Conn]*Client
	rooms       map[string][]*websocket.Conn
	mutex       sync.RWMutex
}

// Client represents a WebSocket client
type Client struct {
	Conn    *websocket.Conn
	UserID  uuid.UUID
	OrderID uuid.UUID
	Send    chan []byte
	Handler *WebSocketHandler
}

// Message represents a WebSocket message
type Message struct {
	Type      string      `json:"type"`
	OrderID   string      `json:"order_id"`
	UserID    string      `json:"user_id"`
	Content   string      `json:"content"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data,omitempty"`
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{
		chatService: service.NewChatService(),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins in development
			},
		},
		clients: make(map[*websocket.Conn]*Client),
		rooms:   make(map[string][]*websocket.Conn),
	}
}

// HandleWebSocket handles WebSocket connections
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.CreateErrorResponse(
			"unauthorized",
			"User ID not found in context",
			nil,
		))
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
			"internal_error",
			"Invalid user ID type",
			nil,
		))
		return
	}

	// Get order ID from query
	orderIDStr := c.Query("order_id")
	if orderIDStr == "" {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"order_id_required",
			"Order ID is required",
			nil,
		))
		return
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"invalid_order_id",
			"Invalid order ID format",
			nil,
		))
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// Create client
	client := &Client{
		Conn:    conn,
		UserID:  userUUID,
		OrderID: orderID,
		Send:    make(chan []byte, 256),
		Handler: h,
	}

	// Register client
	h.mutex.Lock()
	h.clients[conn] = client
	h.rooms[orderID.String()] = append(h.rooms[orderID.String()], conn)
	h.mutex.Unlock()

	// Start goroutines
	go client.writePump()
	go client.readPump()
}

// readPump handles reading messages from WebSocket
func (c *Client) readPump() {
	defer func() {
		c.Handler.unregisterClient(c)
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, messageBytes, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Parse message
		var message Message
		if err := json.Unmarshal(messageBytes, &message); err != nil {
			log.Printf("Message parsing error: %v", err)
			continue
		}

		// Handle message
		c.Handler.handleMessage(c, &message)
	}
}

// writePump handles writing messages to WebSocket
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage handles incoming messages
func (h *WebSocketHandler) handleMessage(client *Client, message *Message) {
	switch message.Type {
	case "chat":
		h.handleChatMessage(client, message)
	case "typing":
		h.handleTypingMessage(client, message)
	case "ping":
		h.handlePingMessage(client, message)
	default:
		log.Printf("Unknown message type: %s", message.Type)
	}
}

// handleChatMessage handles chat messages
func (h *WebSocketHandler) handleChatMessage(client *Client, message *Message) {
	// Create chat message request
	req := &model.ChatMessageRequest{
		OrderID:    message.OrderID,
		ReceiverID: message.UserID,
		Message:    message.Content,
	}

	// Send message through chat service
	chatMessage, err := h.chatService.SendMessage(context.Background(), client.UserID, req)
	if err != nil {
		log.Printf("Chat message error: %v", err)
		return
	}

	// Broadcast message to room
	h.broadcastToRoom(message.OrderID, Message{
		Type:      "chat",
		OrderID:   message.OrderID,
		UserID:    client.UserID.String(),
		Content:   message.Content,
		Timestamp: time.Now(),
		Data:      chatMessage,
	})
}

// handleTypingMessage handles typing indicators
func (h *WebSocketHandler) handleTypingMessage(client *Client, message *Message) {
	// Broadcast typing indicator to room (excluding sender)
	h.broadcastToRoomExcluding(message.OrderID, client.Conn, Message{
		Type:      "typing",
		OrderID:   message.OrderID,
		UserID:    client.UserID.String(),
		Content:   message.Content,
		Timestamp: time.Now(),
	})
}

// handlePingMessage handles ping messages
func (h *WebSocketHandler) handlePingMessage(client *Client, message *Message) {
	// Send pong response
	client.Send <- []byte(`{"type":"pong","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`)
}

// broadcastToRoom broadcasts a message to all clients in a room
func (h *WebSocketHandler) broadcastToRoom(orderID string, message Message) {
	h.mutex.RLock()
	clients := h.rooms[orderID]
	h.mutex.RUnlock()

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Message marshaling error: %v", err)
		return
	}

	for _, conn := range clients {
		if client, exists := h.clients[conn]; exists {
			select {
			case client.Send <- messageBytes:
			default:
				h.unregisterClient(client)
			}
		}
	}
}

// broadcastToRoomExcluding broadcasts a message to all clients in a room except the sender
func (h *WebSocketHandler) broadcastToRoomExcluding(orderID string, sender *websocket.Conn, message Message) {
	h.mutex.RLock()
	clients := h.rooms[orderID]
	h.mutex.RUnlock()

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Message marshaling error: %v", err)
		return
	}

	for _, conn := range clients {
		if conn != sender {
			if client, exists := h.clients[conn]; exists {
				select {
				case client.Send <- messageBytes:
				default:
					h.unregisterClient(client)
				}
			}
		}
	}
}

// unregisterClient removes a client from the handler
func (h *WebSocketHandler) unregisterClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// Remove from clients map
	delete(h.clients, client.Conn)

	// Remove from room
	orderID := client.OrderID.String()
	if clients, exists := h.rooms[orderID]; exists {
		for i, conn := range clients {
			if conn == client.Conn {
				h.rooms[orderID] = append(clients[:i], clients[i+1:]...)
				break
			}
		}
	}

	// Close send channel
	close(client.Send)
}

// GetConnectedUsers returns the number of connected users
func (h *WebSocketHandler) GetConnectedUsers() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return len(h.clients)
}

// GetRoomUsers returns the number of users in a specific room
func (h *WebSocketHandler) GetRoomUsers(orderID string) int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return len(h.rooms[orderID])
}
