package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
)

// Config holds the application configuration
type Config struct {
	MinProfitThreshold float64         `mapstructure:"minProfitThreshold"`
	Exchanges          ExchangesConfig `mapstructure:"exchanges"`
	Logging            LoggingConfig   `mapstructure:"logging"`
	Opportunities      OpportunitiesConfig `mapstructure:"opportunities"`
}

// ExchangesConfig holds the configuration for all exchanges
type ExchangesConfig struct {
	Binance BinanceConfig `mapstructure:"binance"`
	Kraken  KrakenConfig  `mapstructure:"kraken"`
}

// BinanceConfig holds Binance-specific configuration
type BinanceConfig struct {
	Enabled    bool    `mapstructure:"enabled"`
	PairSymbol string  `mapstructure:"pairSymbol"`
	TakerFee   float64 `mapstructure:"takerFee"`
}

// KrakenConfig holds Kraken-specific configuration
type KrakenConfig struct {
	Enabled    bool    `mapstructure:"enabled"`
	PairSymbol string  `mapstructure:"pairSymbol"`
	TakerFee   float64 `mapstructure:"takerFee"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	File   string `mapstructure:"file"`
	Format string `mapstructure:"format"`
}

// OpportunitiesConfig holds configuration for opportunity logging
type OpportunitiesConfig struct {
	LogFile   string `mapstructure:"logFile"`
	LogFormat string `mapstructure:"logFormat"`
}

// LoadConfig loads the application configuration from config.yaml
func LoadConfig() (*Config, error) {
	// Create a new Viper instance
	v := viper.New()
	
	// Set default configuration values
	setDefaultConfig(v)
	
	// Set configuration file name and paths
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	
	// Read configuration file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Warning("No config file found, using default values")
		} else {
			return nil, fmt.Errorf("error reading config file: %v", err)
		}
	}
	
	// Override with environment variables if present
	v.SetEnvPrefix("ARBITRAGE")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	
	// Create the data directory if it doesn't exist
	os.MkdirAll("data", 0755)
	
	// Parse the configuration into our struct
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %v", err)
	}
	
	// Set log level based on config
	setLogLevel(config.Logging.Level)
	
	return &config, nil
}

// setDefaultConfig sets default configuration values
func setDefaultConfig(v *viper.Viper) {
	// Minimum profit threshold
	v.SetDefault("minProfitThreshold", 0.001) // 0.1%
	
	// Binance defaults
	v.SetDefault("exchanges.binance.enabled", true)
	v.SetDefault("exchanges.binance.pairSymbol", "BTCUSDT")
	v.SetDefault("exchanges.binance.takerFee", 0.001) // 0.1%
	
	// Kraken defaults
	v.SetDefault("exchanges.kraken.enabled", true)
	v.SetDefault("exchanges.kraken.pairSymbol", "XBT/USDT")
	v.SetDefault("exchanges.kraken.takerFee", 0.0026) // 0.26%
	
	// Logging defaults
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.file", "data/arbitrage.log")
	v.SetDefault("logging.format", "text")
	
	// Opportunities logging defaults
	v.SetDefault("opportunities.logFile", "data/opportunities.csv")
	v.SetDefault("opportunities.logFormat", "csv")
}

// setLogLevel sets the appropriate log level
func setLogLevel(level string) {
	switch strings.ToLower(level) {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}
