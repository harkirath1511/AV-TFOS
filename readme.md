# AV-TFOS

## Prerequisites Installation

```bash
# Install Redis Server
winget install Redis

# Install NATS Server (Windows)
# 1. Download the latest NATS server release from:
# https://github.com/nats-io/nats-server/releases
# Download the Windows zip file (nats-server-vX.X.X-windows-amd64.zip)

# 2. Extract the zip file
# 3. Add the extracted folder location to your PATH environment variable
# 4. Verify installation
nats-server --version
```

## Start Required Services

```bash
# Start Redis Server (in a new terminal)
redis-server

# Start NATS Server (in another terminal)
nats-server
```

## Initialize Redis with Sample Data

```bash
# Open Redis CLI
redis-cli

# Add sample intersection data
> GEOADD intersections:locations -118.2437 34.0522 intersection_001
> SADD intersections:active intersection_001
> HSET intersection:intersection_001 current_phase north_south emergency_active false
```

## Setup Go Backend

```bash
# Navigate to project directory
cd Tfos

# Initialize Go modules
go mod init flowsync-backend
go mod tidy

# Install Go dependencies
go get github.com/nats-io/nats.go
go get github.com/redis/go-redis/v9
go get github.com/gin-gonic/gin
go get github.com/gorilla/websocket
```

## Run the Backend

```bash
# From the Tfos directory
go run main.go
```

## Verify Services

```bash
# Test simulation start
curl -X POST http://localhost:8080/simulation/start

# Test metrics endpoint
curl http://localhost:8080/metrics
```

The backend server will be available at:
- WebSocket: `ws://localhost:8080/ws`
- REST API: `http://localhost:8080`