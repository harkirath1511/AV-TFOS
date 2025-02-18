### **Deployment Steps**

1. **Install Dependencies**:
```bash
# NATS Server
go install github.com/nats-io/nats-server/v2@latest

# Redis with RedisGraph
docker run -p 6379:6379 redislabs/redismod
```

2. **Build & Run**:
```bash
go mod init flowsync
go mod tidy
go build -o flowsync
./flowsync
```

3. **Start Simulation**:
```bash
curl -X POST http://localhost:8080/simulation/start -d '{
  "vehicle_count": 1000,
  "av_percentage": 0.3,
  "emergency_interval": "5m"
}'
```

4. **Trigger Emergency**:
```bash
curl -X POST http://localhost:8080/emergency -d '{
  "ambulance_id": "medic-1",
  "path": [[37.7749,-122.4194], [37.7813,-122.4168]],
  "priority": 1
}'
```