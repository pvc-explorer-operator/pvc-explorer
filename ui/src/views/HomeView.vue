<template>
  <main>
    <h1 class="sr-only">Explorers</h1>
    <!-- Toolbar: view toggle, sort, page size, pagination -->
    <div class="list-toolbar mb-4">
      <div class="toolbar-left">
        <div class="view-toggle btn-group">
          <button
            :class="['btn-icon', { active: viewMode === 'cards' }]"
            title="Card view"
            aria-label="Card view"
            @click="viewMode = 'cards'"
          >
            <i class="pi pi-th-large" aria-hidden="true"></i>
          </button>
          <button
            :class="['btn-icon', { active: viewMode === 'list' }]"
            title="List view"
            aria-label="List view"
            @click="viewMode = 'list'"
          >
            <i class="pi pi-list" aria-hidden="true"></i>
          </button>
        </div>

        <div class="toolbar-separator"></div>

        <label class="toolbar-label" for="sort-select">Sort:</label>
        <select id="sort-select" v-model="sortBy" class="toolbar-select" @change="page = 1">
          <option value="name">Name</option>
          <option value="namespace">Namespace</option>
          <option value="phase">Phase</option>
          <option value="pvcName">PVC</option>
          <option value="createdAt">Created</option>
        </select>
        <button class="btn-icon btn-sort-dir" @click="toggleSortDir" :title="sortDir === 'asc' ? 'Ascending' : 'Descending'" :aria-label="sortDir === 'asc' ? 'Sort ascending' : 'Sort descending'">
          <i :class="['pi', sortDir === 'asc' ? 'pi-sort-amount-up-alt' : 'pi-sort-amount-down']" aria-hidden="true"></i>
        </button>

        <div class="toolbar-separator"></div>

        <label class="toolbar-label" for="page-size-select">Items per page:</label>
        <select id="page-size-select" v-model="pageSize" class="toolbar-select toolbar-select--small" @change="page = 1">
          <option :value="10">10</option>
          <option :value="20">20</option>
          <option :value="50">50</option>
        </select>
      </div>

      <div class="toolbar-right">
        <span class="pagination-info">
          {{ paginationRange.start + 1 }}–{{ paginationRange.end }} of {{ sorted.length }}
        </span>
        <div class="pagination-btns btn-group">
          <button class="btn-icon" :disabled="page <= 1" @click="page--" title="Previous page" aria-label="Previous page">
            <i class="pi pi-chevron-left" aria-hidden="true"></i>
          </button>
          <button class="btn-icon" :disabled="page >= totalPages" @click="page++" title="Next page" aria-label="Next page">
            <i class="pi pi-chevron-right" aria-hidden="true"></i>
          </button>
        </div>
      </div>
    </div>

    <!-- Cards or List -->
    <AppCardGrid v-if="viewMode === 'cards'" :explorers="paginated" :loading="loading" />
    <AppListView v-else :explorers="paginated" />
  </main>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useExplorerStore } from '../stores/explorerStore'
import { useWebSocket } from '../composables/useWebSocket'
import AppCardGrid from '../components/agents/AppCardGrid.vue'
import AppListView from '../components/agents/AppListView.vue'

const store = useExplorerStore()
const { connect } = useWebSocket({
  onIdleTick: ({ namespace, name, remainingSeconds }) => {
    store.setIdleRemaining(namespace, name, remainingSeconds)
  },
  onAgentReady: ({ namespace, name }) => {
    store.updatePhase(namespace, name, 'Running')
  },
})

const explorers = computed(() => store.explorers)
const loading = ref(true)

/* ── View mode ── */
const viewMode = ref<'cards' | 'list'>('cards')

/* ── Sort ── */
const sortBy = ref('name')
const sortDir = ref<'asc' | 'desc'>('asc')

function toggleSortDir() {
  sortDir.value = sortDir.value === 'asc' ? 'desc' : 'asc'
}

/* ── Pagination ── */
const pageSize = ref(10)
const page = ref(1)

const sorted = computed(() => {
  const list = filtered.value
  const dir = sortDir.value === 'asc' ? 1 : -1
  const by = sortBy.value
  return [...list].sort((a, b) => {
    let cmp = 0
    if (by === 'name') cmp = (a.name ?? '').localeCompare(b.name ?? '')
    else if (by === 'namespace') cmp = (a.namespace ?? '').localeCompare(b.namespace ?? '')
    else if (by === 'phase') cmp = (a.phase ?? '').localeCompare(b.phase ?? '')
    else if (by === 'pvcName') cmp = (a.pvcName ?? '').localeCompare(b.pvcName ?? '')
    else if (by === 'createdAt') cmp = ((a.createdAt ?? '') < (b.createdAt ?? '') ? -1 : 1)
    return cmp * dir
  })
})

const totalPages = computed(() => Math.max(1, Math.ceil(sorted.value.length / pageSize.value)))

