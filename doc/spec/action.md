## IAM actions

### User

|          Method          |        Action         | Dependencies |
|--------------------------|-----------------------|--------------|
| **Create user**          | iam:CreateUser        | None         |
| **Delete user**          | iam:DeleteUser        | iam:GetUser  |
| **Get user**             | iam:GetUser           | None         |
| **List users**           | iam:ListUsers         | None         |
| **Update user**          | iam:UpdateUser        | iam:GetUser  |
| **List groups for user** | iam:ListGroupsForUser | iam:GetUser  |


### Group

|              Method              |            Action             |        Dependencies         |
|----------------------------------|-------------------------------|-----------------------------|
| **Create group**                 | iam:CreateGroup               | None                        |
| **Delete group**                 | iam:DeleteGroup               | iam:GetGroup                |
| **Get group**                    | iam:GetGroup                  | None                        |
| **List groups**                  | iam:ListGroups                | None                        |
| **Update group**                 | iam:UpdateGroup               | iam:GetGroup                |
| **List members**                 | iam:ListMembers               | iam:GetGroup                |
| **Add member**                   | iam:AddMember                 | iam:GetGroup, iam:GetUser   |
| **Remove member**                | iam:RemoveMember              | iam:GetGroup, iam:GetUser   |
| **Attach group policy**          | iam:AttachGroupPolicy         | iam:GetGroup, iam:GetPolicy |
| **Detach group policy**          | iam:DetachGroupPolicy         | iam:GetGroup, iam:GetPolicy |
| **List attached group policies** | iam:ListAttachedGroupPolicies | iam:GetGroup                |

### Policy

|          Method          |         Action         | Dependencies  |
|--------------------------|------------------------|---------------|
| **Create policy**        | iam:CreatePolicy       | None          |
| **Delete policy**        | iam:DeletePolicy       | iam:GetPolicy |
| **Get policy**           | iam:GetPolicy          | None          |
| **Update policy**        | iam:UpdatePolicy       | iam:GetPolicy |
| **List policies**        | iam:ListPolicies       | None          |
| **List attached groups** | iam:ListAttachedGroups | iam:GetPolicy |

## Proxy Resources

|          Method          |         Action             | Dependencies         |
|--------------------------|----------------------------|----------------------|
| **Create Proxy Resource**| iam:CreateProxyResource    | None                 |
| **Delete Proxy Resource**| iam:DeleteProxyResource    | iam:GetProxyResource |
| **Get Proxy Resource**   | iam:GetProxyResource       | None                 |
| **Update Proxy Resource**| iam:UpdateProxyResource    | iam:GetProxyResource |
| **List Proxy Resources** | iam:ListProxyResources     | None                 |

## OIDC Provider

|          Method          |         Action         | Dependencies         |
|--------------------------|------------------------|----------------------|
| **Create OIDC Provider** | auth:CreateOidcProvider| None                 |
| **Delete OIDC Provider** | auth:DeleteOidcProvider| auth:GetOidcProvider |
| **Get OIDC Provider**    | auth:GetOidcProvider   | None                 |
| **Update OIDC Providers**| auth:UpdateOidcProvider| auth:GetOidcProvider |
| **List OIDC Provider**   | auth:ListOidcProviders | None                 |


### Additional info

The dependencies are directly related to the action, for example in AddMember we need permissions to get the group (iam:GetGroup) and the user (iam:GetUser). 
So as described in the table, to use iam:AddMember you need iam:GetGroup + iam:GetUser.

```
Example:
- Add Member (user1) to (group1)
- Dependencies are: iam:GetGroup (group1), iam:GetUser (user1)
```
