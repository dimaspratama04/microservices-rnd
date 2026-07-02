# Microservices R&D Project

This project is a microservices-based architecture research and development environment featuring multiple services built with different technologies (Java, Go, TypeScript). It demonstrates CRUD operations, inter-service communication, message queuing, and distributed tracing.

## Architecture Overview

| Service                  | Language/Runtime | Framework       | Database                  | Others                                |
| ------------------------ | ---------------- | --------------- | ------------------------- | ------------------------------------- |
| **Product Service**      | Java 17          | Spring Boot 3.2 | PostgreSQL, H2 (test)     | Maven, OpenTelemetry                  |
| **Order Service**        | Go 1.25          | Fiber v2        | PostgreSQL, SQLite (test) | GORM, OpenTelemetry                   |
| **Payment Service**      | TypeScript (Bun) | Elysia          | PostgreSQL                | Drizzle ORM, RabbitMQ, OpenTelemetry  |
| **Notification Service** | Java 17          | Spring Boot 4.1 | -                         | Gradle, RabbitMQ, OpenTelemetry       |

- **PostgreSQL:** Shared database instance with separate databases for each service. All financial data (`price`, `total`, `amount`) is consistently stored as `DECIMAL(12, 2)`.
- **RabbitMQ:** Message broker for asynchronous communication between Payment and Notification services.
- **Jaeger (OpenTelemetry):** Distributed tracing implemented across all services. Every API response includes an `X-Request-Id` header mapping to the active Trace ID.

## Prerequisites

- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [Curl](https://curl.se/) or [Postman](https://www.postman.com/) for testing

## Getting Started

1. **Clone the repository:**

   ```bash
   git clone <repository-url>
   cd microservices-rnd
   ```

2. **Start the services:**

   ```bash
   docker compose up -d --build
   ```

3. **Verify services are running:**
   ```bash
   docker ps
   ```

4. **Access Jaeger UI (Tracing):**
   Navigate to [http://localhost:16686](http://localhost:16686) in your browser to view traces.

---

## API Documentation & Testing

All API endpoints return standard JSON responses formatted as `{"message": "...", "data": ...}` and include an `X-Request-Id` header for tracing.

### 1. Product Service (Port 8081)

**Create Product:**
```bash
curl -X POST http://localhost:8081/products \
-H "Content-Type: application/json" \
-d '{
  "name": "MacBook Pro M3",
  "price": 2499.99
}'
```

**Get All Products:**
```bash
curl http://localhost:8081/products
```

**Update Product:**
```bash
curl -X PUT http://localhost:8081/products/1 \
-H "Content-Type: application/json" \
-d '{
  "name": "MacBook Pro M3 Max",
  "price": 3499.99
}'
```

**Delete Product:**
```bash
curl -X DELETE http://localhost:8081/products/1
```

---

### 2. Order Service (Port 8080)

**Create Order:**
*(Automatically fetches the product price from Product Service and calculates the total)*
```bash
curl -X POST http://localhost:8080/orders \
-H "Content-Type: application/json" \
-d '{
  "product_id": 1,
  "quantity": 2
}'
```

**Get All Orders:**
```bash
curl http://localhost:8080/orders
```

**Update Order Status by Invoice ID:**
```bash
curl -X PATCH http://localhost:8080/orders/invoice/INV-12345/status \
-H "Content-Type: application/json" \
-d '{
  "status": "Processing"
}'
```

**Delete Order:**
```bash
curl -X DELETE http://localhost:8080/orders/1
```

---

### 3. Payment Service (Port 8082)

**Process Payment:**
*(Validates the amount against Order Service. On success, it triggers a RabbitMQ notification and automatically updates the order status to SHIPPED).*
```bash
curl -X POST http://localhost:8082/payments \
-H "Content-Type: application/json" \
-d '{
  "invoice_id": "INV-12345",
  "amount": 4999.98
}'
```

**Get Payment Status:**
*(Retrieves the status of the payment and order total).*
```bash
curl http://localhost:8082/payments/invoice/INV-12345/status
```

*Example Response:*
```json
{
  "message": "Payment status retrieved successfully",
  "data": {
    "invoice_id": "INV-12345",
    "total": 4999.98,
    "payment_status": "COMPLETED"
  }
}
```

---

## Development & Testing

- **Go Tests (Order):** `cd order-service && go test -v`
- **Java Tests (Product):** `cd product-service && ./mvnw test`
- **Java Tests (Notification):** `cd notification-service && ./gradlew test`

## Environment Variables

Check `docker-compose.yml` for database connections, service discovery URLs, and OpenTelemetry configurations.
