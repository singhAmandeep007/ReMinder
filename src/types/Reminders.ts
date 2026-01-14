import type {
  Reminder,
  ReminderGroup,
  ReminderFocusSession,
  ReminderState as ReminderStateEnum,
  CreateReminderRequest,
  UpdateReminderRequest,
  CreateReminderGroupRequest,
  UpdateReminderGroupRequest,
  ReminderResponse,
  RemindersResponse,
  ReminderGroupResponse,
  ReminderGroupsResponse,
  MessageResponse,
  ErrorResponse,
} from "../generated/api";

// Re-export ReminderState enum (value export)
export { ReminderState } from "../generated/api";

// Re-export types (type-only exports for isolatedModules compatibility)
export type {
  // Core Models
  Reminder,
  ReminderGroup,
  ReminderFocusSession,
  // Request types
  CreateReminderRequest,
  UpdateReminderRequest,
  CreateReminderGroupRequest,
  UpdateReminderGroupRequest,
  // Response types
  ReminderResponse,
  RemindersResponse,
  ReminderGroupResponse,
  ReminderGroupsResponse,
  MessageResponse,
  ErrorResponse,
} from "../generated/api";

// Legacy type aliases
export type TReminderState = ReminderStateEnum;
export type TReminderFocusSession = ReminderFocusSession;
export type TReminder = Reminder;
export type TReminderGroup = ReminderGroup;

// Generic response wrappers (kept for backward compatibility)
export type TMessageResponsePayload = MessageResponse;
export type TDataResponsePayload<T> = { data: T };

// Error
export type THTTPError = ErrorResponse;

// Reminder request/response types
export type TGetRemindersRequestQueryParams = Partial<{
  groupId: TReminderGroup["id"];
  state: TReminder["state"];
}>;
export type TGetRemindersResponsePayload = RemindersResponse;
export type TGetReminderRequestPathParams = TReminder["id"];
export type TGetReminderResponsePayload = ReminderResponse;
export type TCreateReminderRequestPayload = CreateReminderRequest;
export type TCreateReminderResponsePayload = ReminderResponse;
export type TUpdateReminderRequestPayload = UpdateReminderRequest & { id: string };
export type TUpdateReminderResponsePayload = ReminderResponse;
export type TDeleteReminderRequestPathParams = TReminder["id"];
export type TDeleteReminderResponsePayload = MessageResponse;

// Reminder Group request/response types
export type TGetReminderGroupsResponsePayload = ReminderGroupsResponse;
export type TGetReminderGroupRequestPathParams = TReminderGroup["id"];
export type TGetReminderGroupResponsePayload = ReminderGroupResponse;
export type TCreateReminderGroupRequestPayload = CreateReminderGroupRequest;
export type TCreateReminderGroupResponsePayload = ReminderGroupResponse;
export type TUpdateReminderGroupRequestPayload = UpdateReminderGroupRequest & { id: string };
export type TUpdateReminderGroupResponsePayload = ReminderGroupResponse;
export type TDeleteReminderGroupRequestPathParams = TReminderGroup["id"];
export type TDeleteReminderGroupResponsePayload = MessageResponse;
