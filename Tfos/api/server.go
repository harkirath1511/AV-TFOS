// api/server.go - REST API and WebSocket
package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
)

type Server struct {
	router *gin.Engine
	nc     *nats.Conn
}

func NewServer() *Server {
	r := gin.Default()
	r.Use(cors.Default())

	s := &Server{
		router: r,
		nc:     nats.Cleanup(),
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
	c.JSON(200, gin.H{"status": "emergency_processed"})
}

func (s *Server) getMetrics(c *gin.Context) {
	// Return system metrics
	c.JSON(200, gin.H{
		"vehicles_processed":   getCounter("vehicles"),
		"emergencies_active":   getCounter("emergencies"),
		"messages_processed":   getCounter("nats_messages"),
		"optimization_cycles":  getCounter("optimizations"),
	})
}