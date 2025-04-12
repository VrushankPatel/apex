# Real-Time Arbitrage Detector: Setup Guide

This guide will walk you through setting up and running the Real-Time Arbitrage Detector from scratch.

## System Requirements

- Go 1.18 or higher
- Git (for version control)
- Internet connection (for accessing exchange APIs)
- Modern web browser (for the UI)

## Installation Steps

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/arbitrage-detector.git
cd arbitrage-detector
```

### 2. Install Dependencies

```bash
go mod tidy
```

This will download and install all required Go packages.

### 3. Configure Environment Variables

```bash
# Create an environment file from the example
cp .env.example .env
```

Edit the `.env` file with your preferred text editor and configure:

- Exchange API credentials (if using real-time data)
- Simulation mode setting (`true` for testing, `false` for real data)
- Minimum profit threshold
- Logging level

### 4. Build the Application

```bash
# Optional: Build a binary
go build -o arbitrage-detector
```

Or you can run directly from source code in the next step.

### 5. Run the Application

```bash
# Using the binary
./arbitrage-detector

# Or directly from source
go run main.go
```

### 6. Access the Web Interface

Open your web browser and navigate to:

```
http://localhost:8080
```

You should see the Arbitrage Detector dashboard.

## Configuration Options

### Simulation Mode

By default, the application runs in simulation mode, which generates synthetic market data instead of connecting to real exchanges. This is useful for testing and demonstration purposes.

To use real exchange data:

1. Set `SIMULATION_MODE=false` in your `.env` file
2. Configure your exchange API keys (see [Exchange API Setup](./EXCHANGE_API_SETUP.md))

### Trading Pairs

The application monitors specific trading pairs for arbitrage opportunities. By default, it includes common cryptocurrency pairs like BTC/USDT and ETH/USDT.

To add or modify monitored pairs:

1. Update the `marketConfig` in `web/static/js/app.js` for the UI
2. Add the corresponding pairs to your exchange configuration in the Go code

### Minimum Profit Threshold

Set the minimum profit percentage for considering an arbitrage opportunity valid:

```
MIN_PROFIT_THRESHOLD=0.5
```

This filters out opportunities with potential profits below the specified percentage.

### Logging Level

Control the verbosity of application logs:

```
LOG_LEVEL=info
```

Options include:
- `debug`: Detailed debugging information
- `info`: General operational information
- `warn`: Warning events
- `error`: Error events only

## Additional Configuration

### Web Server Port

By default, the web server runs on port 8080. To change this:

1. Modify the port number in `pkg/server/server.go`
2. Update the workflow configuration if using Replit

### Exchange-Specific Settings

Each exchange connection can be individually configured:

- API endpoints
- Rate limiting
- Timeout settings
- Trading fees

See the exchange implementation files in `pkg/exchange` for details.

## Monitoring and Maintenance

### Log Files

Application logs are printed to the console and can be redirected to a file:

```bash
./arbitrage-detector > arbitrage.log 2>&1
```

### Performance Considerations

- **Network Bandwidth**: The application consumes bandwidth proportional to the number of exchanges and trading pairs monitored
- **Memory Usage**: Increases with the number of arbitrage opportunities stored in history
- **CPU Usage**: Generally low, with spikes during arbitrage calculations

## Troubleshooting

### Common Issues

1. **WebSocket Connection Errors**:
   - Check your internet connection
   - Verify API credentials
   - Ensure the exchange is operational

2. **No Arbitrage Opportunities Detected**:
   - Verify that exchanges are connected
   - Check if the minimum profit threshold is set too high
   - Confirm that the monitored trading pairs exist on multiple exchanges

3. **High CPU Usage**:
   - Reduce the number of monitored trading pairs
   - Increase the refresh interval

4. **Application Won't Start**:
   - Check for error messages in the console
   - Verify that the required port is available
   - Ensure Go is properly installed

### Getting Help

If you encounter issues not covered in this guide:

1. Check the application logs for error messages
2. Review the [GitHub repository](https://github.com/yourusername/arbitrage-detector) for existing issues or to report a new one
3. Consult the [API documentation](./API.md) for reference

## Upgrading

To update to the latest version:

```bash
git pull
go mod tidy
go build
```

Always check the release notes for breaking changes before upgrading.

## Running as a Service

### Linux (systemd)

Create a systemd service file:

```bash
sudo nano /etc/systemd/system/arbitrage-detector.service
```

Add the following content:

```
[Unit]
Description=Real-Time Arbitrage Detector
After=network.target

[Service]
User=yourusername
WorkingDirectory=/path/to/arbitrage-detector
ExecStart=/path/to/arbitrage-detector/arbitrage-detector
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Enable and start the service:

```bash
sudo systemctl enable arbitrage-detector
sudo systemctl start arbitrage-detector
```

### Docker Support

For containerized deployment, you can build a Docker image using the provided Dockerfile:

```bash
docker build -t arbitrage-detector .
docker run -p 8080:8080 --env-file .env arbitrage-detector
```

## Next Steps

Once you have the system running, consider:

1. [Setting up exchange API connections](./EXCHANGE_API_SETUP.md) for real-time data
2. [Learning about arbitrage trading](./ARBITRAGE_GUIDE.md) to better understand the opportunities
3. [Contributing to the project](./CONTRIBUTING.md) with improvements or fixes

---

With this setup guide, you should be able to get the Real-Time Arbitrage Detector up and running quickly. The system is designed to be both a learning tool for understanding arbitrage concepts and a practical tool for identifying real market opportunities.