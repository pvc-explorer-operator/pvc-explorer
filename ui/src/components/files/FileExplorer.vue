<template>
  <div class="fe-shell" @click="closeCtx" @keydown.escape="closeCtx()" @keydown="onShellKeydown">

    <!-- ── Toolbar ────────────────────────────────────────────────────────── -->
    <div class="fe-toolbar">
      <button class="fe-btn fe-btn-ghost fe-btn-icon" :disabled="!currentPath" title="Go to parent" @click.stop="goUp">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="fe-ico"><polyline points="15 18 9 12 15 6"/></svg>
      </button>
      <FileBreadcrumb :path="currentPath" @navigate="emit('navigate', $event)" />
      <div class="fe-toolbar-right">
        <button class="fe-btn fe-btn-ghost fe-btn-icon" title="Toggle view" @click.stop="toggleView">
          <svg v-if="viewMode === 'list'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="fe-ico"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/></svg>
          <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="fe-ico"><line x1="8" y1="6" x2="21" y2="6"/><line x1="8" y1="12" x2="21" y2="12"/><line x1="8" y1="18" x2="21" y2="18"/><line x1="3" y1="6" x2="3.01" y2="6"/><line x1="3" y1="12" x2="3.01" y2="12"/><line x1="3" y1="18" x2="3.01" y2="18"/></svg>
        </button>
        <button class="fe-btn fe-btn-primary" :disabled="readonly" title="Create a new folder" @click.stop="openNewFolder">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="fe-ico"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/><line x1="12" y1="11" x2="12" y2="17"/><line x1="9" y1="14" x2="15" y2="14"/></svg>
          New Folder
        </button>
        <button class="fe-btn fe-btn-ghost" :disabled="readonly" title="Create a new file" @click.stop="openNewFile">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="fe-ico"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="12" y1="18" x2="12" y2="12"/><line x1="9" y1="15" x2="15" y2="15"/></svg>
          New File
        </button>
        <UploadZone :readonly="readonly" :uploading="uploading" @upload="emit('upload', $event)" />
      </div>
    </div>

    <!-- ── Selection bar ──────────────────────────────────────────────────── -->
    <Transition name="fe-slide">
      <div v-if="selectedNames.length > 0" class="fe-sel-bar">
        <span class="fe-sel-count">{{ selectedNames.length }} selected</span>
        <button class="fe-btn fe-btn-ghost fe-btn-sm" title="Deselect all items" @click.stop="clearSelection">Deselect all</button>
        <div class="fe-sel-actions">
          <button class="fe-btn fe-btn-ghost fe-btn-sm" title="Download selected files" @click.stop="downloadSelected">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="fe-ico"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
            Download
          </button>
          <button class="fe-btn fe-btn-ghost fe-btn-sm" :disabled="readonly || selectedNames.length !== 1" title="Rename selected item" @click.stop="renameSelected">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="fe-ico"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
            Rename
          </button>
          <button class="fe-btn fe-btn-danger fe-btn-sm" :disabled="readonly" title="Delete selected items" @click.stop="deleteSelected">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="fe-ico"><polyline points="3 6 5 6 21 6"/><path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6"/></svg>
            Delete
          </button>
        </div>
      </div>
    </Transition>

    <!-- ── File listing ───────────────────────────────────────────────────── -->
    <div class="fe-table-area">
      <FileTable
        :entries="entries"
        :selected-names="selectedNames"
        :view-mode="viewMode"
        :loading="loading"
        :readonly="readonly"
        @toggle-select="toggleSelect"
        @navigate="onNavigate"
        @open="emit('file-opened', $event)"
        @delete="onDeleteEntry"
        @rename="onRenameEntry"
        @download="emit('download', $event)"
        @context-menu="onCtxMenu"
        @clear-selection="clearSelection"
        @new-file="openNewFile"
        @new-folder="openNewFolder"
      />
    </div>

    <!-- ── Status bar ─────────────────────────────────────────────────────── -->
    <div class="fe-status">
      <span class="fe-status-dot" />
      <span>{{ currentPath || '/' }}</span>
      <span v-if="readonly" class="fe-status-badge">read-only</span>
      <span class="fe-status-right">{{ entries.length }} item{{ entries.length !== 1 ? 's' : '' }}</span>
    </div>

    <!-- ── Context menu (position:fixed) ─────────────────────────────────── -->
    <div
      v-if="ctxVisible"
      class="fe-ctx"
      :style="{ position: 'fixed', top: ctxY + 'px', left: ctxX + 'px' }"
      @click.stop
    >
      <button class="fe-ctx-item" @click="ctxOpen">
        <svg v-if="ctxEntry?.isDir" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="fe-ico"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/></svg>
        <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="fe-ico"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
        {{ ctxEntry?.isDir ? 'Open' : 'Download' }}
      </button>
      <button v-if="!readonly" class="fe-ctx-item" @click="ctxRename">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="fe-ico"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
        Rename
      </button>
      <div class="fe-ctx-sep" />
      <button v-if="!readonly" class="fe-ctx-item fe-ctx-item--danger" @click="ctxDelete">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="fe-ico"><polyline points="3 6 5 6 21 6"/><path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6"/></svg>
        Delete
      </button>
    </div>

    <!-- ── Modal ─────────────────────────────────────────────────────────── -->
    <dialog v-if="modal" ref="modalDialogEl" class="fe-modal" @close="onDialogClose">
      <div class="fe-modal-title">{{ modal.title }}</div>
      <div v-if="modal.message" class="fe-modal-msg">{{ modal.message }}</div>
      <input
        v-if="modal.hasInput"
        ref="modalInputEl"
        v-model="modalInputValue"
        class="fe-modal-input"
        :placeholder="modal.placeholder ?? ''"
        @keyup.enter="confirmModal"
      />
      <div class="fe-modal-actions">
        <button class="fe-btn fe-btn-ghost" @click="closeModal">Cancel</button>
        <button
          class="fe-btn"
          :class="modal.danger ? 'fe-btn-danger' : 'fe-btn-primary'"
          @click="confirmModal"
        >{{ modal.confirmLabel }}</button>
      </div>
    </dialog>

    <!-- ── Toast ─────────────────────────────────────────────────────────── -->
    <Transition name="fe-toast">
      <div v-if="toast" class="fe-toast" :class="'fe-toast--' + toast.type">
        <svg v-if="toast.type === 'success'" class="fe-toast-ico" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>
        <svg v-else class="fe-toast-ico" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
        {{ toast.msg }}
      </div>
    </Transition>

  </div>
