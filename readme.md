# AV-TFOS

## Start Required Services

```bash
nats-server -js
```

```bash
docker run --name redis -p 6379:6379 -d redis
```

## Initialize Redis with Sample Data

## Setup Go Backend

```bash
# Navigate to project directory
cd Tfos

go mod tidy

# Install Go dependencies
go get github.com/nats-io/nats.go@latest
go get github.com/redis/go-redis/v9@latest
go get github.com/gin-gonic/gin@latest
go get github.com/gorilla/websocket@latest
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