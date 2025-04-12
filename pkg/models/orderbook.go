package models

import (
        "strconv"
        "time"
)

// OrderBook represents the current best bid and ask prices for an exchange
// @author VrushankPatel
// @description Struct representing the order book data for a specific trading pair on an exchange
type OrderBook struct {
        Exchange      string    `json:"exchange"`      // Name of the exchange (e.g., "Binance", "Kraken")
        Symbol        string    `json:"symbol"`        // Trading symbol in exchange-specific format (e.g., "BTCUSDT", "XBT/USDT")
        BaseCurrency  string    `json:"base_currency"` // Base currency of the trading pair (e.g., "BTC")
        QuoteCurrency string    `json:"quote_currency"`// Quote currency of the trading pair (e.g., "USDT")
        Bid           float64   `json:"bid"`           // Current best bid price (highest buy offer)
        Ask           float64   `json:"ask"`           // Current best ask price (lowest sell offer)
        LastUpdate    time.Time `json:"last_update"`   // Timestamp of the last update to this order book
}

// ArbitrageOpportunity represents a potential arbitrage opportunity between exchanges
// @author VrushankPatel
// @description Struct representing an arbitrage opportunity detected between two exchanges
type ArbitrageOpportunity struct {
        Timestamp        time.Time `json:"timestamp"`       // Time when the opportunity was detected
        BaseCurrency     string    `json:"base_currency"`   // Base currency of the trading pair (e.g., "BTC")
        QuoteCurrency    string    `json:"quote_currency"`  // Quote currency of the trading pair (e.g., "USDT")
        BuyExchange      string    `json:"buy_exchange"`    // Exchange where the asset should be bought
        SellExchange     string    `json:"sell_exchange"`   // Exchange where the asset should be sold
        BuyPrice         float64   `json:"buy_price"`       // Price to buy at on the buy exchange
        SellPrice        float64   `json:"sell_price"`      // Price to sell at on the sell exchange
        ProfitPercentage float64   `json:"profit_percentage"`// Profit as a percentage (e.g., 1.5 means 1.5%)
        NetProfit        float64   `json:"net_profit"`      // Net profit in quote currency (e.g., USDT)
}

// TradingPair represents a cryptocurrency trading pair
// @author VrushankPatel
// @description Struct representing a trading pair in a standardized format
type TradingPair struct {
        BaseCurrency  string // The base currency (e.g., BTC, ETH)
        QuoteCurrency string // The quote currency (e.g., USDT, USD)
}

// GetSymbol returns the formatted symbol for a trading pair based on the exchange format
// @author VrushankPatel
// @description Converts a standard trading pair to the specific format required by different exchanges
// @param exchange The name of the exchange (e.g., "Binance", "Kraken")
// @return The trading pair formatted according to the exchange's requirements
func (tp TradingPair) GetSymbol(exchange string) string {
        switch exchange {
        case "Binance":
                return tp.BaseCurrency + tp.QuoteCurrency // BTCUSDT
        case "Kraken":
                base := tp.BaseCurrency
                if base == "BTC" {
                        base = "XBT" // Kraken uses XBT instead of BTC
                }
                return base + "/" + tp.QuoteCurrency // XBT/USDT
        case "Coinbase":
                return tp.BaseCurrency + "-" + tp.QuoteCurrency // BTC-USDT
        default:
                return tp.BaseCurrency + "/" + tp.QuoteCurrency // Default format
        }
}

// ParseFloat safely parses a string to float64
// @author VrushankPatel
// @description Utility function that safely converts a string to a float64 value
// @param s The string to parse
// @return The parsed float64 value and an error if parsing fails
func ParseFloat(s string) (float64, error) {
        return strconv.ParseFloat(s, 64)
}
