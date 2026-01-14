import { z } from "zod";

/**
 * Transform string '0' or '1' to boolean
 */
const booleanString = z
  .enum(["0", "1"])
  .transform((val) => val === "1")
  .default("0");

/**
 * Environment validation schema using Zod
 * This validates all environment variables at application startup
 */
export const envSchema = z
  .object({
    // ===================
    // Application
    // ===================
    NODE_ENV: z.enum(["development", "production", "test"]).default("development"),
    PORT: z.coerce.number().default(3000),
    APP_NAME: z.string().default("reminder-api"),
    APP_VERSION: z.string().default("1.0.0"),

    // ===================
    // Database
    // ===================
    DB_TYPE: z.enum(["mongodb", "firebase"]).default("mongodb"),

    // MongoDB
    MONGO_DB_URI: z.string().optional(),
    MONGO_DB_NAME: z.string().default("reminder"),

    // Firebase
    FIREBASE_PROJECT_ID: z.string().optional(),
    USE_FIREBASE_EMULATOR: booleanString,
    FIREBASE_EMULATOR_HOST: z.string().optional(),
    FIREBASE_GOOGLE_APP_CREDENTIALS: z.string().optional(),

    // ===================
    // Observability - Logging
    // ===================
    LOG_LEVEL: z.enum(["error", "warn", "info", "debug", "verbose"]).default("info"),
    LOG_FORMAT: z.enum(["json", "pretty"]).default("json"),

    // ===================
    // Observability - Metrics
    // ===================
    ENABLE_METRICS: booleanString,
    METRICS_PORT: z.coerce.number().default(9090),

    // ===================
    // Observability - Tracing (OpenTelemetry)
    // ===================
    ENABLE_TRACING: booleanString,
    OTEL_EXPORTER_OTLP_ENDPOINT: z.string().optional(),
    OTEL_SERVICE_NAME: z.string().default("reminder-api"),

    // ===================
    // Security
    // ===================
    RATE_LIMIT_TTL: z.coerce.number().default(60), // seconds
    RATE_LIMIT_MAX: z.coerce.number().default(100), // requests per TTL
    ENABLE_HELMET: booleanString.default("1"),

    // ===================
    // Caching
    // ===================
    ENABLE_CACHE: booleanString,
    REDIS_HOST: z.string().default("localhost"),
    REDIS_PORT: z.coerce.number().default(6379),
    REDIS_PASSWORD: z.string().optional(),
    CACHE_TTL: z.coerce.number().default(300), // seconds

    // ===================
    // CORS
    // ===================
    CORS_ORIGIN: z.string().default("http://localhost:3000"),

    // ===================
    // Database Seeding
    // ===================
    ENABLE_DB_SEEDING: booleanString,
  })
  .refine(
    (data) => {
      // MongoDB requires URI when selected
      if (data.DB_TYPE === "mongodb" && !data.MONGO_DB_URI) {
        return false;
      }
      return true;
    },
    {
      message: "MONGO_DB_URI is required when DB_TYPE is mongodb",
      path: ["MONGO_DB_URI"],
    }
  )
  .refine(
    (data) => {
      // Firebase requires project ID when selected
      if (data.DB_TYPE === "firebase" && !data.FIREBASE_PROJECT_ID) {
        return false;
      }
      return true;
    },
    {
      message: "FIREBASE_PROJECT_ID is required when DB_TYPE is firebase",
      path: ["FIREBASE_PROJECT_ID"],
    }
  );

/**
 * Inferred type from the Zod schema
 */
export type EnvConfig = z.infer<typeof envSchema>;
