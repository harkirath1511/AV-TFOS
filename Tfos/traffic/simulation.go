// traffic/simulation.go - Traffic Simulation
package traffic

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/nats-io/nats.go"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Update the SimulationConfig struct
type SimulationConfig struct {
	VehicleCount      int     `json:"vehicle_count"`
	AVPercentage      float64 `json:"av_percentage"`
	EmergencyInterval string  `json:"emergency_interval"`
	UpdateFrequency   string  `json:"update_freq,omitempty"` // New field
}

// Modify the RunSimulationEngine function
func RunSimulationEngine(ctx context.Context) {
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Close()

	sub, _ := nc.Subscribe("simulation.control", func(m *nats.Msg) {
		var config SimulationConfig
		if err := json.Unmarshal(m.Data, &config); err != nil {
			log.Printf("Error parsing simulation config: %v", err)
			return
		}

		// Set default values if not provided
		if config.VehicleCount == 0 {
			config.VehicleCount = 500
		}
		if config.AVPercentage == 0 {
			config.AVPercentage = 0.2
		}
		if config.UpdateFrequency == "" {
			config.UpdateFrequency = "1s" // Default update frequency
		}

		go runSimulation(ctx, config)
	})

	<-ctx.Done()
	sub.Unsubscribe()
}

// Update the runSimulation function
func runSimulation(ctx context.Context, config SimulationConfig) {
	// Parse update frequency
	updateFreq, err := time.ParseDuration(config.UpdateFrequency)
	if err != nil {
		updateFreq = 1 * time.Second // fallback to default if parsing fails
	}

	ticker := time.NewTicker(updateFreq)
	defer ticker.Stop()

	for i := 0; i < config.VehicleCount; i++ {
		go simulateVehicle(ctx, i, config, ticker.C)
	}
}

// formatVehicleID generates a unique vehicle identifier
func formatVehicleID(id int) string {
	return fmt.Sprintf("vehicle_%04d", id)
}

// randomPosition generates a random GPS position within simulation bounds
func randomPosition() [2]float64 {
	return [2]float64{
		34.0522 + (rand.Float64()-0.5)*0.1,  // Los Angeles latitude ±0.05°
		118.2437 + (rand.Float64()-0.5)*0.1, // Los Angeles longitude ±0.05°
	}
}

// updatePosition simulates vehicle movement
func updatePosition(current [2]float64) [2]float64 {
	return [2]float64{
		current[0] + (rand.Float64()-0.5)*0.001, // Small random movement
		current[1] + (rand.Float64()-0.5)*0.001,
	}
}

// generateIntent determines vehicle's next action
func generateIntent(avPercentage float64) string {
	intents := []string{
		"maintain_lane",
		"turn_left",
		"turn_right",
		"stop",
	}
	return intents[rand.Intn(len(intents))]
}

// Update simulateVehicle to use nc for publishing
func simulateVehicle(ctx context.Context, id int, config SimulationConfig, updates <-chan time.Time) {
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Close()

	vehicleID := formatVehicleID(id)
	position := randomPosition()
	intent := "maintain_lane"

	for {
		select {
		case <-updates:
			position = updatePosition(position)
			intent = generateIntent(config.AVPercentage)

			data, _ := json.Marshal(map[string]interface{}{
				"vehicle_id": vehicleID,
				"position":   position,
				"speed":      rand.Float64()*100 + 20,
				"intent":     intent,
				"timestamp":  time.Now().Unix(),
			})

			nc.Publish("av.telemetry."+vehicleID, data)

		case <-ctx.Done():
			return
		}
	}
}
