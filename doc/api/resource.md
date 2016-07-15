## <a name="resource-authorize">Resource</a>


Resource API

### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **resourcesAllowed** | *array* | List of allowed resources | `["urn:ews:product:instance:example/resource1"]` |

### Resource authorized

Get authorized resources according selected action and resources

```
POST /api/v1/authorize
```

#### Required Parameters

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **action** | *string* | Action applied over the resources | `"example:Read"` |
| **resources** | *array* | List of resources | `["urn:ews:product:instance:example/resource1"]` |



#### Curl Example

```bash
$ curl -n -X POST /api/v1/authorize \
  -d '{
  "action": "example:Read",
  "resources": [
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
  "resourcesAllowed": [
    "urn:ews:product:instance:example/resource1"
  ]
}
```


