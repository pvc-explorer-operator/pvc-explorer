<template>
  <div class="ft-wrap" @click="$emit('clear-selection')">

    <!-- ── LIST VIEW ──────────────────────────────────────────────────────── -->
    <template v-if="viewMode === 'list'">
      <div class="ft-header">
        <div />
        <button class="ft-col-head" @click.stop="doSort('name')">
          Name<span v-if="sortField === 'name'" class="ft-sort">{{ sortDir === 'asc' ? ' ↑' : ' ↓' }}</span>
        </button>
        <button class="ft-col-head" @click.stop="doSort('size')">
          Size<span v-if="sortField === 'size'" class="ft-sort">{{ sortDir === 'asc' ? ' ↑' : ' ↓' }}</span>
        </button>
        <button class="ft-col-head" @click.stop="doSort('modTime')">
          Modified<span v-if="sortField === 'modTime'" class="ft-sort">{{ sortDir === 'asc' ? ' ↑' : ' ↓' }}</span>
        </button>
        <div />
      </div>

      <VirtualScroller
        v-if="sorted.length"
        :items="sorted"
        :item-size="36"
        class="ft-list-body"
        style="height: 100%"
      >
        <template #item="{ item: e, options }">
          <div
            class="ft-row"
            :class="{ 'ft-row--sel': selSet.has(e.name) }"
            :style="{ animationDelay: `${Math.min(options.index, 10) * 0.025}s` }"
            @click.stop="e.isDir ? $emit('navigate', e) : $emit('open', e)"
            @contextmenu.prevent.stop="$emit('context-menu', $event, e)"
          >
            <div class="ft-cell-chk">
              <input
                type="checkbox"
                class="ft-chk"
                :checked="selSet.has(e.name)"
                @click.stop="$emit('toggle-select', e.name)"
              />
            </div>
            <div class="ft-cell-name">
              <span class="ft-icon" :class="e.isDir ? 'ft-icon--dir' : 'ft-icon--file'">
                <svg v-if="e.isDir" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M10 4H4c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V8c0-1.1-.9-2-2-2h-8l-2-2z"/>
                </svg>
                <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                  <polyline points="14 2 14 8 20 8"/>
                </svg>
              </span>
              <span class="ft-name">{{ e.name }}</span>
            </div>
            <div class="ft-cell-meta">{{ e.isDir ? '—' : fmtSize(e.size) }}</div>
            <div class="ft-cell-meta">{{ fmtDate(e.modTime) }}</div>
            <div class="ft-cell-actions" @click.stop>
              <button v-if="!e.isDir" class="ft-act-btn" title="Download" @click="$emit('download', e)">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                  <polyline points="7 10 12 15 17 10"/>
                  <line x1="12" y1="15" x2="12" y2="3"/>
                </svg>
              </button>
              <button v-if="!readonly" class="ft-act-btn" title="Rename" @click="$emit('rename', e)">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                  <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                </svg>
              </button>
              <button v-if="!readonly" class="ft-act-btn ft-act-btn--danger" title="Delete" @click="$emit('delete', e)">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="3 6 5 6 21 6"/>
                  <path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6"/>
                </svg>
              </button>
            </div>
          </div>
        </template>
      </VirtualScroller>

      <div v-if="loading && !entries.length" class="ft-list-body">
        <div v-for="i in 10" :key="i" class="ft-row ft-row--sk">
          <div class="ft-cell-chk"><Skeleton width="13px" height="13px" /></div>
          <div class="ft-cell-name"><Skeleton width="1rem" height="1rem" /><Skeleton width="60%" height="0.875rem" /></div>
          <div class="ft-cell-meta"><Skeleton width="3rem" height="0.875rem" /></div>
          <div class="ft-cell-meta"><Skeleton width="5rem" height="0.875rem" /></div>
          <div class="ft-cell-actions" style="opacity:1"><Skeleton width="1.75rem" height="1.25rem" /></div>
        </div>
      </div>

      <div v-if="!entries.length && !loading" class="ft-empty">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
          <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
        </svg>
        <span>Empty folder</span>
        <div v-if="!readonly" class="ft-empty-actions">
          <button class="ft-empty-btn" title="Create a new folder" @click.stop="$emit('new-folder')">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/><line x1="12" y1="11" x2="12" y2="17"/><line x1="9" y1="14" x2="15" y2="14"/></svg>
            New Folder
          </button>
          <button class="ft-empty-btn" title="Create a new file" @click.stop="$emit('new-file')">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="12" y1="18" x2="12" y2="12"/><line x1="9" y1="15" x2="15" y2="15"/></svg>
            New File
          </button>
        </div>
      </div>
    </template>

    <!-- ── GRID VIEW ──────────────────────────────────────────────────────── -->
    <div v-else class="ft-grid">
      <div
        v-for="(e, i) in sorted"
        :key="e.name"
        class="ft-grid-item"
        :class="[e.isDir ? 'ft-grid-item--dir' : 'ft-grid-item--file', { 'ft-grid-item--sel': selSet.has(e.name) }]"
        :style="{ animationDelay: `${Math.min(i, 10) * 0.025}s` }"
        @click.stop="e.isDir ? $emit('navigate', e) : $emit('open', e)"
        @contextmenu.prevent.stop="$emit('context-menu', $event, e)"
      >
        <span class="ft-icon ft-icon--xl" :class="e.isDir ? 'ft-icon--dir' : 'ft-icon--file'">
          <svg v-if="e.isDir" viewBox="0 0 24 24" fill="currentColor">
            <path d="M10 4H4c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V8c0-1.1-.9-2-2-2h-8l-2-2z"/>
          </svg>
          <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
            <polyline points="14 2 14 8 20 8"/>
          </svg>
        </span>
        <span class="ft-grid-name">{{ e.name }}</span>
        <span class="ft-grid-size">{{ e.isDir ? '—' : fmtSize(e.size) }}</span>
      </div>

      <template v-if="loading && !entries.length">
        <div v-for="i in 12" :key="i" class="ft-grid-item ft-grid-item--sk">
          <Skeleton width="2.25rem" height="2.25rem" borderRadius="4px" />
          <Skeleton width="70%" height="0.6875rem" />
          <Skeleton width="40%" height="0.625rem" />
        </div>
      </template>

      <div v-if="!entries.length && !loading" class="ft-empty ft-empty--grid">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
            <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
          </svg>
          <span>Empty folder</span>
          <div v-if="!readonly" class="ft-empty-actions">
            <button class="ft-empty-btn" title="Create a new folder" @click.stop="$emit('new-folder')">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/><line x1="12" y1="11" x2="12" y2="17"/><line x1="9" y1="14" x2="15" y2="14"/></svg>
              New Folder
            </button>
            <button class="ft-empty-btn" title="Create a new file" @click.stop="$emit('new-file')">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="12" y1="18" x2="12" y2="12"/><line x1="9" y1="15" x2="15" y2="15"/></svg>
              New File
            </button>
          </div>
        </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { FileEntry } from '../../api/files'
