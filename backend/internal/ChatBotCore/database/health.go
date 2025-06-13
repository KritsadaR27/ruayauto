package database

import (
	"context"
	"log"
	"sync"
	"time"
)

// HealthMonitor monitors database health and provides status information
type HealthMonitor struct {
	db             *DB
	interval       time.Duration
	timeout        time.Duration
	isHealthy      bool
	lastCheck      time.Time
	lastError      error
	mu             sync.RWMutex
	stopCh         chan struct{}
	stoppedCh      chan struct{}
	running        bool
}

// HealthStatus represents the current health status
type HealthStatus struct {
	IsHealthy     bool      `json:"is_healthy"`
	LastCheck     time.Time `json:"last_check"`
	LastError     string    `json:"last_error,omitempty"`
	PoolStats     *PoolStats `json:"pool_stats"`
	ResponseTime  time.Duration `json:"response_time"`
}

// PoolStats represents connection pool statistics
type PoolStats struct {
	AcquiredConns     int32 `json:"acquired_conns"`
	ConstructingConns int32 `json:"constructing_conns"`
	IdleConns         int32 `json:"idle_conns"`
	MaxConns          int32 `json:"max_conns"`
	TotalConns        int32 `json:"total_conns"`
}

// NewHealthMonitor creates a new health monitor
func NewHealthMonitor(db *DB, interval, timeout time.Duration) *HealthMonitor {
	return &HealthMonitor{
		db:        db,
		interval:  interval,
		timeout:   timeout,
		isHealthy: false,
		stopCh:    make(chan struct{}),
		stoppedCh: make(chan struct{}),
	}
}

// Start begins health monitoring in a background goroutine
func (hm *HealthMonitor) Start() {
	hm.mu.Lock()
	if hm.running {
		hm.mu.Unlock()
		return
	}
	hm.running = true
	hm.mu.Unlock()

	go hm.monitor()
	log.Printf("Database health monitor started (interval: %v, timeout: %v)", hm.interval, hm.timeout)
}

// Stop stops the health monitoring
func (hm *HealthMonitor) Stop() {
	hm.mu.Lock()
	if !hm.running {
		hm.mu.Unlock()
		return
	}
	hm.mu.Unlock()

	close(hm.stopCh)
	<-hm.stoppedCh

	hm.mu.Lock()
	hm.running = false
	hm.mu.Unlock()

	log.Printf("Database health monitor stopped")
}

// monitor runs the health check loop
func (hm *HealthMonitor) monitor() {
	defer close(hm.stoppedCh)

	ticker := time.NewTicker(hm.interval)
	defer ticker.Stop()

	// Perform initial health check
	hm.performHealthCheck()

	for {
		select {
		case <-ticker.C:
			hm.performHealthCheck()
		case <-hm.stopCh:
			return
		}
	}
}

// performHealthCheck executes a health check and updates status
func (hm *HealthMonitor) performHealthCheck() {
	start := time.Now()
	
	ctx, cancel := context.WithTimeout(context.Background(), hm.timeout)
	defer cancel()

	err := hm.db.HealthCheck(ctx)
	responseTime := time.Since(start)

	hm.mu.Lock()
	hm.lastCheck = time.Now()
	hm.lastError = err
	hm.isHealthy = (err == nil)
	hm.mu.Unlock()

	if err != nil {
		log.Printf("Database health check failed (took %v): %v", responseTime, err)
	} else {
		log.Printf("Database health check passed (took %v)", responseTime)
	}
}

// GetStatus returns the current health status
func (hm *HealthMonitor) GetStatus() HealthStatus {
	hm.mu.RLock()
	defer hm.mu.RUnlock()

	stats := hm.db.Stats()
	poolStats := &PoolStats{
		AcquiredConns:     stats.AcquiredConns(),
		ConstructingConns: stats.ConstructingConns(),
		IdleConns:         stats.IdleConns(),
		MaxConns:          stats.MaxConns(),
		TotalConns:        stats.TotalConns(),
	}

	status := HealthStatus{
		IsHealthy:    hm.isHealthy,
		LastCheck:    hm.lastCheck,
		PoolStats:    poolStats,
	}

	if hm.lastError != nil {
		status.LastError = hm.lastError.Error()
	}

	return status
}

// IsHealthy returns whether the database is currently healthy
func (hm *HealthMonitor) IsHealthy() bool {
	hm.mu.RLock()
	defer hm.mu.RUnlock()
	return hm.isHealthy
}

// WaitForHealth waits for the database to become healthy or times out
func (hm *HealthMonitor) WaitForHealth(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if hm.IsHealthy() {
				return nil
			}
		}
	}
}
