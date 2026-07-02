import { Elysia } from "elysia";
import { eq } from "drizzle-orm";
import { db } from "../db/index.js";
import { payments } from "../db/schema.js";
import { channel } from "../queue/rabbitmq.js";
import { successResponse, errorResponse } from "../dto/web_response.js";
import { PaymentRequestDTO, PaymentStatusResponseDTO } from "../dto/payment_dto.js";

const ORDER_SERVICE_URL = process.env.ORDER_SERVICE_URL || "http://localhost:8080";
const PRODUCT_SERVICE_URL = process.env.PRODUCT_SERVICE_URL || "http://localhost:8081";

export const paymentRoutes = new Elysia()
  .post("/payments", async ({ body, set }) => {
    const { order_id, amount } = body as PaymentRequestDTO;

    // Validate Order Existence
    let orderData: any;
    try {
      const orderResp = await fetch(`${ORDER_SERVICE_URL}/orders/${order_id}`);
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
        order_id,
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
      await fetch(`${ORDER_SERVICE_URL}/orders/${order_id}/status`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ status: "SHIPPED" }),
      });
      console.log(`Updated order ${order_id} status to SHIPPED`);
    } catch (error) {
      console.error(`Failed to update order ${order_id} status:`, error);
    }

    return successResponse("Payment processed successfully", body);
  })
  .get("/payments/:id/status", async ({ params: { id }, set }) => {
    const paymentId = Number(id);
    const paymentRecords = await db.select().from(payments).where(eq(payments.id, paymentId));
    if (paymentRecords.length === 0) {
      set.status = 404;
      return errorResponse("Payment not found");
    }
    const payment = paymentRecords[0];

    let orderData: any;
    try {
      const orderResp = await fetch(`${ORDER_SERVICE_URL}/orders/${payment.order_id}`);
      if (orderResp.status === 200) {
        orderData = await orderResp.json();
      }
    } catch (e) {
      console.error(e);
    }

    let productName = "Unknown Product";
    const orderId = orderData.data.id;
    const productId = orderData.data.product_id;
    const paymentTotal = orderData.data.total;

    if (orderData && orderData.data && orderData.data.product_id) {
      try {
        const productResp = await fetch(`${PRODUCT_SERVICE_URL}/products/${productId}`);
        if (productResp.status === 200) {
          const productData = (await productResp.json()) as any;
          productName = productData.data.name;
        }
      } catch (e) {
        console.error(e);
      }
    }

    const statusData: PaymentStatusResponseDTO = {
      order_id: orderId,
      product_name: productName,
      total: paymentTotal,
    };
    return successResponse("Payment status retrieved successfully", statusData);
  });