</template>

<script setup lang="ts">
import { ref, watch, nextTick, computed } from 'vue'
import FileBreadcrumb from './FileBreadcrumb.vue'
import FileTable from './FileTable.vue'
import UploadZone from './UploadZone.vue'
import type { FileEntry } from '../../api/files'
import { useFileBrowserShortcuts } from '../../composables/useFileBrowserShortcuts'

interface ModalState {
  title: string
  message?: string
  hasInput: boolean
  placeholder?: string
  danger: boolean
  confirmLabel: string
  onConfirm: (value: string) => void
}

const props = defineProps<{
  entries: FileEntry[]
  currentPath: string
  loading: boolean
  uploading: boolean
  readonly: boolean
}>()

const emit = defineEmits<{
  (e: 'navigate', path: string): void
  (e: 'file-opened', entry: FileEntry): void
  (e: 'delete', path: string): void
  (e: 'delete-many', paths: string[]): void
  (e: 'rename', from: string, to: string): void
  (e: 'upload', files: File[]): void
  (e: 'new-file', path: string): void
  (e: 'new-folder', path: string): void
  (e: 'download', entry: FileEntry): void
}>()

/* ── Selection ─────────────────────────────────────────────────────────────── */
const selectedNames = ref<string[]>([])

function toggleSelect(name: string) {
  const idx = selectedNames.value.indexOf(name)
  selectedNames.value = idx === -1
    ? [...selectedNames.value, name]
    : selectedNames.value.filter(n => n !== name)
}

function onShellKeydown(e: KeyboardEvent) {
  if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'a') {
    e.preventDefault()
    selectAll()
  }
}

function clearSelection() { selectedNames.value = [] }

function selectAll() {
  selectedNames.value = props.entries.filter(e => !e.isDir).map(e => e.name)
}

// Clear selection whenever directory entries change
watch(() => props.entries, () => clearSelection())

/* ── View mode ─────────────────────────────────────────────────────────────── */
const viewMode = ref<'list' | 'grid'>('list')
function toggleView() { viewMode.value = viewMode.value === 'list' ? 'grid' : 'list' }

/* ── Navigation ────────────────────────────────────────────────────────────── */
function goUp() {
  if (!props.currentPath) return
  const segs = props.currentPath.split('/')
  segs.pop()
  emit('navigate', segs.join('/'))
}

function onNavigate(entry: FileEntry) {
  emit('navigate', props.currentPath ? `${props.currentPath}/${entry.name}` : entry.name)
}

/* ── Context menu ──────────────────────────────────────────────────────────── */
const ctxVisible = ref(false)
const ctxX = ref(0)
const ctxY = ref(0)
const ctxEntry = ref<FileEntry | null>(null)

