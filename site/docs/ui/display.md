# UI Display

This section is focused on what users see in the product UI.

For implementation details, data flow, and source-level behavior, jump to [UI Contributing](./contributing.md).

## Screens and experience

### 1. Login

- Route: `/login`
- What users do: authenticate and enter the app.
<!-- Screenshot: /login page — show the animated background canvas, username/password fields, and Login button. Capture a clean empty state before any input. -->
- Tech details: [Auth and route guard flow](./flows.md#route-and-access-flow)

### 2. Explorer dashboard

- Route: `/`
- What users do:
  - switch card/list views
  - sort and paginate explorers
  - watch live state updates
<!-- Screenshot: / dashboard in card view — show the full grid of explorer cards, the top toolbar (view toggle, sort, pagination), and the left sidebar partially visible. Aim for a state with several explorers in different phases (Running, ScaledToZero, Waking). -->
- Tech details:
  - [Dashboard flow](./flows.md#explorer-dashboard-flow)
  - [Explorer list components](./components.md#explorer-listing-and-filtering)

### 3. Sidebar filters

- Context: dashboard sidebar
- What users do: narrow explorers by phase, namespace, scope, mount state, labels, and more.
<!-- Screenshot: dashboard with the filter sidebar expanded — show the Phase filter group open with at least one checkbox ticked, and one or two other filter groups visible below it (Namespace, Scope). The explorer grid on the right should reflect the filtered result. -->
- Tech details: [Filter panel component](./components.md#filter-panel)

### 4. Explorer detail

- Route: `/explorers/:ns/:name`
- What users do:
  - inspect status, mount state, labels, conditions, consumers
  - connect/disconnect
  - open file browser
<!-- Screenshot: /explorers/:ns/:name — show the full detail page for a Running explorer. Include: phase badge, mount state banner, labels as chips, action buttons (Connect / Disconnect / Browse Files / Refresh), and the conditions table below. -->
- Tech details:
  - [Explorer detail flow](./flows.md#explorer-detail-flow)
  - [Explorer detail components](./components.md#explorer-detail-and-actions)

### 5. File browser

- Route: `/explorers/:ns/:name/files`
- What users do:
  - browse tree and directories
  - preview or edit files
  - upload/rename/delete
  - reconnect after idle sleep
<!-- Screenshot: /explorers/:ns/:name/files — show the file browser with the tree panel on the left (at least two directory levels), the file list in the centre, and an open text file in the EditorPanel on the right. The idle timer indicator should be visible. -->
- Tech details:
  - [File browser flow](./flows.md#file-browser-flow)
  - [File browser components](./components.md#file-browser-workspace)

### 6. Scope list and detail

- Routes:
  - `/scopes`
  - `/scopes/:name`
- What users do:
  - inspect managed namespaces and explorer counts
  - review defaults and generated manifest
- Screenshots:
  <!-- Screenshot: /scopes — show the scope list in card view with at least two scope cards. Each card should show the scope name, namespace count, and explorer count. -->
  <!-- Screenshot: /scopes/:name — show the scope detail page with namespace selector section, discovery defaults, and the linked explorers list at the bottom. -->
- Tech details:
  - [Scope flow](./flows.md#scope-flow)
  - [Scope components](./components.md#scope-management)

### 7. Scope creation

- Route: `/scopes/create`
- What users do:
  - configure scope form
  - verify live YAML preview
  - copy/download manifest
<!-- Screenshot: /scopes/create — show the form on the left (identity + namespaces sections filled in) and the live YAML preview panel on the right with syntax highlighting. -->
- Tech details:
  - [Scope flow](./flows.md#scope-flow)
  - [Create scope components](./components.md#scope-create)

### 8. Settings

- Route: `/settings`
- What users do: toggle theme, view version and role, sign out.
<!-- Screenshot: /settings — show role/version info, the dark mode toggle in its current state, and the Sign Out button. -->
- Tech details: [Settings component details](./components.md#admin-and-settings)

## Screenshot notes

To capture screenshots, run the UI in mock mode (`cd ui && npm run dev`) and open the app in your browser.
Replace each comment above with a markdown image: `![description](../public/images/ui/filename.png)`.
Store images under `site/docs/public/images/ui/`.
