# RemindersApi

All URIs are relative to *http://localhost:8080/api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**createReminder**](#createreminder) | **POST** /reminders | Create a reminder|
|[**deleteReminder**](#deletereminder) | **DELETE** /reminders/{id} | Delete a reminder|
|[**getReminders**](#getreminders) | **GET** /reminders | List reminders|
|[**updateReminder**](#updatereminder) | **PATCH** /reminders/{id} | Update a reminder|

# **createReminder**
> ReminderResponse createReminder(createReminderRequest)

Create a new reminder

### Example

```typescript
import {
    RemindersApi,
    Configuration,
    CreateReminderRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new RemindersApi(configuration);

let createReminderRequest: CreateReminderRequest; //

const { status, data } = await apiInstance.createReminder(
    createReminderRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **createReminderRequest** | **CreateReminderRequest**|  | |


### Return type

**ReminderResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Reminder created successfully |  -  |
|**400** | Invalid request body |  -  |
|**404** | Reminder group not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteReminder**
> MessageResponse deleteReminder()

Delete an existing reminder by ID

### Example

```typescript
import {
    RemindersApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RemindersApi(configuration);

let id: string; //Reminder ID (default to undefined)

const { status, data } = await apiInstance.deleteReminder(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] | Reminder ID | defaults to undefined|


### Return type

**MessageResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Reminder deleted successfully |  -  |
|**404** | Reminder not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getReminders**
> RemindersResponse getReminders()

Retrieve a list of reminders with optional filtering by group and state

### Example

```typescript
import {
    RemindersApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RemindersApi(configuration);

let groupId: string; //Filter by reminder group ID (optional) (default to undefined)
let state: ReminderState; //Filter by reminder state (optional) (default to undefined)

const { status, data } = await apiInstance.getReminders(
    groupId,
    state
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **groupId** | [**string**] | Filter by reminder group ID | (optional) defaults to undefined|
| **state** | **ReminderState** | Filter by reminder state | (optional) defaults to undefined|


### Return type

**RemindersResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Successful operation |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateReminder**
> ReminderResponse updateReminder(updateReminderRequest)

Update an existing reminder by ID

### Example

```typescript
import {
    RemindersApi,
    Configuration,
    UpdateReminderRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new RemindersApi(configuration);

let id: string; //Reminder ID (default to undefined)
let updateReminderRequest: UpdateReminderRequest; //

const { status, data } = await apiInstance.updateReminder(
    id,
    updateReminderRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **updateReminderRequest** | **UpdateReminderRequest**|  | |
| **id** | [**string**] | Reminder ID | defaults to undefined|


### Return type

**ReminderResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Reminder updated successfully |  -  |
|**404** | Reminder not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

