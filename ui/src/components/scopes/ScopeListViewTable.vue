<template>
  <div class="sc-table-wrap">
    <table v-if="scopes.length" class="sc-table">
      <thead>
        <tr>
          <th>Status</th>
          <th>Name</th>
          <th>Namespaces</th>
          <th>Explorers</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="scope in scopes"
          :key="scope.name"
          @click="$router.push(`/scopes/${encodeURIComponent(scope.name)}`)"
        >
          <td>
            <span :class="['phase-dot', phaseClass(scope)]" />
            <span :class="['sc-phase-tag', phaseClass(scope)]">{{ scope.phase }}</span>
          </td>
          <td class="font-semibold sc-name">{{ scope.name }}</td>
          <td class="text-muted-color">{{ scope.namespaceCount }}</td>
          <td class="text-muted-color">{{ scope.explorerCount }}</td>
          <td class="text-right"><i class="pi pi-chevron-right text-muted-color" /></td>
        </tr>
      </tbody>
    </table>
    <div v-else class="empty-state">
      <i class="pi pi-inbox empty-icon" />
      <div class="empty-msg">No scopes found.</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Scope } from '../../stores/explorerStore'

defineProps<{ scopes: Scope[] }>()

function phaseClass(s: Scope) {
  const p = s.phase?.toLowerCase() ?? ''
  if (p === 'ready' || p === 'true') return 'phase-ok'
  if (p === 'failed') return 'phase-fail'
  return 'phase-warn'
}
</script>

<style scoped>
.sc-table {
  width: 100%;
  border-collapse: collapse;
}
.sc-table th {
  text-align: left;
  padding: 0.65rem 0.75rem;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--text-color-secondary);
  border-bottom: 1px solid var(--surface-border);
}
.sc-table td {
  padding: 0.65rem 0.75rem;
  font-size: 0.875rem;
  border-bottom: 1px solid var(--surface-border);
  vertical-align: middle;
}
.sc-table tbody tr {
  cursor: pointer;
  transition: background-color 0.1s;
}
.sc-table tbody tr:hover {
  background-color: var(--surface-hover);
}

.phase-dot {
  display: inline-block;
  width: 0.55rem;
  height: 0.55rem;
  border-radius: 50%;
  margin-right: 0.45rem;
  vertical-align: middle;
}
.phase-ok   { background: #22c55e; }
.phase-warn { background: #f59e0b; }
.phase-fail { background: #ef4444; }

.sc-phase-tag {
  font-size: 0.75rem;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 12px;
  vertical-align: middle;
}
.sc-phase-tag.phase-ok   { background: rgba(34,197,94,0.12); color: #22c55e; }
.sc-phase-tag.phase-warn { background: rgba(245,158,11,0.12); color: #f59e0b; }
.sc-phase-tag.phase-fail { background: rgba(239,68,68,0.12);  color: #ef4444; }

.sc-name {
  font-family: 'JetBrains Mono', monospace;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 220px;
  color: var(--text-color-secondary);
  opacity: 0.85;
}
.empty-icon {
  font-size: 2.5rem;
  margin-bottom: 0.7rem;
}
.empty-msg {
  font-size: 1.08rem;
}
</style>
