# Specification

The authorization domain is composed of a set of elements that are depicted in the following diagram:

![Image of IAM](https://docs.google.com/drawings/d/1h82ER9BSRMD_cTSeYOjLSNbeAJqFOE4mjnINhQuhbz0/pub?w=960&h=720)

- Character “*” means one or more elements.
- Character “1” means EXACTLY one.

## <a name="resource"></a>Resource
Resource is the element that needs to be authorized/denied. It has an unique identifier. 
For now, there are two resource types:
- Internal resources (IAM) for self-management
- External resources

### IAM resources
IAM resources allow you to manage groups, organizations, users and policies.
This is a representation of a generic resource with its elements:

```
urn:iws:iam:org:genericresource/pathname
```

- urn: uniform resource name.
- iws: internal web service.
- iam: identity access management.
- org: organization, not apply to IAM users  (google, facebook, coreos, tecsisa, etc.)
- genericresource/pathname: type and unique name for this resource.

In this system we have some representations of users, groups and policies as resources.

- __IAM user__: `urn:iws:iam::user/pathnameuser`
- __IAM group__: `urn:iws:iam:org:group/pathnamegroup`
- __IAM policy__: `urn:iws:iam:org:policy/pathnamepolicy`

Google user account resource example:
```
urn:iws:iam::user/gapps/gmail/user123456
```
Google group resource example:
```
urn:iws:iam:google:group/gapps/gmail/dev-mail
```

### External resources
IAM urns are reserved for AuthZ self-management, so, in order to prevent conflicts, external resources must have different names. This is the representation of an external resource:

```
urn:ews:product:instance:resource/resourcepath
```

- urn: uniform resource name
- ews: external web service (googlews, facebookws, etc).
- product: product name.
- instance: instance of your product.
- resource/resourcepath: unique pathname for your resource.

Facebook API resource to access to user profile 123456:

```
urn:facebookws:socialnet:v123456:resource/user/profile/123456
```

## IAM Entities

### User
User is the basic element to represent a Security Principal that might have access to some resources.
Users could be members of one or more groups. A user might join any group regardless the organization that group belongs.
Go to [User API](../api/user.md) for more information about this entity.

### Organization
Organization is a container of groups and policies. It doesn't have a representation, you don't need to create organizations.

### Group
Group is a collection of users, which belongs to ONLY ONE organization.
According to this draft, a user is granted access to resources by attaching policies to the groups he belongs to.
Group names are unique inside the same organization.
Go to [Group API](../api/group.md) for more information about this entity.

### Policy
A policy is a specification of permissions defined in terms of statements that declare what actions are allowed or denied to be performed on resources.
These policies might be attached to groups in order to restrict their application scope. Policies can’t be attached to users.
Policy names are unique inside the same organization.
Go to [Policy API](../api/policy.md) for more information about this entity.

## Permission definition

The way to define your permissions is using statements inside policies. 
A statement is composed of its `effect`(allow or deny), the `resources` list, and the `actions` you want to allow or deny.
 
Prefixes are allowed in resources and actions. Therefore wildcards (*) in the middle of the string are not allowed, only at the end
E.g:

```
- OK 	→ urn:facebookws:socialnet:v123456:*
- WRONG	→ urn:facebookws:*:socialnet:v123456:someUser
```

#### Default behaviour
When there are some policies that apply to same action and resource for a user, system select effect in this way:

- __If there is an explicit deny, system returns a deny.__
- __If there is an allow and no explicit deny, system returns an allow.__
- __If there isn’t a policy for that resource and action, system returns a deny by default.__

### IAM Policies
IAM policies define system permissions for its internal resources. Each resource type has its own actions predefined by prefix “iam”. This actions are defined in [Action doc](action.md) with its dependencies. When you start the system at first time, you have a system admin user with a password. This user doesn’t have limitations and can’t be assigned to a group.
__Best practice__: don’t use this admin account to manage your system. Create an user with admin rights and use it. Therefore a policy to manage all your IAM system could be:

```json
{
    "id": "01234567-89ab-cdef-0123-456789abcdef",
    "name": "policyAdmin",
    "path": "/admin/",
    "createdAt": "2015-01-01T12:00:00Z",
    "urn": "urn:iws:iam:orgName:policy/admin/policyAdmin",
    "org": "orgName",
    "statements": [
      {
        "effect": "allow",
        "action": [
          "iam:*"
        ],
        "resources": [
          "urn:iws:iam:*"
        ]
      }
    ]
}
```

You could create a group named AdminIAMgroup, create your user and assign it to this group. Now, you can manage all IAM system with your account as system admin.

### External Resource Policies
Policies for external resources depends on your application logic, which defines what actions are associated to these resources. Therefore your system has to define these things and asks to AuthZ if that user has access or not. Your resources should be defined following the [Resource](#resource) specification. 

An example of how Google could manage its application resources with a policy could be next:

```json
{
    "id": "01234567-89ab-cdef-0123-456789abcdef",
    "name": "GmailReadEmailUser123456",
    "path": "/gmail/",
    "createdAt": "2015-01-01T12:00:00Z",
    "urn": "urn:iws:iam:orgName:policy/gmail/GmailReadEmailUser123456",
    "org": "orgName",
    "statements": [
      {
        "effect": "allow",
        "action": [
          "gmail:*"
        ],
        "resources": [
          "urn:googlews:gmail:v123456:resource/user123456"
        ]
      }
    ]
}
```

This policy allows to do whatever action in user123456 account of product gmail for web services of google in instance v123456.

## Admin user
Admin user uses Basic Authentication scheme and it could be set at server start. This user doesn’t follow [Authorization Flow](authorization.md) and can’t be added to a group.

