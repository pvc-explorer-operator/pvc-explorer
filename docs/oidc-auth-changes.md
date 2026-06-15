# OIDC Authentication & RBAC Changes

## Overview

Added OIDC-based single sign-on with group-to-role mapping, a user profile page, a local Dex-based OIDC provider for kind cluster testing, role-based UI restrictions, and login audit logging.

## Architecture

```
Browser ──OIDC──▶ Dex (kind) ──ID Token──▶ Controller ──groups──▶ RBAC ──role──▶ Session
```

- **Generic OIDC** — works with any standards-compliant provider (Dex, Keycloak, Okta, etc.)
- **No Bearer token auth** — browser-based OIDC login only
- **Session store** — cookie-based, extended with OIDC claims (email, groups, subject)
- **ArgoCD-style RBAC** — `policy.csv` with `g, <group>, <role>` mappings

## Roles

| Role | Rank | Permissions |
|------|------|-------------|
| `admin` | 3 | Full access: create/update/delete explorers, manage scopes, edit files, settings |
| `user` | 2 | Create explorers, connect/disconnect, browse files (read-only) |
| `viewer` | 1 | Read-only: view explorers, scopes, browse/download files |

## What Changed

### Backend

#### `internal/auth/oidc.go`
- OIDC types (`OIDCConfig`, `RBACConfig`, `OIDCProvider`, `OIDCClaims`)
- Config loading from ConfigMaps (`pvc-explorer-config`, `pvc-explorer-rbac`)
- Token validation with `go-oidc/v3`
- Group-to-role mapping via `parsePolicyCSV` (strips quotes from ArgoCD-style entries)
- `oidc.externalIssuer` — controller uses in-cluster DNS for discovery, `externalIssuer` for browser redirects
- `oidc.skipTLSVerify` — custom `http.Client` with `InsecureSkipVerify` for self-signed certs

#### `internal/auth/session.go`
- Extended `sessionEntry` with `Email`, `Groups`, `Subject` fields
- `Get()` returns full `sessionEntry` (not just username/role)

#### `internal/auth/auth.go`
- Added `RoleUser` constant alongside `RoleAdmin` and `RoleViewer`

#### `internal/api/auth.go`
- Role hierarchy: `roleRank` map (`admin:3`, `user:2`, `viewer:1`)
- `isAllowed()` compares role ranks against route permissions
- Route permissions:
  - `POST /api/v1/explorers/` — requires `user`
  - `PUT/DELETE /api/v1/explorers/` — requires `admin`
  - `POST/PUT/DELETE /api/v1/scopes` — requires `admin`
  - Non-GET `/proxy/api/` — requires `admin`

#### `internal/api/handlers.go`
- `handleLogin` — logs method, username, role, success/failure
- `handleOIDCCallback` — logs method, username, email, subject, groups, role, success/failure
- OIDC start/callback endpoints
- Enriched `loginResponse` with `email`, `groups`, `subject`
- `handleMe` returns full user profile

#### `cmd/main.go`
- OIDC provider initialization from ConfigMap at startup
- Retry logic: 5 attempts with 3s backoff (handles Dex startup race)
- Uses `ctrl.SetupSignalHandler()` for context

#### `go.mod`
- Added `github.com/coreos/go-oidc/v3 v3.18.0`
- Promoted `golang.org/x/oauth2` to direct dependency

### Frontend

#### `ui/src/stores/authStore.ts`
- `oidcEnabled` flag
- `loginWithOIDC()` — redirects to OIDC start endpoint
- `UserProfile` interface with `email`, `groups`, `subject`
- `isAdmin` computed property

#### `ui/src/views/ProfileView.vue` (new)
- User profile page with avatar, name, role badge, email, groups (Tag badges), subject
- Accessible via clickable username in sidebar

#### `ui/src/views/AuthCallbackView.vue` (new)
- Simple redirect component for OIDC callback

#### `ui/src/views/FileBrowserView.vue`
- `effectiveReadonly` combines agent PVC state (`config.readonly`) with user role (`!authStore.isAdmin`)
- Non-admin users see file write buttons disabled + "read-only" badge

#### `ui/src/views/LoginView.vue`
- SSO button visible when `oidcEnabled`

