import { Module } from "@nestjs/common";
import { ReminderGroupsController } from "./reminder-groups.controller";
import { ReminderGroupsService } from "./reminder-groups.service";

/**
 * Reminder Groups feature module
 */
@Module({
  controllers: [ReminderGroupsController],
  providers: [ReminderGroupsService],
  exports: [ReminderGroupsService],
})
export class ReminderGroupsModule {}
