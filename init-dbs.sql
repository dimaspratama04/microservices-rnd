-- Create Databases
CREATE DATABASE order_db;
CREATE DATABASE payment_db;
CREATE DATABASE product_db;
CREATE DATABASE notification_db;

-- Product Service Schema
\c product_db
CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255),
    price DECIMAL(12, 2)
);

-- Order Service Schema
\c order_db
CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    product_id INTEGER,
    quantity INTEGER,
    total DECIMAL(12, 2),
    status VARCHAR(255)
);
CREATE INDEX IF NOT EXISTS idx_orders_deleted_at ON orders(deleted_at);

-- Payment Service Schema
\c payment_db
CREATE TABLE IF NOT EXISTS payments (
    id BIGSERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL,
    amount DECIMAL(12, 2) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Notification Service Schema
\c notification_db
CREATE TABLE IF NOT EXISTS notifications (
    id BIGSERIAL PRIMARY KEY,
    body TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
