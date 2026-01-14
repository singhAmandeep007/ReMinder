import { Module, MiddlewareConsumer, NestModule } from "@nestjs/common";
import { ThrottlerModule, ThrottlerGuard } from "@nestjs/throttler";
import { APP_GUARD, APP_FILTER, APP_INTERCEPTOR } from "@nestjs/core";
import { ConfigModule } from "./config";
import { HealthModule } from "./infrastructure/health/health.module";
import { LoggingModule } from "./infrastructure/logging/logging.module";
import { MongodbModule } from "./database/mongodb/mongodb.module";
import { RemindersModule } from "./modules/reminders/reminders.module";
import { ReminderGroupsModule } from "./modules/reminder-groups/reminder-groups.module";
import { HttpExceptionFilter } from "./common/filters/http-exception.filter";
import { TransformInterceptor } from "./common/interceptors/transform.interceptor";
import { LoggingInterceptor } from "./common/interceptors/logging.interceptor";
import { CorrelationIdMiddleware } from "./common/middleware/correlation-id.middleware";
import { ConfigService } from "@nestjs/config";

@Module({
  imports: [
    // Configuration (global)
    ConfigModule,

    // Rate limiting
    ThrottlerModule.forRootAsync({
      inject: [ConfigService],
      useFactory: (config: ConfigService) => [
        {
          ttl: config.get<number>("security.rateLimit.ttl") ?? 60000,
          limit: config.get<number>("security.rateLimit.limit") ?? 100,
        },
      ],
    }),

    // Infrastructure
    LoggingModule,
    HealthModule,

    // Database (MongoDB)
    MongodbModule,

    // Feature modules
    RemindersModule,
    ReminderGroupsModule,
  ],
  providers: [
    // Global rate limit guard
    {
      provide: APP_GUARD,
      useClass: ThrottlerGuard,
    },
    // Global exception filter
    {
      provide: APP_FILTER,
      useClass: HttpExceptionFilter,
    },
    // Global response transformer
    {
      provide: APP_INTERCEPTOR,
      useClass: TransformInterceptor,
    },
    // Global logging interceptor
    {
      provide: APP_INTERCEPTOR,
      useClass: LoggingInterceptor,
    },
  ],
})
export class AppModule implements NestModule {
  configure(consumer: MiddlewareConsumer) {
    consumer.apply(CorrelationIdMiddleware).forRoutes("*");
  }
}
