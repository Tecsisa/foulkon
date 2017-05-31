#!/usr/bin/env bash
prmd doc group.json > ../doc/api/group.md
prmd doc user.json > ../doc/api/user.md
prmd doc policy.json > ../doc/api/policy.md
prmd doc proxy_resource.json > ../doc/api/proxy_resource.md
prmd doc resource.json > ../doc/api/resource.md
prmd doc oidc_provider.json > ../doc/api/oidc_provider.md