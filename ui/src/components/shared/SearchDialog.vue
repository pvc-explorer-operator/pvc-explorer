<template>
  <dialog ref="dlgRef" class="search-dialog" aria-label="Search pages" @close="close">
    <div class="sd-header">
      <!-- eslint-disable-next-line vue/a11y-autofocus -->
      <input
        ref="inputRef"
        v-model="query"
        type="text"
        class="sd-input"
        placeholder="Search pages…"
        aria-label="Search query"
        @keydown="onKeydown"
      />
      <kbd class="sd-hint">{{ metaKey }}</kbd>
    </div>

    <div v-if="results.length" class="sd-body" role="listbox" :aria-activedescendant="activeId">
      <button
        v-for="(item, i) in results"
        :key="item.id"
        :id="item.id"
        :ref="(el) => { if (i === activeIndex) activeEl = el as HTMLElement }"
        role="option"
        :aria-selected="i === activeIndex"
        class="sd-result"
        :class="{ 'sd-result--active': i === activeIndex }"
        @click="go(item)"
        @mouseenter="activeIndex = i"
      >
        <i v-if="item.icon" :class="item.icon" aria-hidden="true" />
        <span class="sd-result-label">{{ item.label }}</span>
        <span v-if="item.hint" class="sd-result-hint">{{ item.hint }}</span>
      </button>
    </div>

    <div v-else-if="query && !results.length" class="sd-empty">
      No results for <strong>{{ query }}</strong>
    </div>
  </dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { searchDialogOpen } from '@/composables/useSearchDialog'
import { useExplorerStore } from '@/stores/explorerStore'

/* ── Detect Cmd vs Ctrl label ── */
const metaKey = computed(() => {
  if (typeof navigator === 'undefined') return '⌘'
  return /Mac|iP(hone|ad|od)/.test(navigator.platform) ? '⌘' : 'Ctrl'
})

/* ── Dialog refs ── */
const dlgRef = ref<HTMLDialogElement | null>(null)
const inputRef = ref<HTMLInputElement | null>(null)

/* ── Search state ── */
const query = ref('')
const activeIndex = ref(0)
const activeEl = ref<HTMLElement | null>(null)
const activeId = computed(() => results.value[activeIndex.value]?.id ?? undefined)
const router = useRouter()
const explorerStore = useExplorerStore()

/* ── Searchable items ── */
interface SearchItem {
  id: string
  label: string
  hint?: string
  icon?: string
  route?: string
  action?: () => void
}

const NAV_ITEMS: SearchItem[] = [
  { id: 's-home',         label: 'Explorers',       hint: '/',                    icon: 'pi pi-home',         route: '/' },
  { id: 's-scopes',       label: 'Scopes',           hint: '/scopes',              icon: 'pi pi-list',         route: '/scopes' },
  { id: 's-scope-create', label: 'Create Scope',     hint: '/scopes/create',       icon: 'pi pi-plus',         route: '/scopes/create' },
  { id: 's-settings',     label: 'Settings',         hint: '/settings',            icon: 'pi pi-cog',          route: '/settings' },
  { id: 's-about',        label: 'About',            hint: '/about',               icon: 'pi pi-info-circle',  route: '/about' },
]

const items = computed<SearchItem[]>(() => {
  const explorerItems: SearchItem[] = explorerStore.explorers.map((e) => ({
    id: `e-${e.namespace}-${e.name}`,
    label: e.name,
    hint: e.pvcName ? `pvc: ${e.pvcName}` : e.namespace,
    icon: 'pi pi-database',
    route: `/explorers/${encodeURIComponent(e.namespace)}/${encodeURIComponent(e.name)}`,
  }))
  const scopeItems: SearchItem[] = explorerStore.scopes.map((s) => ({
    id: `sc-${s.name}`,
    label: s.name,
    hint: 'scope',
    icon: 'pi pi-list',
    route: `/scopes/${encodeURIComponent(s.name)}`,
  }))
  return [...NAV_ITEMS, ...explorerItems, ...scopeItems]
})

watch(searchDialogOpen, async (val) => {
  if (val) {
    query.value = ''
    activeIndex.value = 0
    await nextTick()
    dlgRef.value?.showModal()
    inputRef.value?.focus()
  } else {
    dlgRef.value?.close()
  }
})

