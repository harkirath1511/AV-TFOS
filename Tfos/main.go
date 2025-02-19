// main.go - Entry point
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"flowsync-backend/api"
	"flowsync-backend/nats"
	"flowsync-backend/traffic"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize components
	nats.Initialize()
	traffic.InitializeOptimizer(ctx)
	apiServer := api.NewServer()

	// Start services
	go nats.StartAVIngest(ctx)
	go nats.HandleEmergencies(ctx)
	go traffic.RunSimulationEngine(ctx)
	go apiServer.Start()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down...")
	cancel()
	time.Sleep(2 * time.Second) // Allow cleanup
}