function onCtxMenu(ev: MouseEvent, entry: FileEntry) {
  ctxEntry.value = entry
  ctxX.value = Math.min(ev.clientX, window.innerWidth - 186)
  ctxY.value = Math.min(ev.clientY, window.innerHeight - 130)
  ctxVisible.value = true
}

function closeCtx() { ctxVisible.value = false }

function ctxOpen() {
  closeCtx()
  if (ctxEntry.value?.isDir) onNavigate(ctxEntry.value)
  else if (ctxEntry.value) emit('download', ctxEntry.value)
}

function ctxRename() { closeCtx(); if (ctxEntry.value) openRename(ctxEntry.value) }
function ctxDelete()  { closeCtx(); if (ctxEntry.value) onDeleteEntry(ctxEntry.value) }

/* ── Modal ─────────────────────────────────────────────────────────────────── */
const modal = ref<ModalState | null>(null)
const modalInputValue = ref('')
const modalInputEl = ref<HTMLInputElement | null>(null)
const modalDialogEl = ref<HTMLDialogElement | null>(null)

function openModal(config: ModalState, initialInput = '') {
  modal.value = config
  modalInputValue.value = initialInput
  nextTick(() => {
    modalDialogEl.value?.showModal()
    if (config.hasInput) modalInputEl.value?.select()
  })
}

function closeModal() {
  modalDialogEl.value?.close()
}

function onDialogClose() {
  modal.value = null
}

function confirmModal() {
  modal.value?.onConfirm(modalInputValue.value)
  closeModal()
}

/* ── File operations ───────────────────────────────────────────────────────── */
function fullPath(name: string) {
  return props.currentPath ? `${props.currentPath}/${name}` : name
}

function onDeleteEntry(entry: FileEntry) {
  openModal({
    title: `Delete "${entry.name}"?`,
    message: 'This action is not reversible. The file will be permanently deleted.',
    hasInput: false,
    danger: true,
    confirmLabel: 'Delete',
    onConfirm: () => emit('delete', fullPath(entry.name)),
  })
}

function onRenameEntry(entry: FileEntry) { openRename(entry) }

function openRename(entry: FileEntry) {
  openModal({
    title: `Rename "${entry.name}"`,
    hasInput: true,
    placeholder: 'New name',
    danger: false,
    confirmLabel: 'Rename',
    onConfirm: (newName: string) => {
      const n = newName.trim()
      if (n && n !== entry.name) emit('rename', fullPath(entry.name), fullPath(n))
    },
  }, entry.name)
}

function openNewFolder() {
  openModal({
    title: 'New Folder',
    hasInput: true,
    placeholder: 'folder-name',
    danger: false,
    confirmLabel: 'Create',
    onConfirm: (name: string) => {
      const n = name.trim()
      if (n) emit('new-folder', fullPath(`${n}/.keep`))
    },
  }, 'New Folder')
}

function openNewFile() {
  openModal({
    title: 'New File',
    hasInput: true,
    placeholder: 'filename.txt',
    danger: false,
    confirmLabel: 'Create',
    onConfirm: (name: string) => {
      const n = name.trim()
      if (n) emit('new-file', fullPath(n))
    },
  }, '')
}

function downloadSelected() {
  for (const name of selectedNames.value) {
    const entry = props.entries.find(e => e.name === name)
    if (entry && !entry.isDir) emit('download', entry)
  }
  showToast('Downloading selected files')
}

function renameSelected() {
  const entry = props.entries.find(e => e.name === selectedNames.value[0])
  if (entry) openRename(entry)
}

function deleteSelected() {
  const n = selectedNames.value.length
  openModal({
    title: `Delete ${n} item${n > 1 ? 's' : ''}?`,
    message: 'This action is not reversible. The selected item(s) will be permanently deleted.',
    hasInput: false,
    danger: true,
    confirmLabel: 'Delete',
    onConfirm: () => {
      const paths = selectedNames.value.map(name => fullPath(name))
      clearSelection()
      if (paths.length === 1) emit('delete', paths[0])
      else emit('delete-many', paths)
    },
  })
}

const readonlyRef = computed(() => props.readonly)

useFileBrowserShortcuts({
  selectedNames,
  readonly: readonlyRef,
  selectAll,
  clearSelection,
  deleteSelected,
  downloadSelected,
  openNewFile,
  openNewFolder,
})

/* ── Toast ─────────────────────────────────────────────────────────────────── */
const toast = ref<{ msg: string; type: 'success' | 'danger' } | null>(null)
let toastTimer: ReturnType<typeof setTimeout> | null = null

