import { pgTable, serial, integer, decimal, varchar, timestamp } from "drizzle-orm/pg-core";

export const payments = pgTable("payments", {
    id: serial("id").primaryKey(),
    order_id: integer("order_id").notNull(),
    amount: decimal("amount", { precision: 12, scale: 2 }).notNull(),
    status: varchar("status", { length: 50 }).notNull(),
    created_at: timestamp("created_at").defaultNow(),
});
