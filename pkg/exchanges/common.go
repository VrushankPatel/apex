package exchanges

import (
	"context"
	"sync"
	"time"

	"arbitrage-detector/pkg/models"
)

// Exchange defines the interface that all exchange implementations must satisfy
type Exchange interface {
	// Name returns the exchange name
	Name() string
	
	// Connect establishes a websocket connection and starts streaming order book data
	Connect(ctx context.Context, orderBooks map[string]*models.OrderBook, mu *sync.RWMutex)
	
	// GetOrderBook returns the current orderbook snapshot for the exchange
	GetOrderBook() *models.OrderBook
	
	// Close closes the websocket connection
	Close() error
}

// BaseExchange contains common fields and methods for exchanges
type BaseExchange struct {
	name       string
	symbol     string
	orderBook  *models.OrderBook
	lastUpdate time.Time
}

// Name returns the exchange name
func (b *BaseExchange) Name() string {
	return b.name
}

// GetOrderBook returns the current orderbook snapshot
func (b *BaseExchange) GetOrderBook() *models.OrderBook {
	return b.orderBook
}

// updateOrderBookMap safely updates the shared orderbook map
func updateOrderBookMap(exchange string, book *models.OrderBook, orderBooks map[string]*models.OrderBook, mu *sync.RWMutex) {
	mu.Lock()
	defer mu.Unlock()
	orderBooks[exchange] = book
}
