# APEX API Documentation

## Overview

APEX provides a comprehensive set of APIs for integrating with cryptocurrency exchanges and managing arbitrage detection. This document outlines the available APIs, their usage, and implementation details.

## WebSocket API

### Connection

```
ws://localhost:8080/ws
```

The WebSocket API provides real-time updates for:
- Order book changes
- Arbitrage opportunities
- Market statistics

### Message Types

#### Market Data Update
```json
{
  "type": "market_update",
  "data": {
    "exchange": "string",
    "symbol": "string",
    "bid": "float",
    "ask": "float",
    "timestamp": "ISO8601"
  }
}
```

#### Arbitrage Opportunity
```json
{
  "type": "opportunity",
  "data": {
    "timestamp": "ISO8601",
    "base_currency": "string",
    "quote_currency": "string",
    "buy_exchange": "string",
    "sell_exchange": "string",
    "buy_price": "float",
    "sell_price": "float",
    "profit_percentage": "float",
    "net_profit": "float"
  }
}
```

## REST API

### Base URL
```
http://localhost:8080/api
```

### Endpoints

#### Get Market Status
```
GET /status
```

Response:
```json
{
  "exchanges": [
    {
      "name": "string",
      "connected": "boolean",
      "last_update": "ISO8601"
    }
  ]
}
```

#### Get Recent Opportunities
```
GET /opportunities
```

Query Parameters:
- `limit` (optional): Number of opportunities to return (default: 100)
- `min_profit` (optional): Minimum profit percentage filter

Response:
```json
{
  "opportunities": [
    {
      "timestamp": "ISO8601",
      "buy_exchange": "string",
      "sell_exchange": "string",
      "buy_price": "float",
      "sell_price": "float",
      "profit_percentage": "float",
      "net_profit": "float"
    }
  ]
}
```

#### Get Exchange Configuration
```
GET /config/exchanges
```

Response:
```json
{
  "exchanges": [
    {
      "name": "string",
      "enabled": "boolean",
      "taker_fee": "float",
      "trading_pairs": ["string"]
    }
  ]
}
```

## Exchange Integration API

### Interface Definition

The `Exchange` interface that all exchange implementations must satisfy:

```go
type Exchange interface {
    Name() string
    Connect(ctx context.Context, orderBooks map[string]*models.OrderBook, mu *sync.RWMutex)
    GetOrderBook() *models.OrderBook
    Close() error
    GetFormattedSymbol(pair models.TradingPair) string
    GetTakerFee() float64
}
```

### Implementing New Exchanges

To add support for a new exchange:

1. Create a new file in `pkg/exchanges`
2. Implement the `Exchange` interface
3. Add configuration options in `pkg/config`
4. Register the exchange in `main.go`

Example implementation structure:

```go
type NewExchange struct {
    BaseExchange
    wsURL  string
    symbol string
    conn   *websocket.Conn
}

func NewExchangeClient(symbol string) (*NewExchange, error) {
    // Implementation
}

func (e *NewExchange) Connect(ctx context.Context, orderBooks map[string]*models.OrderBook, mu *sync.RWMutex) {
    // Implementation
}
```

## Error Handling

All API endpoints follow a standard error response format:

```json
{
  "error": {
    "code": "string",
    "message": "string",
    "details": "object"
  }
}
```

Common error codes:
- `invalid_request`: Malformed request
- `not_found`: Resource not found
- `exchange_error`: Exchange-specific error
- `internal_error`: Server-side error

## Rate Limiting

- WebSocket: No rate limits on subscribed data
- REST API: 100 requests per minute per IP
- Exchange APIs: Respects exchange-specific rate limits

## Security

### Authentication

Currently, the API is designed for local use and doesn't require authentication. For production deployment, implement appropriate authentication mechanisms:

```go
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Implement authentication logic
    })
}
```

### Best Practices

1. Use HTTPS in production
2. Implement API key authentication for remote access
3. Rate limit by IP and/or API key
4. Validate all input data
5. Sanitize response data

## Testing

The API includes comprehensive test suites:

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./pkg/api

# Run with coverage
go test -cover ./...
```

## Examples

### WebSocket Client Example

```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    switch(data.type) {
        case 'market_update':
            handleMarketUpdate(data);
            break;
        case 'opportunity':
            handleOpportunity(data);
            break;
    }
};
```

### REST API Example

```python
import requests

# Get market status
response = requests.get('http://localhost:8080/api/status')
status = response.json()

# Get recent opportunities with minimum profit filter
params = {'min_profit': 0.5, 'limit': 10}
response = requests.get('http://localhost:8080/api/opportunities', params=params)
opportunities = response.json()
```

## Support

For API support:
1. Check the [GitHub Issues](https://github.com/VrushankPatel/apex/issues)
2. Review the [Troubleshooting Guide](SETUP_GUIDE.md#troubleshooting)
3. Contact the development team

---

This API documentation provides a comprehensive overview of APEX's integration capabilities. For specific implementation details, refer to the source code and inline documentation. 