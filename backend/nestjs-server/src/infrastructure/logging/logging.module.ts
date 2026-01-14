import { Module, Global } from "@nestjs/common";
import { WinstonModule } from "nest-winston";
import { ConfigService } from "@nestjs/config";
import * as winston from "winston";

/**
 * Logging module using Winston
 * Provides structured logging with JSON or pretty format based on environment
 */
@Global()
@Module({
  imports: [
    WinstonModule.forRootAsync({
      inject: [ConfigService],
      useFactory: (configService: ConfigService) => {
        const logLevel = configService.get<string>("observability.logging.level") ?? "info";
        const logFormat = configService.get<string>("observability.logging.format") ?? "json";
        const isProduction = configService.get<boolean>("app.isProduction");

        const formats: winston.Logform.Format[] = [winston.format.timestamp(), winston.format.errors({ stack: true })];

        if (logFormat === "json") {
          formats.push(winston.format.json());
        } else {
          formats.push(
            winston.format.colorize({ all: true }),
            winston.format.printf(({ timestamp, level, message, context, ...meta }) => {
              const metaStr = Object.keys(meta).length ? ` ${JSON.stringify(meta)}` : "";
              return `${timestamp as string} [${(context as string) || "Application"}] ${level}: ${message as string}${metaStr}`;
            })
          );
        }

        return {
          level: logLevel,
          format: winston.format.combine(...formats),
          transports: [
            new winston.transports.Console({
              stderrLevels: ["error"],
            }),
            // Add file transport for production
            ...(isProduction
              ? [
                  new winston.transports.File({
                    filename: "logs/error.log",
                    level: "error",
                    maxsize: 10 * 1024 * 1024, // 10MB
                    maxFiles: 5,
                  }),
                  new winston.transports.File({
                    filename: "logs/combined.log",
                    maxsize: 10 * 1024 * 1024, // 10MB
                    maxFiles: 5,
                  }),
                ]
              : []),
          ],
        };
      },
    }),
  ],
  exports: [WinstonModule],
})
export class LoggingModule {}