import VirtualScroller from 'primevue/virtualscroller'
import Skeleton from 'primevue/skeleton'

const props = defineProps<{
  entries: FileEntry[]
  selectedNames: string[]
  viewMode: 'list' | 'grid'
  loading: boolean
  readonly: boolean
}>()

defineEmits<{
  (e: 'toggle-select', name: string): void
  (e: 'navigate', entry: FileEntry): void
  (e: 'open', entry: FileEntry): void
  (e: 'delete', entry: FileEntry): void
  (e: 'rename', entry: FileEntry): void
  (e: 'download', entry: FileEntry): void
  (e: 'context-menu', ev: MouseEvent, entry: FileEntry): void
  (e: 'clear-selection'): void
  (e: 'new-file'): void
  (e: 'new-folder'): void
}>()

const sortField = ref<'name' | 'size' | 'modTime'>('name')
const sortDir   = ref<'asc' | 'desc'>('asc')

const selSet = computed(() => new Set(props.selectedNames))

const sorted = computed(() => {
  const cmp = (a: FileEntry, b: FileEntry) => {
    const f = sortField.value as keyof FileEntry
    const va = (a[f] ?? '') as string | number
    const vb = (b[f] ?? '') as string | number
    let c: number
    if (f === 'size') c = ((va as number) || 0) - ((vb as number) || 0)
    else c = String(va).localeCompare(String(vb))
    return sortDir.value === 'asc' ? c : -c
  }
  return [
    ...[...props.entries.filter(e => e.isDir)].sort(cmp),
    ...[...props.entries.filter(e => !e.isDir)].sort(cmp),
  ]
})

function doSort(f: 'name' | 'size' | 'modTime') {
  if (sortField.value === f) sortDir.value = sortDir.value === 'asc' ? 'desc' : 'asc'
  else { sortField.value = f; sortDir.value = 'asc' }
}

function fmtSize(b: number): string {
  if (!b) return '0\u00a0B'
  const u = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(b) / Math.log(1024))
  return (b / Math.pow(1024, i)).toFixed(i ? 1 : 0) + '\u00a0' + u[i]
}

function fmtDate(iso: string): string {
  if (!iso) return '—'
  return new Date(iso).toLocaleDateString('en-GB', { day: '2-digit', month: 'short', year: 'numeric' })
}
</script>

