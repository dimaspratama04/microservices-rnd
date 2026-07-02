import { Elysia } from "elysia";
import { eq } from "drizzle-orm";
import { db } from "../db/index.js";
import { payments } from "../db/schema.js";
import { channel } from "../queue/rabbitmq.js";
import { successResponse, errorResponse } from "../dto/web_response.js";
import { PaymentRequestDTO, PaymentStatusResponseDTO } from "../dto/payment_dto.js";
import { context, propagation } from "@opentelemetry/api";

const ORDER_SERVICE_URL = process.env.ORDER_SERVICE_URL || "http://localhost:8080";
const PRODUCT_SERVICE_URL = process.env.PRODUCT_SERVICE_URL || "http://localhost:8081";

export const paymentRoutes = new Elysia()
  .post("/payments", async ({ body, set }) => {
    const { invoice_id, amount } = body as PaymentRequestDTO;

    // Validate Order Existence
    let orderData: any;
    try {
      const headers = {};
      propagation.inject(context.active(), headers);
      const orderResp = await fetch(`${ORDER_SERVICE_URL}/orders/invoice/${invoice_id}`, { headers });
      if (orderResp.status !== 200) {
        set.status = 400;
        return errorResponse("Order not found or order service unavailable");
      }
      orderData = await orderResp.json();
    } catch (error) {
      set.status = 400;
      return errorResponse("Failed to connect to order service");
    }

    const expectedTotal = orderData.data.total;
    if (amount < expectedTotal) {
      set.status = 400;
      return errorResponse("ammount not insufficant");
    }

    console.log("Processing payment...", body);

    // Persist to DB
    try {
      await db.insert(payments).values({
        invoice_id,
        amount: amount.toString(),
        status: "COMPLETED",
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

    // Update Order Status to SHIPPED
    try {
      const headers: Record<string, string> = { "Content-Type": "application/json" };
      propagation.inject(context.active(), headers);
      await fetch(`${ORDER_SERVICE_URL}/orders/invoice/${invoice_id}/status`, {
        method: "PATCH",
        headers,
        body: JSON.stringify({ status: "SHIPPED" }),
      });
      console.log(`Updated order invoice ${invoice_id} status to SHIPPED`);
    } catch (error) {
      console.error(`Failed to update order status:`, error);
    }

    return successResponse("Payment processed successfully", body);
  })
  .get("/payments/invoice/:invoiceId/status", async ({ params: { invoiceId }, set }) => {
    const paymentRecords = await db.select().from(payments).where(eq(payments.invoice_id, invoiceId));
    if (paymentRecords.length === 0) {
      set.status = 404;
      return errorResponse("Payment not found");
    }
    const payment = paymentRecords[0];

    let orderData: any;
    try {
      const headers = {};
      propagation.inject(context.active(), headers);
      const orderResp = await fetch(`${ORDER_SERVICE_URL}/orders/invoice/${payment.invoice_id}`, { headers });
      if (orderResp.status === 200) {
        orderData = await orderResp.json();
      }
    } catch (e) {
      console.error(e);
    }

    const statusData: PaymentStatusResponseDTO = {
      invoice_id: invoiceId,
      total: orderData.data.total,
      payment_status: payment.status,
    };
    return successResponse("Payment status retrieved successfully", statusData);
  });
