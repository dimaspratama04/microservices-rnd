import { Elysia } from "elysia";
import { db } from "../db/index.js";
import { payments } from "../db/schema.js";
import { channel } from "../queue/rabbitmq.js";

const ORDER_SERVICE_URL = process.env.ORDER_SERVICE_URL || "http://localhost:8080";

export const paymentRoutes = new Elysia().post("/payments", async ({ body, set }) => {
    const { order_id, amount } = body as { order_id: number; amount: number };

    // Validate Order Existence
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
    }

    // Send message to queue
    if (channel) {
        const message = JSON.stringify({ event: "PAYMENT_SUCCESS", data: body });
        channel.sendToQueue("notifications", Buffer.from(message));
        console.log("Sent notification event to queue");
    }

    return { message: "Payment processed successfully", data: body };
});
