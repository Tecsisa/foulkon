## <a name="resource-order1_statement">Statement</a>


Policy statement

### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **actions** | *array* | Operations over resources | `["iam:getUser","iam:*"]` |
| **effect** | *string* | allow/deny resources | `"allow"` |
| **resources** | *array* | resources | `["urn:everything:*"]` |


## <a name="resource-order2_policy">Policy</a>


Policy API

### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **createdAt** | *date-time* | Policy creation date | `"2015-01-01T12:00:00Z"` |
| **id** | *uuid* | Unique policy identifier | `"01234567-89ab-cdef-0123-456789abcdef"` |
| **name** | *string* | Policy name | `"policy1"` |
| **org** | *string* | Policy organization | `"tecsisa"` |
| **path** | *string* | Policy location | `"/example/admin/"` |
| **statements** | *array* | Policy statements | `[{"effect":"allow","actions":["iam:getUser","iam:*"],"resources":["urn:everything:*"]}]` |
| **updateAt** | *date-time* | The date timestamp of the last update | `"2015-01-01T12:00:00Z"` |
| **urn** | *string* | Policy's Uniform Resource Name | `"urn:iws:iam:org1:policy/example/admin/policy1"` |

### Policy Create

Create a new policy.

```
POST /api/v1/organizations/{organization_id}/policies
```

#### Required Parameters

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **name** | *string* | Policy name | `"policy1"` |
| **path** | *string* | Policy location | `"/example/admin/"` |
| **statements** | *array* | Policy statements | `[{"effect":"allow","actions":["iam:getUser","iam:*"],"resources":["urn:everything:*"]}]` |



#### Curl Example

```bash
$ curl -n -X POST /api/v1/organizations/$ORGANIZATION_ID/policies \
  -d '{
  "name": "policy1",
  "path": "/example/admin/",
  "statements": [
    {
      "effect": "allow",
      "actions": [
        "iam:getUser",
        "iam:*"
      ],
      "resources": [
        "urn:everything:*"
      ]
    }
  ]
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
  "name": "policy1",
  "path": "/example/admin/",
  "createdAt": "2015-01-01T12:00:00Z",
  "updateAt": "2015-01-01T12:00:00Z",
  "urn": "urn:iws:iam:org1:policy/example/admin/policy1",
  "org": "tecsisa",
  "statements": [
    {
      "effect": "allow",
      "actions": [
        "iam:getUser",
        "iam:*"
      ],
      "resources": [
        "urn:everything:*"
      ]
    }
  ]
}
```

### Policy Update

Update an existing policy.

```
PUT /api/v1/organizations/{organization_id}/policies/{policy_name}
```

#### Required Parameters

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **name** | *string* | Policy name | `"policy1"` |
| **path** | *string* | Policy location | `"/example/admin/"` |
| **statements** | *array* | Policy statements | `[{"effect":"allow","actions":["iam:getUser","iam:*"],"resources":["urn:everything:*"]}]` |



#### Curl Example

```bash
$ curl -n -X PUT /api/v1/organizations/$ORGANIZATION_ID/policies/$POLICY_NAME \
  -d '{
  "name": "policy1",
  "path": "/example/admin/",
  "statements": [
    {
      "effect": "allow",
      "actions": [
        "iam:getUser",
        "iam:*"
      ],
      "resources": [
        "urn:everything:*"
      ]
    }
  ]
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
  "name": "policy1",
  "path": "/example/admin/",
  "createdAt": "2015-01-01T12:00:00Z",
  "updateAt": "2015-01-01T12:00:00Z",
  "urn": "urn:iws:iam:org1:policy/example/admin/policy1",
  "org": "tecsisa",
  "statements": [
    {
      "effect": "allow",
      "actions": [
        "iam:getUser",
        "iam:*"
      ],
      "resources": [
        "urn:everything:*"
      ]
    }
  ]
}
```

### Policy Delete

Delete an existing policy.

```
DELETE /api/v1/organizations/{organization_id}/policies/{policy_name}
```


#### Curl Example

```bash
$ curl -n -X DELETE /api/v1/organizations/$ORGANIZATION_ID/policies/$POLICY_NAME \
  -H "Content-Type: application/json" \
  -H "Authorization: Basic or Bearer XXX"
```


#### Response Example

```
HTTP/1.1 202 Accepted
```


### Policy Get

Get an existing policy.

