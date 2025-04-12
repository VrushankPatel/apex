// Global variables
let socket;
let opportunities = [];
let marketData = {};
let currentMarket = 'crypto';
let currentPage = 1;
let itemsPerPage = 10;
let totalPages = 1;
let filteredOpportunities = [];
let activeExchanges = new Set(['Binance', 'Kraken']);
let refreshInterval = 2000;
let simulationInterval;

let statsData = {
    totalOpportunities: 0,
    totalProfit: 0,
    avgProfit: 0
};

// Market types and assets
const marketConfig = {
    'crypto': {
        pairs: ['BTC/USDT', 'ETH/USDT', 'SOL/USDT', 'ADA/USDT', 'XRP/USDT'],
        exchanges: ['Binance', 'Kraken', 'Coinbase', 'KuCoin']
    },
    'stock-us': {
        pairs: ['AAPL', 'MSFT', 'GOOGL', 'AMZN', 'TSLA', 'META'],
        exchanges: ['NASDAQ', 'NYSE']
    },
    'stock-in': {
        pairs: ['RELIANCE', 'TCS', 'INFY', 'HDFC', 'WIPRO'],
        exchanges: ['NSE', 'BSE']
    },
    'forex': {
        pairs: ['EUR/USD', 'GBP/USD', 'USD/JPY', 'USD/CHF'],
        exchanges: ['FX Markets']
    },
    'commodities': {
        pairs: ['GOLD', 'SILVER', 'OIL', 'NATURAL_GAS'],
        exchanges: ['COMEX', 'NYMEX']
    }
};

// Stock market simulated data
const simulatedStockData = {
    'AAPL': {
        'NASDAQ': { bid: 178.92, ask: 178.97 },
        'NYSE': { bid: 178.90, ask: 179.02 }
    },
    'MSFT': {
        'NASDAQ': { bid: 408.53, ask: 408.61 },
        'NYSE': { bid: 408.49, ask: 408.68 }
    },
    'GOOGL': {
        'NASDAQ': { bid: 152.08, ask: 152.15 },
        'NYSE': { bid: 152.04, ask: 152.18 }
    },
    'AMZN': {
        'NASDAQ': { bid: 181.12, ask: 181.19 },
        'NYSE': { bid: 181.08, ask: 181.24 }
    }
};

// Indian market simulated data
const simulatedIndianStockData = {
    'RELIANCE': {
        'NSE': { bid: 2875.45, ask: 2876.30 },
        'BSE': { bid: 2874.90, ask: 2876.85 }
    },
    'TCS': {
        'NSE': { bid: 3645.20, ask: 3646.50 },
        'BSE': { bid: 3644.75, ask: 3647.00 }
    },
    'INFY': {
        'NSE': { bid: 1499.90, ask: 1500.80 },
        'BSE': { bid: 1499.50, ask: 1501.20 }
    },
    'HDFC': {
        'NSE': { bid: 1632.40, ask: 1633.10 },
        'BSE': { bid: 1632.00, ask: 1633.50 }
    }
};

// DOM elements - will be initialized in init()
let statusIndicator;
let statusText;
let exchangeDataContainer;
let opportunitiesTableBody;
let totalOpportunitiesElement;
let totalProfitElement;
let avgProfitElement;
let lastUpdateElement;
let currentTimeElement;
let notification;
let notificationMessage;
let notificationCloseButton;
let marketSelect;
let refreshRateSelect;
let pairFilterSelect;
let minProfitInput;
let viewAllButton;
let prevPageButton;
let nextPageButton;
let pageInfoElement;
let exchangeBadges;