function close() {
  searchDialogOpen.value = false
}

/* ── Filter logic ── */
const results = computed(() => {
  const q = query.value.toLowerCase().trim()
  if (!q) return NAV_ITEMS
  return items.value.filter(
    (item) =>
      item.label.toLowerCase().includes(q) ||
      item.hint?.toLowerCase().includes(q)
  )
})

/* ── Keyboard navigation ── */
function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    close()
    return
  }

  if (e.key === 'ArrowDown') {
    e.preventDefault()
    activeIndex.value = Math.min(activeIndex.value + 1, results.value.length - 1)
    activeEl.value?.scrollIntoView({ block: 'nearest' })
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    activeIndex.value = Math.max(activeIndex.value - 1, 0)
    activeEl.value?.scrollIntoView({ block: 'nearest' })
  } else if (e.key === 'Enter') {
    e.preventDefault()
    const item = results.value[activeIndex.value]
    if (item) go(item)
  }
}

function go(item: SearchItem) {
  if (item.route) {
    router.push(item.route)
  } else if (item.action) {
    item.action()
  }
  close()
}
</script>

<style scoped>
.search-dialog {
  background: var(--surface-card);
  border: 1px solid var(--surface-border);
  border-radius: 12px;
  box-shadow: 0 24px 64px color-mix(in srgb, #000 40%, transparent);
  width: min(480px, calc(100vw - 2rem));
  max-height: min(480px, calc(100vh - 4rem));
  display: flex;
  flex-direction: column;
  overflow: hidden;
  padding: 0;

  /* Transition */
  opacity: 0;
  transform: scale(0.95) translateY(-8px);
  transition: opacity 0.15s ease, transform 0.15s ease,
    display 0.15s allow-discrete, overlay 0.15s allow-discrete;
  transition-behavior: allow-discrete;
}

.search-dialog[open] {
  opacity: 1;
  transform: scale(1) translateY(0);
  @starting-style {
    opacity: 0;
    transform: scale(0.95) translateY(-8px);
  }
}

.search-dialog::backdrop {
  background: rgba(0, 0, 0, 0);
  transition: display 0.15s allow-discrete, overlay 0.15s allow-discrete,
    background 0.15s ease;
  transition-behavior: allow-discrete;
}

.search-dialog[open]::backdrop {
  background: rgba(0, 0, 0, 0.45);
  @starting-style {
    background: rgba(0, 0, 0, 0);
  }
}

@media (prefers-reduced-motion: reduce) {
  .search-dialog { transform: none; transition-duration: 0.05s; }
}

/* ── Header / Input ── */
.sd-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  border-bottom: 1px solid var(--surface-border);
}

.sd-input {
  flex: 1;
  border: none;
  background: none;
  outline: none;
  font-size: 0.9375rem;
  color: var(--text-color);
  font-family: inherit;
}

.sd-input::placeholder {
  color: var(--text-color-secondary);
}

.sd-hint {
  display: inline-block;
  padding: 0.15em 0.45em;
  font-size: 0.72rem;
  font-family: ui-monospace, 'Cascadia Code', 'Fira Mono', monospace;
  color: var(--text-color-secondary);
  background: var(--surface-ground);
  border: 1px solid var(--surface-border);
  border-radius: 4px;
  line-height: 1.4;
  flex-shrink: 0;
}

/* ── Results ── */
.sd-body {
  overflow-y: auto;
  padding: 0.35rem 0.5rem;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.sd-result {
  display: flex;
  align-items: center;
  gap: 0.6rem;
  padding: 0.5rem 0.6rem;
  border: none;
  background: none;
  border-radius: 6px;
  cursor: pointer;
  color: var(--text-color);
  text-align: left;
  font-size: 0.875rem;
  font-family: inherit;
  width: 100%;
  transition: background 0.1s;
}

.sd-result:hover,
.sd-result--active {
  background: var(--surface-hover);
}

.sd-result .pi {
  font-size: 0.9rem;
  color: var(--text-color-secondary);
  flex-shrink: 0;
}

.sd-result-label {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.sd-result-hint {
  font-size: 0.75rem;
  color: var(--text-color-secondary);
  flex-shrink: 0;
}

.sd-empty {
  padding: 2rem 1rem;
  text-align: center;
  color: var(--text-color-secondary);
  font-size: 0.875rem;
}
</style>