function showToast(msg: string, type: 'success' | 'danger' = 'success') {
  if (toastTimer) clearTimeout(toastTimer)
  toast.value = { msg, type }
  toastTimer = setTimeout(() => { toast.value = null }, 2200)
}
</script>

<style scoped>
/* ── Theme properties — mapped to PrimeVue/app tokens ───────────────────── */
.fe-shell {
  --bg:      var(--surface-ground);
  --surface: var(--surface-card);
  --surface2: var(--surface-hover);
  --border:  var(--surface-border);
  --text:    var(--text-color);
  --muted:   var(--text-color-secondary);
  --accent:  var(--p-primary-500);
  --accent-h: var(--p-primary-400);
  --danger:  var(--p-red-500);
  --danger-h: var(--p-red-400);
  --warn:    var(--p-amber-500);
  --file:    var(--text-color-secondary);
  --sel-bg:  color-mix(in srgb, var(--p-primary-500), transparent 92%);
  --sel-bd:  color-mix(in srgb, var(--p-primary-500), transparent 68%);
  --r:       6px;
  --font:    Lato, sans-serif;

  display: flex;
  flex-direction: column;
  background: var(--surface);
  color: var(--text);
  font-family: var(--font);
  overflow: hidden;
  position: relative;
  height: 100%;
}

/* ── Toolbar ── */
.fe-toolbar {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  background: var(--surface);
  border-bottom: 1px solid var(--border);
  flex-shrink: 0;
}
.fe-toolbar-right {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-left: auto;
}

