## <a name="resource-order1_group">Group</a>


Group API

### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **createdAt** | *date-time* | Group creation date | `"2015-01-01T12:00:00Z"` |
| **id** | *uuid* | Unique group identifier | `"01234567-89ab-cdef-0123-456789abcdef"` |
| **name** | *string* | Group name | `"group1"` |
| **org** | *string* | Group organization | `"tecsisa"` |
| **path** | *string* | Group location | `"/example/admin/"` |
| **urn** | *string* | Group's Uniform Resource Name | `"urn:iws:iam:tecsisa:group/example/admin/group1"` |

### Group Create

Create a new group

```
POST /api/v1/organizations/{organization_id}/groups
```

#### Required Parameters

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **name** | *string* | Group name | `"group1"` |
| **path** | *string* | Group location | `"/example/admin/"` |



#### Curl Example

```bash
$ curl -n -X POST /api/v1/organizations/$ORGANIZATION_ID/groups \
  -d '{
  "name": "group1",
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
  "name": "group1",
  "path": "/example/admin/",
  "createdAt": "2015-01-01T12:00:00Z",
  "urn": "urn:iws:iam:tecsisa:group/example/admin/group1",
  "org": "tecsisa"
}
```

### Group Update

Update an existing group

```
PUT /api/v1/organizations/{organization_id}/groups/{group_name}
```

#### Required Parameters

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **name** | *string* | Group name | `"group1"` |
| **path** | *string* | Group location | `"/example/admin/"` |



#### Curl Example

```bash
$ curl -n -X PUT /api/v1/organizations/$ORGANIZATION_ID/groups/$GROUP_NAME \
  -d '{
  "name": "group1",
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
  "name": "group1",
  "path": "/example/admin/",
  "createdAt": "2015-01-01T12:00:00Z",
  "urn": "urn:iws:iam:tecsisa:group/example/admin/group1",
  "org": "tecsisa"
}
```

### Group Delete

Delete an existing group

```
DELETE /api/v1/organizations/{organization_id}/groups/{group_name}
```


#### Curl Example

```bash
$ curl -n -X DELETE /api/v1/organizations/$ORGANIZATION_ID/groups/$GROUP_NAME \
  -H "Content-Type: application/json" \
  -H "Authorization: Basic or Bearer XXX"
```


#### Response Example

```
HTTP/1.1 202 Accepted
```


### Group Get

Get an existing group

```
GET /api/v1/organizations/{organization_id}/groups/{group_name}
```


#### Curl Example

```bash
$ curl -n /api/v1/organizations/$ORGANIZATION_ID/groups/$GROUP_NAME \
  -H "Authorization: Basic or Bearer XXX"
```


#### Response Example

```
HTTP/1.1 200 OK
```

```json
{
  "id": "01234567-89ab-cdef-0123-456789abcdef",
  "name": "group1",
  "path": "/example/admin/",
  "createdAt": "2015-01-01T12:00:00Z",
  "urn": "urn:iws:iam:tecsisa:group/example/admin/group1",
  "org": "tecsisa"
}
```


## <a name="resource-order2_groupReference"></a>




### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **[groups/name](#resource-order1_group)** | *string* | Group name | `"group1"` |
| **[groups/org](#resource-order1_group)** | *string* | Group organization | `"tecsisa"` |

###  List all organization's groups

List all groups by organization.

```
GET /api/v1/organizations/{organization_id}/groups?PathPrefix={optional_path_prefix}
```


#### Curl Example

```bash
$ curl -n /api/v1/organizations/$ORGANIZATION_ID/groups?PathPrefix=$OPTIONAL_PATH_PREFIX \
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
      "name": "group1"
    }
  ]
}
```

###  List all groups

List all groups.

```
GET /api/v1/groups?PathPrefix={optional_path_prefix}
```


#### Curl Example

```bash
$ curl -n /api/v1/groups?PathPrefix=$OPTIONAL_PATH_PREFIX \
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
      "name": "group1"
    }
  ]
}
```


## <a name="resource-order3_members">Member</a>


Group members

### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **members** | *array* | Identifier of user | `["member1"]` |

### Member Add

Add member to a group.

```
POST /api/v1/organizations/{organization_id}/groups/{group_name}/users/{user_id}
```


#### Curl Example

```bash
$ curl -n -X POST /api/v1/organizations/$ORGANIZATION_ID/groups/$GROUP_NAME/users/$USER_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Basic or Bearer XXX"
```


#### Response Example

```
HTTP/1.1 202 Accepted
```


### Member Remove

Remove member from a group

```
DELETE /api/v1/organizations/{organization_id}/groups/{group_name}/users/{user_id}
```


#### Curl Example

```bash
$ curl -n -X DELETE /api/v1/organizations/$ORGANIZATION_ID/groups/$GROUP_NAME/users/$USER_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Basic or Bearer XXX"
```


#### Response Example

```
HTTP/1.1 202 Accepted
```


### Member List

List members of a group

```
GET /api/v1/organizations/{organization_id}/groups/{group_name}/users
```


#### Curl Example

```bash
$ curl -n /api/v1/organizations/$ORGANIZATION_ID/groups/$GROUP_NAME/users \
  -H "Authorization: Basic or Bearer XXX"
```


#### Response Example

```
HTTP/1.1 200 OK
```

```json
{
  "members": [
    "member1"
  ]
}
```


## <a name="resource-order4_attachedPolicies">Group Policies</a>


Attached Policies

### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **policies/name** | *string* | Name of policy | `"policy1"` |
| **policies/org** | *string* | Organization of policy | `"tecsisa"` |

### Group Policies Attach

Attach policy to group

```
POST /api/v1/organizations/{organization_id}/groups/{group_name}/policies/{policy_id}
```


#### Curl Example

```bash
$ curl -n -X POST /api/v1/organizations/$ORGANIZATION_ID/groups/$GROUP_NAME/policies/$POLICY_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Basic or Bearer XXX"
```


#### Response Example

```
HTTP/1.1 202 Accepted
```


### Group Policies Detach

Detach policy from group

```
DELETE /api/v1/organizations/{organization_id}/groups/{group_name}/policies/{policy_id}
```


#### Curl Example

```bash
$ curl -n -X DELETE /api/v1/organizations/$ORGANIZATION_ID/groups/$GROUP_NAME/policies/$POLICY_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Basic or Bearer XXX"
```


#### Response Example

```
HTTP/1.1 202 Accepted
```


### Group Policies List

List attach policies

```
GET /api/v1/organizations/{organization_id}/groups/{group_name}/policies
```


#### Curl Example

```bash
$ curl -n /api/v1/organizations/$ORGANIZATION_ID/groups/$GROUP_NAME/policies \
  -H "Authorization: Basic or Bearer XXX"
```


#### Response Example

```
HTTP/1.1 200 OK
```

```json
{
  "policies": [
    {
      "org": "tecsisa",
      "name": "policy1"
    }
  ]
}
```


