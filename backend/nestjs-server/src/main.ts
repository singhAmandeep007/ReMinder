import { NestFactory } from "@nestjs/core";
import { ValidationPipe, VersioningType, Logger } from "@nestjs/common";
import { SwaggerModule, OpenAPIObject } from "@nestjs/swagger";
import { ConfigService } from "@nestjs/config";
import * as fs from "fs";
import * as yaml from "yaml";
import * as path from "path";
import helmet from "helmet";
import { AppModule } from "./app.module";

async function bootstrap() {
  const logger = new Logger("Bootstrap");

  const app = await NestFactory.create(AppModule, {
    logger:
      process.env.NODE_ENV === "production" ? ["error", "warn", "log"] : ["error", "warn", "log", "debug", "verbose"],
  });

  const configService = app.get(ConfigService);

  // ===================
  // Security - Helmet
  // ===================
  if (configService.get<boolean>("security.helmet.enabled")) {
    app.use(
      helmet({
        contentSecurityPolicy: configService.get<string>("app.nodeEnv") === "production",
      })
    );
    logger.log("Helmet security middleware enabled");
  }

  // ===================
  // API Versioning + Global Prefix
  // ===================
  app.setGlobalPrefix("api");
  app.enableVersioning({
    type: VersioningType.URI,
    defaultVersion: "1",
  });

  // ===================
  // Validation Pipe
  // ===================
  app.useGlobalPipes(
    new ValidationPipe({
      whitelist: true,
      transform: true,
      forbidNonWhitelisted: true,
      transformOptions: {
        enableImplicitConversion: true,
      },
    })
  );

  // ===================
  // CORS
  // ===================
  const corsOrigin = configService.get<string>("cors.origin");
  app.enableCors({
    origin: corsOrigin,
    credentials: true,
  });
  logger.log(`CORS enabled for origin: ${corsOrigin}`);

  // ===================
  // Swagger - Load existing OpenAPI spec (NO DUPLICATION)
  // ===================
  const nodeEnv = configService.get<string>("app.nodeEnv");
  if (nodeEnv !== "production") {
    try {
      // Use process.cwd() which is nestjs-server directory, then go up to project root
      const openApiPath = path.resolve(process.cwd(), "..", "..", "openapi", "api.yaml");
      logger.log(`Looking for OpenAPI spec at: ${openApiPath}`);

      if (fs.existsSync(openApiPath)) {
        const openApiSpec = yaml.parse(fs.readFileSync(openApiPath, "utf8")) as OpenAPIObject;
        SwaggerModule.setup("api/docs", app, openApiSpec);
        logger.log("Swagger UI available at /api/docs");
      } else {
        logger.warn(`OpenAPI spec not found at ${openApiPath}`);
      }
    } catch (error) {
      logger.error("Failed to load OpenAPI spec:", error);
    }
  }

  // ===================
  // Start Server
  // ===================
  const port = configService.get<number>("app.port") ?? 3000;
  const appName = configService.get<string>("app.name") ?? "reminder-api";

  await app.listen(port);

  logger.log(`🚀 ${appName} running on: http://localhost:${port}`);
  logger.log(`📚 Swagger docs: http://localhost:${port}/api/docs`);
  logger.log(`❤️  Health check: http://localhost:${port}/api/v1/health`);
  logger.log(`🌍 Environment: ${nodeEnv}`);
}

bootstrap().catch((error) => {
  console.error("❌ Error starting application:", error);
  process.exit(1);
});
