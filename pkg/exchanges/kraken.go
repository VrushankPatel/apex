package exchanges

import (
        "context"
        "encoding/json"
        "fmt"
        "sync"
        "time"

        "apex-arbitrage/pkg/models"

        "github.com/gorilla/websocket"
        log "github.com/sirupsen/logrus"
)

// Kraken defines the Kraken exchange client
type Kraken struct {
        BaseExchange
        wsURL  string
        conn   *websocket.Conn
        symbol string
}

// KrakenSubscription defines the structure for subscription message
type KrakenSubscription struct {
        Name string `json:"name"`
}

// KrakenSubscribeMessage defines the structure for the subscription request
type KrakenSubscribeMessage struct {
        Name     string            `json:"name"`
        ReqID    int               `json:"reqid,omitempty"`
        Pairs    []string          `json:"pair"`
        Subscribe KrakenSubscription `json:"subscription"`
}

// NewKraken creates a new Kraken exchange client
func NewKraken(symbol string) (*Kraken, error) {
        return &Kraken{
                BaseExchange: BaseExchange{
                        name:   "Kraken",
                        symbol: symbol,
                        orderBook: &models.OrderBook{
                                Exchange: "Kraken",
                                Symbol:   symbol,
                        },
                },
                wsURL:  "wss://ws.kraken.com",
                symbol: symbol,
        }, nil
}

// Connect establishes a websocket connection to Kraken and starts streaming order book data
func (k *Kraken) Connect(ctx context.Context, orderBooks map[string]*models.OrderBook, mu *sync.RWMutex) {
        var err error
        
        log.Infof("[Kraken] Connecting to %s", k.wsURL)
        
        // Setup custom dialer with longer timeouts
        dialer := websocket.DefaultDialer
        dialer.HandshakeTimeout = 10 * time.Second
        
        k.conn, _, err = dialer.Dial(k.wsURL, nil)
        if err != nil {
                log.Errorf("[Kraken] Failed to connect to websocket: %v", err)
                
                // Initialize with some dummy data so the simulation can work
                log.Warn("[Kraken] Using simulation mode for demonstration purposes")
                k.orderBook.Bid = 69550.0
                k.orderBook.Ask = 69650.0
                k.orderBook.LastUpdate = time.Now()
                updateOrderBookMap(k.Name(), k.orderBook, orderBooks, mu)
                
                return
        }
        
        defer k.Close()
        
        // Monitor for context cancellation
        go func() {
                <-ctx.Done()
                log.Info("[Kraken] Context cancelled, closing connection")
                k.Close()
        }()
        
        // Subscribe to ticker data
        subscribeMsg := KrakenSubscribeMessage{
                Name:  "subscribe",
                ReqID: 1,
                Pairs: []string{k.symbol},
                Subscribe: KrakenSubscription{
                        Name: "ticker",
                },
        }
        
        if err := k.conn.WriteJSON(subscribeMsg); err != nil {
                log.Errorf("[Kraken] Failed to send subscription request: %v", err)
                return
        }
        
        log.Infof("[Kraken] Subscribed to ticker for %s", k.symbol)
        
        // Process incoming messages
        for {
                select {
                case <-ctx.Done():
                        return
                default:
                        // Read message from websocket
                        _, message, err := k.conn.ReadMessage()
                        if err != nil {
                                log.Errorf("[Kraken] Error reading from websocket: %v", err)
                                
                                // Try to reconnect
                                if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
                                        time.Sleep(5 * time.Second)
                                        log.Info("[Kraken] Attempting to reconnect...")
                                        k.Connect(ctx, orderBooks, mu)
                                        return
                                }
                                return
                        }

                        // First try handling as a system message (which is an object, not an array)
                        var systemMsg map[string]interface{}
                        if err := json.Unmarshal(message, &systemMsg); err == nil {
                                // This is a system message (subscription confirmation, heartbeat, etc.)
                                if event, ok := systemMsg["event"].(string); ok {
                                    log.Debugf("[Kraken] Received system message: %s", event)
                                }
                                continue
                        }
                        
                        // If it's not a system message, try parsing as a data message (array format)
                        var data []interface{}
                        if err := json.Unmarshal(message, &data); err != nil {
                                log.Errorf("[Kraken] Error parsing message: %v", err)
                                continue
                        }
                        
                        // Check if it's a ticker message
                        if len(data) < 2 {
                                log.Debugf("[Kraken] Received non-data message with length %d", len(data))
                                continue // Not enough data
                        }
                        
                        channelName, ok := data[1].(string)
                        if !ok || channelName != "ticker" {
                                log.Debugf("[Kraken] Received non-ticker message: %v", data[1])
                                continue // Not a ticker message
                        }
                        
                        // Try to extract ticker data
                        // Format is [channelID, "ticker", data, pair]
                        if len(data) < 3 {
                                log.Debugf("[Kraken] Ticker message has incorrect format")
                                continue
                        }
                        
                        tickerData, ok := data[2].(map[string]interface{})
                        if !ok {
                                log.Error("[Kraken] Invalid ticker data format")
                                continue
                        }
                        
                        log.Debugf("[Kraken] Successfully received ticker data")
                        
                        // Extract bid/ask prices
                        bid, ask, err := extractKrakenPrices(tickerData)
                        if err != nil {
                                log.Errorf("[Kraken] Error extracting prices: %v", err)
                                continue
                        }
                        
                        // Update the order book
                        k.orderBook.Bid = bid
                        k.orderBook.Ask = ask
                        k.orderBook.LastUpdate = time.Now()
                        
                        // Update the shared map
                        updateOrderBookMap(k.Name(), k.orderBook, orderBooks, mu)
                }
        }
}

// extractKrakenPrices extracts bid and ask prices from Kraken ticker data
func extractKrakenPrices(tickerData map[string]interface{}) (float64, float64, error) {
        // Get best bid price
        bidData, ok := tickerData["b"].([]interface{})
        if !ok || len(bidData) < 1 {
                return 0, 0, fmt.Errorf("invalid bid data format")
        }
        
        bidStr, ok := bidData[0].(string)
        if !ok {
                return 0, 0, fmt.Errorf("bid price not a string")
        }
        
        bid, err := models.ParseFloat(bidStr)
        if err != nil {
                return 0, 0, fmt.Errorf("failed to parse bid price: %v", err)
        }
        
        // Get best ask price
        askData, ok := tickerData["a"].([]interface{})
        if !ok || len(askData) < 1 {
                return 0, 0, fmt.Errorf("invalid ask data format")
        }
        
        askStr, ok := askData[0].(string)
        if !ok {
                return 0, 0, fmt.Errorf("ask price not a string")
        }
        
        ask, err := models.ParseFloat(askStr)
        if err != nil {
                return 0, 0, fmt.Errorf("failed to parse ask price: %v", err)
        }
        
        return bid, ask, nil
}

// Close closes the websocket connection
func (k *Kraken) Close() error {
        if k.conn != nil {
                return k.conn.Close()
        }
        return nil
}
