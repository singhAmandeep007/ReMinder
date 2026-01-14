# UpdateReminderRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**title** | **string** | Updated title | [optional] [default to undefined]
**state** | [**ReminderState**](ReminderState.md) |  | [optional] [default to undefined]
**isPinned** | **boolean** | Whether the reminder is pinned | [optional] [default to undefined]
**dueDate** | **string** | Due date and time | [optional] [default to undefined]
**repeatTimes** | **number** | Number of times to repeat | [optional] [default to undefined]
**repeatInterval** | **number** | Repeat interval in milliseconds | [optional] [default to undefined]
**currentFocusSession** | [**ReminderFocusSession**](ReminderFocusSession.md) |  | [optional] [default to undefined]
**focusSessions** | [**Array&lt;ReminderFocusSession&gt;**](ReminderFocusSession.md) |  | [optional] [default to undefined]
**groupId** | **string** | ID of the reminder group to move this reminder to | [optional] [default to undefined]

## Example

```typescript
import { UpdateReminderRequest } from './api';

const instance: UpdateReminderRequest = {
    title,
    state,
    isPinned,
    dueDate,
    repeatTimes,
    repeatInterval,
    currentFocusSession,
    focusSessions,
    groupId,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