#### `ui/src/router/index.ts`
- Added `/auth/callback` and `/profile` routes

#### `ui/src/layout/AppSidebar.vue`
- Username and avatar are `<router-link to="/profile">`

#### `ui/src/layout/AppTopbar.vue`
- Added "Profile" to breadcrumb navigation

#### `ui/src/layout/layout.css`
- Added `text-decoration: none` and hover effect for sidebar footer links

#### Already role-gated (no changes needed)
- `CreateAgentView.vue` — `v-if="authStore.isAdmin"` with "Admin access required" fallback
- `CreateScopeView.vue` — `v-if="authStore.isAdmin"` with "Admin access required" fallback
- `ScopeListView.vue` — "Create Scope" button hidden for non-admin
- `SettingsView.vue` — `v-if="authStore.isAdmin"`
- `AppMenu.vue` — Admin section hidden for non-admin

### Local Development (kind)

#### `kind/setup.sh`
- Deploys Dex v2.45.1 (required for `groups` in `staticPasswords`)
- Self-signed TLS cert for Dex
- OIDC config ConfigMaps and Secret
- In-cluster Dex issuer + external issuer for browser redirects

#### `kind/cluster.yaml`
- Port mapping `30556 → 5556` for Dex

#### `kind/dex/generate-cert.sh`
- Generates self-signed TLS cert for Dex

### Configuration

#### ConfigMap `pvc-explorer-config`
```yaml
oidc.enabled: "true"
oidc.issuer: "https://dex.pvc-explorer-system.svc.cluster.local:5556"
oidc.externalIssuer: "https://localhost:5556"
oidc.clientID: "pvc-explorer"
oidc.clientSecret: "$pvc-explorer-oidc.clientSecret"
oidc.redirectURI: "http://localhost:8080/api/v1/auth/oidc/callback"
oidc.scopes: "openid,profile,email,groups"
oidc.groupClaim: "groups"
oidc.skipTLSVerify: "true"
```

#### ConfigMap `pvc-explorer-rbac`
```yaml
policy.default: "viewer"
policy.csv: |
  g, "platform-admins", admin
  g, "platform-users", user
  g, "platform-viewers", viewer
```

### Tests

#### `internal/auth/oidc_test.go` (new)
- `MapGroupsToRole` — group mapping, default role, empty groups
- `LoadOIDCConfig` — enabled/disabled, missing fields
- `LoadRBACConfig` — defaults, custom config
- `IsOIDCConfigured` — true/false cases
- `SkipTLSVerify` — flag parsing

#### `internal/auth/auth_test.go` (updated)
- Session tests use new `Create(username, role, email, groups, subject)` signature

### Login Audit Logging

Successful logins emit structured logs:

```
Login successful method=oidc username=admin email=admin@pvc-explorer.local subject=... groups=[platform-admins] role=admin
Login successful method=local username=admin role=admin
Login failed method=local username=wronguser reason="invalid credentials"
```

Failed OIDC logins include error details:

```
OIDC token exchange failed err="..."
OIDC token verification failed err="..."
```

### Key Decisions

- **Generic OIDC, not Azure AD specific** — treats OIDC as a standard protocol
- **No Bearer token auth** — browser-based SSO avoids manual user provisioning
- **Three roles (admin/user/viewer)** — role hierarchy with rank comparison
- **ConfigMap-based config** — chosen over CRD or env vars for simplicity
- **ArgoCD-style policy.csv** — group-to-role mapping uses `g, <group>, <role>` format
- **Dex v2.45.1** — required for `groups` support in `staticPasswords` (PR #4456)
- **`oidc.externalIssuer` pattern** — in-cluster Dex URL for API, localhost for browser redirects
- **HTTP redirect URI** — dashboard serves on `http://localhost:8080`, not HTTPS
- **File write ops admin-only** — users see disabled buttons, backend blocks non-GET proxy requests
- **Session store carries OIDC claims** — profile page displays email, groups, subject

### Test Users (kind)

| Email | Password | Role |
|-------|----------|------|
| admin@pvc-explorer.local | admin123 | admin |
| user@pvc-explorer.local | user123 | user |
| viewer@pvc-explorer.local | viewer123 | viewer |

Local login (non-OIDC): `admin` / `admin`