/* ── Button primitives — match global .p-button semi-transparent style ── */
.fe-btn {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 5px 10px;
  border-radius: var(--r);
  border: 1px solid rgba(255,255,255,.2);
  background: rgba(255,255,255,.08);
  color: var(--text-color);
  cursor: pointer;
  font-family: Lato, sans-serif;
  font-size: 0.8125rem;
  transition: background .15s, border-color .15s, color .15s, transform .1s;
  white-space: nowrap;
  line-height: 1.4;
}
.fe-btn:hover:not(:disabled) { background: rgba(255,255,255,.15); border-color: rgba(255,255,255,.35); color: #fff; }
.fe-btn:active:not(:disabled) { transform: scale(.97); }
.fe-btn:disabled { opacity: .35; cursor: not-allowed; }

.fe-btn-icon { padding: 5px; }
.fe-btn-primary { background: var(--p-primary-500); border-color: var(--p-primary-500); color: #fff; }
.fe-btn-primary:hover:not(:disabled) { background: var(--p-primary-400); border-color: var(--p-primary-400); color: #fff; }
.fe-btn-ghost { background: transparent; border-color: transparent; }
.fe-btn-ghost:hover:not(:disabled) { background: rgba(255,255,255,.08); border-color: rgba(255,255,255,.15); color: var(--text-color); }
.fe-btn-danger { background: rgba(239,68,68,.15); border-color: rgba(239,68,68,.4); color: var(--p-red-400); }
.fe-btn-danger:hover:not(:disabled) { background: rgba(239,68,68,.25); border-color: rgba(239,68,68,.6); color: var(--p-red-300); }
.fe-btn-sm { padding: 3px 8px; font-size: 0.75rem; }

.fe-ico { width: 13px; height: 13px; flex-shrink: 0; }

/* ── Selection bar ── */
.fe-sel-bar {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 6px 12px;
  background: color-mix(in srgb, var(--p-primary-500), transparent 94%);
  border-bottom: 1px solid var(--sel-bd);
  font-size: 0.8125rem;
  flex-shrink: 0;
}
.fe-sel-count  { color: var(--p-primary-400); font-weight: 600; }
.fe-sel-actions { margin-left: auto; display: flex; gap: 6px; }

.fe-slide-enter-active { transition: all .2s ease; }
.fe-slide-leave-active { transition: all .15s ease; }
.fe-slide-enter-from, .fe-slide-leave-to { opacity: 0; transform: translateY(-6px); }

/* ── Table area ── */
.fe-table-area { flex: 1; overflow: hidden; min-height: 0; }

/* ── Status bar ── */
.fe-status {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 5px 14px;
  background: var(--bg);
  border-top: 1px solid var(--border);
  color: var(--muted);
  font-size: 0.75rem;
  flex-shrink: 0;
  user-select: none;
}
.fe-status-dot { width: 6px; height: 6px; border-radius: 50%; background: var(--p-primary-500); flex-shrink: 0; }
.fe-status-badge {
  background: rgba(245,158,11,.12);
  color: var(--warn);
  border: 1px solid rgba(245,158,11,.25);
  border-radius: 3px;
  padding: 1px 6px;
  font-size: 0.6875rem;
}
.fe-status-right { margin-left: auto; }

/* ── Context menu ── */
.fe-ctx {
  background: var(--surface-overlay, var(--surface-card));
  border: 1px solid var(--surface-border);
  border-radius: var(--content-border-radius, var(--r));
  padding: 4px;
  min-width: 168px;
  box-shadow: 0 12px 32px rgba(0,0,0,.5);
  z-index: 1000;
  animation: fe-pop .12s ease;
}
@keyframes fe-pop { from { opacity: 0; transform: scale(.95) } to { opacity: 1; transform: none } }

.fe-ctx-item {
  display: flex; align-items: center; gap: 8px;
  width: 100%; background: none; border: none;
  color: var(--text-color); padding: 7px 10px;
  border-radius: calc(var(--content-border-radius, var(--r)) - 2px); cursor: pointer;
  font-family: Lato, sans-serif; font-size: 0.8125rem;
  transition: background .1s; text-align: left;
}
.fe-ctx-item:hover { background: var(--surface-hover); }
.fe-ctx-item--danger { color: var(--p-red-400); }
.fe-ctx-sep { height: 1px; background: var(--surface-border); margin: 3px 6px; }

/* ── Modal (native <dialog>) ── */
.fe-modal {
  opacity: 0;
  transform: scale(0.95);
  transition: opacity 0.2s ease, transform 0.2s ease,
              display 0.2s allow-discrete, overlay 0.2s allow-discrete;
  transition-behavior: allow-discrete;
  background: var(--surface-card);
  border: 1px solid var(--surface-border);
  border-radius: var(--content-border-radius, var(--r));
  padding: 24px;
  min-width: 320px; max-width: 440px; width: 90%;
  box-shadow: 0 24px 48px rgba(0,0,0,.5);
}
.fe-modal[open] {
  opacity: 1;
  transform: scale(1);
  @starting-style { opacity: 0; transform: scale(0.95); }
}
.fe-modal::backdrop {
  background: rgb(0 0 0 / 0);
  transition: display 0.2s allow-discrete, overlay 0.2s allow-discrete,
              background 0.2s ease;
  transition-behavior: allow-discrete;
  backdrop-filter: blur(2px);
}
.fe-modal[open]::backdrop {
  background: rgb(0 0 0 / 50%);
  @starting-style { background: rgb(0 0 0 / 0); }
}
@media (prefers-reduced-motion: reduce) {
  .fe-modal { transform: none; transition-duration: 0.1s; }
}

.fe-modal-title { font-size: 1rem; font-weight: 600; margin-bottom: 8px; color: var(--text-color); }
.fe-modal-msg   { font-size: 0.8125rem; color: var(--text-color-secondary); margin-bottom: 16px; line-height: 1.5; }

.fe-modal-input {
  width: 100%;
  background: var(--p-form-field-background, rgba(255,255,255,.08));
  border: 1px solid var(--p-form-field-border-color, rgba(255,255,255,.2));
  border-radius: var(--r);
  padding: 8px 10px;
  color: var(--p-form-field-color, var(--text-color));
  font-family: Lato, sans-serif;
  font-size: 0.875rem;
  margin-bottom: 16px;
  outline: none;
  transition: border-color .15s;
  box-sizing: border-box;
}
.fe-modal-input:focus { border-color: var(--p-primary-500); }

.fe-modal-actions { display: flex; justify-content: flex-end; gap: 8px; }

/* ── Toast ── */
.fe-toast {
  position: fixed; bottom: 24px; right: 24px;
  background: var(--surface-card);
  border: 1px solid var(--surface-border);
  border-radius: var(--content-border-radius, var(--r));
  padding: 10px 16px;
  font-size: 0.8125rem;
  font-family: Lato, sans-serif;
  box-shadow: 0 8px 24px rgba(0,0,0,.35);
  z-index: 3000;
  display: flex; align-items: center; gap: 8px;
  color: var(--text-color);
  max-width: 320px;
}
.fe-toast--success { border-color: rgba(34,197,94,.4); }
.fe-toast--success .fe-toast-ico { color: var(--p-green-400, #4ade80); }
.fe-toast--danger  { border-color: rgba(239,68,68,.4); }
.fe-toast--danger  .fe-toast-ico { color: var(--p-red-400, #f87171); }
.fe-toast-ico { width: 14px; height: 14px; flex-shrink: 0; }

.fe-toast-enter-active { transition: all .2s ease; }
.fe-toast-leave-active { transition: all .15s ease; }
.fe-toast-enter-from, .fe-toast-leave-to { opacity: 0; transform: translateY(8px); }
</style>
