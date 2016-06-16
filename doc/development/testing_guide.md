## How to test API Method ##

TODO: INGLES
### TODO ###

Se utilizan test case con mapas que tiene que cumplir este criterio:

- Los casos de error empezaran por: ErrorCaseDescription (ErrorCaseUserNotFound)
- Los casos que no sean de error empezaran por: OkCaseDescription (OkCaseAdminRequest)
- El orden de los test será: OkCase, ErrorCase según el orden de ejecución en el método.

### Estructura TestCase ###

```go
    // API Method args
    authUser  AuthenticatedUser
    userID    string
    org       string
    groupName string
    // Expected result
    wantError     *Error
    expectedGroup *Group
    // Manager Results
    getGroupsByUserIDResult   []Group
    getPoliciesAttachedResult []Policy
    getUserByExternalIDResult *User
    getGroupByNameResult      *Group
    isMemberOfGroupResult     bool
    // Manager Errors
    getUserByExternalIDMethodErr error
    getGroupByNameMethodErr      error
    addMemberMethodErr           error
    isMemberOfGroupMethodErr     error
```

Debería seguir la asignación el mismo orden que la declaración

Un buen ejemplo de test [aquí](../../api/group_test.go)
