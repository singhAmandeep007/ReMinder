import { Module } from "@nestjs/common";
import { RemindersController } from "./reminders.controller";
import { RemindersService } from "./reminders.service";

/**
 * Reminders feature module
 */
@Module({
  controllers: [RemindersController],
  providers: [RemindersService],
  exports: [RemindersService],
})
export class RemindersModule {}
