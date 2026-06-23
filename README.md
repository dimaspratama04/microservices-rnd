# Microservices R&D Project

This project is a microservices-based architecture research and development environment featuring multiple services built with different technologies (Java, Go, TypeScript). It demonstrates CRUD operations, inter-service communication, and message queuing.

## Architecture Overview

| Service                  | Language/Runtime | Framework       | Database                  | Others                |
| ------------------------ | ---------------- | --------------- | ------------------------- | --------------------- |
| **Product Service**      | Java 17          | Spring Boot 3.2 | PostgreSQL, H2 (test)     | Maven                 |
| **Order Service**        | Go 1.25          | Fiber v2        | PostgreSQL, SQLite (test) | GORM                  |
| **Payment Service**      | TypeScript (Bun) | Elysia          | PostgreSQL                | Drizzle ORM, RabbitMQ |
| **Notification Service** | Java 17          | Spring Boot 4.1 | -                         | Gradle, RabbitMQ      |

- **PostgreSQL:** Shared database instance with separate databases for each service.
- **RabbitMQ:** Message broker for asynchronous communication between Payment and Notification services.

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

## Testing CRUD Operations

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

```bash
curl -X POST http://localhost:8080/orders \
-H "Content-Type: application/json" \
-d '{
  "product_id": 1,
  "quantity": 2,
  "total": 4999.98
}'
```

**Get All Orders:**

```bash
curl http://localhost:8080/orders
```

**Update Order Status:**

```bash
curl -X PUT http://localhost:8080/orders/1 \
-H "Content-Type: application/json" \
-d '{
  "status": "Shipped"
}'
```

**Delete Order:**

```bash
curl -X DELETE http://localhost:8080/orders/1
```

---

### 3. Payment Service (Port 8082)

**Process Payment:**
_(This will also trigger a notification event to RabbitMQ)_

```bash
curl -X POST http://localhost:8082/payments \
-H "Content-Type: application/json" \
-d '{
  "order_id": 1,
  "amount": 4999.98,
  "currency": "USD"
}'
```

## Development & Testing

- **Go Tests:** `cd order-service && go test -v`
- **Java Tests (Product):** `cd product-service && ./mvnw test`
- **Java Tests (Notification):** `cd notification-service && ./gradlew test`

## Environment Variables

Check `docker-compose.yml` for database connections and service discovery URLs.
