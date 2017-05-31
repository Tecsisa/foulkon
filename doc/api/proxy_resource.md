## <a name="resource-order1_resource_entity">Resource</a>


Entity with the external resource information

### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **action** | *string* | Action related to this resource | `"example:get"` |
| **host** | *string* | Scheme + registered name (hostname) or IP address | `"https://httpbin.org"` |
| **method** | *string* | HTTP Method definition | `"GET"` |
| **path** | *string* | Relative path for destination host. | `"/example"` |
| **urn** | *string* | Uniform Resource Name for this resource | `"urn:examplews:application:v1:resource/get"` |


## <a name="resource-order2_proxy_resource">Proxy Resource</a>


Proxy Resource API

### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **createdAt** | *date-time* | Proxy resource creation date | `"2015-01-01T12:00:00Z"` |
| **id** | *uuid* | Unique proxy resource identifier | `"01234567-89ab-cdef-0123-456789abcdef"` |
| **name** | *string* | Proxy resource name | `"Example"` |
| **org** | *string* | Proxy resource organization | `"tecsisa"` |
| **path** | *string* | Proxy resource location | `"/example/admin/"` |
| **[resource:action](#resource-order1_resource_entity)** | *string* | Action related to this resource | `"example:get"` |
| **[resource:host](#resource-order1_resource_entity)** | *string* | Scheme + registered name (hostname) or IP address | `"https://httpbin.org"` |
| **[resource:method](#resource-order1_resource_entity)** | *string* | HTTP Method definition | `"GET"` |
| **[resource:path](#resource-order1_resource_entity)** | *string* | Relative path for destination host. | `"/example"` |
| **[resource:urn](#resource-order1_resource_entity)** | *string* | Uniform Resource Name for this resource | `"urn:examplews:application:v1:resource/get"` |
| **updateAt** | *date-time* | The date timestamp of the last update | `"2015-01-01T12:00:00Z"` |
| **urn** | *string* | Uniform Resource Name | `"urn:iws:iam:org:proxy/example/admin"` |

### Proxy Resource Create

Create a new proxy resource.

```
POST /api/v1/organizations/{organization_id}/proxy-resources
```

#### Required Parameters

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **name** | *string* | Proxy resource name | `"Example"` |
| **path** | *string* | Proxy resource location | `"/example/admin/"` |
| **resource:action** | *string* | Action related to this resource | `"example:get"` |
| **resource:host** | *string* | Scheme + registered name (hostname) or IP address | `"https://httpbin.org"` |
| **resource:method** | *string* | HTTP Method definition | `"GET"` |
| **resource:path** | *string* | Relative path for destination host. | `"/example"` |
| **resource:urn** | *string* | Uniform Resource Name for this resource | `"urn:examplews:application:v1:resource/get"` |



#### Curl Example

```bash
$ curl -n -X POST /api/v1/organizations/$ORGANIZATION_ID/proxy-resources \
  -d '{
  "name": "Example",
  "path": "/example/admin/",
  "resource": {
    "host": "https://httpbin.org",
    "path": "/example",
    "method": "GET",
    "urn": "urn:examplews:application:v1:resource/get",
    "action": "example:get"
  }
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
  "urn": "urn:iws:iam:org:proxy/example/admin",
  "org": "tecsisa",
  "resource": {
    "host": "https://httpbin.org",
    "path": "/example",
    "method": "GET",
    "urn": "urn:examplews:application:v1:resource/get",
    "action": "example:get"
  }
}
```

### Proxy Resource Update

Update an existing proxy resource.

```
PUT /api/v1/organizations/{organization_id}/proxy-resources/{proxy_resource_name}
```

#### Required Parameters

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **name** | *string* | Proxy resource name | `"Example"` |
| **path** | *string* | Proxy resource location | `"/example/admin/"` |
| **resource:action** | *string* | Action related to this resource | `"example:get"` |
| **resource:host** | *string* | Scheme + registered name (hostname) or IP address | `"https://httpbin.org"` |
| **resource:method** | *string* | HTTP Method definition | `"GET"` |
| **resource:path** | *string* | Relative path for destination host. | `"/example"` |
| **resource:urn** | *string* | Uniform Resource Name for this resource | `"urn:examplews:application:v1:resource/get"` |



#### Curl Example

```bash
$ curl -n -X PUT /api/v1/organizations/$ORGANIZATION_ID/proxy-resources/$PROXY_RESOURCE_NAME \
  -d '{
  "name": "Example",
  "path": "/example/admin/",
  "resource": {
    "host": "https://httpbin.org",
    "path": "/example",
    "method": "GET",
    "urn": "urn:examplews:application:v1:resource/get",
    "action": "example:get"
  }
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
  "urn": "urn:iws:iam:org:proxy/example/admin",
  "org": "tecsisa",
  "resource": {
    "host": "https://httpbin.org",
    "path": "/example",
    "method": "GET",
    "urn": "urn:examplews:application:v1:resource/get",
    "action": "example:get"
  }
}
```

### Proxy Resource Delete

Delete an existing proxy resource.

```
DELETE /api/v1/organizations/{organization_id}/proxy-resources/{proxy_resource_name}
```


#### Curl Example

```bash
$ curl -n -X DELETE /api/v1/organizations/$ORGANIZATION_ID/proxy-resources/$PROXY_RESOURCE_NAME \
  -H "Content-Type: application/json" \
  -H "Authorization: Basic or Bearer XXX"
```


#### Response Example

```
HTTP/1.1 202 Accepted
```


### Proxy Resource Get

Get an existing proxy resource.

```
GET /api/v1/organizations/{organization_id}/proxy-resources/{proxy_resource_name}
```


#### Curl Example

```bash
$ curl -n /api/v1/organizations/$ORGANIZATION_ID/proxy-resources/$PROXY_RESOURCE_NAME \
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
  "urn": "urn:iws:iam:org:proxy/example/admin",
  "org": "tecsisa",
  "resource": {
    "host": "https://httpbin.org",
    "path": "/example",
    "method": "GET",
    "urn": "urn:examplews:application:v1:resource/get",
    "action": "example:get"
  }
}
```


## <a name="resource-order3_ProxyResourceReference">Organization's proxy resources</a>




### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **limit** | *integer* | The maximum number of items in the response (as set in the query or by default) | `20` |
| **offset** | *integer* | The offset of the items returned (as set in the query or by default) | `0` |
| **resources** | *array* | List of proxy resources | `["ProxyResourceName1, ProxyResourceName2"]` |
| **total** | *integer* | The total number of items available to return | `2` |

### Organization's proxy resources List

List all proxy resources by organization.

```
GET /api/v1/organizations/{organization_id}/proxy-resources?PathPrefix={optional_path_prefix}&Offset={optional_offset}&Limit={optional_limit}&OrderBy={columnName-desc}
```


#### Curl Example

```bash
$ curl -n /api/v1/organizations/$ORGANIZATION_ID/proxy-resources?PathPrefix=$OPTIONAL_PATH_PREFIX&Offset=$OPTIONAL_OFFSET&Limit=$OPTIONAL_LIMIT&OrderBy=$COLUMNNAME-DESC \
  -H "Authorization: Basic or Bearer XXX"
```


#### Response Example

```
HTTP/1.1 200 OK
```

```json
{
  "resources": [
    "ProxyResourceName1, ProxyResourceName2"
  ],
  "offset": 0,
  "limit": 20,
  "total": 2
}
```


