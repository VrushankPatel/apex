package models

import (
	"strconv"
	"time"
)

// OrderBook represents the current best bid and ask prices for an exchange
type OrderBook struct {
	Exchange   string    `json:"exchange"`
	Symbol     string    `json:"symbol"`
	Bid        float64   `json:"bid"`
	Ask        float64   `json:"ask"`
	LastUpdate time.Time `json:"last_update"`
}

// ArbitrageOpportunity represents a potential arbitrage opportunity between exchanges
type ArbitrageOpportunity struct {
	Timestamp        time.Time `json:"timestamp"`
	BuyExchange      string    `json:"buy_exchange"`
	SellExchange     string    `json:"sell_exchange"`
	BuyPrice         float64   `json:"buy_price"`
	SellPrice        float64   `json:"sell_price"`
	ProfitPercentage float64   `json:"profit_percentage"`
	NetProfit        float64   `json:"net_profit"`
}

// ParseFloat safely parses a string to float64
func ParseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}