// Initialize the application
function init() {
    // Get DOM elements
    statusIndicator = document.getElementById('status-indicator');
    statusText = document.getElementById('status-text');
    exchangeDataContainer = document.getElementById('exchange-data');
    opportunitiesTableBody = document.getElementById('opportunities-table-body');
    totalOpportunitiesElement = document.getElementById('total-opportunities');
    totalProfitElement = document.getElementById('total-profit');
    avgProfitElement = document.getElementById('avg-profit');
    lastUpdateElement = document.getElementById('last-update');
    currentTimeElement = document.getElementById('current-time');
    notification = document.getElementById('notification');
    notificationMessage = document.getElementById('notification-message');
    notificationCloseButton = document.getElementById('notification-close');
    marketSelect = document.getElementById('market-select');
    refreshRateSelect = document.getElementById('refresh-rate');
    pairFilterSelect = document.getElementById('pair-filter');
    minProfitInput = document.getElementById('min-profit');
    viewAllButton = document.getElementById('view-all-btn');
    prevPageButton = document.getElementById('prev-page');
    nextPageButton = document.getElementById('next-page');
    pageInfoElement = document.getElementById('page-info');
    exchangeBadges = document.querySelectorAll('.exchange-badge');
    
    // Set up WebSocket connection
    connectWebSocket();
    
    // Set up notification close button
    notificationCloseButton.addEventListener('click', () => {
        notification.classList.add('hidden');
    });
    
    // Add event listeners for market settings
    marketSelect.addEventListener('change', handleMarketChange);
    refreshRateSelect.addEventListener('change', handleRefreshRateChange);
    pairFilterSelect.addEventListener('change', handleFilterChange);
    minProfitInput.addEventListener('input', handleFilterChange);
    viewAllButton.addEventListener('click', showAllOpportunities);
    prevPageButton.addEventListener('click', () => goToPage(currentPage - 1));
    nextPageButton.addEventListener('click', () => goToPage(currentPage + 1));
    
    // Add event listeners for exchange badges
    exchangeBadges.forEach(badge => {
        badge.addEventListener('click', toggleExchange);
    });
    
    // Update the current time every second
    setInterval(updateTime, 1000);
    updateTime();
    
    // Start simulating data for stock markets
    startMarketSimulation();
}

// WebSocket connection
function connectWebSocket() {
    // Get the current host and determine WebSocket URL
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    
    // Build the WebSocket URL using the current host
    // For Replit, we need to ensure we're using the right port (5000)
    let host = window.location.host;
    if (host.includes('replit.dev') || host.includes('replit.app')) {
        // Use the same host for Replit deployments
        host = window.location.host;
    } else if (host.includes(':')) {
        // Local development - preserve the port
        host = window.location.host;
    } else {
        // Make sure we use port 8080 if no port is specified
        host = host + ':8080';
    }
    
    const wsUrl = `${protocol}//${host}/ws`;
    console.log('Connecting to WebSocket at:', wsUrl);
    
    // Create WebSocket connection
    socket = new WebSocket(wsUrl);
    
    // Connection opened
    socket.addEventListener('open', () => {
        setConnectionStatus(true);
        console.log('WebSocket connection established');
    });
    
    // Listen for messages
    socket.addEventListener('message', (event) => {
        try {
            const data = JSON.parse(event.data);
            console.log('Received data:', data);
            processData(data);
        } catch (error) {
            console.error('Error parsing WebSocket message:', error);
        }
    });
    
    // Connection closed
    socket.addEventListener('close', () => {
        setConnectionStatus(false);
        console.log('WebSocket connection closed');
        
        // Try to reconnect after 3 seconds
        setTimeout(connectWebSocket, 3000);
    });
    
    // Error handling
    socket.addEventListener('error', (error) => {
        console.error('WebSocket error:', error);
        setConnectionStatus(false);
    });
}

// Process incoming WebSocket data
function processData(data) {
    switch (data.type) {
        case 'market':
            marketData = data.data;
            updateExchangeGrid();
            break;
        case 'opportunities':
            opportunities = data.data;
            updateStats();
            updateOpportunitiesTable();
            break;
        case 'opportunity':
            // Add new opportunity to the beginning of the array
            opportunities.unshift(data.data);
            
            // Keep only the last 100 opportunities
            if (opportunities.length > 100) {
                opportunities = opportunities.slice(0, 100);
            }
            
            // Show notification
            showNotification(data.data);
            
            // Update UI
            updateStats();
            updateOpportunitiesTable();
            break;
        default:
            console.log('Unknown message type:', data.type);
    }
    
    // Update the last update time
    lastUpdateElement.textContent = formatTime(new Date());
}

