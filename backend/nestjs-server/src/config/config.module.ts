import { Module, Global } from "@nestjs/common";
import { ConfigModule as NestConfigModule } from "@nestjs/config";
import configuration from "./configuration";

/**
 * Global configuration module
 * Loads and validates environment variables at application startup
 */
@Global()
@Module({
  imports: [
    NestConfigModule.forRoot({
      isGlobal: true,
      load: [configuration],
      // Environment file loading based on NODE_ENV
      envFilePath: [`.env.${process.env.NODE_ENV ?? "development"}`, ".env"],
    }),
  ],
  exports: [NestConfigModule],
})
export class ConfigModule {}
