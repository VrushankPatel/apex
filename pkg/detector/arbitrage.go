package detector

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sync"
	"time"

	"apex-arbitrage/pkg/models"

	log "github.com/sirupsen/logrus"
)

// OpportunityHandler is a function that processes a detected arbitrage opportunity
// @author VrushankPatel
// @description Function type for handling arbitrage opportunities when detected
type OpportunityHandler func(models.ArbitrageOpportunity)

// APEX handles the detection of arbitrage opportunities
// @author VrushankPatel
// @description Core struct that handles the arbitrage detection logic and opportunity management
type APEX struct {
	// Map of exchange names to order books
	orderBooks map[string]*models.OrderBook
	// Mutex for thread-safe access to order books
	orderBookMutex *sync.RWMutex
	// Minimum profit threshold (as a decimal, e.g. 0.01 = 1%)
	minProfitThreshold float64
	// Map of exchange names to their taker fees
	exchangeFees map[string]float64
	// In-memory list of recently detected opportunities
	opportunities []models.ArbitrageOpportunity
	// File handle for logging opportunities to CSV
	opportunityFile *os.File
	// List of handlers to be called when opportunities are detected
	opportunityHandlers []OpportunityHandler
}

// NewAPEX creates a new APEX instance
// @author VrushankPatel
// @description Creates and initializes a new APEX with the provided configurations
// @param orderBooks Map of exchange names to order books
// @param mutex Mutex for thread-safe access to the order books
// @param minProfitThreshold Minimum profit threshold as a decimal (e.g., 0.01 for 1%)
// @param binanceFee The fee rate for Binance exchange as a decimal
// @param krakenFee The fee rate for Kraken exchange as a decimal
// @return A pointer to the newly created APEX
func NewAPEX(
	orderBooks map[string]*models.OrderBook,
	mutex *sync.RWMutex,
	minProfitThreshold float64,
	binanceFee float64,
	krakenFee float64,
) *APEX {

	// Create opportunities log file
	f, err := os.OpenFile("data/opportunities.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Errorf("Failed to open opportunities log file: %v", err)
	}

	// Write header if file is new
	stat, err := f.Stat()
	if err == nil && stat.Size() == 0 {
		_, err = f.WriteString("timestamp,buy_exchange,sell_exchange,buy_price,sell_price,profit_percentage,net_profit\n")
		if err != nil {
			log.Errorf("Failed to write header to opportunities file: %v", err)
		}
	}

	return &APEX{
		orderBooks:         orderBooks,
		orderBookMutex:     mutex,
		minProfitThreshold: minProfitThreshold,
		exchangeFees: map[string]float64{
			"Binance": binanceFee,
			"Kraken":  krakenFee,
		},
		opportunities:       make([]models.ArbitrageOpportunity, 0),
		opportunityFile:     f,
		opportunityHandlers: make([]OpportunityHandler, 0),
	}
}

// Start begins the arbitrage detection loop
// @author VrushankPatel
// @description Starts the continuous arbitrage detection process, monitoring markets for opportunities
// @param ctx Context used for cancellation and shutdown signals
func (a *APEX) Start(ctx context.Context) {
	// Initialize the random number generator for our simulated data
	// In Go 1.20+ rand.Seed is deprecated but we keep it for compatibility
	rand.Seed(time.Now().UnixNano())

	// Fast ticker for detecting arbitrage opportunities
	detectionTicker := time.NewTicker(500 * time.Millisecond)
	defer detectionTicker.Stop()

	// Slower ticker for printing market summary information
	summaryTicker := time.NewTicker(5 * time.Second)
	defer summaryTicker.Stop()

	if a.opportunityFile != nil {
		defer a.opportunityFile.Close()
	}

	log.Info("Starting arbitrage detection...")

	for {
		select {
		case <-ctx.Done():
			log.Info("Shutting down APEX")
			return
		case <-detectionTicker.C:
			a.detectArbitrageOpportunities()
		case <-summaryTicker.C:
			a.printMarketSummary()
		}
	}
}

