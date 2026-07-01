import { pgTable, serial, integer, doublePrecision, varchar, timestamp } from "drizzle-orm/pg-core";

export const payments = pgTable("payments", {
    id: serial("id").primaryKey(),
    order_id: integer("order_id").notNull(),
    amount: doublePrecision("amount").notNull(),
    status: varchar("status", { length: 50 }).notNull(),
    created_at: timestamp("created_at").defaultNow(),
});