const paginated = computed(() => {
  const start = (page.value - 1) * pageSize.value
  return sorted.value.slice(start, start + pageSize.value)
})

const paginationRange = computed(() => {
  const start = (page.value - 1) * pageSize.value
  const end = Math.min(start + pageSize.value, sorted.value.length)
  return { start, end }
})

/* ── Filters (from sidebar) ── */
const filtered = computed(() => {
  let list = explorers.value
  const f = store.sidebarFilters as Record<string, unknown>
  const search = (f.search as string) ?? ''
  const phases = (f.phases as string[]) ?? []
  const namespaces = (f.namespaces as string[]) ?? []
  const mountStates = (f.mountStates as string[]) ?? []
  const scopes = (f.scopes as string[]) ?? []
  const accessModes = (f.accessModes as string[]) ?? []
  const consumers = (f.consumers as string) ?? ''
  const created = (f.created as string) ?? ''
  const labels = (f.labels as string[]) ?? []

  if (phases.length) {
    const hasInUse = phases.includes('InUse')
    const otherPhases = phases.filter(p => p !== 'InUse')
    list = list.filter(e => {
      if (otherPhases.length && otherPhases.includes(e.phase)) return true
      if (hasInUse && (e.consumerCount ?? 0) > 0) return true
      return false
    })
  }
  if (namespaces.length)  list = list.filter(e => namespaces.includes(e.namespace))
  if (mountStates.length) list = list.filter(e => mountStates.includes(e.mountState))
  if (scopes.length)      list = list.filter(e => e.scope !== undefined && scopes.includes(e.scope))
  if (accessModes.length) list = list.filter(e => {
    const m = e.accessMode || e.mode
    return m !== undefined && accessModes.includes(m)
  })
  if (consumers === 'has')  list = list.filter(e => (e.consumerCount ?? 0) > 0)
  if (consumers === 'none') list = list.filter(e => !(e.consumerCount ?? 0))
  if (created) {
    const now = Date.now()
    list = list.filter(e => {
      if (!e.createdAt) return false
      const age = now - new Date(e.createdAt).getTime()
      if (created === '24h')   return age < 86_400_000
      if (created === '7d')    return age < 604_800_000
      if (created === '30d')   return age < 2_592_000_000
      if (created === 'older') return age >= 2_592_000_000
      return true
    })
  }
  if (labels.length) {
    list = list.filter(e => {
      if (!e.labels?.length) return false
      return labels.every(l => e.labels!.includes(l))
    })
  }
  if (search) {
    const q = search.toLowerCase()
    list = list.filter(e =>
      e.name.toLowerCase().includes(q) ||
      e.namespace.toLowerCase().includes(q) ||
      e.pvcName.toLowerCase().includes(q)
    )
  }
  return list
})

onMounted(async () => {
  loading.value = true
  try {
    await store.fetchExplorers()
  } finally {
    loading.value = false
  }
  connect()
})
</script>

<style scoped>
/* ── Toolbar ── */
.list-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 0.5rem;
}
.toolbar-left,
.toolbar-right {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}
.toolbar-label {
  font-size: 0.8125rem;
  color: var(--text-color-secondary);
  white-space: nowrap;
}
.toolbar-select {
  height: 1.9rem;
  padding: 0 0.5rem;
  border: 1px solid rgba(255,255,255,0.2);
  border-radius: 6px;
  background: rgba(255,255,255,0.08);
  color: var(--text-color);
  font-size: 0.8125rem;
  cursor: pointer;
  outline: none;
}
.toolbar-select:focus {
  border-color: var(--primary-color);
}
.toolbar-select--small {
  width: 4.5rem;
}
.toolbar-separator {
  width: 1px;
  height: 1.4rem;
  background: var(--surface-border);
}
.btn-group {
  display: flex;
  border: 1px solid rgba(255,255,255,0.2);
  border-radius: 6px;
  overflow: hidden;
}
.btn-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  border: none;
  background: transparent;
  color: var(--text-color-secondary);
  cursor: pointer;
  font-size: 0.875rem;
  transition: background 0.1s, color 0.1s;
}
.btn-icon:hover {
  background: var(--surface-hover);
  color: var(--text-color);
}
.btn-icon:disabled {
  opacity: 0.4;
  cursor: default;
}
.btn-icon.active {
  background: var(--primary-color);
  color: var(--primary-color-text);
}
.btn-group .btn-icon + .btn-icon {
  border-left: 1px solid var(--surface-border);
}
.btn-group .btn-icon.active + .btn-icon {
  border-left-color: transparent;
}
.btn-sort-dir {
  border: 1px solid var(--surface-border);
  border-radius: var(--border-radius);
  width: 2rem;
  height: 2rem;
}
.pagination-info {
  font-size: 0.8125rem;
  color: var(--text-color-secondary);
  white-space: nowrap;
}
.pagination-btns .btn-icon {
  width: 2rem;
  height: 2rem;
}
.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border-width: 0;
}
</style>
