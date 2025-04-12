# Setting Up Exchange API Connections

This guide explains how to configure your Arbitrage Detector to use real-time data from cryptocurrency exchanges instead of simulated data.

## Overview

The Arbitrage Detector can connect to multiple cryptocurrency exchanges via their APIs to access real-time market data. While the system includes a simulation mode for testing, connecting to actual exchanges provides the most accurate arbitrage detection.

## Exchange API Keys

### What Are API Keys?

API keys are credentials that allow applications to access exchange data and services programmatically. For arbitrage detection, we only need read-only access to market data.

### Supported Exchanges

The system currently supports:

1. **Binance**
2. **Kraken**
3. **Coinbase**

## Step-by-Step API Key Setup

### Binance API Keys

1. **Create an account** on [Binance](https://www.binance.com) if you don't have one.
2. **Complete verification** (KYC) as required by Binance.
3. **Generate API keys**:
   - Log into your Binance account
   - Navigate to "API Management" (usually under your profile or account settings)
   - Click "Create API"
   - Set a label for your API key (e.g., "Arbitrage Detector")
   - For security, restrict API access to only read-only permissions
   - Copy and securely store your API Key and Secret Key

### Kraken API Keys

1. **Create an account** on [Kraken](https://www.kraken.com) if you don't have one.
2. **Complete verification** (KYC) as required by Kraken.
3. **Generate API keys**:
   - Log into your Kraken account
   - Navigate to "Security" > "API"
   - Click "Add New Key"
   - Set a key description (e.g., "Arbitrage Detector")
   - Under permissions, select only "Query Funds" and "Query Open Orders & Trades"
   - Click "Generate Key"
   - Copy and securely store your API Key and Private Key

### Coinbase API Keys

1. **Create an account** on [Coinbase Advanced](https://advanced.coinbase.com/) if you don't have one.
2. **Complete verification** (KYC) as required by Coinbase.
3. **Generate API keys**:
   - Log into your Coinbase Advanced account
   - Navigate to "API" (found under your profile)
   - Click "New API Key"
   - Select appropriate portfolios
   - Under permissions, select only "View" permissions
   - Set a passphrase
   - Click "Create API Key"
   - Copy and securely store your API Key, Secret, and Passphrase

## Configuring the Arbitrage Detector

1. **Create a .env file**:
   ```bash
   cp .env.example .env
   ```

2. **Edit the .env file** with your API keys:
   ```
   # Binance API credentials
   BINANCE_API_KEY=your_binance_api_key_here
   BINANCE_API_SECRET=your_binance_api_secret_here
   
   # Kraken API credentials
   KRAKEN_API_KEY=your_kraken_api_key_here
   KRAKEN_API_SECRET=your_kraken_api_secret_here
   
   # Coinbase API credentials
   COINBASE_API_KEY=your_coinbase_api_key_here
   COINBASE_API_SECRET=your_coinbase_api_secret_here
   COINBASE_PASSPHRASE=your_coinbase_passphrase_here
   
   # Configuration
   SIMULATION_MODE=false
   LOG_LEVEL=info
   MIN_PROFIT_THRESHOLD=0.5
   ```

3. **Switch to Real-Time Mode**:
   - Make sure `SIMULATION_MODE=false` in your .env file
   - Adjust the minimum profit threshold as needed
   - Save the file

4. **Restart the Application**:
   ```bash
   go run main.go
   ```

## Security Best Practices

1. **Use Read-Only Keys**: For security, only enable view/read permissions, not trading
2. **IP Restrictions**: If possible, restrict API access to your server's IP address
3. **Secure Storage**: Keep your .env file secure and never commit it to version control
4. **Regular Rotation**: Periodically rotate your API keys
5. **Monitoring**: Regularly check your exchange account for unauthorized activity

## Troubleshooting

### Common API Connection Issues

1. **Rate Limiting**: Exchanges impose rate limits on API calls. The system is designed to respect these limits, but if you exceed them, you may see temporary connection failures.

2. **Network Issues**: Ensure your server has reliable internet connectivity.

3. **Invalid Credentials**: Double-check your API keys and secrets for accuracy.

4. **API Changes**: Exchange APIs may change over time. Check the application logs for any indication of API compatibility issues.

### Connection Status Check

The web interface displays the connection status for each exchange. If you see "Disconnected" for an exchange:

1. Check the application logs for error messages
2. Verify your API keys are correct
3. Ensure the exchange is operational (check their status page)
4. Restart the application

## Adding New Exchanges

The system is designed to be extensible. To add a new exchange:

1. Implement the exchange interface in `pkg/exchange`
2. Add configuration options in `pkg/config`
3. Update the web interface to display the new exchange

Refer to the existing exchange implementations as examples.

---

With proper API configuration, your Arbitrage Detector will transition from simulation mode to monitoring real market conditions, potentially identifying actual arbitrage opportunities in the cryptocurrency market.