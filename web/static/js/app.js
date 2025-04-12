// DOM Elements
const connectionStatus = document.getElementById('connection-status');
const activeExchanges = document.getElementById('active-exchanges');
const tradingPairs = document.getElementById('trading-pairs');
const totalOpportunities = document.getElementById('total-opportunities');
const potentialProfit = document.getElementById('potential-profit');
const opportunitiesTableBody = document.getElementById('opportunities-table-body');
const exchangeGrid = document.getElementById('exchange-grid');
const currentTime = document.getElementById('current-time');
const exchangeCardTemplate = document.getElementById('exchange-card-template');

// Global State
let socket = null;
let dataStore = {
    opportunities: [],
    orderBooks: {},
    exchangeSet: new Set(),
    pairSet: new Set(),
    totalProfit: 0
};

// WebSocket Connection
function connectWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws`;
    
    socket = new WebSocket(wsUrl);
    
    socket.onopen = function() {
        setConnectionStatus(true);
        console.log('WebSocket connection established');
    };
    
    socket.onclose = function() {
        setConnectionStatus(false);
        console.log('WebSocket connection closed');
        
        // Try to reconnect after a delay
        setTimeout(connectWebSocket, 5000);
    };
    
    socket.onerror = function(error) {
        console.error('WebSocket error:', error);
        setConnectionStatus(false);
    };
    
    socket.onmessage = function(event) {
        const data = JSON.parse(event.data);
        processData(data);
    };
}

// Process received data
function processData(data) {
    // Process opportunities
    if (data.opportunities && data.opportunities.length > 0) {
        data.opportunities.forEach(opp => {
            // Only add if it's not already in our list
            if (!dataStore.opportunities.some(existing => 
                existing.timestamp === opp.timestamp && 
                existing.buyExchange === opp.buyExchange && 
                existing.sellExchange === opp.sellExchange)) {
                
                dataStore.opportunities.push(opp);
                dataStore.totalProfit += opp.netProfit;
                
                // Keep only the most recent 100 opportunities
                if (dataStore.opportunities.length > 100) {
                    dataStore.opportunities = dataStore.opportunities.slice(-100);
                }
            }
        });
    }
    
    // Process order books
    if (data.orderBooks) {
        dataStore.orderBooks = data.orderBooks;
        
        // Extract unique exchanges and pairs
        Object.values(data.orderBooks).forEach(book => {
            dataStore.exchangeSet.add(book.exchange);
            const pairKey = `${book.baseCurrency}/${book.quoteCurrency}`;
            dataStore.pairSet.add(pairKey);
        });
    }
    
    // Update the UI
    updateUI();
}

// Update all UI elements
function updateUI() {
    updateStats();
    updateOpportunitiesTable();
    updateExchangeGrid();
    updateTime();
}

// Update statistics cards
function updateStats() {
    activeExchanges.textContent = dataStore.exchangeSet.size;
    tradingPairs.textContent = dataStore.pairSet.size;
    totalOpportunities.textContent = dataStore.opportunities.length;
    potentialProfit.textContent = formatCurrency(dataStore.totalProfit);
}

// Update opportunities table
function updateOpportunitiesTable() {
    // Clear any "empty" messages
    if (opportunitiesTableBody.querySelector('.empty-row')) {
        opportunitiesTableBody.innerHTML = '';
    }
    
    // Sort opportunities by timestamp (newest first)
    const sortedOpportunities = [...dataStore.opportunities].sort((a, b) => 
        new Date(b.timestamp) - new Date(a.timestamp)
    );
    
    // Only display the 10 most recent opportunities
    const recentOpportunities = sortedOpportunities.slice(0, 10);
    
    // Create rows for each opportunity
    recentOpportunities.forEach(opp => {
        // Check if row already exists
        const existingRow = document.getElementById(`opp-${opp.timestamp}-${opp.buyExchange}-${opp.sellExchange}`);
        if (existingRow) return;
        
        const row = document.createElement('tr');
        row.id = `opp-${opp.timestamp}-${opp.buyExchange}-${opp.sellExchange}`;
        
        const time = new Date(opp.timestamp);
        const pairName = opp.baseCurrency && opp.quoteCurrency 
            ? `${opp.baseCurrency}/${opp.quoteCurrency}`
            : 'BTC/USDT'; // Default fallback
        
        row.innerHTML = `
            <td>${formatTime(time)}</td>
            <td>${pairName}</td>
            <td>${opp.buyExchange}</td>
            <td>${opp.sellExchange}</td>
            <td>${formatCurrency(opp.buyPrice)}</td>
            <td>${formatCurrency(opp.sellPrice)}</td>
            <td class="profit-positive">+${opp.profitPercentage.toFixed(2)}%</td>
            <td class="profit-positive">${formatCurrency(opp.netProfit)}</td>
        `;
        
        opportunitiesTableBody.insertBefore(row, opportunitiesTableBody.firstChild);
    });
    
    // Show empty message if no opportunities
    if (dataStore.opportunities.length === 0) {
        opportunitiesTableBody.innerHTML = `
            <tr class="empty-row">
                <td colspan="8">No opportunities detected yet</td>
            </tr>
        `;
    }
}

// Update exchange grid
function updateExchangeGrid() {
    // Clear any "empty" messages
    if (exchangeGrid.querySelector('.empty-exchanges')) {
        exchangeGrid.innerHTML = '';
    }
    
    // Skip if no order books
    if (Object.keys(dataStore.orderBooks).length === 0) {
        exchangeGrid.innerHTML = '<div class="empty-exchanges">No exchange data available</div>';
        return;
    }
    
    // Group order books by pair
    const pairGroups = {};
    
    Object.values(dataStore.orderBooks).forEach(book => {
        const pairKey = book.baseCurrency && book.quoteCurrency
            ? `${book.baseCurrency}/${book.quoteCurrency}`
            : book.symbol || 'Unknown';
        
        if (!pairGroups[pairKey]) {
            pairGroups[pairKey] = [];
        }
        
        pairGroups[pairKey].push(book);
    });
    
    // Create/update exchange cards for each book
    Object.entries(pairGroups).forEach(([pair, books]) => {
        books.forEach(book => {
            const cardId = `exchange-card-${book.exchange}-${pair}`;
            let card = document.getElementById(cardId);
            
            if (!card) {
                // Create new card from template
                const template = exchangeCardTemplate.content.cloneNode(true);
                card = template.querySelector('.exchange-card');
                card.id = cardId;
                
                // Set exchange name and pair
                card.querySelector('.exchange-name').textContent = book.exchange;
                card.querySelector('.pair-name').textContent = pair;
                
                exchangeGrid.appendChild(card);
            }
            
            // Update prices
            card.querySelector('.bid .value').textContent = formatCurrency(book.bid);
            card.querySelector('.ask .value').textContent = formatCurrency(book.ask);
            
            // Update time
            const lastUpdate = new Date(book.lastUpdate);
            card.querySelector('.last-update').textContent = `Updated: ${formatTime(lastUpdate)}`;
        });
    });
}

// Set connection status indicator
function setConnectionStatus(isConnected) {
    if (isConnected) {
        connectionStatus.className = 'connection-online';
        connectionStatus.innerHTML = '<i class="fas fa-circle"></i> Online';
    } else {
        connectionStatus.className = 'connection-offline';
        connectionStatus.innerHTML = '<i class="fas fa-circle"></i> Offline';
    }
}

// Update current time
function updateTime() {
    const now = new Date();
    currentTime.textContent = formatTime(now);
}

// Format time (HH:MM:SS)
function formatTime(date) {
    const hours = date.getHours().toString().padStart(2, '0');
    const minutes = date.getMinutes().toString().padStart(2, '0');
    const seconds = date.getSeconds().toString().padStart(2, '0');
    return `${hours}:${minutes}:${seconds}`;
}

// Format currency
function formatCurrency(value) {
    return new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'USD',
        minimumFractionDigits: 2,
        maximumFractionDigits: 2
    }).format(value);
}

// Initialize
function init() {
    // Update the time every second
    setInterval(updateTime, 1000);
    
    // Connect to WebSocket
    connectWebSocket();
    
    // Update UI initially
    updateUI();
    
    // Fetch initial data if WebSocket isn't available
    fetch('/api/status')
        .then(response => response.json())
        .then(data => processData(data))
        .catch(error => console.error('Error fetching initial data:', error));
}

// Start the application
document.addEventListener('DOMContentLoaded', init);