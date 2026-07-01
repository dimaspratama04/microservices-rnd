import * as amqplib from "amqplib";

const RABBITMQ_URL = process.env.RABBITMQ_URL || "amqp://localhost";
export let channel: amqplib.Channel | null = null;

export async function connectQueue() {
    try {
        const connection = await amqplib.connect(RABBITMQ_URL);
        channel = await connection.createChannel();
        await channel.assertQueue("notifications");
        console.log("Connected to RabbitMQ");
        return channel;
    } catch (error) {
        console.error("Failed to connect to RabbitMQ", error);
        return null;
    }
}
