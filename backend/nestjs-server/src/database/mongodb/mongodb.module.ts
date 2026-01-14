import { Module, Global } from "@nestjs/common";
import { MongooseModule } from "@nestjs/mongoose";
import { ConfigService } from "@nestjs/config";
import { Reminder, ReminderSchema } from "./schemas/reminder.schema";
import { ReminderGroup, ReminderGroupSchema } from "./schemas/reminder-group.schema";
import { MongoReminderRepository } from "./repositories/reminder.repository";
import { MongoReminderGroupRepository } from "./repositories/reminder-group.repository";
import { REMINDER_REPOSITORY } from "../interfaces/reminder-repository.interface";
import { REMINDER_GROUP_REPOSITORY } from "../interfaces/reminder-group-repository.interface";

/**
 * MongoDB module
 * Provides MongoDB connection and repositories
 * Global module - exports are available in all modules
 */
@Global()
@Module({
  imports: [
    MongooseModule.forRootAsync({
      inject: [ConfigService],
      useFactory: (configService: ConfigService) => ({
        uri: configService.get<string>("database.mongodb.uri"),
        dbName: configService.get<string>("database.mongodb.dbName"),
      }),
    }),
    MongooseModule.forFeature([
      { name: Reminder.name, schema: ReminderSchema },
      { name: ReminderGroup.name, schema: ReminderGroupSchema },
    ]),
  ],
  providers: [
    MongoReminderRepository,
    MongoReminderGroupRepository,
    {
      provide: REMINDER_REPOSITORY,
      useExisting: MongoReminderRepository,
    },
    {
      provide: REMINDER_GROUP_REPOSITORY,
      useExisting: MongoReminderGroupRepository,
    },
  ],
  exports: [REMINDER_REPOSITORY, REMINDER_GROUP_REPOSITORY, MongoReminderRepository, MongoReminderGroupRepository],
})
export class MongodbModule {}
