# AV-TFOS

# Prerequisites Installation

```bash
# Install Redis Server
winget install Redis

# Install NATS Server
winget install nats-io.nats-server

# Clone the repository (if not done already)
git clone <your-repo-url>
cd AV-TFOS
```

#  Start Required Services

```bash
# Start Redis Server (in a new terminal)
redis-server

# Start NATS Server (in another terminal)
nats-server
```

# Initialize Redis with Sample Data

```bash
# Open Redis CLI
redis-cli

# Add sample intersection data
> GEOADD intersections:locations -118.2437 34.0522 intersection_001
> SADD intersections:active intersection_001
> HSET intersection:intersection_001 current_phase north_south emergency_active false
```

# Setup Go Backend

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

# Run the Backend

```bash
# From the Tfos directory
go run main.go
```