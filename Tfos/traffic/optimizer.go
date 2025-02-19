// traffic/optimizer.go - Traffic Light Control
package traffic

import (
	"context"
	"encoding/json"
	"math"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	rdb          = redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	optimizeInterval = 5 * time.Second
)

type IntersectionConfig struct {
	ID              string
	DefaultCycle    map[string]time.Duration
	CurrentPhase    string
	PhaseStartTime  time.Time
	EmergencyActive bool
}

func InitializeOptimizer(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(optimizeInterval)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				optimizeAllIntersections()
			case <-ctx.Done():
				return
			}
		}
	}()
}

func optimizeAllIntersections() {
	// Get all intersections from Redis
	ids, _ := rdb.SMembers(context.Background(), "intersections:active").Result()

	for _, id := range ids {
		cfg := getIntersectionConfig(id)
		newPhase := calculateOptimalPhase(cfg)
		applyPhase(cfg, newPhase)
	}
}

func calculateOptimalPhase(cfg IntersectionConfig) map[string]time.Duration {
	// Simple round-robin implementation
	// Replace with actual optimization logic
	return map[string]time.Duration{
		"north_south": 30 * time.Second,
		"east_west":   25 * time.Second,
	}
}

func applyPhase(cfg IntersectionConfig, phase map[string]time.Duration) {
	data, _ := json.Marshal(phase)
	rdb.Publish(context.Background(), "trafficlight:"+cfg.ID, data)
}