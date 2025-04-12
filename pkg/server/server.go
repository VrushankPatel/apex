package server

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"arbitrage-detector/pkg/models"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// WebSocketMessage represents a message sent over WebSocket
type WebSocketMessage struct {
	Type        string      `json:"type"`
	Data        interface{} `json:"data"`
	Timestamp   int64       `json:"timestamp"`
}

// WebServer handles HTTP requests and WebSocket connections
type WebServer struct {
	port             string
	orderBooks       map[string]*models.OrderBook
	orderBookMutex   *sync.RWMutex
	clients          map[*websocket.Conn]bool
	clientsMutex     sync.Mutex
	opportunities    []models.ArbitrageOpportunity
	opportunitiesMutex sync.Mutex
	upgrader         websocket.Upgrader
}

// NewWebServer creates a new web server instance
func NewWebServer(port string, orderBooks map[string]*models.OrderBook, orderBookMutex *sync.RWMutex) *WebServer {
	return &WebServer{
		port:             port,
		orderBooks:       orderBooks,
		orderBookMutex:   orderBookMutex,
		clients:          make(map[*websocket.Conn]bool),
		opportunities:    make([]models.ArbitrageOpportunity, 0),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for WebSocket connections
			},
		},
	}
}

// Start initializes and runs the web server
func (s *WebServer) Start() error {
	// Create a file server to serve static files
	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/", fs)

	// WebSocket endpoint
	http.HandleFunc("/ws", s.handleWebSocket)

	// API endpoints
	http.HandleFunc("/api/opportunities", s.handleOpportunitiesAPI)
	http.HandleFunc("/api/market", s.handleMarketAPI)

	// Start market data broadcast
	go s.broadcastMarketData()

	// Start the server
	log.Infof("Starting web server on port %s", s.port)
	return http.ListenAndServe("0.0.0.0:"+s.port, nil)
}

// AddOpportunity adds a new arbitrage opportunity and broadcasts it to clients
func (s *WebServer) AddOpportunity(opportunity models.ArbitrageOpportunity) {
	s.opportunitiesMutex.Lock()
	defer s.opportunitiesMutex.Unlock()

	// Add to opportunities list
	s.opportunities = append(s.opportunities, opportunity)

	// Keep only the last 100 opportunities
	if len(s.opportunities) > 100 {
		s.opportunities = s.opportunities[len(s.opportunities)-100:]
	}

	// Broadcast to connected clients
	s.broadcastOpportunity(opportunity)
}

// handleWebSocket handles WebSocket connections
func (s *WebServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("Failed to upgrade connection to WebSocket: %v", err)
		return
	}

	// Register the new client
	s.clientsMutex.Lock()
	s.clients[conn] = true
	s.clientsMutex.Unlock()

	log.Infof("New WebSocket client connected: %s", conn.RemoteAddr())

	// Handle disconnection
	defer func() {
		s.clientsMutex.Lock()
		delete(s.clients, conn)
		s.clientsMutex.Unlock()
		conn.Close()
		log.Infof("WebSocket client disconnected: %s", conn.RemoteAddr())
	}()

	// Send initial data
	s.sendInitialData(conn)

	// Handle incoming messages (though we don't expect many)
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Errorf("WebSocket read error: %v", err)
			}
			break
		}
	}
}

// sendInitialData sends the initial data to a new WebSocket client
func (s *WebServer) sendInitialData(conn *websocket.Conn) {
	// Send market data
	s.orderBookMutex.RLock()
	marketData := make(map[string]*models.OrderBook)
	for k, v := range s.orderBooks {
		marketData[k] = v
	}
	s.orderBookMutex.RUnlock()

	marketMsg := WebSocketMessage{
		Type:      "market",
		Data:      marketData,
		Timestamp: time.Now().Unix(),
	}

	if err := conn.WriteJSON(marketMsg); err != nil {
		log.Errorf("Failed to send market data: %v", err)
	}

	// Send recent opportunities
	s.opportunitiesMutex.Lock()
	opportunities := make([]models.ArbitrageOpportunity, len(s.opportunities))
	copy(opportunities, s.opportunities)
	s.opportunitiesMutex.Unlock()

	oppMsg := WebSocketMessage{
		Type:      "opportunities",
		Data:      opportunities,
		Timestamp: time.Now().Unix(),
	}

	if err := conn.WriteJSON(oppMsg); err != nil {
		log.Errorf("Failed to send opportunities data: %v", err)
	}
}

// broadcastMarketData periodically broadcasts market data to all connected clients
func (s *WebServer) broadcastMarketData() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		s.orderBookMutex.RLock()
		marketData := make(map[string]*models.OrderBook)
		for k, v := range s.orderBooks {
			marketData[k] = v
		}
		s.orderBookMutex.RUnlock()

		// Skip if no data or no clients
		if len(marketData) == 0 {
			continue
		}

		s.clientsMutex.Lock()
		if len(s.clients) == 0 {
			s.clientsMutex.Unlock()
			continue
		}

		// Create message
		msg := WebSocketMessage{
			Type:      "market",
			Data:      marketData,
			Timestamp: time.Now().Unix(),
		}

		// Send to all clients
		for client := range s.clients {
			if err := client.WriteJSON(msg); err != nil {
				log.Errorf("Failed to send market data: %v", err)
				client.Close()
				delete(s.clients, client)
			}
		}
		s.clientsMutex.Unlock()
	}
}

// broadcastOpportunity broadcasts an arbitrage opportunity to all connected clients
func (s *WebServer) broadcastOpportunity(opp models.ArbitrageOpportunity) {
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()

	// Skip if no clients
	if len(s.clients) == 0 {
		return
	}

	// Create message
	msg := WebSocketMessage{
		Type:      "opportunity",
		Data:      opp,
		Timestamp: time.Now().Unix(),
	}

	// Send to all clients
	for client := range s.clients {
		if err := client.WriteJSON(msg); err != nil {
			log.Errorf("Failed to send opportunity data: %v", err)
			client.Close()
			delete(s.clients, client)
		}
	}
}

// handleOpportunitiesAPI handles API requests for arbitrage opportunities
func (s *WebServer) handleOpportunitiesAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	s.opportunitiesMutex.Lock()
	opportunities := make([]models.ArbitrageOpportunity, len(s.opportunities))
	copy(opportunities, s.opportunities)
	s.opportunitiesMutex.Unlock()

	if err := json.NewEncoder(w).Encode(opportunities); err != nil {
		log.Errorf("Failed to encode opportunities: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// handleMarketAPI handles API requests for market data
func (s *WebServer) handleMarketAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	s.orderBookMutex.RLock()
	marketData := make(map[string]*models.OrderBook)
	for k, v := range s.orderBooks {
		marketData[k] = v
	}
	s.orderBookMutex.RUnlock()

	if err := json.NewEncoder(w).Encode(marketData); err != nil {
		log.Errorf("Failed to encode market data: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}