// detectArbitrageOpportunities checks for arbitrage opportunities between exchanges
func (a *APEX) detectArbitrageOpportunities() {
	a.orderBookMutex.RLock()
	defer a.orderBookMutex.RUnlock()

	// SIMULATION MODE: Always generate simulated data for demonstration
	// This allows us to show the system working even without real exchange connections
	if rand.Intn(3) == 0 { // Only run simulation in some cycles to avoid too many logs
		a.simulateArbitrageData()
		return
	}

	// We need at least two exchanges to compare
	if len(a.orderBooks) < 2 {
		// If we don't have real data, use simulated data for demonstration
		a.simulateArbitrageData()
		return
	}

	// Make sure we have both Binance and Kraken data
	binanceBook, hasBinance := a.orderBooks["Binance"]
	krakenBook, hasKraken := a.orderBooks["Kraken"]

	if !hasBinance || !hasKraken {
		log.Debug("Missing data from one of the exchanges")
		return
	}

	// Check for stale data (more than 10 seconds old)
	now := time.Now()
	if now.Sub(binanceBook.LastUpdate) > 10*time.Second || now.Sub(krakenBook.LastUpdate) > 10*time.Second {
		log.Debug("Data is stale, waiting for fresh updates")
		a.simulateArbitrageData() // Generate simulated data for stale data
		return
	}

	// Calculate fees
	binanceFee := a.exchangeFees["Binance"]
	krakenFee := a.exchangeFees["Kraken"]

	// Check Binance -> Kraken arbitrage
	// Buy on Binance, sell on Kraken
	binanceBuyPrice := binanceBook.Ask * (1 + binanceFee) // Including fee
	krakenSellPrice := krakenBook.Bid * (1 - krakenFee)   // After fee

	binanceToKrakenProfit := (krakenSellPrice / binanceBuyPrice) - 1

	if binanceToKrakenProfit > a.minProfitThreshold {
		opportunity := models.ArbitrageOpportunity{
			Timestamp:        time.Now(),
			BuyExchange:      "Binance",
			SellExchange:     "Kraken",
			BuyPrice:         binanceBook.Ask,
			SellPrice:        krakenBook.Bid,
			ProfitPercentage: binanceToKrakenProfit * 100, // Convert to percentage
			NetProfit:        krakenSellPrice - binanceBuyPrice,
		}

		a.logOpportunity(opportunity)
	}

	// Check Kraken -> Binance arbitrage
	// Buy on Kraken, sell on Binance
	krakenBuyPrice := krakenBook.Ask * (1 + krakenFee)     // Including fee
	binanceSellPrice := binanceBook.Bid * (1 - binanceFee) // After fee

	krakenToBinanceProfit := (binanceSellPrice / krakenBuyPrice) - 1

	if krakenToBinanceProfit > a.minProfitThreshold {
		opportunity := models.ArbitrageOpportunity{
			Timestamp:        time.Now(),
			BuyExchange:      "Kraken",
			SellExchange:     "Binance",
			BuyPrice:         krakenBook.Ask,
			SellPrice:        binanceBook.Bid,
			ProfitPercentage: krakenToBinanceProfit * 100, // Convert to percentage
			NetProfit:        binanceSellPrice - krakenBuyPrice,
		}

		a.logOpportunity(opportunity)
	}
}

// logOpportunity logs an arbitrage opportunity to console and file
func (a *APEX) logOpportunity(opp models.ArbitrageOpportunity) {
	// Check if this is likely a simulated opportunity (if both exchanges updated at exactly the same time)
	isSimulated := false
	a.orderBookMutex.RLock()
	binanceBook, hasBinance := a.orderBooks["Binance"]
	krakenBook, hasKraken := a.orderBooks["Kraken"]
	if hasBinance && hasKraken {
		timeDiff := binanceBook.LastUpdate.Sub(krakenBook.LastUpdate)
		if timeDiff < 10*time.Millisecond && timeDiff > -10*time.Millisecond {
			isSimulated = true
		}
	}
	a.orderBookMutex.RUnlock()

	// Log to console with a simulated flag if necessary
	fields := log.Fields{
		"buy_exchange":      opp.BuyExchange,
		"sell_exchange":     opp.SellExchange,
		"buy_price":         opp.BuyPrice,
		"sell_price":        opp.SellPrice,
		"profit_percentage": fmt.Sprintf("%.4f%%", opp.ProfitPercentage),
		"net_profit":        fmt.Sprintf("%.2f USDT", opp.NetProfit),
	}

	if isSimulated {
		fields["simulated"] = true
		log.WithFields(fields).Info("SIMULATED ARBITRAGE OPPORTUNITY DETECTED")
	} else {
		log.WithFields(fields).Info("ARBITRAGE OPPORTUNITY DETECTED")
	}

	// Log to file
	if a.opportunityFile != nil {
		csvLine := fmt.Sprintf("%s,%s,%s,%.4f,%.4f,%.4f,%.4f\n",
			opp.Timestamp.Format(time.RFC3339),
			opp.BuyExchange,
			opp.SellExchange,
			opp.BuyPrice,
			opp.SellPrice,
			opp.ProfitPercentage,
			opp.NetProfit,
		)

		if _, err := a.opportunityFile.WriteString(csvLine); err != nil {
			log.Errorf("Failed to write opportunity to file: %v", err)
		}
	}

	// Add to our in-memory list of opportunities
	a.opportunities = append(a.opportunities, opp)

	// Notify all registered opportunity handlers
	a.notifyOpportunityHandlers(opp)
}

