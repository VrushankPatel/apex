package config

import (
        "os"
        "strconv"

        "github.com/joho/godotenv"
        log "github.com/sirupsen/logrus"
)

// Trading pair represents a market pair like BTC/USDT
type TradingPair struct {
        BaseCurrency  string
        QuoteCurrency string
}

// ExchangeConfig stores configuration for a specific exchange
type ExchangeConfig struct {
        Enabled  bool
        TakerFee float64
        MakerFee float64
        APIKey   string
        APISecret string
}

// ExchangesConfig holds configuration for all exchanges
type ExchangesConfig struct {
        Binance  ExchangeConfig
        Kraken   ExchangeConfig
        Coinbase ExchangeConfig
}

// Config stores all configuration of the application
type Config struct {
        // Exchange API keys
        BinanceAPIKey     string
        BinanceAPISecret  string
        KrakenAPIKey      string
        KrakenAPISecret   string
        CoinbaseAPIKey    string
        CoinbaseAPISecret string
        CoinbasePassphrase string

        // Application configuration
        SimulationMode     bool
        MinProfitThreshold float64
        LogLevel           string
        
        // Trading pairs to monitor
        TradingPairs []TradingPair
        
        // Exchange-specific configurations
        Exchanges ExchangesConfig
}

// LoadConfig reads configuration from .env file or environment variables
func LoadConfig() (*Config, error) {
        // Load .env file if it exists
        _ = godotenv.Load()

        config := &Config{
                // Exchange API keys
                BinanceAPIKey:     getEnv("BINANCE_API_KEY", ""),
                BinanceAPISecret:  getEnv("BINANCE_API_SECRET", ""),
                KrakenAPIKey:      getEnv("KRAKEN_API_KEY", ""),
                KrakenAPISecret:   getEnv("KRAKEN_API_SECRET", ""),
                CoinbaseAPIKey:    getEnv("COINBASE_API_KEY", ""),
                CoinbaseAPISecret: getEnv("COINBASE_API_SECRET", ""),
                CoinbasePassphrase: getEnv("COINBASE_PASSPHRASE", ""),

                // Application configuration
                SimulationMode:     getBoolEnv("SIMULATION_MODE", true),
                MinProfitThreshold: getFloatEnv("MIN_PROFIT_THRESHOLD", 0.1),
                LogLevel:           getEnv("LOG_LEVEL", "info"),
                
                // Default trading pairs
                TradingPairs: []TradingPair{
                        {BaseCurrency: "BTC", QuoteCurrency: "USDT"},
                        {BaseCurrency: "ETH", QuoteCurrency: "USDT"},
                },
                
                // Exchange configurations
                Exchanges: ExchangesConfig{
                        Binance: ExchangeConfig{
                                Enabled:   true,
                                TakerFee:  getFloatEnv("BINANCE_TAKER_FEE", 0.001), // 0.1%
                                MakerFee:  getFloatEnv("BINANCE_MAKER_FEE", 0.0008), // 0.08%
                                APIKey:    getEnv("BINANCE_API_KEY", ""),
                                APISecret: getEnv("BINANCE_API_SECRET", ""),
                        },
                        Kraken: ExchangeConfig{
                                Enabled:   true,
                                TakerFee:  getFloatEnv("KRAKEN_TAKER_FEE", 0.0026), // 0.26%
                                MakerFee:  getFloatEnv("KRAKEN_MAKER_FEE", 0.0016), // 0.16%
                                APIKey:    getEnv("KRAKEN_API_KEY", ""),
                                APISecret: getEnv("KRAKEN_API_SECRET", ""),
                        },
                        Coinbase: ExchangeConfig{
                                Enabled:   getBoolEnv("COINBASE_ENABLED", false),
                                TakerFee:  getFloatEnv("COINBASE_TAKER_FEE", 0.0025), // 0.25%
                                MakerFee:  getFloatEnv("COINBASE_MAKER_FEE", 0.0015), // 0.15%
                                APIKey:    getEnv("COINBASE_API_KEY", ""),
                                APISecret: getEnv("COINBASE_API_SECRET", ""),
                        },
                },
        }

        return config, nil
}

// Helper function to read an environment variable or return a default value
func getEnv(key, defaultValue string) string {
        if value, exists := os.LookupEnv(key); exists {
                return value
        }
        return defaultValue
}

// Helper function to read a boolean environment variable
func getBoolEnv(key string, defaultValue bool) bool {
        if valueStr, exists := os.LookupEnv(key); exists {
                value, err := strconv.ParseBool(valueStr)
                if err == nil {
                        return value
                }
                log.Warnf("Invalid boolean value for %s: %s", key, valueStr)
        }
        return defaultValue
}

// Helper function to read a float environment variable
func getFloatEnv(key string, defaultValue float64) float64 {
        if valueStr, exists := os.LookupEnv(key); exists {
                value, err := strconv.ParseFloat(valueStr, 64)
                if err == nil {
                        return value
                }
                log.Warnf("Invalid float value for %s: %s", key, valueStr)
        }
        return defaultValue
}