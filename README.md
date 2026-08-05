# ⚙️ Centralized Dynamic Config Management Service

![CI Status](https://github.com/goesbams/config-management-service/actions/workflows/ci.yml/badge.svg)
![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=flat&logo=go)
![Architecture](https://img.shields.io/badge/Architecture-Clean%20%2F%20Modular-blue)
![License](https://img.shields.io/badge/License-MIT-green)

High-performance, concurrency-safe Centralized Configuration & Version Control Service built in Go. Allows microservices to create, fetch, update, rollback, and audit dynamic feature flags and environment configurations in real-time with full version history.

---

## 🏗️ System Architecture & Data Flow

```mermaid
graph TD
    Client["🌐 Microservice / Client App"] -->|1. REST API Requests| Router["🚦 HTTP Mux Router (:8090)"]
    Router -->|2. Route to Handler| Handler["🎮 Config Handlers (/handlers)"]
    
    subgraph CoreEngine ["Concurrency-Safe Version Control Engine"]
        Handler -->|3. Mutex-Lock State Ops| Store["📦 In-Memory Version Store (/models)"]
        Store -->|4. Maintain Immutable History| VersionMap[("Audit Version Map")]
    end

    subgraph SwaggerUI ["OpenAPI / Documentation"]
        OpenAPI["📄 OpenAPI Specification"] -->|5. Serve Specs| Swagger["📊 Swagger UI (:8080)"]
    end
```

---

## 📊 Sequence Diagram: Dynamic Configuration & Atomic Rollback Flow

```mermaid
sequenceDiagram
    autonumber
    actor Client as Microservice / Admin Client
    participant Router as HTTP Router (:8090)
    participant Handler as Config Handler
    participant Engine as Version Control Engine
    participant Store as State Store Map

    note over Client,Store: Scenario 1: Creating & Updating Dynamic Configuration (v1 -> v2)
    Client->>Router: POST /config (Name: DB_CONFIG, Type: DATABASE, v1)
    activate Router
    Router->>Handler: CreateConfig(w, r)
    activate Handler
    Handler->>Engine: Validate & Insert New Config
    activate Engine
    Engine->>Store: Store[DB_CONFIG] = Version 1
    Store-->>Engine: State Saved
    deactivate Engine
    Handler-->>Router: 200 OK (Config Created)
    deactivate Handler
    Router-->>Client: 200 OK (Status: Created, Version: 1)
    deactivate Router

    Client->>Router: POST /config/update (Name: DB_CONFIG, v2: max_limit=2000)
    activate Router
    Router->>Handler: UpdateConfig(w, r)
    activate Handler
    Handler->>Engine: Append Version 2
    activate Engine
    Engine->>Store: Append Version 2 to DB_CONFIG History
    Store-->>Engine: Version 2 Appended
    deactivate Engine
    Handler-->>Router: 200 OK (Config Updated to v2)
    deactivate Handler
    Router-->>Client: 200 OK (Status: Updated, Active Version: 2)
    deactivate Router

    note over Client,Store: Scenario 2: Atomic Rollback to Past Version (v2 -> v1)
    Client->>Router: POST /config/rollback (Name: DB_CONFIG, Target Version: 1)
    activate Router
    Router->>Handler: RollbackConfig(w, r)
    activate Handler
    Handler->>Engine: Retrieve Historical Version 1 & Clone to New Head (v3)
    activate Engine
    Engine->>Store: Create Version 3 with Content of Version 1
    Store-->>Engine: Atomic Rollback Saved
    deactivate Engine
    Handler-->>Router: 200 OK (Rollback Executed)
    deactivate Handler
    Router-->>Client: 200 OK (Active Version: 3, Content: Rolled back to v1)
    deactivate Router
```

---

## 💡 Key Features & Business Capabilities

- **🔄 Versioning per Configuration**: Every change appends a new immutable version to the configuration audit history.
- **⚡ Instant Rollback Engine**: Atomic rollback to any historical version without restarting microservices.
- **🛡️ Type Validation**: Enforces valid config types (`DATABASE`, `API`, `FEATURE_FLAG`, etc.) and payload schemas.
- **📊 Real-time Fetch & Audit**: Query the active version or audit the complete version lifecycle (`GET /config/versions`).
- **🚀 High Concurrency**: Thread-safe memory operations designed for low-latency P99 reads.

---

## 🛠️ REST API Specification & Endpoint Matrix

| Method | Route | Description | Sample Payload | Response |
| :---: | :--- | :--- | :--- | :---: |
| `POST` | `/config` | Create a new configuration entry | `{"name":"DB_CONFIG","type":"DATABASE","versions":[{"version":1,"property":{"max_limit":1000}}]}` | `200 OK` |
| `POST` | `/config/update` | Append a new version to existing config | `{"name":"DB_CONFIG","type":"DATABASE","versions":[{"version":2,"property":{"max_limit":2000}}]}` | `200 OK` |
| `POST` | `/config/rollback` | Rollback configuration to a target version | `{"name":"DB_CONFIG","version":1}` | `200 OK` |
| `GET` | `/config/fetch` | Fetch latest or specific version | `GET /config/fetch?name=DB_CONFIG&version=1` | `200 OK` |
| `GET` | `/config/versions` | List all historical versions of a config | `GET /config/versions?name=DB_CONFIG` | `200 OK` |

---

## 🚀 Quick Start Guide

### Prerequisites
- Go `1.22+`
- Docker & Docker Compose
- `make` utility

---

### Running Locally with Go

```bash
# 1. Clone repository
git clone https://github.com/goesbams/config-management-service.git
cd config-management-service

# 2. Run unit tests (All tests pass)
make test

# 3. Start the service
make run
```

The REST API service will listen on `http://localhost:8090`.

---

### Running via Docker & OpenAPI Swagger

```bash
# 1. Build and start service container
make docker-build
make docker-up

# 2. View Swagger API documentation
make docker-openapi
```

Open your browser at `http://localhost:8080` to view the interactive **OpenAPI / Swagger UI**.

---

## 🧪 Local Test Verification

The unit test suite has been verified locally and passes 100%:

```bash
go test -v ./...
```

```text
=== RUN   TestCreateConfig
--- PASS: TestCreateConfig (0.00s)
=== RUN   TestUpdateConfig
--- PASS: TestUpdateConfig (0.00s)
=== RUN   TestRollbackConfig
--- PASS: TestRollbackConfig (0.00s)
=== RUN   TestFetchConfig
--- PASS: TestFetchConfig (0.00s)
=== RUN   TestListVersionsHandler
--- PASS: TestListVersionsHandler (0.00s)
PASS
ok  	github.com/goesbams/config-management-service/tests	0.415s
```

---
<sub>Designed & engineered with precision by Ando | Bambang Handoko.</sub>
