/**
 * Re-exports from the auto-generated OpenAPI types
 *
 * This file provides convenient access to the types generated from openapi/api.yaml
 * Use these types throughout the application to ensure consistency with the API contract.
 *
 * To regenerate: npm run generate:types
 */

import type { components, operations, paths } from "./api.types";

// ===================
// Schema Types (Models)
// ===================

/** Current state of a reminder */
export type ReminderState = components["schemas"]["ReminderState"];

/** A focus session for a reminder */
export type ReminderFocusSession = components["schemas"]["ReminderFocusSession"];

/** A group/category for organizing reminders */
export type ReminderGroup = components["schemas"]["ReminderGroup"];

/** A reminder/task item */
export type Reminder = components["schemas"]["Reminder"];

// ===================
// Request Types
// ===================

/** Create reminder request body */
export type CreateReminderRequest = components["schemas"]["CreateReminderRequest"];

/** Update reminder request body */
export type UpdateReminderRequest = components["schemas"]["UpdateReminderRequest"];

/** Create reminder group request body */
export type CreateReminderGroupRequest = components["schemas"]["CreateReminderGroupRequest"];

/** Update reminder group request body */
export type UpdateReminderGroupRequest = components["schemas"]["UpdateReminderGroupRequest"];

// ===================
// Response Types
// ===================

/** Single reminder response */
export type ReminderResponse = components["schemas"]["ReminderResponse"];

/** Multiple reminders response */
export type RemindersResponse = components["schemas"]["RemindersResponse"];

/** Single reminder group response */
export type ReminderGroupResponse = components["schemas"]["ReminderGroupResponse"];

/** Multiple reminder groups response */
export type ReminderGroupsResponse = components["schemas"]["ReminderGroupsResponse"];

/** Generic message response */
export type MessageResponse = components["schemas"]["MessageResponse"];

/** Error response */
export type ErrorResponse = components["schemas"]["ErrorResponse"];

// ===================
// Query Parameters
// ===================

/** Query parameters for listing reminders */
export type GetRemindersQuery = operations["getReminders"]["parameters"]["query"];

// ===================
// Constants (derived from enum)
// ===================

/** All valid reminder states */
export const REMINDER_STATES = ["INACTIVE", "INPROGRESS", "COMPLETED", "ARCHIVED"] as const;

// ===================
// Re-export full types for advanced use cases
// ===================
export type { components, operations, paths };
