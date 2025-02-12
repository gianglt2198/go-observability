# Coffee Shop API

## Overview
The Coffee Shop API is a simple HTTP server built using the Fiber framework in Go. It provides endpoints for managing coffee orders, allowing users to create and retrieve orders.

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

## API Endpoints
- `POST /orders`: Create a new coffee order.
- `GET /orders`: Retrieve all coffee orders.

## Usage
You can use tools like Postman or curl to interact with the API endpoints. Make sure to send the appropriate JSON payload when creating orders.

## License
This project is licensed under the MIT License.