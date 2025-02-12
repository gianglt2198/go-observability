# go-observability

## Overview
The Coffee Shop API is a simple HTTP server built using the Fiber framework in Go. It provides endpoints for managing coffee orders, allowing users to create and retrieve orders.

The main purpose is for implementing observability.
Tech stack includes:
- Database: postgres, grafan *(time-series database)*, loki *(non-sql object storage)*
- Language Programing: go
   - Logger: zap
   - Metric: opentelemetry
   - Tracer: opentelemetry
- Collector: vector, otel-collector
- Visualizatoin: Grafana

## Features
- Create coffee orders
- Retrieve existing coffee orders

## Setup Instructions
1. Clone the repository:
   ```
   git clone <repository-url>
   cd coffee-shop-api
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Run the server:
   ```
   go run cmd/server/main.go
   ```
4. Generate docs
   ```
   make install-dependencies
   # then run command for generating
   make generate-swag
   ```
5. Create a new migration
   ```
   make migrate-add
   ```
6. Run migration 
   ```
   make migrate-up
   ```


## API Endpoints
- `POST /coffee`: Create a new coffee drink.
- `GET /coffees`: Retrieve all coffee menu.

## Usage
You can use tools like Postman or curl to interact with the API endpoints. Make sure to send the appropriate JSON payload when creating orders.

## License
This project is licensed under the MIT License.