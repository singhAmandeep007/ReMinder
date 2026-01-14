/**
 * Reminder Group Repository Interface
 *
 * Uses types from OpenAPI-generated types (single source of truth).
 */

import type { ReminderGroup, CreateReminderGroupRequest, UpdateReminderGroupRequest } from "../../generated";
import type { IBaseRepository } from "./base-repository.interface";

/**
 * Injection token for the reminder group repository
 */
export const REMINDER_GROUP_REPOSITORY = Symbol("REMINDER_GROUP_REPOSITORY");

/**
 * Filter options for querying reminder groups
 */
export interface IReminderGroupFilter {
  name?: string;
}

/**
 * Reminder group repository interface
 * Extends base repository with group-specific operations
 */
export interface IReminderGroupRepository
  extends IBaseRepository<ReminderGroup, CreateReminderGroupRequest, UpdateReminderGroupRequest, IReminderGroupFilter> {
  /**
   * Find a reminder group by name
   */
  findByName(name: string): Promise<ReminderGroup | null>;
}