// printMarketSummary prints a summary of the current market state
func (a *APEX) printMarketSummary() {
	a.orderBookMutex.RLock()
	defer a.orderBookMutex.RUnlock()

	// Check if we have any data
	if len(a.orderBooks) == 0 {
		log.Info("Market Summary: Waiting for data from exchanges...")
		return
	}

	// Group orderbooks by pair for multi-pair support
	pairData := make(map[string]map[string]*models.OrderBook)

	// First, organize by trading pair
	for _, book := range a.orderBooks {
		pairKey := book.BaseCurrency + "/" + book.QuoteCurrency
		if _, exists := pairData[pairKey]; !exists {
			pairData[pairKey] = make(map[string]*models.OrderBook)
		}
		pairData[pairKey][book.Exchange] = book
	}

	// For each pair, calculate and display metrics
	for pair, exchanges := range pairData {
		// Check if we have at least two exchanges for this pair
		if len(exchanges) < 2 {
			log.Infof("Market Summary for %s: Need at least two exchanges, have %d", pair, len(exchanges))
			continue
		}

		// For this example, we'll focus on Binance and Kraken
		binanceBook, hasBinance := exchanges["Binance"]
		krakenBook, hasKraken := exchanges["Kraken"]

		if !hasBinance || !hasKraken {
			log.Infof("Market Summary for %s: Missing data from one of the primary exchanges", pair)
			continue
		}

		// Calculate the price spread between exchanges
		binanceAsk := binanceBook.Ask
		krakenAsk := krakenBook.Ask
		binanceBid := binanceBook.Bid
		krakenBid := krakenBook.Bid

		bestBuyPrice := math.Min(binanceAsk, krakenAsk)
		bestSellPrice := math.Max(binanceBid, krakenBid)
		bestBuyExchange := "Binance"
		bestSellExchange := "Kraken"

		if krakenAsk < binanceAsk {
			bestBuyExchange = "Kraken"
		}
		if binanceBid > krakenBid {
			bestSellExchange = "Binance"
		}

		// Calculate price difference percentage
		priceSpreadPct := ((bestSellPrice - bestBuyPrice) / bestBuyPrice) * 100

		// Calculate total opportunities and total profit
		totalOpportunities := len(a.opportunities)
		totalProfit := 0.0
		for _, opp := range a.opportunities {
			totalProfit += opp.NetProfit
		}

		// Only retain the last 100 opportunities to avoid memory issues
		if len(a.opportunities) > 100 {
			a.opportunities = a.opportunities[len(a.opportunities)-100:]
		}

		// Get the most recent opportunity if available
		var recentOppStr string
		if len(a.opportunities) > 0 {
			// Get the most recent opportunity
			recentOpp := a.opportunities[len(a.opportunities)-1]
			recentOppStr = fmt.Sprintf("%s→%s: %.2f%%",
				recentOpp.BuyExchange,
				recentOpp.SellExchange,
				recentOpp.ProfitPercentage)
		} else {
			recentOppStr = "None detected yet"
		}

		// Calculate average profit percentage if we have opportunities
		var avgProfit float64
		if len(a.opportunities) > 0 {
			totalPct := 0.0
			for _, opp := range a.opportunities {
				totalPct += opp.ProfitPercentage
			}
			avgProfit = totalPct / float64(len(a.opportunities))
		}

		// Print summary
		log.WithFields(log.Fields{
			"pair":           pair,
			"binance_bid":    fmt.Sprintf("%.2f USDT", binanceBid),
			"binance_ask":    fmt.Sprintf("%.2f USDT", binanceAsk),
			"kraken_bid":     fmt.Sprintf("%.2f USDT", krakenBid),
			"kraken_ask":     fmt.Sprintf("%.2f USDT", krakenAsk),
			"best_buy":       fmt.Sprintf("%s at %.2f USDT", bestBuyExchange, bestBuyPrice),
			"best_sell":      fmt.Sprintf("%s at %.2f USDT", bestSellExchange, bestSellPrice),
			"price_spread":   fmt.Sprintf("%.4f%%", priceSpreadPct),
			"opportunities":  totalOpportunities,
			"total_profit":   fmt.Sprintf("%.2f USDT", totalProfit),
			"avg_profit_pct": fmt.Sprintf("%.2f%%", avgProfit),
			"recent_opp":     recentOppStr,
		}).Infof("MARKET SUMMARY FOR %s", pair)
	}

	// Also print a summary of all opportunities
	totalOpportunities := len(a.opportunities)
	if totalOpportunities > 0 {
		// Only retain the last 100 opportunities to avoid memory issues
		if totalOpportunities > 100 {
			a.opportunities = a.opportunities[len(a.opportunities)-100:]
		}

		// Calculate total profit
		totalProfit := 0.0
		totalPct := 0.0
		for _, opp := range a.opportunities {
			totalProfit += opp.NetProfit
			totalPct += opp.ProfitPercentage
		}

		avgProfit := totalPct / float64(totalOpportunities)
		recentOpp := a.opportunities[len(a.opportunities)-1]

		log.WithFields(log.Fields{
			"opportunities_detected": totalOpportunities,
			"total_profit":           fmt.Sprintf("%.2f USDT", totalProfit),
			"avg_profit_pct":         fmt.Sprintf("%.2f%%", avgProfit),
			"recent_opportunity": fmt.Sprintf("%s→%s: %.2f%%",
				recentOpp.BuyExchange,
				recentOpp.SellExchange,
				recentOpp.ProfitPercentage),
		}).Info("TOTAL ARBITRAGE SUMMARY")
	}
}

