# Real-Time Crypto Arbitrage Detector

![Arbitrage Detector](assets/arbitrage-banner.png)

[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/arbitrage-detector)](https://goreportcard.com/report/github.com/yourusername/arbitrage-detector)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A high-performance, real-time cryptocurrency arbitrage detector built in Go. Monitor multiple exchanges simultaneously to identify profitable trading opportunities where price differences between exchanges create potential for profit.

## 🚀 Features

- **Multi-Exchange Support**: Currently supports Binance, Kraken, and Coinbase with easy extensibility
- **Real-Time Monitoring**: WebSocket connections for live market data
- **Multi-Pair Trading**: Monitor BTC, ETH, SOL, and more against USDT/USD
- **Fee-Aware Calculations**: Accounts for exchange trading fees when calculating potential profits
- **Opportunity Logging**: Records all detected arbitrage opportunities for analysis
- **High Performance**: Written in Go for speed and efficiency
- **Web UI**: Simple dashboard to visualize opportunities (optional)
- **Simulation Mode**: Test the system without real exchange connections

## 📊 How It Works

1. **Connect to Exchanges**: Establishes WebSocket connections to supported cryptocurrency exchanges
2. **Market Data Analysis**: Continuously monitors the bid-ask spreads across exchanges
3. **Arbitrage Detection**: Identifies situations where buying on one exchange and selling on another would yield profit
4. **Fee Calculation**: Accounts for trading fees to ensure opportunities are genuinely profitable
5. **Opportunity Logging**: Records profitable arbitrage opportunities with timestamps and profit metrics

## 📋 Requirements

- Go 1.19 or higher
- Internet connection for live exchange data
- API credentials for exchanges (optional, simulation mode works without them)

## 🔧 Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/arbitrage-detector.git
cd arbitrage-detector

# Install dependencies
go mod tidy

# Run the application
go run main.go
```

## ⚙️ Configuration

Configuration is handled via the `config.yaml` file:

```yaml
# Minimum profit threshold (as a decimal, e.g., 0.001 = 0.1%)
minProfitThreshold: 0.001

# Trading pairs to monitor across exchanges
tradingPairs:
  - baseCurrency: BTC
    quoteCurrency: USDT
  - baseCurrency: ETH
    quoteCurrency: USDT
  - baseCurrency: SOL
    quoteCurrency: USDT

# Exchange-specific configurations
exchanges:
  binance:
    enabled: true
    takerFee: 0.001  # 0.1%
  
  kraken:
    enabled: true
    takerFee: 0.0026  # 0.26%
  
  coinbase:
    enabled: false
    takerFee: 0.006  # 0.6%
```

## 📊 Sample Output

```
time="2025-04-12 00:13:44" level=info msg="Real-Time Arbitrage Detector Started"
time="2025-04-12 00:13:44" level=info msg="Monitoring pairs: BTC/USDT, ETH/USDT, SOL/USDT"
time="2025-04-12 00:13:44" level=info msg="Min profit threshold: 0.10%"

time="2025-04-12 00:13:54" level=info msg="MARKET SUMMARY FOR BTC/USDT" avg_profit_pct="16.25%" best_buy="Binance at 70925.14 USDT" best_sell="Binance at 70868.51 USDT" binance_ask="70925.14 USDT" binance_bid="70868.51 USDT" kraken_ask="70946.20 USDT" kraken_bid="70851.72 USDT" opportunities=11 pair=BTC/USDT price_spread="-0.0798%" 

time="2025-04-12 00:13:54" level=info msg="SIMULATED ARBITRAGE OPPORTUNITY DETECTED" buy_exchange=Binance buy_price=63526.14 sell_exchange=Kraken sell_price=74113.83 profit_percentage="16.25%" net_profit="10331.47 USDT" simulated=true
```

## 🔍 Web UI

The project includes an optional web-based dashboard for visualizing arbitrage opportunities in real-time. To access it:

1. Run the application: `go run main.go`
2. Open your browser to: `http://localhost:8080`

![Web UI Preview](assets/web-ui-preview.png)

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## ⚠️ Disclaimer

This software is for educational purposes only. Cryptocurrency trading carries significant risk. This tool does not constitute financial advice. Always do your own research before trading.