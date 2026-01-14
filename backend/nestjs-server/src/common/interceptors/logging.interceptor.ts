import { Injectable, NestInterceptor, ExecutionContext, CallHandler, Logger } from "@nestjs/common";
import { Observable } from "rxjs";
import { tap } from "rxjs/operators";
import { Request, Response } from "express";

/**
 * Logging interceptor
 * Logs all incoming requests and their response times
 */
@Injectable()
export class LoggingInterceptor implements NestInterceptor {
  private readonly logger = new Logger("HTTP");

  intercept(context: ExecutionContext, next: CallHandler): Observable<unknown> {
    const ctx = context.switchToHttp();
    const request = ctx.getRequest<Request>();
    const response = ctx.getResponse<Response>();

    const { method, url, body } = request as { method: string; url: string; body: unknown };
    const correlationId = request.headers["x-correlation-id"] as string;
    const userAgent = request.get("user-agent") ?? "";
    const startTime = Date.now();

    // Log request
    this.logger.log(`→ ${method} ${url} [${correlationId}] - ${userAgent}`);

    // Log request body in debug mode (except for sensitive endpoints)
    if (process.env.NODE_ENV !== "production" && Object.keys(body || {}).length > 0) {
      this.logger.debug(`Request body: ${JSON.stringify(body)}`);
    }

    return next.handle().pipe(
      tap({
        next: () => {
          const duration = Date.now() - startTime;
          const { statusCode } = response;
          this.logger.log(`← ${method} ${url} [${correlationId}] - ${statusCode} - ${duration}ms`);
        },
        error: (error: { status?: number; message?: string }) => {
          const duration = Date.now() - startTime;
          const statusCode = error.status ?? 500;
          this.logger.error(`← ${method} ${url} [${correlationId}] - ${statusCode} - ${duration}ms - ${error.message}`);
        },
      })
    );
  }
}
