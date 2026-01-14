import { Injectable, NestMiddleware } from "@nestjs/common";
import { Request, Response, NextFunction } from "express";
import { randomUUID } from "crypto";

/**
 * Correlation ID middleware
 * Adds a unique correlation ID to each request for distributed tracing
 * If the request already has a correlation ID header, it uses that value
 */
@Injectable()
export class CorrelationIdMiddleware implements NestMiddleware {
  use(req: Request, res: Response, next: NextFunction) {
    const correlationId = (req.headers["x-correlation-id"] as string) || randomUUID();

    // Set correlation ID on request headers
    req.headers["x-correlation-id"] = correlationId;

    // Set correlation ID on response headers
    res.setHeader("x-correlation-id", correlationId);

    next();
  }
}