// Show notification for new opportunity
function showNotification(opportunity) {
    const message = `
        Buy on ${opportunity.BuyExchange} at ${formatCurrency(opportunity.BuyPrice)}
        Sell on ${opportunity.SellExchange} at ${formatCurrency(opportunity.SellPrice)}
        Profit: ${opportunity.ProfitPercentage.toFixed(4)}% (${formatCurrency(opportunity.NetProfit)})
    `;
    
    notificationMessage.textContent = message;
    notification.classList.remove('hidden');
    
    // Auto-hide after 5 seconds
    setTimeout(() => {
        notification.classList.add('hidden');
    }, 5000);
}

// Update the statistics
function updateStats() {
    if (opportunities.length > 0) {
        // Calculate statistics
        const totalOpportunities = opportunities.length;
        let totalProfit = 0;
        let totalProfitPct = 0;
        
        opportunities.forEach(opp => {
            totalProfit += opp.NetProfit;
            totalProfitPct += opp.ProfitPercentage;
        });
        
        const avgProfit = totalProfitPct / totalOpportunities;
        
        // Update the stats data
        statsData = {
            totalOpportunities,
            totalProfit,
            avgProfit
        };
        
        // Update UI
        totalOpportunitiesElement.textContent = totalOpportunities;
        totalProfitElement.textContent = formatCurrency(totalProfit);
        avgProfitElement.textContent = avgProfit.toFixed(2) + '%';
    }
}

// Update the opportunities table
function updateOpportunitiesTable() {
    if (opportunities.length === 0) {
        opportunitiesTableBody.innerHTML = `
            <tr>
                <td colspan="7" class="table-empty-message">Waiting for arbitrage opportunities...</td>
            </tr>
        `;
        return;
    }
    
    // Clear the table
    opportunitiesTableBody.innerHTML = '';
    
    // Add opportunities to the table
    opportunities.forEach((opp, index) => {
        const row = document.createElement('tr');
        
        // Highlight the newest opportunity
        if (index === 0) {
            row.classList.add('highlight');
        }
        
        // Format the timestamp
        const timestamp = new Date(opp.Timestamp);
        
        row.innerHTML = `
            <td>${formatTime(timestamp)}</td>
            <td>${opp.BuyExchange}</td>
            <td>${opp.SellExchange}</td>
            <td>${formatCurrency(opp.BuyPrice)}</td>
            <td>${formatCurrency(opp.SellPrice)}</td>
            <td>${opp.ProfitPercentage.toFixed(4)}%</td>
            <td>${formatCurrency(opp.NetProfit)}</td>
        `;
        
        opportunitiesTableBody.appendChild(row);
    });
}

// Update the exchange grid
function updateExchangeGrid() {
    if (Object.keys(marketData).length === 0) {
        exchangeDataContainer.innerHTML = `
            <div class="exchange-grid-row">
                <div class="exchange-col">Waiting for market data...</div>
                <div class="exchange-col"></div>
                <div class="exchange-col"></div>
                <div class="exchange-col"></div>
                <div class="exchange-col"></div>
            </div>
        `;
        return;
    }
    
    // Clear the existing data
    exchangeDataContainer.innerHTML = '';
    
    // Group by trading pairs
    const pairs = {};
    
    // Organize by trading pair
    Object.values(marketData).forEach(book => {
        const pairKey = `${book.BaseCurrency}/${book.QuoteCurrency}`;
        if (!pairs[pairKey]) {
            pairs[pairKey] = {};
        }
        pairs[pairKey][book.Exchange] = book;
    });
    
    // Add each exchange and pair to the grid
    Object.entries(pairs).forEach(([pair, exchanges]) => {
        Object.entries(exchanges).forEach(([exchange, data]) => {
            const row = document.createElement('div');
            row.className = 'exchange-grid-row';
            
            // Calculate the spread percentage
            const spreadPercentage = ((data.Ask - data.Bid) / data.Bid) * 100;
            
            row.innerHTML = `
                <div class="exchange-col">${pair}</div>
                <div class="exchange-col">${exchange}</div>
                <div class="exchange-col">${formatCurrency(data.Bid)}</div>
                <div class="exchange-col">${formatCurrency(data.Ask)}</div>
                <div class="exchange-col">${spreadPercentage.toFixed(4)}%</div>
            `;
            
            exchangeDataContainer.appendChild(row);
        });
    });
}

