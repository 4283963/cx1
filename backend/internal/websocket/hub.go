package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"smart-home/internal/services"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Send chan []byte
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

var hub = &Hub{
	broadcast:  make(chan []byte, 256),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("New client connected. Total clients: %d", len(h.clients))
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("Client disconnected. Total clients: %d", len(h.clients))
		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func StartHub() {
	go hub.run()
	log.Println("WebSocket hub started")
}

func BroadcastMessage(msgType string, data interface{}) {
	message := map[string]interface{}{
		"type": msgType,
		"data": data,
	}
	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling broadcast message: %v", err)
		return
	}
	hub.broadcast <- jsonData
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	hub.register <- client

	defer func() {
		hub.unregister <- client
		conn.Close()
	}()

	go client.writePump()
	client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		hub.unregister <- c
		c.Conn.Close()
	}()

	systemService := services.NewSystemService()

	for {
		_, msgBytes, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}

		var msg struct {
			Type string                 `json:"type"`
			Data map[string]interface{} `json:"data"`
		}

		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			log.Printf("WebSocket message parse error: %v", err)
			continue
		}

		switch msg.Type {
		case "force_mode":
			if enabled, ok := msg.Data["enabled"].(bool); ok {
				systemService.SetForceMode(enabled)
				log.Printf("收到客户端强启模式请求: enabled=%v", enabled)

				response := map[string]interface{}{
					"type": "force_mode_result",
					"data": map[string]interface{}{
						"success": true,
						"enabled": enabled,
					},
				}
				responseBytes, _ := json.Marshal(response)
				c.Send <- responseBytes

				BroadcastMessage("force_mode_update", map[string]interface{}{
					"enabled": enabled,
				})
			}
		}
	}
}

func (c *Client) writePump() {
	defer func() {
		c.Conn.Close()
	}()

	for message := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("WebSocket write error: %v", err)
			break
		}
	}
}
