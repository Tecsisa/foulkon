## <a name="resource-order1_resource_entity">OIDC Client</a>


Entity with the OIDC Client configuration to use in Authentication Middleware

### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **[name](#resource-order1_oidc_client)** | *string* | Identifier associated to this OIDC Client for the OIDC Provider | `"client-api-identifier"` |


## <a name="resource-order2_oidc_provider">OIDC Provider</a>


Entity with the OIDC Provider configuration to use in Authentication Middleware

### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **clients** | *array* | OIDC Clients associated | `[{"name":"client-api-identifier"}]` |
| **createdAt** | *date-time* | OIDC Provider creation date | `"2015-01-01T12:00:00Z"` |
| **id** | *uuid* | Unique OIDC Provider identifier | `"01234567-89ab-cdef-0123-456789abcdef"` |
| **issuerUrl** | *string* | The issuer URL which issues the tokens | `"https://accounts.google.com"` |
| **name** | *string* | OIDC Provider name | `"Example"` |
| **path** | *string* | OIDC Provider location | `"/example/admin/"` |
| **updateAt** | *date-time* | The date timestamp of the last update | `"2015-01-01T12:00:00Z"` |
| **urn** | *string* | Uniform Resource Name | `"urn:iws:auth::oidc/example/admin/Example"` |

### OIDC Provider Create

Create a new OIDC Provider.

```
POST /api/v1/admin/auth/oidc/providers
```

#### Required Parameters

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **clients** | *array* | OIDC Client identifiers associated | `["client-api-identifier"]` |
| **issuerUrl** | *string* | The issuer URL which issues the tokens | `"https://accounts.google.com"` |
| **name** | *string* | OIDC Provider name | `"Example"` |
| **path** | *string* | OIDC Provider location | `"/example/admin/"` |



#### Curl Example

```bash
$ curl -n -X POST /api/v1/admin/auth/oidc/providers \
  -d '{
  "name": "Example",
  "path": "/example/admin/",
  "issuerUrl": "https://accounts.google.com",
  "clients": [
    "client-api-identifier"
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
  "name": "Example",
  "path": "/example/admin/",
  "createdAt": "2015-01-01T12:00:00Z",
  "updateAt": "2015-01-01T12:00:00Z",
  "issuerUrl": "https://accounts.google.com",
  "urn": "urn:iws:auth::oidc/example/admin/Example",
  "clients": [
    {
      "name": "client-api-identifier"
    }
  ]
}
```

### OIDC Provider Update

Update an existing OIDC Provider.

```
PUT /api/v1/admin/auth/oidc/providers/{oidc_provider_name}
```

#### Required Parameters

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **clients** | *array* | OIDC Client identifiers associated | `["client-api-identifier"]` |
| **issuerUrl** | *string* | The issuer URL which issues the tokens | `"https://accounts.google.com"` |
| **name** | *string* | OIDC Provider name | `"Example"` |
| **path** | *string* | OIDC Provider location | `"/example/admin/"` |



#### Curl Example

```bash
$ curl -n -X PUT /api/v1/admin/auth/oidc/providers/$OIDC_PROVIDER_NAME \
  -d '{
  "name": "Example",
  "path": "/example/admin/",
  "issuerUrl": "https://accounts.google.com",
  "clients": [
    "client-api-identifier"
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
  "name": "Example",
  "path": "/example/admin/",
  "createdAt": "2015-01-01T12:00:00Z",
  "updateAt": "2015-01-01T12:00:00Z",
  "issuerUrl": "https://accounts.google.com",
  "urn": "urn:iws:auth::oidc/example/admin/Example",
  "clients": [
    {
      "name": "client-api-identifier"
    }
  ]
}
```

### OIDC Provider Delete

Delete an existing OIDC Provider.

```
DELETE /api/v1/admin/auth/oidc/providers/{oidc_provider_name}
```


#### Curl Example

```bash
$ curl -n -X DELETE /api/v1/admin/auth/oidc/providers/$OIDC_PROVIDER_NAME \
  -H "Content-Type: application/json" \
  -H "Authorization: Basic or Bearer XXX"
```


#### Response Example

```
HTTP/1.1 202 Accepted
```


### OIDC Provider Get

Get an existing OIDC Provider.

```
GET /api/v1/admin/auth/oidc/providers/{oidc_provider_name}
```


#### Curl Example

```bash
$ curl -n /api/v1/admin/auth/oidc/providers/$OIDC_PROVIDER_NAME \
  -H "Authorization: Basic or Bearer XXX"
```


#### Response Example

```
HTTP/1.1 200 OK
```

```json
{
  "id": "01234567-89ab-cdef-0123-456789abcdef",
  "name": "Example",
  "path": "/example/admin/",
  "createdAt": "2015-01-01T12:00:00Z",
  "updateAt": "2015-01-01T12:00:00Z",
  "issuerUrl": "https://accounts.google.com",
  "urn": "urn:iws:auth::oidc/example/admin/Example",
  "clients": [
    {
      "name": "client-api-identifier"
    }
  ]
}
```


## <a name="resource-order3_OidcProviderReference"></a>




### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **limit** | *integer* | The maximum number of items in the response (as set in the query or by default) | `20` |
| **offset** | *integer* | The offset of the items returned (as set in the query or by default) | `0` |
| **providers** | *array* | OIDC Provider identifiers | `["google","keycloak"]` |
| **total** | *integer* | The total number of items available to return | `2` |

###  OIDC Provider List All

List all OIDC Providers, using optional query parameters.

```
GET /api/v1/admin/auth/oidc/providers?PathPrefix={optional_path_prefix}&Offset={optional_offset}&Limit={optional_limit}&OrderBy={columnName-desc}
```


#### Curl Example

```bash
$ curl -n /api/v1/admin/auth/oidc/providers?PathPrefix=$OPTIONAL_PATH_PREFIX&Offset=$OPTIONAL_OFFSET&Limit=$OPTIONAL_LIMIT&OrderBy=$COLUMNNAME-DESC \
  -H "Authorization: Basic or Bearer XXX"
```


#### Response Example

```
HTTP/1.1 200 OK
```

```json
{
  "providers": [
    "google",
    "keycloak"
  ],
  "offset": 0,
  "limit": 20,
  "total": 2
}
```


