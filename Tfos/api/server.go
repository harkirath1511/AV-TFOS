// api/server.go - REST API and WebSocket
package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	router *gin.Engine
	nc     *nats.Conn
	rdb    *redis.Client
}

func NewServer() *Server {
	r := gin.Default()
	r.Use(cors.Default())

	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("Failed to connect to NATS:", err)
	}

	// Initialize Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	s := &Server{
		router: r,
		nc:     nc,
		rdb:    rdb,
	}

	s.routes()
	return s
}

func (s *Server) routes() {
	s.router.GET("/ws", s.handleWebSocket)
	s.router.POST("/simulation/start", s.startSimulation)
	s.router.POST("/emergency", s.handleEmergency)
	s.router.GET("/metrics", s.getMetrics)
}

func (s *Server) Start() {
	log.Fatal(s.router.Run(":8080"))
}

func (s *Server) handleWebSocket(c *gin.Context) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer ws.Close()

	// Subscribe to traffic updates
	sub, _ := s.nc.Subscribe("trafficlight.*", func(m *nats.Msg) {
		// Add logging before sending the message
		log.Printf("WebSocket sending data: %s", string(m.Data))

		if err := ws.WriteMessage(websocket.TextMessage, m.Data); err != nil {
			log.Printf("WebSocket write error: %v", err)
		}
	})

	defer sub.Unsubscribe()

	// Keep connection alive
	for {
		if _, _, err := ws.NextReader(); err != nil {
			log.Printf("WebSocket connection closed: %v", err)
			break
		}
	}
}

func (s *Server) startSimulation(c *gin.Context) {
	type simulationConfig struct {
		VehicleCount      int     `json:"vehicle_count"`
		AVPercentage      float64 `json:"av_percentage"`
		EmergencyInterval string  `json:"emergency_interval"`
		UpdateFrequency   string  `json:"update_freq"`
	}

	var config simulationConfig

	// Try to parse the JSON body, use defaults if not provided
	if err := c.BindJSON(&config); err != nil {
		log.Printf("Using default simulation config: %v", err)
		config = simulationConfig{
			VehicleCount:      500,
			AVPercentage:      0.2,
			EmergencyInterval: "10m",
			UpdateFrequency:   "1s",
		}
	}

	data, _ := json.Marshal(config)
	log.Printf("Starting simulation with config: %s", string(data))
	s.nc.Publish("simulation.control", data)
	c.JSON(200, gin.H{"status": "simulation_started", "config": config})
}

func (s *Server) handleEmergency(c *gin.Context) {
	var req struct {
		Route [][2]float64 `json:"route"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid_request"})
		return
	}

	data, _ := json.Marshal(req)
	s.nc.Publish("emergency.route", data)
	s.incrementCounter("emergencies") // New line to track emergency events
	c.JSON(200, gin.H{"status": "emergency_processed"})
}

func (s *Server) getCounter(key string) int64 {
	val, err := s.rdb.Get(context.Background(), "counter:"+key).Int64()
	if err != nil {
		return 0
	}
	return val
}

func (s *Server) getMetrics(c *gin.Context) {
	c.JSON(200, gin.H{
		"vehicles_processed":  s.getCounter("vehicles"),
		"emergencies_active":  s.getCounter("emergencies"),
		"messages_processed":  s.getCounter("nats_messages"),
		"optimization_cycles": s.getCounter("optimizations"),
	})
}

func (s *Server) incrementCounter(key string) {
	s.rdb.Incr(context.Background(), "counter:"+key)
}

func (s *Server) Cleanup() {
	if s.nc != nil {
		s.nc.Close()
	}
	if s.rdb != nil {
		s.rdb.Close()
	}
}
