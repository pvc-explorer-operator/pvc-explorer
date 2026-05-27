# UI Flows

This page describes the UI user flows and data flows, including auth, list and detail navigation, and the file browser lifecycle.

## Route and access flow

```mermaid
flowchart TD
  A["App start"] --> B["authStore.init"]
  B --> C{Authenticated?}
  C -- No --> D["/login"]
  C -- Yes --> E["/ or protected route"]
  E --> F{adminOnly route?}
  F -- No --> G["Render view"]
  F -- Yes and admin --> G
  F -- Yes and not admin --> H["Redirect to /"]
```

Implementation references:

- Route guards: `ui/src/router/index.ts`
- Auth state: `ui/src/stores/authStore.ts`

## Explorer dashboard flow

```mermaid
flowchart LR
  A["HomeView mounted"] --> B["fetchExplorers REST"]
  A --> C["connect WebSocket"]
  C --> D["snapshot event"]
  D --> E["store.setSnapshot"]
  C --> F["explorer.updated/deleted"]
  F --> G["store.upsert/remove explorer"]
  E --> H["FilterSidebar emits filters"]
  H --> I["store.setSidebarFilters"]
  I --> J["HomeView computed filter/sort/page"]
  J --> K["Card or List render"]
```

Implementation references:

- View: `ui/src/views/HomeView.vue`
- Store: `ui/src/stores/explorerStore.ts`
- Realtime: `ui/src/composables/useWebSocket.ts`
- Filters: `ui/src/components/filters/FilterSidebar.vue`

## Explorer detail flow

```mermaid
sequenceDiagram
  participant User
  participant AgentDetailView
  participant ExplorerStore
  participant API

  User->>AgentDetailView: Open /explorers/:ns/:name
  AgentDetailView->>ExplorerStore: fetchExplorer(ns, name)
  ExplorerStore->>API: GET /api/v1/explorers/:ns/:name
  API-->>ExplorerStore: explorer payload
  ExplorerStore-->>AgentDetailView: explorer model
  User->>AgentDetailView: Connect/Disconnect/Refresh
  AgentDetailView->>ExplorerStore: wakeExplorer/sleepExplorer/fetchExplorer
```

Implementation references:

- View: `ui/src/views/AgentDetailView.vue`
- Store methods: `ui/src/stores/explorerStore.ts`

## File browser flow

```mermaid
flowchart TD
  A["Open /explorers/:ns/:name/files"] --> B["FileBrowserView onMounted"]
  B --> C["fetch agent config"]
  B --> D["heartbeat POST"]
  B --> E["connect WebSocket idle events"]
  B --> F["render FileExplorerApp"]

  F --> G["navigate tree/list"]
  G --> H["fetchFiles path"]

  F --> I["open file"]
  I --> J{previewable or too large?}
  J -- Yes --> K["show FilePreview"]
  J -- No --> L["fetchContent + EditorPanel"]

  F --> M["save/rename/delete/upload/create"]
  M --> N["file API operation"]
  N --> O["reload entries"]

  E --> P["idle.warning/idle.expired"]
  P --> Q["warning banner or disconnected overlay"]
  Q --> R["Reconnect action"]
```

Implementation references:

- View: `ui/src/views/FileBrowserView.vue`
- API adapter: `ui/src/api/files.ts`
- Workspace container: `ui/src/components/files/FileExplorerApp.vue`

## Scope flow

```mermaid
flowchart LR
  A["Open /scopes"] --> B["fetchScopes"]
  B --> C["Sort + paginate"]
  C --> D["Scope card/list view"]
  D --> E["Open /scopes/:name"]
  E --> F["Render scope detail + related explorers"]
  A --> G["Admin clicks Create Scope"]
  G --> H["/scopes/create"]
  H --> I["Fill form"]
  I --> J["Live YAML preview"]
  J --> K["Copy or download manifest"]
```

Implementation references:

- Scope list: `ui/src/views/ScopeListView.vue`
- Scope detail: `ui/src/views/ScopeDetailView.vue`
- Scope create: `ui/src/views/CreateScopeView.vue`
