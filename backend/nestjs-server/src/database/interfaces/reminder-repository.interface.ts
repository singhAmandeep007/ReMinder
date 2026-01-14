/**
 * Reminder Repository Interface
 *
 * Uses types from OpenAPI-generated types (single source of truth).
 */

import type { Reminder, ReminderState, CreateReminderRequest, UpdateReminderRequest } from "../../generated";
import type { IBaseRepository } from "./base-repository.interface";

/**
 * Injection token for the reminder repository
 */
export const REMINDER_REPOSITORY = Symbol("REMINDER_REPOSITORY");

/**
 * Filter options for querying reminders
 */
export interface IReminderFilter {
  state?: ReminderState;
  isPinned?: boolean;
  groupId?: string;
}

/**
 * Reminder repository interface
 * Extends base repository with reminder-specific operations
 */
export interface IReminderRepository
  extends IBaseRepository<Reminder, CreateReminderRequest, UpdateReminderRequest, IReminderFilter> {
  /**
   * Find reminders by group ID
   */
  findByGroupId(groupId: string): Promise<Reminder[]>;

  /**
   * Find reminders by state
   */
  findByState(state: ReminderState): Promise<Reminder[]>;

  /**
   * Find reminders by group ID and state
   */
  findByGroupIdAndState(groupId: string, state: ReminderState): Promise<Reminder[]>;

  /**
   * Delete all reminders in a group
   */
  deleteByGroupId(groupId: string): Promise<number>;
}
