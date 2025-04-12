# Monitoring Guide

This guide outlines the monitoring and observability setup for the APEX system, including metrics collection, alerting, and debugging procedures.

## Monitoring Architecture

### Components

1. **Metrics Collection**
   - Prometheus for time-series metrics
   - OpenTelemetry for distributed tracing
   - Grafana for visualization
   - Loki for log aggregation

2. **Health Checks**
   - Exchange connectivity status
   - WebSocket connection health
   - Database connectivity
   - Memory/CPU utilization

### Infrastructure Setup

```yaml
# docker-compose.monitoring.yml
version: '3'
services:
  prometheus:
    image: prom/prometheus:latest
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    ports:
      - "9090:9090"

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    depends_on:
      - prometheus

  loki:
    image: grafana/loki:latest
    ports:
      - "3100:3100"
```

## Key Metrics

### System Metrics

1. **Resource Utilization**
   ```go
   // CPU Usage
   var cpuUsage = prometheus.NewGauge(prometheus.GaugeOpts{
       Name: "apex_cpu_usage_percent",
       Help: "Current CPU usage percentage",
   })

   // Memory Usage
   var memoryUsage = prometheus.NewGauge(prometheus.GaugeOpts{
       Name: "apex_memory_usage_bytes",
       Help: "Current memory usage in bytes",
   })
   ```

2. **Network Metrics**
   ```go
   // Network Traffic
   var networkTraffic = prometheus.NewCounterVec(
       prometheus.CounterOpts{
           Name: "apex_network_traffic_bytes",
           Help: "Network traffic in bytes",
       },
       []string{"direction"}, // in/out
   )
   ```

### Business Metrics

1. **Exchange Connectivity**
   ```go
   // Exchange Connection Status
   var exchangeStatus = prometheus.NewGaugeVec(
       prometheus.GaugeOpts{
           Name: "apex_exchange_status",
           Help: "Exchange connection status (1=connected, 0=disconnected)",
       },
       []string{"exchange"},
   )
   ```

2. **Arbitrage Opportunities**
   ```go
   // Opportunity Metrics
   var opportunityMetrics = prometheus.NewHistogramVec(
       prometheus.HistogramOpts{
           Name:    "apex_arbitrage_opportunity",
           Help:    "Arbitrage opportunity details",
           Buckets: prometheus.LinearBuckets(0, 0.5, 20),
       },
       []string{"pair", "exchanges"},
   )
   ```

## Alerting Rules

### Critical Alerts

```yaml
# prometheus/alerts.yml
groups:
  - name: apex_alerts
    rules:
      - alert: HighMemoryUsage
        expr: apex_memory_usage_bytes > 1.8e9  # 1.8GB
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: High memory usage detected

      - alert: ExchangeDisconnected
        expr: apex_exchange_status == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: Exchange connection lost
```

### Warning Alerts

```yaml
      - alert: HighLatency
        expr: rate(apex_arbitrage_detection_latency_seconds_sum[5m]) > 0.01
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: High arbitrage detection latency

      - alert: LowOpportunityRate
        expr: rate(apex_arbitrage_opportunity_total[15m]) < 1
        for: 15m
        labels:
          severity: warning
        annotations:
          summary: Low arbitrage opportunity detection rate
```

## Logging

### Log Levels

1. **Debug Logging**
   ```go
   log.WithFields(log.Fields{
       "exchange": exchange.Name,
       "pair":     pair,
       "latency":  latency,
   }).Debug("Market data update received")
   ```

2. **Error Logging**
   ```go
   log.WithFields(log.Fields{
       "error":    err,
       "exchange": exchange.Name,
       "attempt":  retryCount,
   }).Error("Failed to connect to exchange")
   ```

### Log Aggregation

```yaml
# loki/config.yml
auth_enabled: false

server:
  http_listen_port: 3100

ingester:
  lifecycler:
    ring:
      kvstore:
        store: inmemory
      replication_factor: 1
    final_sleep: 0s
  chunk_idle_period: 5m
  chunk_retain_period: 30s
```

## Dashboards

### System Overview

1. **Resource Dashboard**
   ```json
   {
     "title": "APEX System Overview",
     "panels": [
       {
         "title": "CPU Usage",
         "type": "graph",
         "datasource": "Prometheus",
         "targets": [
           {
             "expr": "rate(apex_cpu_usage_percent[5m])"
           }
         ]
       }
     ]
   }
   ```

2. **Business Metrics Dashboard**
   ```json
   {
     "title": "APEX Business Metrics",
     "panels": [
       {
         "title": "Arbitrage Opportunities",
         "type": "graph",
         "datasource": "Prometheus",
         "targets": [
           {
             "expr": "rate(apex_arbitrage_opportunity_total[5m])"
           }
         ]
       }
     ]
   }
   ```

## Debugging Procedures

### Common Issues

1. **High Memory Usage**
   ```bash
   # Generate heap profile
   curl -s http://localhost:6060/debug/pprof/heap > heap.prof
   go tool pprof heap.prof
   ```

2. **Performance Issues**
   ```bash
   # CPU profile collection
   curl -s http://localhost:6060/debug/pprof/profile > cpu.prof
   go tool pprof cpu.prof
   ```

### Troubleshooting Steps

1. **Exchange Connectivity**
   - Check exchange API status
   - Verify network connectivity
   - Review rate limits
   - Check API credentials

2. **Data Quality**
   - Validate order book consistency
   - Check for stale data
   - Monitor update frequencies
   - Verify price accuracy

## Maintenance

### Regular Tasks

1. **Log Rotation**
   ```bash
   # logrotate configuration
   /var/log/apex/*.log {
       daily
       rotate 7
       compress
       delaycompress
       missingok
       notifempty
   }
   ```

2. **Metric Retention**
   ```yaml
   # prometheus/config.yml
   global:
     scrape_interval: 15s
     evaluation_interval: 15s
     
   storage:
     tsdb:
       retention.time: 15d
       retention.size: 50GB
   ```

### Backup Procedures

1. **Metrics Backup**
   ```bash
   # Prometheus data backup
   tar czf prometheus-data-$(date +%Y%m%d).tar.gz /path/to/prometheus/data
   ```

2. **Log Backup**
   ```bash
   # Log archival
   find /var/log/apex -type f -name "*.gz" -mtime +30 -exec aws s3 cp {} s3://apex-logs/ \;
   ```

## Best Practices

1. **Monitoring**
   - Set appropriate thresholds
   - Implement gradual alerts
   - Monitor all critical paths
   - Keep dashboards simple

2. **Logging**
   - Use structured logging
   - Include relevant context
   - Implement log levels
   - Rotate logs regularly

3. **Alerting**
   - Define clear severity levels
   - Avoid alert fatigue
   - Include actionable information
   - Set up escalation paths

## References

- [Prometheus Best Practices](https://prometheus.io/docs/practices/naming/)
- [Grafana Documentation](https://grafana.com/docs/)
- [OpenTelemetry Guide](https://opentelemetry.io/docs/)
- [Loki Documentation](https://grafana.com/docs/loki/latest/) 