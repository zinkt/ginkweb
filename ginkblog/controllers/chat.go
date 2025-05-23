package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zinkt/ginkweb/gink"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { // 允许跨域
		return true
	},
}

// Message struct to include roomID, UserID, and Content
type Message struct {
	RoomID    string `json:"roomID"`
	UserID    string `json:"userID"`
	Content   string `json:"content"` // Changed from []byte to string for direct JSON marshaling
	Timestamp string `json:"timestamp"`
	UserIP    string `json:"userIP"` // Stores IP:Port
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients within specific rooms.
type Hub struct {
	// Registered clients. clients[roomID][client]
	rooms map[string]map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Mutex to protect rooms map
	mu sync.RWMutex // Use RWMutex for better concurrency on reads
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		rooms:      make(map[string]map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if _, ok := h.rooms[client.roomID]; !ok {
				h.rooms[client.roomID] = make(map[*Client]bool)
			}
			h.rooms[client.roomID][client] = true
			h.mu.Unlock()
			log.Printf("Client %s (%s) registered to room %s", client.userID, client.remoteAddr, client.roomID)
		case client := <-h.unregister:
			h.mu.Lock()
			if roomClients, ok := h.rooms[client.roomID]; ok {
				if _, ok_client := roomClients[client]; ok_client {
					delete(roomClients, client)
					close(client.send)
					if len(roomClients) == 0 {
						delete(h.rooms, client.roomID)
						log.Printf("Room %s closed", client.roomID)
					}
				}
			}
			h.mu.Unlock()
			log.Printf("Client %s (%s) unregistered from room %s", client.userID, client.remoteAddr, client.roomID)
		case message := <-h.broadcast: // message is of type Message
			h.mu.RLock()
			if roomClients, ok := h.rooms[message.RoomID]; ok {
				// log.Printf("Broadcasting message from %s in room %s: %s", message.UserID, message.RoomID, message.Content)
				// Prepare the message to be sent to clients (JSON)
				// The actual marshaling will happen in writePump or before sending to client.send
				for client := range roomClients {
					if client.userID == message.UserID {
						continue
					}
					select {
					// Send the structured Message itself, writePump will marshal it
					case client.send <- message:
					default:
						log.Printf("Failed to send message to client %s in room %s (channel full or closed), removing client.", client.userID, client.roomID)
						// It's safer to trigger unregister from here if the channel is problematic,
						// but for simplicity, we'll just log and expect pong/ping or read errors to handle cleanup.
						// Consider closing client.send from here if consistently failing might indicate a stuck client.
						// However, closing client.send here without removing from h.rooms can lead to send on closed channel.
						// A robust solution might involve a separate cleanup goroutine based on send failures.
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan Message // Changed from chan []byte to chan Message

	// User ID or name
	userID string

	// Room ID
	roomID string

	// Remote address
	remoteAddr string // To store IP:Port
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// readPump pumps messages from the websocket connection to the hub.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { _ = c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, rawMessage, err := c.conn.ReadMessage() // rawMessage is []byte
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message for client %s (%s) in room %s: %v", c.userID, c.remoteAddr, c.roomID, err)
			}
			break
		}
		// Content is now the pure message string
		message := Message{
			RoomID:    c.roomID,
			UserID:    c.userID,
			Content:   string(rawMessage),
			Timestamp: time.Now().Format("15:04:05"), // Format time as HH:MM:SS
			UserIP:    c.remoteAddr,
		}
		c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send: // message is of type Message
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Printf("Send channel closed for client %s (%s) in room %s", c.userID, c.remoteAddr, c.roomID)
				return
			}

			// Marshal the Message struct to JSON
			jsonMessage, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error marshaling message for client %s (%s) in room %s: %v", c.userID, c.remoteAddr, c.roomID, err)
				continue // Or return, depending on desired error handling
			}

			err = c.conn.WriteMessage(websocket.TextMessage, jsonMessage)
			if err != nil {
				log.Printf("error writing json message to client %s (%s) in room %s: %v", c.userID, c.remoteAddr, c.roomID, err)
				return
			}

		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("error sending ping to client %s (%s) in room %s: %v", c.userID, c.remoteAddr, c.roomID, err)
				return
			}
		}
	}
}

// Global hub instance
var hub = newHub()

func init() {
	go hub.run()
}

// Chat handles websocket requests from the peer.
func Chat(ctx *gink.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Req, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	userID := ctx.Query("userID")
	roomID := ctx.Query("roomID")
	remoteAddrStr := conn.RemoteAddr().String() // Get IP:Port

	// You might want to parse out just the IP if needed, but IP:Port is often useful
	// host, _, err := net.SplitHostPort(remoteAddrStr)
	// if err == nil {
	// 	 remoteAddrStr = host
	// }

	if userID == "" {
		userID = "Anonymous"
		log.Println("UserID not provided, defaulting to Anonymous for", remoteAddrStr)
	}
	if roomID == "" {
		// You might want to return an error or assign to a default room
		log.Printf("RoomID not provided for userID %s (%s). Closing connection.", userID, remoteAddrStr)
		_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "roomID is required"))
		_ = conn.Close()
		return
	}

	log.Printf("New WebSocket connection: userID=%s, roomID=%s, remoteAddr=%s", userID, roomID, remoteAddrStr)

	client := &Client{
		hub:        hub,
		conn:       conn,
		send:       make(chan Message, 256), // Changed to chan Message
		userID:     userID,
		roomID:     roomID,
		remoteAddr: remoteAddrStr, // Store the remote address
	}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
