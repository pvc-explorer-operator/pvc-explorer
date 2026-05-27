<template>
  <main>
    <h1 class="sr-only">Scopes</h1>
    <div class="flex justify-end items-center mb-4">
      <Button v-if="authStore.isAdmin" severity="primary" icon="pi pi-plus" label="Create Scope" rounded @click="router.push('/scopes/create')" />
    </div>

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

        <label class="toolbar-label" for="scope-sort-select">Sort:</label>
        <select id="scope-sort-select" v-model="sortBy" class="toolbar-select" @change="page = 1">
          <option value="name">Name</option>
          <option value="phase">Phase</option>
          <option value="namespaceCount">Namespaces</option>
          <option value="explorerCount">Explorers</option>
        </select>
        <button class="btn-icon btn-sort-dir" @click="toggleSortDir" :title="sortDir === 'asc' ? 'Ascending' : 'Descending'" :aria-label="sortDir === 'asc' ? 'Sort ascending' : 'Sort descending'">
          <i :class="['pi', sortDir === 'asc' ? 'pi-sort-amount-up-alt' : 'pi-sort-amount-down']"></i>
        </button>

        <div class="toolbar-separator"></div>

        <label class="toolbar-label" for="scope-page-size-select">Items per page:</label>
        <select id="scope-page-size-select" v-model="pageSize" class="toolbar-select toolbar-select--small" @change="page = 1">
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
    <ScopeCardGrid v-if="viewMode === 'cards'" :scopes="paginated" />
    <ScopeListViewTable v-else :scopes="paginated" />
  </main>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useExplorerStore } from '../stores/explorerStore'
import { useAuthStore } from '../stores/authStore'
import Button from 'primevue/button'
import ScopeCardGrid from '../components/scopes/ScopeCardGrid.vue'
import ScopeListViewTable from '../components/scopes/ScopeListViewTable.vue'

const router = useRouter()
const store = useExplorerStore()
const authStore = useAuthStore()

const scopes = computed(() => store.scopes)

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
  const list = scopes.value
  const dir = sortDir.value === 'asc' ? 1 : -1
  const by = sortBy.value
  return [...list].sort((a, b) => {
    let cmp = 0
    if (by === 'name') cmp = (a.name ?? '').localeCompare(b.name ?? '')
    else if (by === 'phase') cmp = (a.phase ?? '').localeCompare(b.phase ?? '')
    else if (by === 'namespaceCount') cmp = (a.namespaceCount ?? 0) - (b.namespaceCount ?? 0)
    else if (by === 'explorerCount') cmp = (a.explorerCount ?? 0) - (b.explorerCount ?? 0)
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

onMounted(() => store.fetchScopes())
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
