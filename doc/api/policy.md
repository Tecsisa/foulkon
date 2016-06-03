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
| **CreatedAt** | *date-time* | when policy was created | `"2015-01-01T12:00:00Z"` |
| **ID** | *uuid* | Unique identifier of policy | `"01234567-89ab-cdef-0123-456789abcdef"` |
| **Name** | *string* | Name of policy | `"policy1"` |
| **Org** | *string* | Organization of policy | `"tecsisa"` |
| **Path** | *string* | Policy's location | `"/example/admin/"` |
| **Statements** | *array* | Policy statements | `[{"effect":"allow","action":["iam:*"],"resources":["urn:everything:*"]}]` |
| **Urn** | *string* | Uniform Resource Name of policy | `"urn:iws:iam:org1:policy/example/admin/policy1"` |

### Policy Create

Create a new policy.

```
POST /api/v1/organizations/{organization_id}/policies
```

#### Required Parameters

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **Name** | *string* | Name of policy | `"policy1"` |
| **Path** | *string* | Policy's location | `"/example/admin/"` |
| **Statements** | *array* | Policy statements | `[{"effect":"allow","action":["iam:*"],"resources":["urn:everything:*"]}]` |



#### Curl Example

```bash
$ curl -n -X POST /api/v1/organizations/$ORGANIZATION_ID/policies \
  -d '{
  "Name": "policy1",
  "Path": "/example/admin/",
  "Statements": [
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
  "ID": "01234567-89ab-cdef-0123-456789abcdef",
  "Name": "policy1",
  "Path": "/example/admin/",
  "CreatedAt": "2015-01-01T12:00:00Z",
  "Urn": "urn:iws:iam:org1:policy/example/admin/policy1",
  "Org": "tecsisa",
  "Statements": [
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
| **Name** | *string* | Name of policy | `"policy1"` |
| **Path** | *string* | Policy's location | `"/example/admin/"` |
| **Statements** | *array* | Policy statements | `[{"effect":"allow","action":["iam:*"],"resources":["urn:everything:*"]}]` |



#### Curl Example

```bash
$ curl -n -X PUT /api/v1/organizations/$ORGANIZATION_ID/policies/$POLICY_NAME \
  -d '{
  "Name": "policy1",
  "Path": "/example/admin/",
  "Statements": [
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
  "ID": "01234567-89ab-cdef-0123-456789abcdef",
  "Name": "policy1",
  "Path": "/example/admin/",
  "CreatedAt": "2015-01-01T12:00:00Z",
  "Urn": "urn:iws:iam:org1:policy/example/admin/policy1",
  "Org": "tecsisa",
  "Statements": [
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
  "ID": "01234567-89ab-cdef-0123-456789abcdef",
  "Name": "policy1",
  "Path": "/example/admin/",
  "CreatedAt": "2015-01-01T12:00:00Z",
  "Urn": "urn:iws:iam:org1:policy/example/admin/policy1",
  "Org": "tecsisa",
  "Statements": [
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
| **[Name](#resource-order2_policy)** | *string* | Name of policy | `"policy1"` |
| **[Org](#resource-order2_policy)** | *string* | Organization of policy | `"tecsisa"` |

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
    "Org": "tecsisa",
    "Name": "policy1"
  }
]
```


## <a name="resource-order4_attachedGroups"></a>


List attached groups

### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **Name** | *string* | Name of group | `"group1"` |
| **Org** | *string* | Organization of group | `"tecsisa"` |

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
    "Org": "tecsisa",
    "Name": "group1"
  }
]
```


