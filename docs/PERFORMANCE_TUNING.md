# Performance Tuning Guide

This guide provides detailed information on optimizing APEX for maximum performance across different deployment scenarios and workloads.

## Performance Metrics

### Key Performance Indicators (KPIs)

1. **Latency Metrics**
   - Order book update latency: < 1ms
   - Arbitrage detection time: < 5ms
   - End-to-end opportunity notification: < 10ms
   - API response time: < 50ms

2. **Throughput Metrics**
   - Market updates per second: 10,000+
   - Concurrent WebSocket connections: 1,000+
   - API requests per second: 5,000+
   - Opportunities detected per second: 100+

3. **Resource Utilization**
   - CPU usage per instance: < 70%
   - Memory usage per instance: < 2GB
   - Network bandwidth: < 100Mbps
   - Disk I/O: < 1000 IOPS

## System Requirements

### Minimum Requirements
```
CPU: 4 cores
Memory: 8GB RAM
Storage: 100GB SSD
Network: 1Gbps
```

### Recommended Requirements
```
CPU: 8+ cores
Memory: 16GB+ RAM
Storage: 500GB NVMe SSD
Network: 10Gbps
```

## Go Runtime Optimization

### Memory Management

1. **Garbage Collection Tuning**
   ```bash
   # Environment variables for GC tuning
   GOGC=100
   GOMEMLIMIT=2048MiB
   ```

2. **Memory Allocation**
   ```go
   // Use sync.Pool for frequently allocated objects
   var orderBookPool = sync.Pool{
       New: func() interface{} {
           return &OrderBook{
               Bids: make(map[float64]float64, 1000),
               Asks: make(map[float64]float64, 1000),
           }
       },
   }
   ```

### CPU Optimization

1. **Goroutine Management**
   ```go
   // Worker pool implementation
   type WorkerPool struct {
       maxWorkers int
       tasks      chan Task
       workers    []*Worker
   }
   ```

2. **Lock Contention**
   ```go
   // Use atomic operations where possible
   var counter uint64
   atomic.AddUint64(&counter, 1)
   ```

## Network Optimization

### WebSocket Tuning

1. **Connection Settings**
   ```go
   websocket.Upgrader{
       ReadBufferSize:  1024,
       WriteBufferSize: 1024,
       WriteBufferPool: &sync.Pool{},
   }
   ```

2. **Heartbeat Configuration**
   ```go
   const (
       pingPeriod = 25 * time.Second
       pongWait   = 60 * time.Second
   )
   ```

### HTTP Optimization

1. **Server Settings**
   ```go
   server := &http.Server{
       ReadTimeout:  5 * time.Second,
       WriteTimeout: 10 * time.Second,
       IdleTimeout:  120 * time.Second,
   }
   ```

2. **Connection Pooling**
   ```go
   transport := &http.Transport{
       MaxIdleConns:        100,
       MaxIdleConnsPerHost: 10,
       IdleConnTimeout:     90 * time.Second,
   }
   ```

## Database Optimization

### Time Series Database

1. **Batch Processing**
   ```go
   type BatchWriter struct {
       batchSize int
       buffer    []Record
       flushChan chan struct{}
   }
   ```

2. **Index Optimization**
   ```sql
   CREATE INDEX idx_timestamp ON opportunities (timestamp DESC);
   CREATE INDEX idx_profit ON opportunities (profit_percentage DESC);
   ```

## Monitoring and Profiling

### Runtime Profiling

1. **CPU Profiling**
   ```go
   import "runtime/pprof"
   
   f, _ := os.Create("cpu.prof")
   pprof.StartCPUProfile(f)
   defer pprof.StopCPUProfile()
   ```

2. **Memory Profiling**
   ```go
   import "runtime/pprof"
   
   f, _ := os.Create("mem.prof")
   pprof.WriteHeapProfile(f)
   defer f.Close()
   ```

### Metrics Collection

1. **Prometheus Integration**
   ```go
   // Latency histogram
   var detectionLatency = prometheus.NewHistogram(
       prometheus.HistogramOpts{
           Name:    "arbitrage_detection_latency_seconds",
           Help:    "Arbitrage detection latency in seconds",
           Buckets: prometheus.ExponentialBuckets(0.001, 2, 10),
       },
   )
   ```

## Load Testing

### Test Scenarios

1. **Market Data Simulation**
   ```go
   type MarketSimulator struct {
       UpdateRate     int
       VolatilityPct float64
       NumPairs      int
   }
   ```

2. **Client Load Testing**
   ```bash
   # Using k6 for load testing
   k6 run --vus 100 --duration 30s websocket_test.js
   ```

## Scaling Strategies

### Horizontal Scaling

1. **Load Balancer Configuration**
   ```nginx
   upstream apex_servers {
       least_conn;
       server apex1:8080;
       server apex2:8080;
       server apex3:8080;
   }
   ```

2. **Service Discovery**
   ```yaml
   # Docker Compose configuration
   version: '3'
   services:
     apex:
       image: apex:latest
       deploy:
         replicas: 3
   ```

### Vertical Scaling

1. **Resource Allocation**
   ```bash
   # Container resource limits
   docker run -d \
     --cpus=4 \
     --memory=8g \
     apex:latest
   ```

## Troubleshooting

### Common Issues

1. **High Memory Usage**
   - Check for memory leaks
   - Monitor GC frequency
   - Review object allocation patterns

2. **High CPU Usage**
   - Profile hot spots
   - Optimize algorithms
   - Review goroutine patterns

3. **Network Latency**
   - Check network configuration
   - Monitor connection pools
   - Review timeout settings

### Performance Debugging

1. **Logging**
   ```go
   // Performance logging
   func logPerformanceMetrics(metrics *Metrics) {
       log.WithFields(log.Fields{
           "latency":    metrics.Latency,
           "throughput": metrics.Throughput,
           "memory":     metrics.MemoryUsage,
       }).Info("Performance metrics")
   }
   ```

2. **Tracing**
   ```go
   // OpenTelemetry integration
   tracer := otel.Tracer("apex")
   ctx, span := tracer.Start(ctx, "arbitrage_detection")
   defer span.End()
   ```

## Best Practices

1. **Code Optimization**
   - Use efficient data structures
   - Minimize allocations
   - Implement caching strategies

2. **Resource Management**
   - Monitor resource usage
   - Implement graceful degradation
   - Set appropriate limits

3. **Testing**
   - Regular performance testing
   - Continuous monitoring
   - Benchmark critical paths

## References

- [Go Performance](https://golang.org/doc/performance)
- [pprof Documentation](https://golang.org/pkg/runtime/pprof/)
- [WebSocket Best Practices](https://tools.ietf.org/html/rfc6455)
- [System Monitoring Guide](docs/MONITORING.md) 