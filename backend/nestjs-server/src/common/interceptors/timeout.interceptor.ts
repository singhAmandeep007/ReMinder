import { Injectable, NestInterceptor, ExecutionContext, CallHandler, RequestTimeoutException } from "@nestjs/common";
import { Observable, throwError, TimeoutError } from "rxjs";
import { catchError, timeout } from "rxjs/operators";

/**
 * Timeout interceptor
 * Automatically times out requests that take too long
 */
@Injectable()
export class TimeoutInterceptor implements NestInterceptor {
  constructor(private readonly timeoutMs = 30000) {}

  intercept(context: ExecutionContext, next: CallHandler): Observable<unknown> {
    return next.handle().pipe(
      timeout(this.timeoutMs),
      catchError((err: unknown) => {
        if (err instanceof TimeoutError) {
          return throwError(() => new RequestTimeoutException("Request timed out"));
        }
        return throwError(() => err);
      })
    );
  }
}
