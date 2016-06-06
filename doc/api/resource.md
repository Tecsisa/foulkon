## <a name="resource-order1_restriction">Restriction</a>




### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **AllowedFullUrns** | *array* | Locations where urn's are allowed | `["urn:ews:product:instance1:example/resource_path"]` |
| **AllowedUrnPrefixes** | *array* | Locations where prefixes are allowed | `["urn:ews:product:instance2:*"]` |
| **DeniedFullUrns** | *array* | Locations where urn's are denied | `["urn:ews:product:instance2:example2/resource_path"]` |
| **DeniedUrnPrefixes** | *array* | Locations where prefixes are denied | `["urn:ews:product2:*"]` |


## <a name="resource-order2_resource">Resource</a>


Resource API

### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **Effect** | *string* | allow/deny resources | `"allow/deny"` |
| **Restrictions:AllowedFullUrns** | *array* | Locations where urn's are allowed | `["urn:ews:product:instance1:example/resource_path"]` |
| **Restrictions:AllowedUrnPrefixes** | *array* | Locations where prefixes are allowed | `["urn:ews:product:instance2:*"]` |
| **Restrictions:DeniedFullUrns** | *array* | Locations where urn's are denied | `["urn:ews:product:instance2:example2/resource_path"]` |
| **Restrictions:DeniedUrnPrefixes** | *array* | Locations where prefixes are denied | `["urn:ews:product2:*"]` |

### Resource Get Effect

Get user effect to do the action over the resource. If urn is full only return effect else if is a prefix return restrictions

```
GET /api/v1/resources?Action={Action_example}&Urn={Urn_example}
```


#### Curl Example

```bash
$ curl -n /api/v1/resources?Action=$ACTION_EXAMPLE&Urn=$URN_EXAMPLE \
  -H "Authorization: Basic or Bearer XXX"
```


#### Response Example

```
HTTP/1.1 200 OK
```

```json
{
  "Effect": "allow/deny",
  "Restrictions": {
    "AllowedUrnPrefixes": [
      "urn:ews:product:instance2:*"
    ],
    "AllowedFullUrns": [
      "urn:ews:product:instance1:example/resource_path"
    ],
    "DeniedUrnPrefixes": [
      "urn:ews:product2:*"
    ],
    "DeniedFullUrns": [
      "urn:ews:product:instance2:example2/resource_path"
    ]
  }
}
```


