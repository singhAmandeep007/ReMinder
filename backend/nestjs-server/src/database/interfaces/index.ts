/**
 * Database Interfaces Barrel Export
 *
 * All types are derived from OpenAPI-generated types.
 */

// Base repository
export * from "./base-repository.interface";

// Reminder repository
export { REMINDER_REPOSITORY, type IReminderRepository, type IReminderFilter } from "./reminder-repository.interface";

// Reminder group repository
export {
  REMINDER_GROUP_REPOSITORY,
  type IReminderGroupRepository,
  type IReminderGroupFilter,
} from "./reminder-group-repository.interface";

// Re-export commonly used types from generated
export type {
  Reminder,
  ReminderGroup,
  ReminderFocusSession,
  ReminderState,
  CreateReminderRequest,
  UpdateReminderRequest,
  CreateReminderGroupRequest,
  UpdateReminderGroupRequest,
} from "../../generated";
