// Global variables
let socket;
let opportunities = [];
let marketData = {};
let statsData = {
    totalOpportunities: 0,
    totalProfit: 0,
    avgProfit: 0
};

// DOM Elements
const statusIndicator = document.getElementById('status-indicator');
const statusText = document.getElementById('status-text');
const exchangeDataContainer = document.getElementById('exchange-data');
const opportunitiesTableBody = document.getElementById('opportunities-table-body');
const totalOpportunitiesElement = document.getElementById('total-opportunities');
const totalProfitElement = document.getElementById('total-profit');
const avgProfitElement = document.getElementById('avg-profit');
const lastUpdateElement = document.getElementById('last-update');
const currentTimeElement = document.getElementById('current-time');
const notification = document.getElementById('notification');
const notificationMessage = document.getElementById('notification-message');
const notificationCloseButton = document.getElementById('notification-close');

// Initialize the application
function init() {
    // Set up WebSocket connection
    connectWebSocket();
    
    // Set up notification close button
    notificationCloseButton.addEventListener('click', () => {
        notification.classList.add('hidden');
    });
    
    // Update the current time every second
    setInterval(updateTime, 1000);
    updateTime();
}

// WebSocket connection
function connectWebSocket() {
    // Get the current host and determine WebSocket URL
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws`;
    
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

// Initialize the app when the DOM is loaded
document.addEventListener('DOMContentLoaded', init);