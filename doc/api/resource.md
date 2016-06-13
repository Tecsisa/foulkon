## <a name="resource-authorize">Authorize</a>


Authorize API

### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **ResourcesAllowed** | *array* | List of resources allowed | `["urn:ews:product:instance:example/resource1"]` |

### Authorize resources

Authorize user to access resources

```
POST /api/v1/authorize
```

#### Required Parameters

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **Action** | *string* | Action applied over the resources | `"example:Read"` |
| **Resources** | *array* | List of resources | `["urn:ews:product:instance:example/resource1"]` |



#### Curl Example

```bash
$ curl -n -X POST /api/v1/authorize \
  -d '{
  "Action": "example:Read",
  "Resources": [
    "urn:ews:product:instance:example/resource1"
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
  "ResourcesAllowed": [
    "urn:ews:product:instance:example/resource1"
  ]
}
```


