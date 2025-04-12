package models

import (
        "strconv"
        "time"
)

// OrderBook represents the current best bid and ask prices for an exchange
type OrderBook struct {
        Exchange      string    `json:"exchange"`
        Symbol        string    `json:"symbol"`
        BaseCurrency  string    `json:"base_currency"`
        QuoteCurrency string    `json:"quote_currency"`
        Bid           float64   `json:"bid"`
        Ask           float64   `json:"ask"`
        LastUpdate    time.Time `json:"last_update"`
}

// ArbitrageOpportunity represents a potential arbitrage opportunity between exchanges
type ArbitrageOpportunity struct {
        Timestamp        time.Time `json:"timestamp"`
        BaseCurrency     string    `json:"base_currency"`
        QuoteCurrency    string    `json:"quote_currency"`
        BuyExchange      string    `json:"buy_exchange"`
        SellExchange     string    `json:"sell_exchange"`
        BuyPrice         float64   `json:"buy_price"`
        SellPrice        float64   `json:"sell_price"`
        ProfitPercentage float64   `json:"profit_percentage"`
        NetProfit        float64   `json:"net_profit"`
}

// TradingPair represents a cryptocurrency trading pair
type TradingPair struct {
        BaseCurrency  string // The base currency (e.g., BTC, ETH)
        QuoteCurrency string // The quote currency (e.g., USDT, USD)
}

// GetSymbol returns the formatted symbol for a trading pair based on the exchange format
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
func ParseFloat(s string) (float64, error) {
        return strconv.ParseFloat(s, 64)
}
