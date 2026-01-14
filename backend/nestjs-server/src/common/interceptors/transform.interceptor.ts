import { Injectable, NestInterceptor, ExecutionContext, CallHandler } from "@nestjs/common";
import { Observable } from "rxjs";
import { map } from "rxjs/operators";
import { Request } from "express";

/**
 * Standard API response wrapper
 */
export interface ApiResponse<T> {
  data: T;
  meta?: {
    timestamp: string;
    correlationId?: string;
  };
}

/**
 * Transform interceptor
 * Wraps all successful responses in a consistent format
 *
 * Note: This matches the OpenAPI spec response format with `data` wrapper
 */
@Injectable()
export class TransformInterceptor<T> implements NestInterceptor<T, ApiResponse<T>> {
  intercept(context: ExecutionContext, next: CallHandler): Observable<ApiResponse<T>> {
    const request = context.switchToHttp().getRequest<Request>();
    const correlationId = request.headers["x-correlation-id"] as string;

    return next.handle().pipe(
      map((data: T): ApiResponse<T> => {
        const meta = {
          timestamp: new Date().toISOString(),
          correlationId,
        };

        // If data is already wrapped with 'data' key, extract and re-wrap
        if (data && typeof data === "object" && "data" in data) {
          return {
            data: (data as { data: T }).data,
            meta,
          };
        }

        // Wrap response in standard format
        return {
          data,
          meta,
        };
      })
    );
  }
}
