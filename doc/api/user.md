## <a name="resource-order1_user">User</a>


User API

### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **CreatedAt** | *date-time* | When user was created | `"2015-01-01T12:00:00Z"` |
| **ExternalID** | *string* | Identifier of user | `"user1"` |
| **ID** | *uuid* | Unique identifier of user | `"01234567-89ab-cdef-0123-456789abcdef"` |
| **Path** | *string* | User's location | `"/example/admin/"` |
| **Urn** | *string* | Uniform Resource Name of user | `"urn:iws:iam::user/example/admin/user1"` |

### User Create

Create a new user.

```
POST /api/v1/users
```

#### Required Parameters

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **ExternalID** | *string* | Identifier of user | `"user1"` |
| **Path** | *string* | User's location | `"/example/admin/"` |



#### Curl Example

```bash
$ curl -n -X POST /api/v1/users \
  -d '{
  "ExternalID": "user1",
  "Path": "/example/admin/"
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
  "ID": "01234567-89ab-cdef-0123-456789abcdef",
  "ExternalID": "user1",
  "Path": "/example/admin/",
  "CreatedAt": "2015-01-01T12:00:00Z",
  "Urn": "urn:iws:iam::user/example/admin/user1"
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
| **Path** | *string* | User's location | `"/example/admin/"` |



#### Curl Example

```bash
$ curl -n -X PUT /api/v1/users/$USER_EXTERNALID \
  -d '{
  "Path": "/example/admin/"
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
  "ID": "01234567-89ab-cdef-0123-456789abcdef",
  "ExternalID": "user1",
  "Path": "/example/admin/",
  "CreatedAt": "2015-01-01T12:00:00Z",
  "Urn": "urn:iws:iam::user/example/admin/user1"
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
  "ID": "01234567-89ab-cdef-0123-456789abcdef",
  "ExternalID": "user1",
  "Path": "/example/admin/",
  "CreatedAt": "2015-01-01T12:00:00Z",
  "Urn": "urn:iws:iam::user/example/admin/user1"
}
```


## <a name="resource-order2_userReference"></a>




### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **ExternalID** | *array* | Identifier of user | `["User1"]` |

###  User List All

List all users filtered by PathPrefix.

```
GET /api/v1/users?PathPrefix={OptionalPath}
```


#### Curl Example

```bash
$ curl -n /api/v1/users?PathPrefix=$OPTIONALPATH \
  -H "Authorization: Basic or Bearer XXX"
```


#### Response Example

```
HTTP/1.1 200 OK
```

```json
{
  "ExternalID": [
    "User1"
  ]
}
```


## <a name="resource-order3_groupIdentity"></a>




### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **Name** | *string* | Name of group | `"group1"` |
| **Org** | *string* | Organization of group | `"tecsisa"` |

###  List user groups

List all groups that a user is a member.

```
GET /api/v1/users/{user_externalID}/groups
```


#### Curl Example

```bash
$ curl -n /api/v1/users/$USER_EXTERNALID/groups \
  -H "Authorization: Basic or Bearer XXX"
```


#### Response Example

```
HTTP/1.1 200 OK
```

```json
[
  {
    "Org": "tecsisa",
    "Name": "group1"
  }
]
```