// Update connection status indicator
function setConnectionStatus(isConnected) {
    if (isConnected) {
        statusIndicator.classList.add('connected');
        statusText.textContent = 'Connected';
    } else {
        statusIndicator.classList.remove('connected');
        statusText.textContent = 'Disconnected';
    }
}

// Update the current time display
function updateTime() {
    const now = new Date();
    currentTimeElement.textContent = formatTime(now);
}

// Format a timestamp as HH:MM:SS
function formatTime(date) {
    return date.toLocaleTimeString(undefined, {
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        hour12: false
    });
}

// Format a currency value with 2 decimal places
function formatCurrency(value) {
    return value.toFixed(2);
}

// Handle market change
function handleMarketChange(event) {
    currentMarket = event.target.value;
    
    // Update the pair filter dropdown
    updatePairFilter();
    
    // Clear current data
    if (simulationInterval) {
        clearInterval(simulationInterval);
    }
    
    // Start simulation for the selected market
    startMarketSimulation();
}

// Handle refresh rate change
function handleRefreshRateChange(event) {
    refreshInterval = parseInt(event.target.value);
    
    // Update simulation
    if (simulationInterval) {
        clearInterval(simulationInterval);
    }
    startMarketSimulation();
}

// Update pair filter based on current market
function updatePairFilter() {
    // Clear current options
    pairFilterSelect.innerHTML = '<option value="all">All Pairs</option>';
    
    // Add options based on selected market
    const pairs = marketConfig[currentMarket].pairs;
    pairs.forEach(pair => {
        const value = pair.replace('/', '-').toLowerCase();
        pairFilterSelect.innerHTML += `<option value="${value}">${pair}</option>`;
    });
}

// Handle filter change
function handleFilterChange() {
    const pairFilter = pairFilterSelect.value;
    const minProfit = parseFloat(minProfitInput.value);
    
    // Filter opportunities
    filteredOpportunities = opportunities.filter(opp => {
        // Check profit threshold
        if (opp.ProfitPercentage < minProfit) {
            return false;
        }
        
        // Check pair filter
        if (pairFilter !== 'all') {
            const pairKey = `${opp.BaseCurrency}/${opp.QuoteCurrency}`.toLowerCase();
            if (!pairKey.includes(pairFilter)) {
                return false;
            }
        }
        
        return true;
    });
    
    // Update pagination
    totalPages = Math.ceil(filteredOpportunities.length / itemsPerPage);
    currentPage = 1; // Reset to first page
    
    // Update UI
    updateOpportunitiesTable();
    updatePagination();
}

// Show all opportunities
function showAllOpportunities() {
    pairFilterSelect.value = 'all';
    minProfitInput.value = '0.1';
    handleFilterChange();
}

// Go to specific page
function goToPage(page) {
    if (page < 1 || page > totalPages) {
        return;
    }
    
    currentPage = page;
    updateOpportunitiesTable();
    updatePagination();
}

// Update pagination controls
function updatePagination() {
    pageInfoElement.textContent = `Page ${currentPage} of ${totalPages}`;
    prevPageButton.disabled = currentPage <= 1;
    nextPageButton.disabled = currentPage >= totalPages;
}

// Toggle exchange in the active exchanges list
function toggleExchange(event) {
    const badge = event.target;
    const exchange = badge.dataset.exchange;
    
    if (badge.classList.contains('active')) {
        badge.classList.remove('active');
        activeExchanges.delete(exchange);
    } else {
        badge.classList.add('active');
        activeExchanges.add(exchange);
    }
    
    // Update data
    startMarketSimulation();
}

