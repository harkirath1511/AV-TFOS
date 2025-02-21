# AV-TFOS (Autonomous Vehicle Traffic Flow Optimization System)

## Prerequisites

### Install Docker Desktop
```bash
# Download and install from
https://www.docker.com/products/docker-desktop
```

### Install NATS Server
```bash
# Download latest release from
https://github.com/nats-io/nats-server/releases
# Add to PATH and verify
nats-server --version
```

### Install Go
```bash
# Download and install from
https://golang.org/dl/
# Verify installation
go version
```

### Install Node.js
```bash
# Download and install from
https://nodejs.org/
# Verify installation
node --version
npm --version
```

## Starting the Backend

### 1. Start Required Services

```bash
# Start NATS with JetStream enabled
nats-server -js

# Start Redis using Docker (in a new terminal)
docker run --name redis -p 6379:6379 -d redis
```

### 2. Initialize Redis with Sample Data

```bash
# Open Redis CLI in Docker container
docker exec -it redis redis-cli

# Add sample intersection data
GEOADD intersections:locations -118.2437 34.0522 intersection_001
SADD intersections:active intersection_001
HSET intersection:intersection_001 current_phase north_south emergency_active false default_cycle "{\"north_south\":30,\"east_west\":25}"
```

### 3. Setup Go Backend

```bash
# Navigate to backend directory
cd Tfos

# Initialize Go modules
go mod init flowsync-backend
go mod tidy

# Install required dependencies
go get github.com/nats-io/nats.go@latest
go get github.com/redis/go-redis/v9@latest
go get github.com/gin-gonic/gin@latest
go get github.com/gorilla/websocket@latest
```

### 4. Run the Backend

```bash
# From the Tfos directory
go run main.go
```

## Starting the Frontend

```bash
# Navigate to frontend directory
cd clientt

# Install dependencies
npm install

# Start development server
npm run dev
```

## Testing the System

### 1. Start Traffic Simulation
```bash
# Start simulation with default parameters
curl -X POST http://localhost:8080/simulation/start

# Or start with custom configuration
curl -X POST http://localhost:8080/simulation/start \
  -H "Content-Type: application/json" \
  -d '{
    "vehicle_count": 1000,     # Number of vehicles to simulate
    "av_percentage": 0.3,      # Percentage of autonomous vehicles (0.0 to 1.0)
    "emergency_interval": "5m"  # Interval between emergency vehicles
  }'
```

#### Simulation Parameters
- `vehicle_count`: Number of vehicles to simulate (default: 500)
- `av_percentage`: Percentage of vehicles that are autonomous (default: 0.2)
- `emergency_interval`: Time between emergency vehicle spawns (default: "10m")
  - Format: "1h" (hour), "5m" (minutes), "30s" (seconds)

### 2. Test Emergency Route
```bash
# Send emergency vehicle route
curl -X POST http://localhost:8080/emergency \
  -H "Content-Type: application/json" \
  -d '{
    "route": [
      [34.0522, -118.2437],
      [34.0523, -118.2438]
    ]
  }'
```

### Emergency Vehicle Integration

The emergency route endpoint allows dispatching emergency vehicles (ambulances, fire trucks, etc.) through the traffic system.

```bash
# Dispatch an emergency vehicle
curl -X POST http://localhost:8080/emergency \
  -H "Content-Type: application/json" \
  -d '{
    "ambulance_id": "medic-1",    # Unique identifier for the emergency vehicle
    "path": [                      # Array of [latitude, longitude] coordinates
      [37.7749, -122.4194],       # Starting point
      [37.7813, -122.4168]        # Destination
    ],
    "priority": 1                  # Priority level (1: highest, 3: lowest)
  }'
```

#### Emergency Parameters
- `ambulance_id`: Unique identifier for the emergency vehicle
- `path`: Array of coordinates representing the route
  - Each point is [latitude, longitude]
  - Must have at least start and end points
- `priority`: Emergency priority level
  - 1: Critical (e.g., cardiac arrest)
  - 2: Urgent (e.g., serious injury)
  - 3: Non-urgent (e.g., minor injury)

When an emergency route is submitted:
1. Traffic lights along the route are optimized
2. Other vehicles are notified to clear the path
3. Real-time updates are sent via WebSocket

### 3. Monitor System
```bash
# Check metrics
curl http://localhost:8080/metrics

# View real-time updates in frontend
http://localhost:5173
```

## Available Endpoints

### WebSocket
- `ws://localhost:8080/ws` - Real-time traffic updates

### REST API
- `POST /simulation/start` - Start traffic simulation
- `POST /emergency` - Submit emergency vehicle route
- `GET /metrics` - System metrics

## Cleanup

```bash
# Stop the backend (Ctrl+C)
# Stop Redis
docker stop redis
docker rm redis

# Stop NATS (Ctrl+C)
```

## Troubleshooting

### Redis Connection Issues
```bash
# Check Redis container status
docker ps -a
# Restart if needed
docker restart redis
```

### NATS Connection Issues
```bash
# Verify NATS is running
nats-server --jetstream -DV
```

### Backend Startup Issues
```bash
# Check all dependencies
go mod tidy
# Verify Redis and NATS are running
```
`````
