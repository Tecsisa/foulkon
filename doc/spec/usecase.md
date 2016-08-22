# Use case

We have an application (app1) composed by UI, API, and some backend:
```
+--------+   +---------+   +-------------+
|        |   |         |   |             |
|   UI   +--->   API   +--->   Backend   |
|        |   |         |   |             |
+--------+   +---------+   +-------------+
```

The application manages data from several __companies__.
These companies, have __roles or groups__ that contain __users__.

In order to simplify this use case, we are going to define just __two roles: admin and member__.
<br /><br />
##### User definition in Foulkon
```
urn:iws:iam::user/user_admin1
urn:iws:iam::user/user_member1
```

##### Group definition in Foulkon
```
urn:iws:iam:company1:group/app1/member
urn:iws:iam:company1:group/app1/admin
```

Then we should attach the users to the groups. See [group API doc](../api/group.md#add_member)

Note: That being said, keep in mind that a user can be member of several groups (and organizations).

<br /><br />
The API exposes __two resources: AdminResource, and Resource__.

We'd define the resources like this:
```
urn:examplews:app1:v1:resource/company1/AdminResource
urn:examplews:app1:v1:resource/company1/Resource
```
<br />

Now we have to grant admin users access to everything related to `company1`.
To do this, you have to create a policy in the first place:
```json
{
   "Name": "admin_access",
   "Path": "/app1/policies/company1/",
   "Statements": [
       {
           "Effect": "allow",
           "Action": ["*"],
           "Resources": ["urn:examplews:app1:v1:resource/company1/*"]
       }
   ]
}
```
And another policy for members to have access to Resource:
```json
{
   "Name": "member_access_to_Resource",
   "Path": "/app1/policies/company1/",
   "Statements": [
       {
           "Effect": "allow",
           "Action": ["app1:read", "app1:create"],
           "Resources": ["urn:examplews:app1:v1:resource/company1/Resource"]
       }
   ]
}
```

Note: You have to define the possible resource actions in your app (e.g.: read, create, delete)
<br /><br />

With these definitions, if you ask Foulkon ([calling /api/v1/authorize](../api/resource.md#resource_authorized)):
- if user user_member1 asks access to action:"read" over `urn:examplews:app1:v1:resource/company1/AdminResource`, it will respond 403 forbidden.
- if user user_member1 asks access to action:"read" over `urn:examplews:app1:v1:resource/company1/Resource`, it will respond 200 OK.

## Proxy
To avoid implementing permissions logic in the APIs, we created a Foulkon Proxy to manage the API resources access.

```
             +---+
             |   |
+--------+   | P |   +---------+   +-------------+
|        |   | r |   |         |   |             |
|   UI   +---> o +--->   API   +--->   Backend   |
|        |   | x |   |         |   |             |
+--------+   | y |   +---------+   +-------------+
             |   |
             +-+-+
               |
         +-----v-----+
         |           |
         |  Foulkon  |
         |           |
         +-----------+

```

You can configure this using the proxy toml config file:
```
[[resources]]
    id = "resource1"
    host = "https://app1.company1.com/"
    url = "/resource"
    method = "GET"
    urn = "urn:examplews:app1:v1:resource/company1/Resource"
    action = "app1:read"
```

