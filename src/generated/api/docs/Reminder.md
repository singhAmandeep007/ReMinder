# Reminder

A reminder/task item

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **string** | Unique identifier for the reminder | [default to undefined]
**title** | **string** | Title of the reminder | [default to undefined]
**state** | [**ReminderState**](ReminderState.md) |  | [default to undefined]
**isPinned** | **boolean** | Whether the reminder is pinned to top | [default to false]
**createdAt** | **string** | Creation timestamp (ISO 8601 format) | [default to undefined]
**updatedAt** | **string** | Last update timestamp (ISO 8601 format) | [default to undefined]
**group** | [**ReminderGroup**](ReminderGroup.md) | Optional reminder group this reminder belongs to | [optional] [default to undefined]
**dueDate** | **string** | Due date and time for the reminder (ISO 8601 format) | [optional] [default to undefined]
**repeatTimes** | **number** | Number of times to repeat the reminder | [optional] [default to 1]
**repeatInterval** | **number** | Interval to repeat the reminder (in milliseconds) | [optional] [default to 300000]
**currentFocusSession** | [**ReminderFocusSession**](ReminderFocusSession.md) | Current active focus session for this reminder | [optional] [default to undefined]
**focusSessions** | [**Array&lt;ReminderFocusSession&gt;**](ReminderFocusSession.md) | History of focus sessions for this reminder | [default to undefined]

## Example

```typescript
import { Reminder } from './api';

const instance: Reminder = {
    id,
    title,
    state,
    isPinned,
    createdAt,
    updatedAt,
    group,
    dueDate,
    repeatTimes,
    repeatInterval,
    currentFocusSession,
    focusSessions,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
