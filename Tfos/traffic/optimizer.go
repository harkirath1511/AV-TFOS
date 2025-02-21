// traffic/optimizer.go - Traffic Light Control
package traffic

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	rdb              = redis.NewClient(&redis.Options{Addr: "localhost:6379"})
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

func getIntersectionConfig(id string) IntersectionConfig {
	var cfg IntersectionConfig

	// Get intersection data from Redis
	data, err := rdb.HGetAll(context.Background(), "intersection:"+id).Result()
	if err != nil {
		return cfg
	}

	// Parse the basic fields
	cfg.ID = id
	cfg.CurrentPhase = data["current_phase"]
	cfg.EmergencyActive = data["emergency_active"] == "true"

	// Parse phase start time
	if timeStr, ok := data["phase_start_time"]; ok {
		timestamp, err := time.Parse(time.RFC3339, timeStr)
		if err == nil {
			cfg.PhaseStartTime = timestamp
		}
	}

	// Parse default cycle from JSON
	if cycleJSON, ok := data["default_cycle"]; ok {
		var cycle map[string]float64
		if err := json.Unmarshal([]byte(cycleJSON), &cycle); err == nil {
			cfg.DefaultCycle = make(map[string]time.Duration)
			for phase, seconds := range cycle {
				cfg.DefaultCycle[phase] = time.Duration(seconds * float64(time.Second))
			}
		}
	}

	return cfg
}
