import { Elysia } from "elysia";
import * as amqplib from "amqplib";

const RABBITMQ_URL = process.env.RABBITMQ_URL || "amqp://localhost";

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
    .post("/payments", async ({ body }) => {
        console.log("Processing payment...", body);
        
        // Send message to queue
        if (channel) {
            const message = JSON.stringify({ event: "PAYMENT_SUCCESS", data: body });
            channel.sendToQueue("notifications", Buffer.from(message));
            console.log("Sent notification event to queue");
        }
        
        return { status: "Payment processed" };
    })
    .listen(8082);

console.log(`Payment service running at ${app.server?.hostname}:${app.server?.port}`);