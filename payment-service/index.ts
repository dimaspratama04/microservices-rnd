import { Elysia } from "elysia";
import * as amqplib from "amqplib";
import { pgTable, serial, integer, doublePrecision, varchar, timestamp } from "drizzle-orm/pg-core";
import { drizzle } from "drizzle-orm/node-postgres";
import pg from "pg";

const DB_URL = process.env.DB_URL || "postgres://user:password@localhost:5432/payment_db";
const RABBITMQ_URL = process.env.RABBITMQ_URL || "amqp://localhost";

// DB Schema
const payments = pgTable("payments", {
    id: serial("id").primaryKey(),
    order_id: integer("order_id").notNull(),
    amount: doublePrecision("amount").notNull(),
    status: varchar("status", { length: 50 }).notNull(),
    created_at: timestamp("created_at").defaultNow(),
});

// DB Connection
const pool = new pg.Pool({ connectionString: DB_URL });
const db = drizzle(pool);

async function connectQueue() {
    try {
        const connection = await amqplib.connect(RABBITMQ_URL);
        const channel = await connection.createChannel();
        await channel.assertQueue("notifications");
        return channel;
    } catch (error) {
        console.error("Failed to connect to RabbitMQ", error);
        return null;
    }
}

let channel: amqplib.Channel | null = null;
connectQueue().then(ch => {
    channel = ch;
    if (ch) console.log("Connected to RabbitMQ");
});

const app = new Elysia()
    .post("/payments", async ({ body, set }) => {
        const { order_id, amount } = body as { order_id: number, amount: number };

        // Validate Order Existence
        const ORDER_SERVICE_URL = process.env.ORDER_SERVICE_URL || "http://localhost:8080";
        try {
            const orderResp = await fetch(`${ORDER_SERVICE_URL}/orders/${order_id}`);
            if (orderResp.status !== 200) {
                set.status = 400;
                return { error: "Order not found or order service unavailable" };
            }
        } catch (error) {
            set.status = 400;
            return { error: "Failed to connect to order service" };
        }

        console.log("Processing payment...", body);
        
        // Persist to DB
        try {
            await db.insert(payments).values({
                order_id,
                amount,
                status: "COMPLETED"
            });
            console.log("Payment persisted to database");
        } catch (dbError) {
            console.error("Database error:", dbError);
            // We'll continue even if DB fail for this demo, or we could return 500
        }

        // Send message to queue
        if (channel) {
            const message = JSON.stringify({ event: "PAYMENT_SUCCESS", data: body });
            channel.sendToQueue("notifications", Buffer.from(message));
            console.log("Sent notification event to queue");
        }
        
        return { message: "Payment processed successfully", data: body };
    })
    .listen(8082);

console.log(`Payment service running at ${app.server?.hostname}:${app.server?.port}`);