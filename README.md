# ☕ Coffee Shop API

> A production-ready RESTful API for managing a coffee shop, built with **Go** and **Fiber** — showcasing clean architecture, full observability stack, and containerized deployment.

[![Go Version](https://img.shields.io/badge/Go-1.23-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Framework](https://img.shields.io/badge/Framework-Fiber_v2-00ACD7?style=flat)](https://gofiber.io)
[![OTel](https://img.shields.io/badge/OpenTelemetry-Enabled-blueviolet?style=flat)](https://opentelemetry.io)

---

## 📖 Overview

The **Coffee Shop API** is a fully-featured HTTP server built with the [Fiber](https://gofiber.io) framework in Go. It provides endpoints for managing a coffee menu with full **CRUD** operations, and is designed as a reference implementation for **production-grade observability** — covering distributed tracing, metrics, and centralized log aggregation.

---

## ✨ Features

### 🔧 Core API

- **Full CRUD** operations on the coffee menu (Create, Read, Update, Delete)
- **Swagger UI** auto-generated documentation at `/swagger`
- **Health check** endpoint at `/health`
- **Structured JSON** success/error responses with consistent format

### 🏗️ Clean Architecture

- Separated layers: `handlers` → `routes` → `models` / `dto`
- **Generic route handler** (`Usecase[T, R]`) with middleware injection support
- **Unified payload validator** — merges body, query, and path params before validation using `go-playground/validator`
- **Custom error types** (`AError`) with HTTP status mapping

### 🗄️ Database

- **PostgreSQL** with GORM ORM
- **Connection pooling** with configurable max connections and timeout
- **Database migration** tool powered by [Goose](https://github.com/pressly/goose)
- Schema includes: `coffees`, `orders`, `order_coffees`, `payments` tables

### 📊 Observability (OpenTelemetry)

- **Distributed Tracing** — traces every HTTP request and DB operation via OTLP gRPC to Tempo
- **Metrics** — tracks:
  - `http_service_requests_total` (Counter)
  - `http_service_active_requests_total` (UpDown Counter)
  - `http_service_request_duration_milliseconds` (Histogram)
  - `system_memory_heap` (Observable Gauge)
- **Structured Logging** — [Zap](https://github.com/uber-go/zap) logger with:
  - Request ID correlation on every log entry
  - Separate info/error log files with rotation (production mode)
  - JSON-format console logging (development mode)
  - Secret masking utility (`InfoWithMask`)

### 📦 Log Aggregation Pipeline

- **Vector** — ships container logs (Docker mode) or file logs (file mode) to Loki
- **Loki** — stores and indexes logs (31-day retention)
- **Grafana** — unified dashboard for logs, traces, and metrics

### 🔍 Metrics Collection Stack

- **OpenTelemetry Collector** — receives OTLP (gRPC/HTTP), exports to Prometheus and Tempo
- **Prometheus** — scrapes app, collector, Tempo, cAdvisor, Node Exporter, and Vector
- **Grafana Tempo** — distributed trace storage with Memcached caching
- **cAdvisor** — Docker container resource monitoring
- **Node Exporter** — host-level system metrics

### 🐳 Docker & Dev Experience

- Full **Docker Compose** setup with all services in one command
- **Air** hot-reload in development (`Dockerfile.dev`)
- **Makefile** commands for common tasks

### 🧪 Load Testing

- **k6** test scripts with scenario-based load configuration:
  - `coffees` scenario: 5 VUs × 500 iterations
  - `detail_coffee` scenario: 2 VUs × 100 iterations

---

## 🛠️ Tech Stack

| Category         | Technology                 |
| ---------------- | -------------------------- |
| Language         | Go 1.23                    |
| Web Framework    | Fiber v2                   |
| ORM              | GORM                       |
| Database         | PostgreSQL 16              |
| Migration        | Goose                      |
| Logger           | Uber Zap + Lumberjack      |
| Tracing          | OpenTelemetry → Tempo      |
| Metrics          | OpenTelemetry → Prometheus |
| Log Shipping     | Vector                     |
| Log Storage      | Loki                       |
| Visualization    | Grafana                    |
| API Docs         | Swagger (swaggo)           |
| Load Testing     | k6                         |
| Config           | Viper                      |
| Containerization | Docker + Docker Compose    |

---

## 📁 Project Structure

```
coffee-shop-api/
├── cmd/
│   ├── migrate/        # DB migration CLI tool
│   └── server/         # Application entrypoint
├── config/             # Viper config loader
├── database/           # GORM database + custom Zap logger for GORM
├── deployment/         # Infrastructure configs
│   ├── loki/           # Loki log storage config
│   ├── otlp/           # OpenTelemetry Collector config
│   ├── prometheus/     # Prometheus scrape config
│   ├── tempo/          # Tempo trace storage config
│   └── vector/         # Vector log shipper config
├── docs/               # Auto-generated Swagger docs
├── internal/
│   ├── app/            # Fiber app bootstrap & middleware setup
│   ├── common/         # Shared enums, error types
│   ├── dto/            # Request/Response data transfer objects
│   ├── handlers/       # HTTP request handlers
│   └── models/         # GORM database models
├── middlewares/        # Logger, Tracing, Metrics, Validator middlewares
├── migrations/         # SQL migration files
├── monitoring/         # Zap logger + OTel SDK bootstrap
├── routes/             # Generic handler, response helpers, validator
├── testing/            # k6 load test scripts
├── utils/              # Shared utilities
├── config.yml          # App configuration
├── docker-compose.yml  # Full stack Docker Compose
├── Dockerfile.dev      # Dev Docker image with Air hot-reload
└── Makefile            # Developer shortcuts
```

---

## 🚀 Getting Started

### Prerequisites

- [Go 1.23+](https://golang.org/dl/)
- [Docker](https://www.docker.com/) & Docker Compose
- [Make](https://www.gnu.org/software/make/)

### 1. Clone the repository

```bash
git clone <repository-url>
cd coffee-shop-api
```

### 2. Start the full stack (Recommended)

```bash
docker compose up -d
```

This will start all services: the app, PostgreSQL, Prometheus, Tempo, Loki, Grafana, Vector, cAdvisor, Node Exporter, and the OTel Collector.

### 3. Run locally (without Docker)

```bash
# Install dependencies
go mod tidy

# Run database migration
make migrate-up

# Start the server
go run ./cmd/server
```

---

## ⚙️ Configuration

The app is configured via `config.yml`:

```yaml
app:
  name: coffee-shop-api
  env: development # or "production"
  port: 8081

database:
  host: db
  port: 5432
  name: coffee_shop
  user: postgres
  password: postgres
  max_connections: 100
  timeout: 5s

tracing:
  endpoint: collector:4317 # OTel Collector gRPC endpoint
```

> 💡 Set `APP_CONFIG_PATH` environment variable to override the config file path.

---

## 📋 Database Migrations

```bash
# Create a new migration
make migrate-add

# Apply all pending migrations
make migrate-up

# Roll back the latest migration
make migrate-down
```

---

## 📡 API Endpoints

Base URL: `http://localhost:8081`

| Method   | Endpoint          | Description         |
| -------- | ----------------- | ------------------- |
| `GET`    | `/health`         | Health check        |
| `GET`    | `/swagger/*`      | Swagger UI          |
| `GET`    | `/metrics`        | Prometheus metrics  |
| `GET`    | `/api/coffee`     | List all coffees    |
| `POST`   | `/api/coffee`     | Create a new coffee |
| `GET`    | `/api/coffee/:id` | Get a coffee by ID  |
| `PATCH`  | `/api/coffee/:id` | Update a coffee     |
| `DELETE` | `/api/coffee/:id` | Delete a coffee     |

### Example Requests

**Create a Coffee**

```bash
curl -X POST http://localhost:8081/api/coffee \
  -H "Content-Type: application/json" \
  -d '{"name": "Espresso", "price": 3.50, "description": "Strong and bold"}'
```

**List All Coffees**

```bash
curl http://localhost:8081/api/coffee
```

**Update a Coffee**

```bash
curl -X PATCH http://localhost:8081/api/coffee/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Double Espresso", "price": 4.50}'
```

### Response Format

**Success**

```json
{
  "success": true,
  "data": { ... }
}
```

**Error**

```json
{
  "success": false,
  "error": "error message"
}
```

---

## 📚 API Documentation

Swagger UI is available at:

```
http://localhost:8081/swagger/index.html
```

To regenerate Swagger docs from code annotations:

```bash
# Install swaggo CLI
make install-dependencies

# Generate docs
make generate-swag
```

---

## 🔭 Observability Stack

| Service            | URL                    | Description                |
| ------------------ | ---------------------- | -------------------------- |
| **Grafana**        | http://localhost:3000  | Dashboards (admin / `123`) |
| **Prometheus**     | http://localhost:9090  | Metrics browser            |
| **Tempo**          | http://localhost:3200  | Trace storage              |
| **Loki**           | http://localhost:3100  | Log storage                |
| **OTel Collector** | http://localhost:13133 | Health check               |
| **cAdvisor**       | internal               | Container metrics          |
| **Node Exporter**  | internal               | Host metrics               |

### Architecture Diagram

```
App (Fiber)
  │
  ├─── Traces (OTLP gRPC) ──► OTel Collector ──► Tempo ──► Grafana
  │
  ├─── Metrics (OTLP gRPC) ─► OTel Collector ──► Prometheus ──► Grafana
  │
  └─── Logs (stdout/file) ──► Vector ──────────► Loki ──────► Grafana
```

---

## 🧪 Load Testing

The project includes [k6](https://k6.io) load test scripts:

```bash
k6 run testing/k6-testing.js
```

Scenarios defined:

- **`coffees`**: 5 VUs, 500 iterations — tests `GET /api/coffee`
- **`detail_coffee`**: 2 VUs, 100 iterations — tests `GET /api/coffee/:id`

---

## 🛠️ Makefile Commands

```bash
make migrate-add       # Create a new SQL migration file
make migrate-up        # Apply all pending migrations
make migrate-down      # Roll back the latest migration
make install-dependencies  # Install swaggo CLI
make generate-swag     # Regenerate Swagger docs
```

---

## 📄 License

This project is licensed under the **MIT License** — see the [LICENSE](LICENSE) file for details.

---

> Built with ❤️ using Go · Fiber · OpenTelemetry · Grafana Stack
