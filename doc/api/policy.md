## <a name="resource-order1_statement">Statement</a>


Policy statement

### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **action** | *array* | CRUD functions | `["iam:*"]` |
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
| **org** | *string* | Policy's organization | `"tecsisa"` |
| **path** | *string* | Policy's location | `"/example/admin/"` |
| **statements** | *array* | Policy statements | `[{"effect":"allow","action":["iam:*"],"resources":["urn:everything:*"]}]` |
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
| **path** | *string* | Policy's location | `"/example/admin/"` |
| **statements** | *array* | Policy statements | `[{"effect":"allow","action":["iam:*"],"resources":["urn:everything:*"]}]` |



#### Curl Example

```bash
$ curl -n -X POST /api/v1/organizations/$ORGANIZATION_ID/policies \
  -d '{
  "name": "policy1",
  "path": "/example/admin/",
  "statements": [
    {
      "effect": "allow",
      "action": [
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
  "urn": "urn:iws:iam:org1:policy/example/admin/policy1",
  "org": "tecsisa",
  "statements": [
    {
      "effect": "allow",
      "action": [
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
| **path** | *string* | Policy's location | `"/example/admin/"` |
| **statements** | *array* | Policy statements | `[{"effect":"allow","action":["iam:*"],"resources":["urn:everything:*"]}]` |



#### Curl Example

```bash
$ curl -n -X PUT /api/v1/organizations/$ORGANIZATION_ID/policies/$POLICY_NAME \
  -d '{
  "name": "policy1",
  "path": "/example/admin/",
  "statements": [
    {
      "effect": "allow",
      "action": [
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
  "urn": "urn:iws:iam:org1:policy/example/admin/policy1",
  "org": "tecsisa",
  "statements": [
    {
      "effect": "allow",
      "action": [
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
  "urn": "urn:iws:iam:org1:policy/example/admin/policy1",
  "org": "tecsisa",
  "statements": [
    {
      "effect": "allow",
      "action": [
        "iam:*"
      ],
      "resources": [
        "urn:everything:*"
      ]
    }
  ]
}
```


## <a name="resource-order3_policyReference"></a>




### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **[name](#resource-order2_policy)** | *string* | Policy name | `"policy1"` |
| **[org](#resource-order2_policy)** | *string* | Policy's organization | `"tecsisa"` |

###  Policy List

List all policies by organization.

```
GET /api/v1/organizations/{organization_id}/policies
```


#### Curl Example

```bash
$ curl -n /api/v1/organizations/$ORGANIZATION_ID/policies \
  -H "Authorization: Basic or Bearer XXX"
```


#### Response Example

```
HTTP/1.1 200 OK
```

```json
[
  {
    "org": "tecsisa",
    "name": "policy1"
  }
]
```


## <a name="resource-order4_attachedGroups"></a>


List attached groups

### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **name** | *string* | Group's name | `"group1"` |
| **org** | *string* | Group's organization | `"tecsisa"` |

###  Policy Groups List

List attached groups

```
GET /api/v1/organizations/{organization_id}/policies/{policy_name}/groups
```


#### Curl Example

```bash
$ curl -n /api/v1/organizations/$ORGANIZATION_ID/policies/$POLICY_NAME/groups \
  -H "Authorization: Basic or Bearer XXX"
```


#### Response Example

```
HTTP/1.1 200 OK
```

```json
[
  {
    "org": "tecsisa",
    "name": "group1"
  }
]
```


