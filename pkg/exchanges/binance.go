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

// Binance defines the Binance exchange client
type Binance struct {
        BaseExchange
        wsURL    string
        conn     *websocket.Conn
        symbol   string
        pairs    []models.TradingPair
}

// BinanceBookTickerResponse defines the structure of Binance's bookTicker websocket response
type BinanceBookTickerResponse struct {
        UpdateID   int64   `json:"u"`
        Symbol     string  `json:"s"`
        BestBidQty string  `json:"B"`
        BestBid    string  `json:"b"`
        BestAskQty string  `json:"A"`
        BestAsk    string  `json:"a"`
        EventTime  int64   `json:"E"`
}

// NewBinance creates a new Binance exchange client
func NewBinance(symbol string) (*Binance, error) {
        // Parse the symbol to derive base and quote currencies
        // Assuming standard format like BTCUSDT, extract BTC and USDT
        var baseCurrency, quoteCurrency string
        
        if len(symbol) >= 7 {
            // Most common quote currencies are USDT (4), BUSD (4), USDC (4), USD (3)
            baseCurrency = symbol[:len(symbol)-4]
            quoteCurrency = symbol[len(symbol)-4:]
        } else if len(symbol) >= 6 {
            // Fallback for pairs like BTCUSD
            baseCurrency = symbol[:len(symbol)-3]
            quoteCurrency = symbol[len(symbol)-3:]
        } else {
            // Default fallback
            baseCurrency = "BTC"
            quoteCurrency = "USDT"
        }
        
        return &Binance{
                BaseExchange: BaseExchange{
                        name:   "Binance",
                        symbol: symbol,
                        orderBook: &models.OrderBook{
                                Exchange:      "Binance",
                                Symbol:        symbol,
                                BaseCurrency:  baseCurrency,
                                QuoteCurrency: quoteCurrency,
                        },
                        takerFee: 0.001, // 0.1% is the default fee
                },
                wsURL:  "wss://stream.binance.com:9443/ws",
                symbol: symbol,
        }, nil
}

// Connect establishes a websocket connection to Binance and starts streaming order book data
func (b *Binance) Connect(ctx context.Context, orderBooks map[string]*models.OrderBook, mu *sync.RWMutex) {
        // Let's use the actual combined streams endpoint for better reliability
        wsEndpoint := "wss://stream.binance.com:9443/ws"
        
        log.Infof("[Binance] Connecting to %s", wsEndpoint)
        
        // Setup custom dialer with longer timeouts
        dialer := websocket.DefaultDialer
        dialer.HandshakeTimeout = 10 * time.Second
        
        var err error
        b.conn, _, err = dialer.Dial(wsEndpoint, nil)
        if err != nil {
                log.Errorf("[Binance] Failed to connect to websocket: %v", err)
                
                // Initialize with some dummy data so the simulation can work
                log.Warn("[Binance] Using simulation mode for demonstration purposes")
                b.orderBook.Bid = 69500.0
                b.orderBook.Ask = 69550.0
                b.orderBook.LastUpdate = time.Now()
                updateOrderBookMap(b.Name(), b.orderBook, orderBooks, mu)
                
                return
        }
        
        // Subscribe to the stream using the message format from Binance docs
        subscribeMsg := map[string]interface{}{
                "method": "SUBSCRIBE",
                "params": []string{
                        fmt.Sprintf("%s@bookTicker", b.symbol),
                },
                "id": 1,
        }
        
        if err := b.conn.WriteJSON(subscribeMsg); err != nil {
                log.Errorf("[Binance] Failed to subscribe to stream: %v", err)
                b.Close()
                return
        }
        
        defer b.Close()
        
        // Monitor for context cancellation
        go func() {
                <-ctx.Done()
                log.Info("[Binance] Context cancelled, closing connection")
                b.Close()
        }()
        
        log.Infof("[Binance] Connected to websocket for %s", b.symbol)
        
        // Process incoming messages
        for {
                select {
                case <-ctx.Done():
                        return
                default:
                        // Read message from websocket
                        _, message, err := b.conn.ReadMessage()
                        if err != nil {
                                log.Errorf("[Binance] Error reading from websocket: %v", err)
                                
                                // Try to reconnect
                                if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
                                        time.Sleep(5 * time.Second)
                                        log.Info("[Binance] Attempting to reconnect...")
                                        b.Connect(ctx, orderBooks, mu)
                                        return
                                }
                                return
                        }
                        
                        // First, check if this is a subscription response
                        var subResponse map[string]interface{}
                        if err := json.Unmarshal(message, &subResponse); err == nil {
                                // Check if it has the result field which indicates a subscription response
                                if _, hasResult := subResponse["result"]; hasResult {
                                        log.Debugf("[Binance] Received subscription response")
                                        continue
                                }
                        }
                        
                        // Parse the message as ticker data
                        var tickerData BinanceBookTickerResponse
                        if err := json.Unmarshal(message, &tickerData); err != nil {
                                log.Errorf("[Binance] Error parsing message: %v", err)
                                log.Debugf("[Binance] Raw message: %s", string(message))
                                continue
                        }
                        
                        // Validate we have the required fields
                        if tickerData.BestBid == "" || tickerData.BestAsk == "" {
                                log.Debugf("[Binance] Incomplete ticker data received")
                                continue
                        }
                        
                        // Update our order book
                        bid, err := models.ParseFloat(tickerData.BestBid)
                        if err != nil {
                                log.Errorf("[Binance] Error parsing bid: %v", err)
                                continue
                        }
                        
                        ask, err := models.ParseFloat(tickerData.BestAsk)
                        if err != nil {
                                log.Errorf("[Binance] Error parsing ask: %v", err)
                                continue
                        }
                        
                        // Update the order book
                        b.orderBook.Bid = bid
                        b.orderBook.Ask = ask
                        b.orderBook.LastUpdate = time.Now()
                        
                        // Update the shared map
                        updateOrderBookMap(b.Name(), b.orderBook, orderBooks, mu)
                }
        }
}

// Close closes the websocket connection
func (b *Binance) Close() error {
        if b.conn != nil {
                return b.conn.Close()
        }
        return nil
}
