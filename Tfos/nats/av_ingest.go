// nats/av_ingest.go - AV Telemetry Handling
package nats

import (
	"context"
	"encoding/json"
	"log"
	"time"
	"flowsync-backend/traffic"

	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
)

type AVTelemetry struct {
	VehicleID      string     `json:"vehicle_id"`
	Position       [2]float64 `json:"position"` // [latitude, longitude]
	Speed          float64    `json:"speed"`    // in km/h
	Intent         string     `json:"intent"`   // e.g., "turn_left", "maintain_lane"
	IntersectionID string     `json:"intersection_id"`
	Timestamp      int64      `json:"timestamp"`
}

var (
	nc                 *nats.Conn
	js                 nats.JetStreamContext
	rdb                = redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	vehiclePositionTTL = 5 * time.Minute
)

func Initialize() {
	var err error
	nc, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("NATS connection failed:", err)
	}

	js, err = nc.JetStream()
	if err != nil {
		log.Fatal("JetStream init failed:", err)
	}

	configureStreams()
}

func configureStreams() {
	js.AddStream(&nats.StreamConfig{
		Name:      "TRAFFIC",
		Subjects:  []string{"av.telemetry.*", "emergency.*"},
		Retention: nats.InterestPolicy,
	})
}

func StartAVIngest(ctx context.Context) {
	sub, _ := js.Subscribe("av.telemetry.*", func(m *nats.Msg) {
		var telemetry AVTelemetry
		if err := json.Unmarshal(m.Data, &telemetry); err != nil {
			log.Printf("Invalid telemetry data: %v", err)
			return
		}

		processTelemetry(telemetry)
		m.Ack()
	}, nats.ManualAck())

	<-ctx.Done()
	sub.Unsubscribe()
}

func processTelemetry(t AVTelemetry) {
	// Store position
	rdb.GeoAdd(context.Background(), "vehicles:positions", &redis.GeoLocation{
		Name:      t.VehicleID,
		Longitude: t.Position[1],
		Latitude:  t.Position[0],
	}).Err()

	// Update vehicle metadata
	rdb.HSet(context.Background(), "vehicle:"+t.VehicleID,
		"speed", t.Speed,
		"intent", t.Intent,
		"lastSeen", time.Now().Unix(),
	).Err()

	// Predict collisions
	if isCollisionRisk(t) {
		handleCollisionRisk(t)
	}
}

func isCollisionRisk(t AVTelemetry) bool {
	res, _ := rdb.GeoRadius(context.Background(), "vehicles:positions",
		t.Position[0], t.Position[1], &redis.GeoRadiusQuery{
			Radius:      0.05, // 50 meters
			Unit:        "km",
			WithDist:    true,
			WithCoord:   true,
			WithGeoHash: true,
		}).Result()

	return len(res) > 3 // More than 3 vehicles in 50m radius
}

func handleCollisionRisk(t AVTelemetry) {
	// Send alert to vehicle
	js.Publish("av.alerts."+t.VehicleID, []byte(`{"type":"collision_warning"}`))

	// Update traffic lights
	traffic.AdjustIntersection(t.IntersectionID, map[string]interface{}{
		"alert":       "collision_risk",
		"vehicle_id":  t.VehicleID,
		"coordinates": t.Position,
	})
}
