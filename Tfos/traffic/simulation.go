// traffic/simulation.go - Traffic Simulation
package traffic

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/nats-io/nats.go"
)

type SimulationConfig struct {
	VehicleCount     int           `json:"vehicle_count"`
	AVPercentage     float64       `json:"av_percentage"`
	UpdateFrequency  time.Duration `json:"update_freq"`
}

func RunSimulationEngine(ctx context.Context) {
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Close()

	sub, _ := nc.Subscribe("simulation.control", func(m *nats.Msg) {
		var config SimulationConfig
		json.Unmarshal(m.Data, &config)
		go runSimulation(ctx, config)
	})

	<-ctx.Done()
	sub.Unsubscribe()
}

func runSimulation(ctx context.Context, config SimulationConfig) {
	ticker := time.NewTicker(config.UpdateFrequency)
	defer ticker.Stop()

	for i := 0; i < config.VehicleCount; i++ {
		go simulateVehicle(ctx, i, config, ticker.C)
	}
}

func simulateVehicle(ctx context.Context, id int, config SimulationConfig, updates <-chan time.Time) {
	vehicleID := formatVehicleID(id)
	position := randomPosition()
	intent := "maintain_lane"

	for {
		select {
		case <-updates:
			position = updatePosition(position)
			intent = generateIntent(config.AVPercentage)

			data, _ := json.Marshal(map[string]interface{}{
				"vehicle_id":  vehicleID,
				"position":    position,
				"speed":       rand.Float64()*100 + 20,
				"intent":      intent,
				"timestamp":   time.Now().Unix(),
			})

			nats.Publish("av.telemetry."+vehicleID, data)

		case <-ctx.Done():
			return
		}
	}
}