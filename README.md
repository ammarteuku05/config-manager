# Configuration Management
## Setup and Running

### 1. Using Docker Compose

```bash
docker compose up --build
```

### 2. Manual Run

**Terminal 1 (Controller):**
```bash
make run-controller
```

**Terminal 2 (Worker):**
```bash
make run-worker
```

**Terminal 3 (Agent):**
```bash
make run-agent
```

### Running Tests
```bash
make test
make test-coverage
```

## API Endpoints

### Controller (Port 8080)

**1. Update Global Configuration**
```bash
curl -X POST http://localhost:8080/v1/config \
  -H "Content-Type: application/json" \
  -d '{"config":{"url":"https://ifconfig.me"}}'
```

**2. Register Agent (Internal)**
```bash
curl -X POST http://localhost:8080/v1/register \
  -H "Content-Type: application/json" \
  -H "Authorization: agent-secret" \
  -d '{"name":"agent-test"}'
```

**3. Get Latest Configuration (Internal)**
```bash
curl -X GET http://localhost:8080/v1/config
```

### Worker (Port 8082)

**1. Execute Configured Action (Proxy Hit)**
```bash
curl -X GET http://localhost:8082/hit
```

**2. Receive Configuration Push (Internal)**
```bash
curl -X POST http://localhost:8082/v1/config \
  -H "Content-Type: application/json" \
  -d '{"config":{"url":"https://ifconfig.me"}}'
```

