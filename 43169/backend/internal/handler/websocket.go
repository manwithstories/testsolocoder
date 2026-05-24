package handler

import (
	"encoding/json"
	"log"
	"matchmaking-platform/internal/model"
	"matchmaking-platform/internal/service"
	"matchmaking-platform/internal/utils"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	ID   uint
	Conn *websocket.Conn
	Send chan []byte
	mu   sync.Mutex
}

type ChatServer struct {
	clients    map[uint]*Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

var chatServer = &ChatServer{
	clients:    make(map[uint]*Client),
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
}

func (s *ChatServer) Run() {
	for {
		select {
		case client := <-s.register:
			s.mu.Lock()
			s.clients[client.ID] = client
			s.mu.Unlock()
			log.Printf("Client connected: %d, total: %d", client.ID, len(s.clients))

		case client := <-s.unregister:
			s.mu.Lock()
			if _, ok := s.clients[client.ID]; ok {
				close(client.Send)
				delete(s.clients, client.ID)
			}
			s.mu.Unlock()
			log.Printf("Client disconnected: %d", client.ID)

		case message := <-s.broadcast:
			s.mu.RLock()
			for _, client := range s.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(s.clients, client.ID)
				}
			}
			s.mu.RUnlock()
		}
	}
}

func (c *Client) ReadPump() {
	defer func() {
		chatServer.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(65536)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var msg struct {
			Type       string `json:"type"`
			ReceiverID uint   `json:"receiver_id"`
			Content    string `json:"content"`
		}
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		if msg.Type == "chat" {
			handleChatMessage(c, msg.ReceiverID, msg.Content)
		} else if msg.Type == "typing" {
			broadcastToUser(msg.ReceiverID, map[string]interface{}{
				"type":      "typing",
				"sender_id": c.ID,
			})
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(30 * time.Second)
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
			c.mu.Lock()
			c.Conn.WriteMessage(websocket.TextMessage, message)
			c.mu.Unlock()

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			c.mu.Lock()
			c.Conn.WriteMessage(websocket.PingMessage, nil)
			c.mu.Unlock()
		}
	}
}

func handleChatMessage(sender *Client, receiverID uint, content string) {
	svc := service.NewChatService()

	msg, err := svc.SendMessage(sender.ID, receiverID, "text", content)
	if err != nil {
		return
	}

	data, _ := json.Marshal(map[string]interface{}{
		"type":    "message",
		"message": msg,
	})

	broadcastToUser(receiverID, data)
	broadcastToUser(sender.ID, data)
}

func broadcastToUser(userID uint, data interface{}) {
	chatServer.mu.RLock()
	client, ok := chatServer.clients[userID]
	chatServer.mu.RUnlock()

	if ok {
		var msgBytes []byte
		switch v := data.(type) {
		case []byte:
			msgBytes = v
		default:
			msgBytes, _ = json.Marshal(v)
		}
		select {
		case client.Send <- msgBytes:
		default:
		}
	}
}

func WebSocketHandler(c *gin.Context) {
	userIDStr := c.Query("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid user_id"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	client := &Client{
		ID:   uint(userID),
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	chatServer.register <- client

	go client.WritePump()
	go client.ReadPump()
}

func GetOnlineUsers(c *gin.Context) {
	chatServer.mu.RLock()
	var onlineIDs []uint
	for id := range chatServer.clients {
		onlineIDs = append(onlineIDs, id)
	}
	chatServer.mu.RUnlock()

	utils.Success(c, gin.H{
		"online_users":  onlineIDs,
		"online_count":  len(onlineIDs),
	})
}

func RunChatServer() {
	go chatServer.Run()
}

func NotifyUser(userID uint, msgType string, data interface{}) {
	broadcastToUser(userID, map[string]interface{}{
		"type": msgType,
		"data": data,
	})
}

var _ = model.User{}
