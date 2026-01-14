# ReminderGroupsApi

All URIs are relative to *http://localhost:8080/api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**createReminderGroup**](#createremindergroup) | **POST** /reminder-groups | Create a reminder group|
|[**deleteReminderGroup**](#deleteremindergroup) | **DELETE** /reminder-groups/{id} | Delete a reminder group|
|[**getReminderGroup**](#getremindergroup) | **GET** /reminder-groups/{id} | Get a reminder group|
|[**getReminderGroups**](#getremindergroups) | **GET** /reminder-groups | List reminder groups|
|[**updateReminderGroup**](#updateremindergroup) | **PATCH** /reminder-groups/{id} | Update a reminder group|

# **createReminderGroup**
> ReminderGroupResponse createReminderGroup(createReminderGroupRequest)

Create a new reminder group

### Example

```typescript
import {
    ReminderGroupsApi,
    Configuration,
    CreateReminderGroupRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new ReminderGroupsApi(configuration);

let createReminderGroupRequest: CreateReminderGroupRequest; //

const { status, data } = await apiInstance.createReminderGroup(
    createReminderGroupRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **createReminderGroupRequest** | **CreateReminderGroupRequest**|  | |


### Return type

**ReminderGroupResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Reminder group created successfully |  -  |
|**400** | Invalid request body |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteReminderGroup**
> MessageResponse deleteReminderGroup()

Delete an existing reminder group by ID. This also deletes all reminders in the group.

### Example

```typescript
import {
    ReminderGroupsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new ReminderGroupsApi(configuration);

let id: string; //Reminder group ID (default to undefined)

const { status, data } = await apiInstance.deleteReminderGroup(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] | Reminder group ID | defaults to undefined|


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
|**200** | Reminder group deleted successfully |  -  |
|**404** | Reminder group not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getReminderGroup**
> ReminderGroupResponse getReminderGroup()

Retrieve a single reminder group by ID

### Example

```typescript
import {
    ReminderGroupsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new ReminderGroupsApi(configuration);

let id: string; //Reminder group ID (default to undefined)

const { status, data } = await apiInstance.getReminderGroup(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] | Reminder group ID | defaults to undefined|


### Return type

**ReminderGroupResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Successful operation |  -  |
|**404** | Reminder group not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getReminderGroups**
> ReminderGroupsResponse getReminderGroups()

Retrieve all reminder groups

### Example

```typescript
import {
    ReminderGroupsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new ReminderGroupsApi(configuration);

const { status, data } = await apiInstance.getReminderGroups();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**ReminderGroupsResponse**

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

# **updateReminderGroup**
> ReminderGroupResponse updateReminderGroup(updateReminderGroupRequest)

Update an existing reminder group by ID

### Example

```typescript
import {
    ReminderGroupsApi,
    Configuration,
    UpdateReminderGroupRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new ReminderGroupsApi(configuration);

let id: string; //Reminder group ID (default to undefined)
let updateReminderGroupRequest: UpdateReminderGroupRequest; //

const { status, data } = await apiInstance.updateReminderGroup(
    id,
    updateReminderGroupRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **updateReminderGroupRequest** | **UpdateReminderGroupRequest**|  | |
| **id** | [**string**] | Reminder group ID | defaults to undefined|


### Return type

**ReminderGroupResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Reminder group updated successfully |  -  |
|**404** | Reminder group not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