<style scoped>
/* CSS vars cascade in from .fe-shell (FileExplorer.vue) */
.ft-wrap {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
  font-family: var(--font, monospace);
  background: var(--surface, #13162a);
  color: var(--text, #dde3f8);
}

/* ── Header ── */
.ft-header {
  display: grid;
  grid-template-columns: 36px 1fr 90px 130px 88px;
  padding: 6px 10px;
  background: var(--surface, #13162a);
  border-bottom: 1px solid var(--border, #252a42);
  flex-shrink: 0;
  position: sticky;
  top: 0;
  z-index: 1;
}
.ft-col-head {
  background: none;
  border: none;
  color: var(--muted, var(--text-color-secondary));
  cursor: pointer;
  font-family: Lato, sans-serif;
  font-size: 0.6875rem;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  padding: 2px 4px;
  text-align: left;
  transition: color .15s;
}
.ft-col-head:hover { color: var(--text, var(--text-color)); }
.ft-sort { color: var(--accent, var(--p-primary-500)); }

/* ── List body ── */
.ft-list-body {
  overflow-y: auto;
  flex: 1;
  scrollbar-color: light-dark(#aaa, #555) light-dark(#f0f0f0, #1a1d2b);
  scrollbar-width: thin;
}
.ft-row {
  display: grid;
  grid-template-columns: 36px 1fr 90px 130px 88px;
  align-items: center;
  padding: 6px 10px;
  border-bottom: 1px solid var(--border, var(--surface-border));
  cursor: pointer;
  animation: ft-fadein 0.18s ease both;
  transition: background .1s;
}
@keyframes ft-fadein {
  from { opacity: 0; transform: translateX(-4px) }
  to   { opacity: 1; transform: none }
}
.ft-row:hover     { background: var(--surface-hover); }
.ft-row--sel      { background: var(--sel-bg, color-mix(in srgb, var(--p-primary-500), transparent 92%)); border-color: var(--sel-bd); }
.ft-row:last-child { border-bottom: none; }

.ft-cell-chk { display: flex; align-items: center; padding: 0 4px; }
.ft-chk { accent-color: var(--accent, #4f8ef7); cursor: pointer; width: 13px; height: 13px; }

.ft-cell-name {
  display: flex;
  align-items: center;
  gap: 8px;
  overflow: hidden;
  padding: 0 4px;
}
.ft-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: var(--fs-sm, 0.875rem);
}
.ft-icon { display: flex; flex-shrink: 0; }
.ft-icon svg { width: 16px; height: 16px; }
.ft-icon--dir  { color: var(--warn, var(--p-amber-500)); }
.ft-icon--file { color: var(--file, var(--text-color-secondary)); }
.ft-icon--xl svg { width: 36px; height: 36px; }

.ft-cell-meta {
  color: var(--muted, var(--text-color-secondary));
  font-size: var(--fs-xs, 0.8125rem);
  padding: 0 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.ft-cell-actions {
  display: flex;
  align-items: center;
  gap: 3px;
  justify-content: flex-end;
  opacity: 0;
  transition: opacity .15s;
  padding: 0 4px;
}
.ft-row:hover .ft-cell-actions { opacity: 1; }

.ft-act-btn {
  background: none;
  border: none;
  cursor: pointer;
  color: var(--muted, var(--text-color-secondary));
  padding: 3px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  transition: color .1s, background .1s;
}
.ft-act-btn svg { width: 13px; height: 13px; }
.ft-act-btn:hover { color: var(--text-color); background: var(--surface-hover); }
.ft-act-btn--danger:hover { color: var(--p-red-400, #f87171); }

/* ── Empty state ── */
.ft-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 60px 20px;
  color: var(--muted, #5a6490);
  font-size: 13px;
}
.ft-empty svg { width: 40px; height: 40px; opacity: 0.4; }
.ft-empty--grid { grid-column: 1 / -1; }

.ft-empty-actions {
  display: flex;
  gap: 8px;
  margin-top: 4px;
}
.ft-empty-btn {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 5px 12px;
  border-radius: 5px;
  border: 1px solid rgba(255,255,255,.18);
  background: rgba(255,255,255,.07);
  color: var(--text, #dde3f8);
  cursor: pointer;
  font-size: 12px;
  font-family: Lato, sans-serif;
  transition: background .15s, border-color .15s;
}
.ft-empty-btn svg { width: 13px; height: 13px; flex-shrink: 0; }
.ft-empty-btn:hover { background: rgba(255,255,255,.14); border-color: rgba(255,255,255,.32); }

/* ── Grid ── */
.ft-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
  gap: 8px;
  padding: 16px;
  overflow-y: auto;
  flex: 1;
  align-content: start;
  scrollbar-color: light-dark(#aaa, #555) light-dark(#f0f0f0, #1a1d2b);
  scrollbar-width: thin;
}
.ft-grid-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 12px 8px 10px;
  border-radius: var(--r, 8px);
  border: 1px solid transparent;
  cursor: pointer;
  transition: background .15s, border-color .15s;
  text-align: center;
  animation: ft-fadein 0.2s ease both;
  user-select: none;
}
.ft-grid-item:hover { background: var(--surface2, #1c2038); border-color: var(--border, #252a42); }
.ft-grid-item--sel  { background: var(--sel-bg); border-color: var(--sel-bd); }
.ft-grid-name { font-size: 11px; word-break: break-word; line-height: 1.3; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; }
.ft-grid-size { font-size: 10px; color: var(--muted); }
</style>
