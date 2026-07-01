import pg from "pg";
import { drizzle } from "drizzle-orm/node-postgres";

const DB_URL = process.env.DB_URL || "postgres://user:password@localhost:5432/payment_db";

const pool = new pg.Pool({ connectionString: DB_URL });
export const db = drizzle(pool);
