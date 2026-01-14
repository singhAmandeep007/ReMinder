import { Controller, Get, VERSION_NEUTRAL } from "@nestjs/common";
import { ApiTags, ApiOperation, ApiResponse } from "@nestjs/swagger";
import {
  HealthCheck,
  HealthCheckService,
  MongooseHealthIndicator,
  MemoryHealthIndicator,
  DiskHealthIndicator,
} from "@nestjs/terminus";
import { SkipThrottle } from "@nestjs/throttler";

/**
 * Health check controller
 * Provides endpoints for Kubernetes liveness and readiness probes
 */
@ApiTags("Health")
@Controller({
  path: "health",
  version: VERSION_NEUTRAL,
})
@SkipThrottle()
export class HealthController {
  constructor(
    private health: HealthCheckService,
    private mongoose: MongooseHealthIndicator,
    private memory: MemoryHealthIndicator,
    private disk: DiskHealthIndicator
  ) {}

  /**
   * Full health check (liveness probe)
   * Checks database, memory, and disk health
   */
  @Get()
  @HealthCheck()
  @ApiOperation({ summary: "Full health check (liveness probe)" })
  @ApiResponse({ status: 200, description: "Service is healthy" })
  @ApiResponse({ status: 503, description: "Service is unhealthy" })
  async check() {
    return this.health.check([
      // Database health
      () => this.mongoose.pingCheck("database"),
      // Memory health (heap should be under 200MB)
      () => this.memory.checkHeap("memory_heap", 200 * 1024 * 1024),
      // Disk health (should have at least 10% free)
      () =>
        this.disk.checkStorage("storage", {
          path: "/",
          thresholdPercent: 0.9,
        }),
    ]);
  }

  /**
   * Readiness probe
   * Only checks if the service can handle requests (database connection)
   */
  @Get("ready")
  @HealthCheck()
  @ApiOperation({ summary: "Readiness probe" })
  @ApiResponse({ status: 200, description: "Service is ready to receive traffic" })
  @ApiResponse({ status: 503, description: "Service is not ready" })
  async ready() {
    return this.health.check([() => this.mongoose.pingCheck("database")]);
  }

  /**
   * Simple liveness check
   * Just returns OK - used for basic health verification
   */
  @Get("live")
  @ApiOperation({ summary: "Simple liveness check" })
  @ApiResponse({ status: 200, description: "Service is alive" })
  live() {
    return {
      status: "ok",
      timestamp: new Date().toISOString(),
    };
  }
}
