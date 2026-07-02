import { Elysia } from "elysia";
import { opentelemetry } from "@elysiajs/opentelemetry";
import { BatchSpanProcessor } from "@opentelemetry/sdk-trace-node";
import { OTLPTraceExporter } from "@opentelemetry/exporter-trace-otlp-proto";
import { trace, propagation } from "@opentelemetry/api";
import { W3CTraceContextPropagator } from "@opentelemetry/core";
import { connectQueue } from "./src/queue/rabbitmq.js";

propagation.setGlobalPropagator(new W3CTraceContextPropagator());
import { paymentRoutes } from "./src/routes/payment.js";

// Connect to RabbitMQ
connectQueue();

const app = new Elysia()
  .use(
    opentelemetry({
      serviceName: "payment-service",
      spanProcessors: [new BatchSpanProcessor(new OTLPTraceExporter())],
    }),
  )
  .onAfterHandle(({ set }) => {
    const span = trace.getActiveSpan();
    if (span) {
      set.headers["X-Request-Id"] = span.spanContext().traceId;
    }
  })
  .use(paymentRoutes)
  .listen(8082);

console.log(`Payment service running at ${app.server?.hostname}:${app.server?.port}`);