// Start market simulation
function startMarketSimulation() {
    // Stop current simulation if any
    if (simulationInterval) {
        clearInterval(simulationInterval);
    }
    
    // Add simulated market data
    simulateMarketData();
    
    // Start interval for continuous simulation
    simulationInterval = setInterval(simulateMarketData, refreshInterval);
}

// Simulate market data for different markets
function simulateMarketData() {
    let simulatedData = {};
    
    switch (currentMarket) {
        case 'crypto':
            // Crypto data should come from real WebSocket
            // But we'll simulate if no real data available
            if (Object.keys(marketData).length === 0) {
                simulateCryptoData();
            }
            break;
            
        case 'stock-us':
            simulateStockMarketData();
            break;
            
        case 'stock-in':
            simulateIndianStockMarketData();
            break;
            
        case 'forex':
            simulateForexData();
            break;
            
        case 'commodities':
            simulateCommoditiesData();
            break;
    }
    
    // Generate new arbitrage opportunities
    simulateArbitrageOpportunities();
}

// Simulate crypto market data
function simulateCryptoData() {
    // Clear current data
    marketData = {};
    
    // Generate data for BTC/USDT pair
    const btcBase = 60000 + Math.random() * 10000;
    
    marketData['Binance-BTC'] = {
        Exchange: 'Binance',
        Symbol: 'BTCUSDT',
        BaseCurrency: 'BTC',
        QuoteCurrency: 'USDT',
        Bid: btcBase - Math.random() * 50,
        Ask: btcBase + Math.random() * 50,
        LastUpdate: new Date().toISOString()
    };
    
    marketData['Kraken-BTC'] = {
        Exchange: 'Kraken',
        Symbol: 'XBT/USDT',
        BaseCurrency: 'BTC',
        QuoteCurrency: 'USDT',
        Bid: btcBase - Math.random() * 80,
        Ask: btcBase + Math.random() * 80,
        LastUpdate: new Date().toISOString()
    };
    
    // Generate data for ETH/USDT pair
    const ethBase = 3000 + Math.random() * 500;
    
    marketData['Binance-ETH'] = {
        Exchange: 'Binance',
        Symbol: 'ETHUSDT',
        BaseCurrency: 'ETH',
        QuoteCurrency: 'USDT',
        Bid: ethBase - Math.random() * 5,
        Ask: ethBase + Math.random() * 5,
        LastUpdate: new Date().toISOString()
    };
    
    marketData['Kraken-ETH'] = {
        Exchange: 'Kraken',
        Symbol: 'ETH/USDT',
        BaseCurrency: 'ETH',
        QuoteCurrency: 'USDT',
        Bid: ethBase - Math.random() * 8,
        Ask: ethBase + Math.random() * 8,
        LastUpdate: new Date().toISOString()
    };
    
    // Update UI
    updateExchangeGrid();
}

// Simulate US stock market data
function simulateStockMarketData() {
    // Clear current data
    marketData = {};
    
    // Generate data for each stock
    Object.entries(simulatedStockData).forEach(([stock, exchanges]) => {
        Object.entries(exchanges).forEach(([exchange, prices]) => {
            // Add small random variations
            const bid = prices.bid * (1 + (Math.random() * 0.002 - 0.001));
            const ask = prices.ask * (1 + (Math.random() * 0.002 - 0.001));
            
            marketData[`${exchange}-${stock}`] = {
                Exchange: exchange,
                Symbol: stock,
                BaseCurrency: stock,
                QuoteCurrency: 'USD',
                Bid: bid,
                Ask: ask,
                LastUpdate: new Date().toISOString()
            };
        });
    });
    
    // Update UI
    updateExchangeGrid();
}

// Simulate Indian stock market data
function simulateIndianStockMarketData() {
    // Clear current data
    marketData = {};
    
    // Generate data for each stock
    Object.entries(simulatedIndianStockData).forEach(([stock, exchanges]) => {
        Object.entries(exchanges).forEach(([exchange, prices]) => {
            // Add small random variations
            const bid = prices.bid * (1 + (Math.random() * 0.002 - 0.001));
            const ask = prices.ask * (1 + (Math.random() * 0.002 - 0.001));
            
            marketData[`${exchange}-${stock}`] = {
                Exchange: exchange,
                Symbol: stock,
                BaseCurrency: stock,
                QuoteCurrency: 'INR',
                Bid: bid,
                Ask: ask,
                LastUpdate: new Date().toISOString()
            };
        });
    });
    
    // Update UI
    updateExchangeGrid();
}

