## <a name="resource-order1_group">Group</a>


Group API

### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **CreatedAt** | *date-time* | When group was created | `"2015-01-01T12:00:00Z"` |
| **ID** | *uuid* | Unique identifier of group | `"01234567-89ab-cdef-0123-456789abcdef"` |
| **Name** | *string* | Name of group | `"group1"` |
| **Org** | *string* | Organization of group | `"tecsisa"` |
| **Path** | *string* | Group's location | `"/example/admin/"` |
| **Urn** | *string* | Uniform Resource Name of group | `"urn:iws:iam:tecsisa:group/example/admin/group1"` |

### Group Create

Create a new group.

```
POST /api/v1/organizations/{organization_id}/groups
```

#### Required Parameters

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **Name** | *string* | Name of group | `"group1"` |
| **Path** | *string* | Group's location | `"/example/admin/"` |



#### Curl Example

```bash
$ curl -n -X POST /api/v1/organizations/$ORGANIZATION_ID/groups \
  -d '{
  "Name": "group1",
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
  "Name": "group1",
  "Path": "/example/admin/",
  "CreatedAt": "2015-01-01T12:00:00Z",
  "Urn": "urn:iws:iam:tecsisa:group/example/admin/group1",
  "Org": "tecsisa"
}
```

### Group Update

Update an existing group.

```
PUT /api/v1/organizations/{organization_id}/groups/{group_name}
```

#### Required Parameters

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **Name** | *string* | Name of group | `"group1"` |
| **Path** | *string* | Group's location | `"/example/admin/"` |



#### Curl Example

```bash
$ curl -n -X PUT /api/v1/organizations/$ORGANIZATION_ID/groups/$GROUP_NAME \
  -d '{
  "Name": "group1",
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
  "Name": "group1",
  "Path": "/example/admin/",
  "CreatedAt": "2015-01-01T12:00:00Z",
  "Urn": "urn:iws:iam:tecsisa:group/example/admin/group1",
  "Org": "tecsisa"
}
```

### Group Delete

Delete an existing group.

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

Get an existing group.

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
  "ID": "01234567-89ab-cdef-0123-456789abcdef",
  "Name": "group1",
  "Path": "/example/admin/",
  "CreatedAt": "2015-01-01T12:00:00Z",
  "Urn": "urn:iws:iam:tecsisa:group/example/admin/group1",
  "Org": "tecsisa"
}
```


## <a name="resource-order2_groupReference"></a>




### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **[Name](#resource-order1_group)** | *string* | Name of group | `"group1"` |
| **[Org](#resource-order1_group)** | *string* | Organization of group | `"tecsisa"` |

###  Group List All

List all groups by organization.

```
GET /api/v1/organizations/{organization_id}/groups
```


#### Curl Example

```bash
$ curl -n /api/v1/organizations/$ORGANIZATION_ID/groups \
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


## <a name="resource-order3_members">Member</a>


Members of a group.

### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **ExternalID** | *array* | Identifier of user | `["member1"]` |

### Member Add

Add member to a group

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
  "ExternalID": [
    "member1"
  ]
}
```


## <a name="resource-order4_attachedPolicies">Group Policies</a>


Attached Policies

### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **Name** | *string* | Name of policy | `"policy1"` |
| **Org** | *string* | Organization of policy | `"tecsisa"` |

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

Detach policy to group

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
[
  {
    "Org": "tecsisa",
    "Name": "policy1"
  }
]
```


