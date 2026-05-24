package handlers

import (
	"encoding/json"
	"errand-service/internal/models"
	"errand-service/pkg/logger"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WebSocketHandler struct {
	clients   map[string]map[uint]*websocket.Conn
	mu        sync.RWMutex
	db        *gorm.DB
}

func NewWebSocketHandler(db *gorm.DB) *WebSocketHandler {
	return &WebSocketHandler{
		clients: make(map[string]map[uint]*websocket.Conn),
		db:      db,
	}
}

type WSMessage struct {
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
	OrderID uint        `json:"order_id,omitempty"`
	UserID  uint        `json:"user_id,omitempty"`
}

func (h *WebSocketHandler) HandleConnection(c *gin.Context) {
	userID := c.GetUint("user_id")
	orderID := c.Param("orderId")

	var order models.Order
	if err := h.db.Where("id = ? AND (publisher_id = ? OR courier_id = ?)", orderID, userID, userID).
		First(&order).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "Access denied"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Errorf("WebSocket upgrade error: %v", err)
		return
	}

	h.mu.Lock()
	if h.clients[orderID] == nil {
		h.clients[orderID] = make(map[uint]*websocket.Conn)
	}
	h.clients[orderID][userID] = conn
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.clients[orderID], userID)
		if len(h.clients[orderID]) == 0 {
			delete(h.clients, orderID)
		}
		h.mu.Unlock()
		conn.Close()
	}()

	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Errorf("WebSocket read error: %v", err)
			}
			break
		}

		var msg WSMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		switch msg.Type {
		case "chat":
			h.handleChatMessage(orderID, userID, msg)
		case "location":
			h.handleLocationUpdate(orderID, userID, msg)
		case "status":
			h.handleStatusUpdate(orderID, userID, msg)
		}
	}
}

func (h *WebSocketHandler) handleChatMessage(orderID string, userID uint, msg WSMessage) {
	chatMsg := models.ChatMessage{
		OrderID:  0,
		SenderID: userID,
		Content:  msg.Data.(string),
		MsgType:  "text",
	}

	if id, err := parseOrderID(orderID); err == nil {
		chatMsg.OrderID = id
	}

	h.db.Create(&chatMsg)

	h.broadcastToOrder(orderID, WSMessage{
		Type:   "chat",
		Data:   chatMsg,
		UserID: userID,
	})
}

func (h *WebSocketHandler) handleLocationUpdate(orderID string, userID uint, msg WSMessage) {
	locationData, _ := msg.Data.(map[string]interface{})

	track := models.OrderTrack{
		EventType: "location",
	}

	if lat, ok := locationData["latitude"].(float64); ok {
		track.Latitude = lat
	}
	if lng, ok := locationData["longitude"].(float64); ok {
		track.Longitude = lng
	}
	if addr, ok := locationData["address"].(string); ok {
		track.Address = addr
	}

	if id, err := parseOrderID(orderID); err == nil {
		track.OrderID = id
	}

	h.db.Create(&track)

	h.broadcastToOrder(orderID, WSMessage{
		Type:   "location",
		Data:   track,
		UserID: userID,
	})
}

func (h *WebSocketHandler) handleStatusUpdate(orderID string, userID uint, msg WSMessage) {
	h.broadcastToOrder(orderID, WSMessage{
		Type:   "status",
		Data:   msg.Data,
		UserID: userID,
	})
}

func (h *WebSocketHandler) broadcastToOrder(orderID string, msg WSMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	if clients, ok := h.clients[orderID]; ok {
		for _, conn := range clients {
			conn.WriteMessage(websocket.TextMessage, data)
		}
	}
}

func (h *WebSocketHandler) SendNotification(orderID string, notification WSMessage) {
	h.broadcastToOrder(orderID, notification)
}

func parseOrderID(orderID string) (uint, error) {
	var id uint
	_, err := fmt.Sscanf(orderID, "%d", &id)
	return id, err
}