// Simulate forex data
function simulateForexData() {
    // Clear current data
    marketData = {};
    
    // EUR/USD
    const eurusdBase = 1.082 + (Math.random() * 0.004 - 0.002);
    marketData['FX-EURUSD'] = {
        Exchange: 'FX Markets',
        Symbol: 'EUR/USD',
        BaseCurrency: 'EUR',
        QuoteCurrency: 'USD',
        Bid: eurusdBase - 0.0002,
        Ask: eurusdBase + 0.0002,
        LastUpdate: new Date().toISOString()
    };
    
    // GBP/USD
    const gbpusdBase = 1.262 + (Math.random() * 0.006 - 0.003);
    marketData['FX-GBPUSD'] = {
        Exchange: 'FX Markets',
        Symbol: 'GBP/USD',
        BaseCurrency: 'GBP',
        QuoteCurrency: 'USD',
        Bid: gbpusdBase - 0.0003,
        Ask: gbpusdBase + 0.0003,
        LastUpdate: new Date().toISOString()
    };
    
    // Update UI
    updateExchangeGrid();
}

// Simulate commodities data
function simulateCommoditiesData() {
    // Clear current data
    marketData = {};
    
    // Gold
    const goldBase = 2320 + (Math.random() * 20 - 10);
    marketData['COMEX-GOLD'] = {
        Exchange: 'COMEX',
        Symbol: 'GOLD',
        BaseCurrency: 'GOLD',
        QuoteCurrency: 'USD',
        Bid: goldBase - 0.5,
        Ask: goldBase + 0.5,
        LastUpdate: new Date().toISOString()
    };
    
    // Silver
    const silverBase = 27.5 + (Math.random() * 0.4 - 0.2);
    marketData['COMEX-SILVER'] = {
        Exchange: 'COMEX',
        Symbol: 'SILVER',
        BaseCurrency: 'SILVER',
        QuoteCurrency: 'USD',
        Bid: silverBase - 0.05,
        Ask: silverBase + 0.05,
        LastUpdate: new Date().toISOString()
    };
    
    // Update UI
    updateExchangeGrid();
}

