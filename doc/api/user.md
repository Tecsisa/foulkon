## <a name="resource-order1_user">User</a>


User API

### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **createdAt** | *date-time* | User creation date | `"2015-01-01T12:00:00Z"` |
| **externalId** | *string* | User's external identifier | `"user1"` |
| **id** | *uuid* | Unique user identifier | `"01234567-89ab-cdef-0123-456789abcdef"` |
| **path** | *string* | User location | `"/example/admin/"` |
| **updateAt** | *date-time* | The date timestamp of the last update | `"2015-01-01T12:00:00Z"` |
| **urn** | *string* | User's Uniform Resource Name | `"urn:iws:iam::user/example/admin/user1"` |

### User Create

Create a new user.

```
POST /api/v1/users
```

#### Required Parameters

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **externalId** | *string* | User's external identifier | `"user1"` |
| **path** | *string* | User location | `"/example/admin/"` |



#### Curl Example

```bash
$ curl -n -X POST /api/v1/users \
  -d '{
  "externalId": "user1",
  "path": "/example/admin/"
}' \
  -H "Content-Type: application/json" \
  -H "Authorization: Basic or Bearer XXX"
```


#### Response Example

```
HTTP/1.1 201 Created
```

```json
{
  "id": "01234567-89ab-cdef-0123-456789abcdef",
  "externalId": "user1",
  "path": "/example/admin/",
  "createdAt": "2015-01-01T12:00:00Z",
  "updateAt": "2015-01-01T12:00:00Z",
  "urn": "urn:iws:iam::user/example/admin/user1"
}
```

### User Update

Update an existing user.

```
PUT /api/v1/users/{user_externalID}
```

#### Required Parameters

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **path** | *string* | User location | `"/example/admin/"` |



#### Curl Example

```bash
$ curl -n -X PUT /api/v1/users/$USER_EXTERNALID \
  -d '{
  "path": "/example/admin/"
}' \
  -H "Content-Type: application/json" \
  -H "Authorization: Basic or Bearer XXX"
```


#### Response Example

```
HTTP/1.1 200 OK
```

```json
{
  "id": "01234567-89ab-cdef-0123-456789abcdef",
  "externalId": "user1",
  "path": "/example/admin/",
  "createdAt": "2015-01-01T12:00:00Z",
  "updateAt": "2015-01-01T12:00:00Z",
  "urn": "urn:iws:iam::user/example/admin/user1"
}
```

### User Delete

Delete an existing user.

```
DELETE /api/v1/users/{user_externalID}
```


#### Curl Example

```bash
$ curl -n -X DELETE /api/v1/users/$USER_EXTERNALID \
  -H "Content-Type: application/json" \
  -H "Authorization: Basic or Bearer XXX"
```


#### Response Example

```
HTTP/1.1 202 Accepted
```


### User Get

Get an existing user.

```
GET /api/v1/users/{user_externalID}
```


#### Curl Example

```bash
$ curl -n /api/v1/users/$USER_EXTERNALID \
  -H "Authorization: Basic or Bearer XXX"
```


#### Response Example

```
HTTP/1.1 200 OK
```

```json
{
  "id": "01234567-89ab-cdef-0123-456789abcdef",
  "externalId": "user1",
  "path": "/example/admin/",
  "createdAt": "2015-01-01T12:00:00Z",
  "updateAt": "2015-01-01T12:00:00Z",
  "urn": "urn:iws:iam::user/example/admin/user1"
}
```


## <a name="resource-order2_userReference"></a>




### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **limit** | *integer* | The maximum number of items in the response (as set in the query or by default) | `20` |
| **offset** | *integer* | The offset of the items returned (as set in the query or by default) | `0` |
| **total** | *integer* | The total number of items available to return | `50` |
| **users** | *array* | User identifiers | `["User1","User2"]` |

###  User List All

List all users filtered, using optional query parameters.

```
GET /api/v1/users?PathPrefix={optional_path_prefix}&Offset={optional_offset}&Limit={optional_limit}&OrderBy={columnName-desc}
```


#### Curl Example

```bash
$ curl -n /api/v1/users?PathPrefix=$OPTIONAL_PATH_PREFIX&Offset=$OPTIONAL_OFFSET&Limit=$OPTIONAL_LIMIT&OrderBy=$COLUMNNAME-DESC \
  -H "Authorization: Basic or Bearer XXX"
```


#### Response Example

```
HTTP/1.1 200 OK
```

```json
{
  "users": [
    "User1",
    "User2"
  ],
  "offset": 0,
  "limit": 20,
  "total": 50
}
```


## <a name="resource-order3_groupIdentity"></a>




### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **groups/joined** | *date-time* | When relationship was created | `"2015-01-01T12:00:00Z"` |
| **groups/name** | *string* | Group name | `"group1"` |
| **groups/org** | *string* | Group organization | `"tecsisa"` |
| **limit** | *integer* | The maximum number of items in the response (as set in the query or by default) | `20` |
| **offset** | *integer* | The offset of the items returned (as set in the query or by default) | `0` |
| **total** | *integer* | The total number of items available to return | `50` |

###  List user groups

List all groups that a user is a member.

```
GET /api/v1/users/{user_externalId}/groups?Offset={optional_offset}&Limit={optional_limit}&OrderBy={columnName-asc}
```


#### Curl Example

```bash
$ curl -n /api/v1/users/$USER_EXTERNALID/groups?Offset=$OPTIONAL_OFFSET&Limit=$OPTIONAL_LIMIT&OrderBy=$COLUMNNAME-ASC \
  -H "Authorization: Basic or Bearer XXX"
```


#### Response Example

```
HTTP/1.1 200 OK
```

```json
{
  "groups": [
    {
      "org": "tecsisa",
      "name": "group1",
      "joined": "2015-01-01T12:00:00Z"
    }
  ],
  "offset": 0,
  "limit": 20,
  "total": 50
}
```


