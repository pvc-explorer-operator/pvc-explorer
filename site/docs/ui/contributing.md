# UI Contributing

Reference for contributors working on the UI codebase.

## Local development

### Run UI only with mock data

```bash
cd ui
npm install
npm run dev
```

This uses the Vite mock plugin so you can develop without a running cluster.

### Run UI against local controller

Terminal 1:

```bash
cd ui
VITE_DEV_AUTH_BYPASS=false npm run dev
```

Terminal 2:

```bash
make run
```

Then open `/login` in the UI and test the real authentication flow.

## Technical map

### Routing and access control

- Router and guards: `ui/src/router/index.ts`
- Auth store: `ui/src/stores/authStore.ts`
- Flow doc: [Route and access flow](./flows.md#route-and-access-flow)

### State and realtime updates

- Explorer store: `ui/src/stores/explorerStore.ts`
- WebSocket composable: `ui/src/composables/useWebSocket.ts`
- Flow doc: [Explorer dashboard flow](./flows.md#explorer-dashboard-flow)

### File browser integration

- View shell: `ui/src/views/FileBrowserView.vue`
- API adapter: `ui/src/api/files.ts`
- Workspace component: `ui/src/components/files/FileExplorerApp.vue`
- Flow doc: [File browser flow](./flows.md#file-browser-flow)

### Scope management

- List: `ui/src/views/ScopeListView.vue`
- Detail: `ui/src/views/ScopeDetailView.vue`
- Create: `ui/src/views/CreateScopeView.vue`
- Flow doc: [Scope flow](./flows.md#scope-flow)

## Source reference

- [ui/README.md](https://github.com/pvc-explorer-operator/pvc-explorer/blob/main/ui/README.md)