// Simulate arbitrage opportunities
function simulateArbitrageOpportunities() {
    // Find potential arbitrage opportunities
    const pairs = {};
    
    // Organize by trading pair
    Object.values(marketData).forEach(book => {
        const pairKey = `${book.BaseCurrency}/${book.QuoteCurrency}`;
        if (!pairs[pairKey]) {
            pairs[pairKey] = [];
        }
        pairs[pairKey].push(book);
    });
    
    // Check for arbitrage opportunities
    Object.entries(pairs).forEach(([pair, books]) => {
        if (books.length < 2) return; // Need at least 2 exchanges for arbitrage
        
        // Check all pairs of exchanges
        for (let i = 0; i < books.length; i++) {
            for (let j = i + 1; j < books.length; j++) {
                const exchange1 = books[i];
                const exchange2 = books[j];
                
                // Check if these are active exchanges
                if (!activeExchanges.has(exchange1.Exchange) || 
                    !activeExchanges.has(exchange2.Exchange)) {
                    continue;
                }
                
                // Check if buy on exchange1, sell on exchange2 is profitable
                const profitPct1 = (exchange2.Bid - exchange1.Ask) / exchange1.Ask * 100;
                
                // Check if buy on exchange2, sell on exchange1 is profitable
                const profitPct2 = (exchange1.Bid - exchange2.Ask) / exchange2.Ask * 100;
                
                // Add opportunities that are above threshold
                const minThreshold = parseFloat(minProfitInput.value) || 0.1;
                
                if (profitPct1 > minThreshold) {
                    const opportunity = {
                        Timestamp: new Date().toISOString(),
                        BaseCurrency: exchange1.BaseCurrency,
                        QuoteCurrency: exchange1.QuoteCurrency,
                        BuyExchange: exchange1.Exchange,
                        SellExchange: exchange2.Exchange,
                        BuyPrice: exchange1.Ask,
                        SellPrice: exchange2.Bid,
                        ProfitPercentage: profitPct1,
                        NetProfit: (exchange2.Bid - exchange1.Ask) * 1 // Assume 1 unit volume
                    };
                    
                    // Add opportunity
                    opportunities.unshift(opportunity);
                    
                    // Show notification for significant opportunities
                    if (profitPct1 > 1) {
                        showNotification(opportunity);
                    }
                }
                
                if (profitPct2 > minThreshold) {
                    const opportunity = {
                        Timestamp: new Date().toISOString(),
                        BaseCurrency: exchange2.BaseCurrency,
                        QuoteCurrency: exchange2.QuoteCurrency,
                        BuyExchange: exchange2.Exchange,
                        SellExchange: exchange1.Exchange,
                        BuyPrice: exchange2.Ask,
                        SellPrice: exchange1.Bid,
                        ProfitPercentage: profitPct2,
                        NetProfit: (exchange1.Bid - exchange2.Ask) * 1 // Assume 1 unit volume
                    };
                    
                    // Add opportunity
                    opportunities.unshift(opportunity);
                    
                    // Show notification for significant opportunities
                    if (profitPct2 > 1) {
                        showNotification(opportunity);
                    }
                }
            }
        }
    });
    
    // Keep only the last 100 opportunities
    if (opportunities.length > 100) {
        opportunities = opportunities.slice(0, 100);
    }
    
    // Update filtered opportunities
    handleFilterChange();
    
    // Update stats and UI
    updateStats();
    updateOpportunitiesTable();
}

// Update the opportunities table with filtered and paginated data
function updateOpportunitiesTable() {
    if (opportunities.length === 0) {
        opportunitiesTableBody.innerHTML = `
            <tr>
                <td colspan="9" class="table-empty-message">Waiting for arbitrage opportunities...</td>
            </tr>
        `;
        return;
    }
    
    // Filter opportunities
    handleFilterChange();
    
    // Paginate
    const start = (currentPage - 1) * itemsPerPage;
    const end = start + itemsPerPage;
    const paginatedOpportunities = filteredOpportunities.slice(start, end);
    
    // Clear the table
    opportunitiesTableBody.innerHTML = '';
    
    // Check if we have any opportunities after filtering
    if (paginatedOpportunities.length === 0) {
        opportunitiesTableBody.innerHTML = `
            <tr>
                <td colspan="9" class="table-empty-message">No opportunities match your filters</td>
            </tr>
        `;
        return;
    }
    
    // Add opportunities to the table
    paginatedOpportunities.forEach((opp, index) => {
        const row = document.createElement('tr');
        
        // Highlight the newest opportunity
        if (index === 0 && currentPage === 1) {
            row.classList.add('highlight');
        }
        
        // Format the timestamp
        const timestamp = new Date(opp.Timestamp);
        
        // Create asset name
        const asset = opp.BaseCurrency && opp.QuoteCurrency ? 
            `${opp.BaseCurrency}/${opp.QuoteCurrency}` : 
            opp.BaseCurrency || 'Unknown';
        
        row.innerHTML = `
            <td>${formatTime(timestamp)}</td>
            <td>${asset}</td>
            <td>${opp.BuyExchange}</td>
            <td>${opp.SellExchange}</td>
            <td>${formatCurrency(opp.BuyPrice)}</td>
            <td>${formatCurrency(opp.SellPrice)}</td>
            <td>${opp.ProfitPercentage.toFixed(4)}%</td>
            <td>${formatCurrency(opp.NetProfit)}</td>
            <td><button class="action-btn">Details</button></td>
        `;
        
        opportunitiesTableBody.appendChild(row);
    });
    
    // Update pagination
    updatePagination();
}

// Initialize the app when the DOM is loaded
document.addEventListener('DOMContentLoaded', init);