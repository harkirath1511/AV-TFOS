package nats

import (
	"context"
	"encoding/json"
	"flowsync-backend/traffic"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
)

type EmergencyRoute struct {
	Route [][2]float64 `json:"route"`
}

// HandleEmergencies processes emergency vehicle route requests
func HandleEmergencies(ctx context.Context) {
	if js == nil {
		log.Fatal("JetStream context not initialized")
		return
	}

	sub, err := js.Subscribe("emergency.route", func(m *nats.Msg) {
		var route EmergencyRoute
		if err := json.Unmarshal(m.Data, &route); err != nil {
			log.Printf("Invalid emergency route data: %v", err)
			return
		}

		if err := processEmergencyRoute(route); err != nil {
			log.Printf("Failed to process emergency route: %v", err)
			return
		}

		if err := m.Ack(); err != nil {
			log.Printf("Failed to acknowledge message: %v", err)
		}
	}, nats.ManualAck())

	if err != nil {
		log.Printf("Failed to subscribe to emergency routes: %v", err)
		return
	}

	<-ctx.Done()
	if err := sub.Unsubscribe(); err != nil {
		log.Printf("Failed to unsubscribe: %v", err)
	}
}

func processEmergencyRoute(route EmergencyRoute) error {
	if len(route.Route) < 2 {
		return fmt.Errorf("route must contain at least 2 points")
	}

	// Find affected intersections along the route
	for i := 0; i < len(route.Route)-1; i++ {
		intersection := findNearestIntersection(route.Route[i])
		if intersection != "" {
			// Prioritize emergency vehicle passage
			if err := traffic.AdjustIntersection(intersection, map[string]interface{}{
				"alert":      "emergency_vehicle",
				"position":   route.Route[i],
				"next_point": route.Route[i+1],
				"priority":   "high",
			}); err != nil {
				return fmt.Errorf("failed to adjust intersection %s: %v", intersection, err)
			}
		}
	}
	return nil
}

func findNearestIntersection(position [2]float64) string {
	// Find closest intersection using Redis geospatial query
	results, err := rdb.GeoRadius(context.Background(), "intersections:locations",
		position[1], position[0], &redis.GeoRadiusQuery{
			Radius: 0.1, // 100 meters
			Unit:   "km",
			Count:  1, // Get closest intersection only
		}).Result()

	if err != nil || len(results) == 0 {
		return ""
	}

	return results[0].Name
}
