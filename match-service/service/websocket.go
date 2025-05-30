package service

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Client represents a single WebSocket connection.
type Client struct {
	conn    *websocket.Conn
	send    chan []byte // Buffered channel for outgoing messages
	matchID string      // The match this client is interested in
}

// WebSocketHub manages all WebSocket connections and broadcasting.
type WebSocketHub struct {
	clients map[*Client]bool // Registered clients
	// Mapping from matchID to clients interested in that match
	matchSubscriptions map[string]map[*Client]bool
	broadcast          chan []byte // Inbound messages from clients (if any)
	register           chan *Client
	unregister         chan *Client
	mu                 sync.RWMutex // Mutex to protect client maps
}

// NewWebSocketHub creates a new WebSocketHub.
func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		broadcast:          make(chan []byte),
		register:           make(chan *Client),
		unregister:         make(chan *Client),
		clients:            make(map[*Client]bool),
		matchSubscriptions: make(map[string]map[*Client]bool),
	}
}

// Run starts the hub's goroutine for managing clients and messages.
func (h *WebSocketHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			if client.matchID != "" {
				if _, ok := h.matchSubscriptions[client.matchID]; !ok {
					h.matchSubscriptions[client.matchID] = make(map[*Client]bool)
				}
				h.matchSubscriptions[client.matchID][client] = true
			}
			h.mu.Unlock()
			log.Printf("Client registered for match %s. Total clients: %d\n", client.matchID, len(h.clients))
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				if client.matchID != "" {
					if subs, ok := h.matchSubscriptions[client.matchID]; ok {
						delete(subs, client)
						if len(subs) == 0 {
							delete(h.matchSubscriptions, client.matchID) // Clean up empty subscriptions
						}
					}
				}
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("Client unregistered. Total clients: %d\n", len(h.clients))
		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastMatchUpdate sends a match update to all clients subscribed to that match.
// This is called by the MatchService when data changes.
func (h *WebSocketHub) BroadcastMatchUpdate(matchID string, data interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	message, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling match update for broadcast: %v", err)
		return
	}

	if subs, ok := h.matchSubscriptions[matchID]; ok {
		for client := range subs {
			select {
			case client.send <- message:
			default:
				log.Printf("Client send channel full for match %s, unregistering.", matchID)
				// Handle slow client by unregistering or logging
				h.unregister <- client // Send to unregister channel to clean up
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for development, but restrict in production!
		return true
	},
}

// HandleConnections handles new WebSocket connections.
// It expects a 'match_id' query parameter to subscribe clients to specific matches.
// Example: ws://localhost:8080/ws?match_id=match-123
func (h *WebSocketHub) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// Extract match_id from query parameter
	matchID := r.URL.Query().Get("match_id")
	if matchID == "" {
		// Optionally, disconnect if no match_id provided or provide general updates
		log.Println("WebSocket connection without match_id provided. Providing general updates.")
		// For a specific match, match_id is critical
	}

	client := &Client{conn: conn, send: make(chan []byte, 256), matchID: matchID} // Increased buffer size
	h.register <- client

	go client.writePump() // Goroutine to write messages to the client
	client.readPump(h)    // Blocking call to read messages from the client
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Client) readPump(hub *WebSocketHub) {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))                                                           // Set a deadline for pings
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(60 * time.Second)); return nil }) // Extend deadline on pong
	for {
		_, _, err := c.conn.ReadMessage() // Read messages from client (e.g., pings or control messages)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		// If you expect messages from clients, process them here
		// For now, we only expect clients to receive updates.
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second) // Ping interval
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current WebSocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			// Send ping messages to keep the connection alive
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