```
GET /api/v1/organizations/{organization_id}/policies/{policy_name}
```


#### Curl Example

```bash
$ curl -n /api/v1/organizations/$ORGANIZATION_ID/policies/$POLICY_NAME \
  -H "Authorization: Basic or Bearer XXX"
```


#### Response Example

```
HTTP/1.1 200 OK
```

```json
{
  "id": "01234567-89ab-cdef-0123-456789abcdef",
  "name": "policy1",
  "path": "/example/admin/",
  "createdAt": "2015-01-01T12:00:00Z",
  "updateAt": "2015-01-01T12:00:00Z",
  "urn": "urn:iws:iam:org1:policy/example/admin/policy1",
  "org": "tecsisa",
  "statements": [
    {
      "effect": "allow",
      "actions": [
        "iam:getUser",
        "iam:*"
      ],
      "resources": [
        "urn:everything:*"
      ]
    }
  ]
}
```


## <a name="resource-order3_policyReference">Organization's policies</a>




### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **limit** | *integer* | The maximum number of items in the response (as set in the query or by default) | `20` |
| **offset** | *integer* | The offset of the items returned (as set in the query or by default) | `0` |
| **policies** | *array* | List of policies | `["policyName1, policyName2"]` |
| **total** | *integer* | The total number of items available to return | `50` |

### Organization's policies List

List all policies by organization.

```
GET /api/v1/organizations/{organization_id}/policies?PathPrefix={optional_path_prefix}&Offset={optional_offset}&Limit={optional_limit}&OrderBy={columnName-desc}
```


#### Curl Example

```bash
$ curl -n /api/v1/organizations/$ORGANIZATION_ID/policies?PathPrefix=$OPTIONAL_PATH_PREFIX&Offset=$OPTIONAL_OFFSET&Limit=$OPTIONAL_LIMIT&OrderBy=$COLUMNNAME-DESC \
  -H "Authorization: Basic or Bearer XXX"
```


#### Response Example

```
HTTP/1.1 200 OK
```

```json
{
  "policies": [
    "policyName1, policyName2"
  ],
  "offset": 0,
  "limit": 20,
  "total": 50
}
```


## <a name="resource-order4_policyAllReference">All policies</a>




### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **limit** | *integer* | The maximum number of items in the response (as set in the query or by default) | `20` |
| **offset** | *integer* | The offset of the items returned (as set in the query or by default) | `0` |
| **[policies/name](#resource-order2_policy)** | *string* | Policy name | `"policy1"` |
| **[policies/org](#resource-order2_policy)** | *string* | Policy organization | `"tecsisa"` |
| **total** | *integer* | The total number of items available to return | `50` |

### All policies List

List all policies.

```
GET /api/v1/policies?PathPrefix={optional_path_prefix}&Offset={optional_offset}&Limit={optional_limit}&OrderBy={columnName-asc}
```


#### Curl Example

```bash
$ curl -n /api/v1/policies?PathPrefix=$OPTIONAL_PATH_PREFIX&Offset=$OPTIONAL_OFFSET&Limit=$OPTIONAL_LIMIT&OrderBy=$COLUMNNAME-ASC \
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
  ],
  "offset": 0,
  "limit": 20,
  "total": 50
}
```


## <a name="resource-order5_attachedGroups">Attached group</a>


List attached groups

### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **groups/attached** | *date-time* | When relationship was created | `"2015-01-01T12:00:00Z"` |
| **groups/group** | *string* | Group name | `"groupName1"` |
| **limit** | *integer* | The maximum number of items in the response (as set in the query or by default) | `20` |
| **offset** | *integer* | The offset of the items returned (as set in the query or by default) | `0` |
| **total** | *integer* | The total number of items available to return | `50` |

### Attached group List

List attached groups to this policy

```
GET /api/v1/organizations/{organization_id}/policies/{policy_name}/groups?Offset={optional_offset}&Limit={optional_limit}&OrderBy={columnName-desc}
```


#### Curl Example

```bash
$ curl -n /api/v1/organizations/$ORGANIZATION_ID/policies/$POLICY_NAME/groups?Offset=$OPTIONAL_OFFSET&Limit=$OPTIONAL_LIMIT&OrderBy=$COLUMNNAME-DESC \
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
      "group": "groupName1",
      "attached": "2015-01-01T12:00:00Z"
    }
  ],
  "offset": 0,
  "limit": 20,
  "total": 50
}
```


