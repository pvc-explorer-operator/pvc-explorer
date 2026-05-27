# UI Components

This page documents the main UI components, their responsibilities, and their role in the runtime experience.

## Layout shell

### App layout

- Component: `AppLayout`
- Purpose: top-level authenticated shell with topbar, sidebar, route outlet, and footer.
- Source: `ui/src/layout/AppLayout.vue`

### Topbar

- Component: `AppTopbar`
- Purpose: breadcrumb context, manual refresh, auto-refresh interval picker, theme toggle.
- Source: `ui/src/layout/AppTopbar.vue`

### Sidebar

- Component: `AppSidebar`
- Purpose: main navigation + dashboard-only filter panel + user footer actions.
- Source: `ui/src/layout/AppSidebar.vue`

## Authentication components

### Login screen

- Route: `/login`
- View: `LoginView`
- Main component behavior:
  - `LoginForm` behavior is implemented directly in the view with PrimeVue `InputText`, `Password`, and `Button`.
  - Animated background canvas runs in the page.
- Source: `ui/src/views/LoginView.vue`

## Explorer listing and filtering

### Explorer listing views

- Route: `/`
- View: `HomeView`
- Display modes:
  - Card mode: `AppCardGrid`
  - Table mode: `AppListView`
- Toolbar features:
  - view mode toggle, sort field/direction, page size, pagination
- Sources:
  - `ui/src/views/HomeView.vue`
  - `ui/src/components/agents/AppCardGrid.vue`
  - `ui/src/components/agents/AppListView.vue`

### Filter panel

- Component: `FilterSidebar`
- Purpose: multi-dimension explorer filtering from the left sidebar.
- Filter groups:
  - Search
  - Phase
  - Namespace
  - Scope
  - Mount State
  - Access Mode
  - Consumers
  - Created age
  - Labels
- Source: `ui/src/components/filters/FilterSidebar.vue`

## Explorer detail and actions

### Explorer detail

- Route: `/explorers/:ns/:name`
- View: `AgentDetailView`
- Key UI elements:
  - phase indicator + status tags
  - mount-state banner
  - labels as chips
  - action buttons: Connect / Disconnect / Browse Files / Refresh
  - conditions table
  - consumer list
  - wake-up dialog
- Sources:
  - `ui/src/views/AgentDetailView.vue`
  - `ui/src/components/agents/ConditionsTable.vue`
  - `ui/src/components/agents/ConsumerList.vue`
  - `ui/src/components/shared/WakeUpDialog.vue`

## File browser workspace

### File browser page

- Route: `/explorers/:ns/:name/files`
- View: `FileBrowserView`
- Container component: `FileExplorerApp`
- Functional areas:
  - `FileTree`: directory tree and navigation
  - `FileExplorer`: current directory listing, actions, multi-select operations
  - `FilePreview`: image/pdf/binary and non-editable previews
  - `EditorPanel`: Monaco-based text editing
  - idle timer and idle warning banner
  - disconnected overlay with reconnect action
- Sources:
  - `ui/src/views/FileBrowserView.vue`
  - `ui/src/components/files/FileExplorerApp.vue`
  - `ui/src/components/files/FileTree.vue`
  - `ui/src/components/files/FileExplorer.vue`
  - `ui/src/components/files/FilePreview.vue`
  - `ui/src/components/files/EditorPanel.vue`

## Scope management

### Scope list

- Route: `/scopes`
- View: `ScopeListView`
- Modes:
  - card: `ScopeCardGrid`
  - list: `ScopeListViewTable`
- Source:
  - `ui/src/views/ScopeListView.vue`
  - `ui/src/components/scopes/ScopeCardGrid.vue`
  - `ui/src/components/scopes/ScopeListViewTable.vue`

### Scope detail

- Route: `/scopes/:name`
- View: `ScopeDetailView`
- Areas:
  - status badges
  - namespaces and selectors
  - defaults, discovery, deletion policy sections
  - linked explorer list
  - generated manifest panel with copy/download
- Source: `ui/src/views/ScopeDetailView.vue`

### Scope create

- Route: `/scopes/create`
- View: `CreateScopeView`
- Features:
  - form sections for identity, namespaces, discovery, defaults, deletion policy
  - live YAML generation with syntax highlighting
  - copy/download actions
- Source: `ui/src/views/CreateScopeView.vue`

## Admin and settings

### Settings

- Route: `/settings`
- View: `SettingsView`
- Features:
  - current identity/role/version
  - dark mode toggle
  - sign-out
- Source: `ui/src/views/SettingsView.vue`

### Create explorer (legacy/simple form)

- Route: `/explorers/create`
- View: `CreateAgentView`
- Features:
  - namespace + PVC selectors
  - idle timeout and force read-write options
- Source: `ui/src/views/CreateAgentView.vue`
