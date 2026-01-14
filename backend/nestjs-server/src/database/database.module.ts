import { Module, DynamicModule, Logger } from "@nestjs/common";
import { ConfigModule } from "@nestjs/config";
import { MongodbModule } from "./mongodb/mongodb.module";

/**
 * Dynamic database module
 * Selects the appropriate database implementation based on configuration
 */
@Module({})
export class DatabaseModule {
  private static readonly logger = new Logger(DatabaseModule.name);

  static forRoot(): DynamicModule {
    return {
      module: DatabaseModule,
      imports: [ConfigModule],
      providers: [],
      exports: [],
    };
  }

  /**
   * Create dynamic module that imports the correct database module
   * based on the DB_TYPE configuration
   */
  static register(): DynamicModule {
    return {
      module: DatabaseModule,
      imports: [
        // We always import MongoDB for now
        // When Firebase is implemented, add conditional logic here
        MongodbModule,
      ],
      exports: [MongodbModule],
    };
  }
}

// Re-export the simple async module for app.module.ts
@Module({
  imports: [MongodbModule],
  exports: [MongodbModule],
})
export class DatabaseModuleAsync {}
