import { ExceptionFilter, Catch, ArgumentsHost, HttpException, HttpStatus, Logger } from "@nestjs/common";
import { Request, Response } from "express";

interface ErrorResponse {
  statusCode: number;
  message: string | string[];
  error: string;
  timestamp: string;
  path: string;
  correlationId?: string;
}

/**
 * Global HTTP exception filter
 * Catches all HTTP exceptions and formats them consistently
 */
@Catch(HttpException)
export class HttpExceptionFilter implements ExceptionFilter {
  private readonly logger = new Logger(HttpExceptionFilter.name);

  catch(exception: HttpException, host: ArgumentsHost) {
    const ctx = host.switchToHttp();
    const response = ctx.getResponse<Response>();
    const request = ctx.getRequest<Request>();
    const status = exception.getStatus();
    const exceptionResponse = exception.getResponse();

    const errorResponse: ErrorResponse = {
      statusCode: status,
      message:
        typeof exceptionResponse === "object" && "message" in exceptionResponse
          ? (exceptionResponse as { message: string | string[] }).message
          : exception.message,
      error:
        typeof exceptionResponse === "object" && "error" in exceptionResponse
          ? (exceptionResponse as { error: string }).error
          : HttpStatus[status],
      timestamp: new Date().toISOString(),
      path: request.url,
      correlationId: request.headers["x-correlation-id"] as string,
    };

    // Log the error
    this.logger.error(
      `${request.method} ${request.url} - ${status} - ${JSON.stringify(errorResponse.message)}`,
      exception.stack
    );

    response.status(status).json(errorResponse);
  }
}
