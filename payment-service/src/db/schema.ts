import { pgTable, serial, integer, decimal, varchar, timestamp } from "drizzle-orm/pg-core";

export const payments = pgTable("payments", {
    id: serial("id").primaryKey(),
    invoice_id: varchar("invoice_id", { length: 50 }).notNull(),
    amount: decimal("amount", { precision: 12, scale: 2 }).notNull(),
    status: varchar("status", { length: 50 }).notNull(),
    created_at: timestamp("created_at").defaultNow(),
});
