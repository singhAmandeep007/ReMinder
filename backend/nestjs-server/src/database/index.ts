/**
 * Database Module Exports
 *
 * All entity types are derived from OpenAPI-generated types (single source of truth).
 * This module provides repository interfaces and implementations.
 */

// Database module
export * from "./database.module";

// Interfaces (contracts for repositories)
export * from "./interfaces";

// MongoDB implementation (only export module and repositories, not schemas)
export { MongodbModule } from "./mongodb/mongodb.module";
export { MongoReminderRepository } from "./mongodb/repositories/reminder.repository";
export { MongoReminderGroupRepository } from "./mongodb/repositories/reminder-group.repository";
