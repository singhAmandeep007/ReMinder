# Reminders API Documentation

## Types

1. Reminder State

```typescript
type TReminderState = "INACTIVE" | "INPROGRESS" | "COMPLETED" | "ARCHIVED";
```

2. Reminder Focus Session

```typescript
type TReminderFocusSession = {
  startTime: string;
  endTime: string;
};
```

3. Reminder

```typescript
type TReminder = {
  id: string;
  title: string;
  /**
   * Current state of the reminder.
   *
   * @example "IN_ACTIVE" | "IN_PROGRESS" | "COMPLETED"
   */
  state: TReminderState;
  /**
   * Indicates if the reminder is pinned on top.
   */
  isPinned: boolean;
  /**
   * Date and time the reminder was created. (ISO 8601 format)
   */
  createdAt: string;
  /**
   * Date and time the reminder was last updated. (ISO 8601 format)
   */
  updatedAt: string;
  /**
   * Optional reminder group this reminder belongs to.
   */
  group: TReminderGroup | null;
  /**
   * Due date and time for the reminder. (ISO 8601 format)
   *
   * @example "2021-09-30T00:00:00.000Z"
   */
  dueDate: string | null;
  /**
   * Number of times to repeat the reminder.
   *
   * @default 1
   */
  repeatTimes: number | null;
  /**
   * Interval to repeat the reminder. (in milliseconds)
   *
   * @default 300000 (5 minutes)
   */
  repeatInterval: number | null;

  /**
   * Current focus session for this reminder.
   */
  currentFocusSession: {
    /**
     * date and time when focus session was started for this reminder. (ISO 8601 format)
     */
    startTime: string;
    /**
     * date and time the focus session stopped for this reminder. (ISO 8601 format)
     */
    endTime: string;
  } | null;
  /**
   * Array to store history of multiple focus sessions for this reminder.
   * Each session includes start time and end time.
   */
  focusSessions: TReminderFocusSession[];
};
```

## Endpoints

### 1. Get Reminders

- **Method:** GET
- **URL:** `/reminders?groupId={groupId}&state={state}`
- **Description:** Fetch a list of reminders.

#### Request Query Parameters

```typescript
type TGetRemindersRequestQueryParams = {
  groupId?: string; // ID of the reminder group.
  state?: TReminderState; // State of the reminders.
};
```

#### Response

- **Status:** 200 OK
- **Body:**
  ```typescript
  type TGetRemindersResponsePayload = {
    data: TReminder[];
  };
  ```

### 2. Get Reminder

- **Method:** GET
- **URL:** `/reminders/{id}`
- **Description:** Fetch a single reminder by ID.

#### Request Path Parameters

```typescript
type TGetReminderRequestPathParams = string;
```

#### Response

- **Status:** 200 OK
- **Body:**
  ```typescript
  type TGetReminderResponsePayload = {
    data: TReminder;
  };
  ```

### 3. Create Reminder

- **Method:** POST
- **URL:** `/reminders`
- **Description:** Create a new reminder.

#### Request Payload

```typescript
type TCreateReminderRequestPayload = {
  title: string;
  groupId?: string;
};
```

#### Response

- **Status:** 201 Created
- **Body:**
  ```typescript
  type TCreateReminderResponsePayload = {
    data: TReminder;
  };
  ```

### 4. Update Reminder

- **Method:** PUT
- **URL:** `/reminders/{id}`
- **Description:** Update an existing reminder.

#### Request Path Parameters

```typescript
type TUpdateReminderRequestPathParams = string;
```

#### Request Payload

```typescript
type TUpdateReminderRequestPayload = {
  title?: string;
  state?: TReminderState;
  isPinned?: boolean;
  dueDate?: string;
  currentFocusSession?: TReminderFocusSession;
  focusSessions?: TReminderFocusSession[];
  repeatTimes?: number;
  repeatInterval?: number;
  groupId?: string;
};
```

#### Response

- **Status:** 201 Created
- **Body:**
  ```typescript
  type TUpdateReminderResponsePayload = {
    data: TReminder;
  };
  ```

### 5. Delete Reminder

- **Method:** DELETE
- **URL:** `/reminders/{id}`
- **Description:** Delete an existing reminder.

#### Request Path Parameters

```typescript
type TDeleteReminderRequestPathParams = string;
```

#### Response

- **Status:** 200 OK
- **Body:**
  ```typescript
  type TDeleteReminderResponsePayload = {
    message: string;
  };
  ```