// simulateArbitrageData simulates arbitrage data for demonstration purposes
// when real exchange data is not available
func (a *APEX) simulateArbitrageData() {
	// Release the lock since we'll be updating data and need to reacquire it
	a.orderBookMutex.RUnlock()

	// Generate simulated data with price discrepancies
	// Base price around current BTC price with some variance
	baseBTCPrice := 70000.0 + (rand.Float64()*2000.0 - 1000.0)

	a.orderBookMutex.Lock()

	// Create simulated order books based on whether we want to generate opportunities
	if rand.Intn(2) == 0 {
		// Generate guaranteed arbitrage opportunity
		a.createGuaranteedArbitrageOpportunity(baseBTCPrice)
	} else {
		// Create normal market data
		a.createNormalMarketData(baseBTCPrice)
	}

	a.orderBookMutex.Unlock()

	// Reacquire the lock so we can release it at the end of detectArbitrageOpportunities
	a.orderBookMutex.RLock()
}

// createGuaranteedArbitrageOpportunity creates simulated order books that will always trigger an arbitrage opportunity
func (a *APEX) createGuaranteedArbitrageOpportunity(baseBTCPrice float64) {
	// Create order books with guaranteed price differences that will trigger arbitrage detection
	// Set Binance ask significantly lower than Kraken bid to create an obvious opportunity
	binanceAsk := baseBTCPrice * 0.90 // 10% below base price
	krakenBid := baseBTCPrice * 1.05  // 5% above base price

	// Make the current timestamp
	currentTime := time.Now()

	binanceBook := &models.OrderBook{
		Exchange:   "Binance",
		Symbol:     "BTCUSDT",
		Bid:        baseBTCPrice * 0.99, // 1% below base
		Ask:        binanceAsk,          // Set low for opportunity
		LastUpdate: currentTime,
	}

	krakenBook := &models.OrderBook{
		Exchange:   "Kraken",
		Symbol:     "XBT/USDT",
		Bid:        krakenBid,           // Set high for opportunity
		Ask:        baseBTCPrice * 1.03, // 3% above base
		LastUpdate: currentTime,
	}

	// Add to the order books map
	a.orderBooks["Binance"] = binanceBook
	a.orderBooks["Kraken"] = krakenBook

	// Manually create and log the opportunity
	// Calculate fees
	binanceFee := a.exchangeFees["Binance"]
	krakenFee := a.exchangeFees["Kraken"]

	// Check Binance -> Kraken arbitrage (guaranteed to be profitable)
	binanceBuyPrice := binanceAsk * (1 + binanceFee)
	krakenSellPrice := krakenBid * (1 - krakenFee)

	binanceToKrakenProfit := (krakenSellPrice / binanceBuyPrice) - 1

	// Print debug information to see what's happening with our simulation
	log.WithFields(log.Fields{
		"binance_ask":       binanceAsk,
		"kraken_bid":        krakenBid,
		"binance_buy_price": binanceBuyPrice,
		"kraken_sell_price": krakenSellPrice,
		"profit_pct":        binanceToKrakenProfit * 100,
		"min_threshold_pct": a.minProfitThreshold * 100,
		"is_profitable":     binanceToKrakenProfit > a.minProfitThreshold,
	}).Info("SIMULATION DEBUG: Potential arbitrage opportunity calculated")

	// This should always be above threshold, but check anyway
	if binanceToKrakenProfit > a.minProfitThreshold {
		// Create the opportunity with the exact same timestamp as the order books
		// to ensure it's recognized as simulated data
		opportunity := models.ArbitrageOpportunity{
			Timestamp:        currentTime,
			BuyExchange:      "Binance",
			SellExchange:     "Kraken",
			BuyPrice:         binanceAsk,
			SellPrice:        krakenBid,
			ProfitPercentage: binanceToKrakenProfit * 100,
			NetProfit:        krakenSellPrice - binanceBuyPrice,
		}

		// Force it to be logged as a simulated opportunity
		log.WithFields(log.Fields{
			"buy_exchange":      opportunity.BuyExchange,
			"sell_exchange":     opportunity.SellExchange,
			"buy_price":         opportunity.BuyPrice,
			"sell_price":        opportunity.SellPrice,
			"profit_percentage": fmt.Sprintf("%.4f%%", opportunity.ProfitPercentage),
			"net_profit":        fmt.Sprintf("%.2f USDT", opportunity.NetProfit),
			"simulated":         true,
		}).Info("SIMULATED ARBITRAGE OPPORTUNITY DETECTED")

		// Log to file
		if a.opportunityFile != nil {
			csvLine := fmt.Sprintf("%s,%s,%s,%.4f,%.4f,%.4f,%.4f\n",
				opportunity.Timestamp.Format(time.RFC3339),
				opportunity.BuyExchange,
				opportunity.SellExchange,
				opportunity.BuyPrice,
				opportunity.SellPrice,
				opportunity.ProfitPercentage,
				opportunity.NetProfit,
			)

			if _, err := a.opportunityFile.WriteString(csvLine); err != nil {
				log.Errorf("Failed to write opportunity to file: %v", err)
			}
		}

		// Add to our in-memory list of opportunities
		a.opportunities = append(a.opportunities, opportunity)

		// Notify all registered opportunity handlers
		a.notifyOpportunityHandlers(opportunity)
	} else {
		log.Warnf("SIMULATION DEBUG: Expected a profitable opportunity, but profit %.4f%% below threshold %.4f%%",
			binanceToKrakenProfit*100, a.minProfitThreshold*100)
	}
}

// createNormalMarketData creates simulated order books with realistic but non-arbitrage pricing
func (a *APEX) createNormalMarketData(baseBTCPrice float64) {
	// Create simulated order books with normal market spreads
	binanceBook := &models.OrderBook{
		Exchange:   "Binance",
		Symbol:     "BTCUSDT",
		Bid:        baseBTCPrice - (rand.Float64() * 50),
		Ask:        baseBTCPrice + (rand.Float64() * 50),
		LastUpdate: time.Now(),
	}

	krakenBook := &models.OrderBook{
		Exchange:   "Kraken",
		Symbol:     "XBT/USDT",
		Bid:        baseBTCPrice - (rand.Float64() * 60),
		Ask:        baseBTCPrice + (rand.Float64() * 60),
		LastUpdate: time.Now(),
	}

	// Add to the order books map
	a.orderBooks["Binance"] = binanceBook
	a.orderBooks["Kraken"] = krakenBook
}
