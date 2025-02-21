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
		ws.WriteMessage(websocket.TextMessage, m.Data)
	})

	defer sub.Unsubscribe()

	// Keep connection alive
	for {
		if _, _, err := ws.NextReader(); err != nil {
			break
		}
	}
}

func (s *Server) startSimulation(c *gin.Context) {
	// Start simulation logic
	s.nc.Publish("simulation.control", []byte(`{"action":"start"}`))
	c.JSON(200, gin.H{"status": "simulation_started"})
}

// Replace the existing handleEmergency function with this updated version
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
