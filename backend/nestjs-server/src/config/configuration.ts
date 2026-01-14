import { envSchema, EnvConfig } from "./env.schema";

/**
 * Validates environment variables against the Zod schema
 * Throws an error with detailed information if validation fails
 */
export const validateEnv = (): EnvConfig => {
  const result = envSchema.safeParse(process.env);

  if (!result.success) {
    console.error("❌ Invalid environment variables:");
    console.error(JSON.stringify(result.error.format(), null, 2));
    throw new Error("Invalid environment configuration. See above for details.");
  }

  return result.data;
};

/**
 * Configuration factory for @nestjs/config
 * Returns a typed configuration object organized by domain
 */
export default () => {
  const env = validateEnv();

  return {
    // Application settings
    app: {
      nodeEnv: env.NODE_ENV,
      port: env.PORT,
      name: env.APP_NAME,
      version: env.APP_VERSION,
      isProduction: env.NODE_ENV === "production",
      isDevelopment: env.NODE_ENV === "development",
      isTest: env.NODE_ENV === "test",
    },

    // Database configuration
    database: {
      type: env.DB_TYPE,
      mongodb: {
        uri: env.MONGO_DB_URI,
        dbName: env.MONGO_DB_NAME,
      },
      firebase: {
        projectId: env.FIREBASE_PROJECT_ID,
        useEmulator: env.USE_FIREBASE_EMULATOR,
        emulatorHost: env.FIREBASE_EMULATOR_HOST,
        credentialsPath: env.FIREBASE_GOOGLE_APP_CREDENTIALS,
      },
    },

    // Observability settings
    observability: {
      logging: {
        level: env.LOG_LEVEL,
        format: env.LOG_FORMAT,
      },
      metrics: {
        enabled: env.ENABLE_METRICS,
        port: env.METRICS_PORT,
      },
      tracing: {
        enabled: env.ENABLE_TRACING,
        endpoint: env.OTEL_EXPORTER_OTLP_ENDPOINT,
        serviceName: env.OTEL_SERVICE_NAME,
      },
    },

    // Security settings
    security: {
      rateLimit: {
        ttl: env.RATE_LIMIT_TTL * 1000, // Convert to milliseconds for NestJS
        limit: env.RATE_LIMIT_MAX,
      },
      helmet: {
        enabled: env.ENABLE_HELMET,
      },
    },

    // Cache settings
    cache: {
      enabled: env.ENABLE_CACHE,
      redis: {
        host: env.REDIS_HOST,
        port: env.REDIS_PORT,
        password: env.REDIS_PASSWORD,
        ttl: env.CACHE_TTL,
      },
    },

    // CORS settings
    cors: {
      origin: env.CORS_ORIGIN,
    },

    // Database seeding
    seeding: {
      enabled: env.ENABLE_DB_SEEDING,
    },
  };
};

/**
 * Type for the configuration object
 */
export type AppConfig = ReturnType<typeof import("./configuration").default>;
