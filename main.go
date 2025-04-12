package main

import (
        "context"
        "os"
        "os/signal"
        "strings"
        "sync"
        "syscall"
        "time"

        "apex-arbitrage/pkg/config"
        "apex-arbitrage/pkg/detector"
        "apex-arbitrage/pkg/exchanges"
        "apex-arbitrage/pkg/models"
        "apex-arbitrage/pkg/server"
        "apex-arbitrage/pkg/util"

        log "github.com/sirupsen/logrus"
)

func main() {
        // Initialize logger
        util.InitLogger()

        // Load config
        cfg, err := config.LoadConfig()
        if err != nil {
                log.Fatalf("Failed to load config: %v", err)
        }

        // Initialize order book map to store data from exchanges
        orderBooks := make(map[string]*models.OrderBook)
        orderBookMutex := &sync.RWMutex{}

        // Create context with cancellation for graceful shutdown
        ctx, cancel := context.WithCancel(context.Background())
        defer cancel()

        // Setup signal handling for graceful shutdown
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

        // Create a waitgroup to coordinate goroutines
        var wg sync.WaitGroup

        // Initialize exchanges
        exchangeClients := []exchanges.Exchange{}

        // Default trading pair if needed
        defaultPair := "BTCUSDT"
        if len(cfg.TradingPairs) > 0 {
                // Use the first trading pair by default
                defaultPair = cfg.TradingPairs[0].BaseCurrency + cfg.TradingPairs[0].QuoteCurrency
        }

        // Add binance
        if cfg.Exchanges.Binance.Enabled {
                binance, err := exchanges.NewBinance(defaultPair)
                if err != nil {
                        log.Fatalf("Failed to initialize Binance: %v", err)
                }
                exchangeClients = append(exchangeClients, binance)
        }

        // Add kraken
        if cfg.Exchanges.Kraken.Enabled {
                kraken, err := exchanges.NewKraken(defaultPair)
                if err != nil {
                        log.Fatalf("Failed to initialize Kraken: %v", err)
                }
                exchangeClients = append(exchangeClients, kraken)
        }

        // Start exchange websocket connections
        for _, exchange := range exchangeClients {
                wg.Add(1)
                go func(e exchanges.Exchange) {
                        defer wg.Done()
                        e.Connect(ctx, orderBooks, orderBookMutex)
                }(exchange)
        }

        // Initialize arbitrage detector
        arb := detector.NewArbitrageDetector(
                orderBooks,
                orderBookMutex,
                cfg.MinProfitThreshold,
                cfg.Exchanges.Binance.TakerFee,
                cfg.Exchanges.Kraken.TakerFee,
        )

        // Start arbitrage detection loop
        wg.Add(1)
        go func() {
                defer wg.Done()
                arb.Start(ctx)
        }()

        // Print header for the dashboard
        log.Info("APEX: Arbitrage Profit EXplorer Started")
        log.Info("-----------------------------------------")
        
        // Display the trading pairs we're monitoring
        pairs := []string{}
        for _, pair := range cfg.TradingPairs {
            pairs = append(pairs, pair.BaseCurrency+"/"+pair.QuoteCurrency)
        }
        log.Infof("Monitoring pairs: %s", strings.Join(pairs, ", "))
        log.Infof("Min profit threshold: %.2f%%", cfg.MinProfitThreshold*100)
        log.Info("Press Ctrl+C to exit")
        log.Info("--------------------------------------")

        // Initialize and start web server (on port 8080)
        webServer := server.NewWebServer("8080", orderBooks, orderBookMutex)
        
        // Register the opportunity handler to receive detected opportunities
        arb.RegisterOpportunityHandler(func(opp models.ArbitrageOpportunity) {
            webServer.AddOpportunity(opp)
        })
        
        // Start web server in a goroutine
        wg.Add(1)
        go func() {
            defer wg.Done()
            if err := webServer.Start(); err != nil {
                log.Errorf("Web server error: %v", err)
            }
        }()
        
        log.Info("Web UI available at http://localhost:8080")

        // Block until we receive a termination signal
        <-sigChan
        log.Info("Shutdown signal received, closing connections...")
        
        // Cancel the context to notify all goroutines to shut down
        cancel()
        
        // Wait for all goroutines to finish with a timeout
        waitCh := make(chan struct{})
        go func() {
                wg.Wait()
                close(waitCh)
        }()
        
        select {
        case <-waitCh:
                log.Info("Graceful shutdown completed")
        case <-time.After(5 * time.Second):
                log.Warn("Shutdown timed out, forcing exit")
        }
}
