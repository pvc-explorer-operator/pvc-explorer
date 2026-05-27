<template>
  <div class="fea-layout">

    <!-- Disconnected overlay -->
    <div v-if="disconnected" class="fea-disconnected">
      <div class="fea-disconnected__card">
        <i class="pi pi-power-off fea-disconnected__icon" />
        <p class="fea-disconnected__title">Disconnected</p>
        <p class="fea-disconnected__sub">The agent went to sleep due to inactivity.</p>
        <button
          class="fea-disconnected__btn"
          :disabled="reconnecting"
          @click="emit('reconnect')"
        >
          <i v-if="reconnecting" class="pi pi-spin pi-spinner" style="margin-right:6px" />
          {{ reconnecting ? 'Waking up…' : 'Reconnect' }}
        </button>
      </div>
    </div>
    <!-- Left panel: recursive directory tree -->
    <FileTree
      :current-path="currentPath"
      :fetch-files="fetchFiles"
      class="fea-tree"
      @navigate="navigateTo"
      @open-file="handleOpenFileFromTree"
    />

    <!-- Right panel: top-bar + explorer + editor/preview -->
    <div class="fea-main">
      <!-- Top bar: explorer label + idle timer -->
      <div class="fea-topbar">
        <span class="fea-label">{{ explorerLabel }}</span>
        <IdleTimer v-if="remainingSeconds !== null" :remaining-seconds="remainingSeconds" class="ml-auto" />
      </div>

      <!-- Idle warning banner -->
      <div v-if="idleWarning" class="fea-idle-warn">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width:14px;height:14px;flex-shrink:0"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
        <span>Agent will go to sleep soon.</span>
        <button class="fea-keepalive-btn" @click="emit('heartbeat')">Keep Alive</button>
      </div>

      <!-- Directory listing, breadcrumb, toolbar -->
      <FileExplorer
        :entries="entries"
        :current-path="currentPath"
        :loading="loading"
        :uploading="uploading"
        :readonly="readonly"
        class="fea-explorer"
        @navigate="navigateTo"
        @file-opened="openFile"
        @delete="handleDelete"
        @rename="handleRename"
        @upload="handleUpload"
        @new-file="handleNewFile"
        @new-folder="handleNewFolder"
        @download="handleDownload"
        @delete-many="handleDeleteMany"
      />

      <!-- Non-editable / binary / image / pdf preview -->
      <FilePreview
        v-if="activeFile && showPreview"
        :file="activeFile"
        :download-url="resolvedDownloadUrl"
        :too-large="tooLarge"
        class="fea-preview"
        @download="handleDownload"
      />

      <!-- Monaco editor (kept mounted via v-show to preserve tab state) -->
      <div class="fea-editor" v-show="showEditor">
        <EditorPanel
          ref="editorRef"
          :readonly="readonly"
          @save="handleSave"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import IdleTimer from '../shared/IdleTimer.vue'
import FileTree from './FileTree.vue'
import FileExplorer from './FileExplorer.vue'
import FilePreview from './FilePreview.vue'
import EditorPanel from './EditorPanel.vue'
import { isEditable as isEditableExt } from '../../utils/lang'
import type { FileEntry } from '../../api/files'

const props = defineProps<{
  // API function props — components never call fetch directly
  fetchFiles:   (path: string) => Promise<{ entries: FileEntry[] }>
  fetchContent: (path: string) => Promise<string>
  saveFile:     (path: string, content: string) => Promise<void>
  deleteFile:   (path: string) => Promise<void>
  renameFile:   (from: string, to: string) => Promise<void>
  uploadFiles:  (path: string, files: File[]) => Promise<void>
  createFile:   (path: string) => Promise<void>
  downloadUrl:  (path: string) => string
  readonly: boolean
  explorerLabel?: string
  remainingSeconds: number | null
  idleWarning: boolean
  disconnected?: boolean
  reconnecting?: boolean
}>()

const emit = defineEmits<{
  (e: 'heartbeat'): void
  (e: 'reconnect'): void
}>()

// ── Shared state ────────────────────────────────────────────────────────────
const currentPath = ref('')
const entries     = ref<FileEntry[]>([])
const activeFile  = ref<FileEntry | null>(null)
const loading     = ref(false)
const uploading   = ref(false)
const tooLarge    = ref(false)
const editorRef   = ref<InstanceType<typeof EditorPanel> | null>(null)

// ── Computed helpers ─────────────────────────────────────────────────────────
const ext      = computed(() => activeFile.value?.name.split('.').pop()?.toLowerCase() ?? '')
const isImage  = computed(() => ['png','jpg','jpeg','gif','svg','webp'].includes(ext.value))
const isPdf    = computed(() => ext.value === 'pdf')
const canEdit  = computed(() => !!activeFile.value && isEditableExt(activeFile.value.name))

/** Show the preview panel (image, pdf, binary, too-large) */
const showPreview = computed(() =>
  !!activeFile.value && (isImage.value || isPdf.value || tooLarge.value || !canEdit.value)
)
/** Show the editor panel (text files, or empty state when no file is selected) */
const showEditor = computed(() => !showPreview.value)

/** Full proxy URL for the currently active file */
const resolvedDownloadUrl = computed(() => {
  if (!activeFile.value) return ''
  const p = currentPath.value
    ? `${currentPath.value}/${activeFile.value.name}`
    : activeFile.value.name
  return props.downloadUrl(p)
})

// ── Navigation ───────────────────────────────────────────────────────────────
async function navigateTo(path: string) {
  currentPath.value = path
  activeFile.value  = null
  await loadEntries(path)
}

/** Open a file clicked in the tree sidebar */
async function handleOpenFileFromTree(filePath: string, fileName: string) {
  // Navigate to the parent directory if we're not already there
  const parts = filePath.split('/')
  parts.pop()
  const parentDir = parts.join('/')
  if (currentPath.value !== parentDir) {
    currentPath.value = parentDir
    activeFile.value  = null
    await loadEntries(parentDir)
  }
  // Open the file (size unknown from tree — content fetch will determine)
  const entry: FileEntry = { name: fileName, size: 0, modTime: '', isDir: false }
  await openFile(entry)
}

async function loadEntries(path: string) {
  loading.value = true
  try {
    const { entries: e } = await props.fetchFiles(path)
    entries.value = e
  } catch (err) {
    console.error('Failed to load entries', err)
    entries.value = []
  } finally {
    loading.value = false
  }
}

// ── File open ─────────────────────────────────────────────────────────────────
async function openFile(entry: FileEntry) {
  activeFile.value = entry
  tooLarge.value   = entry.size > 1_048_576

  // For images, pdfs, binary or oversized files skip content fetch
  if (tooLarge.value || isImage.value || isPdf.value || !isEditableExt(entry.name)) return

  const filePath = currentPath.value ? `${currentPath.value}/${entry.name}` : entry.name
  const content  = await props.fetchContent(filePath)
  editorRef.value?.openFile({ path: filePath, name: entry.name, content })
}

// ── File operations (call API function props, then refresh) ──────────────────
async function handleSave(path: string, content: string) {
  await props.saveFile(path, content)
}

async function handleDelete(path: string) {
  await props.deleteFile(path)
  // Close active file if it was deleted
  if (activeFile.value) {
    const activePath = currentPath.value
      ? `${currentPath.value}/${activeFile.value.name}`
      : activeFile.value.name
    if (activePath === path) activeFile.value = null
  }
  await loadEntries(currentPath.value)
}

async function handleRename(from: string, to: string) {
  await props.renameFile(from, to)
  await loadEntries(currentPath.value)
}

async function handleUpload(files: File[]) {
  uploading.value = true
  try {
    await props.uploadFiles(currentPath.value, files)
    await loadEntries(currentPath.value)
  } finally {
    uploading.value = false
  }
}

async function handleNewFile(path: string) {
  await props.createFile(path)
  await loadEntries(currentPath.value)
}

async function handleDeleteMany(paths: string[]) {
  for (const path of paths) {
    await props.deleteFile(path)
    if (activeFile.value) {
      const activePath = currentPath.value
        ? `${currentPath.value}/${activeFile.value.name}`
        : activeFile.value.name
      if (activePath === path) activeFile.value = null
    }
  }
  await loadEntries(currentPath.value)
}

async function handleNewFolder(keepPath: string) {
  await props.createFile(keepPath) // .keep file creates the dir
  await loadEntries(currentPath.value)
}

function handleDownload(entry: FileEntry) {
  const path = currentPath.value ? `${currentPath.value}/${entry.name}` : entry.name
  window.open(props.downloadUrl(path))
}

// ── Initial load + retry (agent proxy may take a moment) ─────────────────────
onMounted(async () => {
  await loadEntries('')
  // Retry up to 5× if directory comes back empty (agent not ready yet)
  for (let i = 0; i < 5 && !entries.value.length; i++) {
    await new Promise(r => setTimeout(r, 1500))
    await loadEntries(currentPath.value)
  }
})
</script>

<style scoped>
/* ── Dark theme vars cascade to FileTree, FileExplorer and all children ── */
.fea-layout {
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
  height: calc(100vh - 80px);
  min-height: 500px;
  overflow: hidden;
  background: var(--bg);
  color: var(--text);
  font-family: var(--font);
  position: relative;
}

.fea-tree { flex-shrink: 0; }

.fea-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  overflow: hidden;
}

.fea-topbar {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 5px 12px;
  background: var(--surface-ground);
  border-bottom: 1px solid var(--surface-border);
  flex-shrink: 0;
}

.fea-label {
  font-size: var(--fs-xs, 0.8125rem);
  font-weight: 600;
  color: var(--text-color-secondary);
  font-family: Lato, sans-serif;
  letter-spacing: 0.04em;
}

/* ── Idle warning banner ── */
.fea-idle-warn {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 7px 14px;
  background: rgba(245,158,11,.1);
  border-bottom: 1px solid rgba(245,158,11,.3);
  color: var(--warn);
  font-family: var(--font);
  font-size: 12px;
  flex-shrink: 0;
}
.fea-keepalive-btn {
  margin-left: auto;
  background: rgba(245,158,11,.15);
  border: 1px solid rgba(245,158,11,.4);
  color: var(--warn);
  border-radius: 4px;
  padding: 3px 10px;
  font-family: var(--font);
  font-size: 11px;
  cursor: pointer;
  transition: background .15s;
}
.fea-keepalive-btn:hover { background: rgba(245,158,11,.25); }

.fea-explorer {
  flex: 0 0 42%;
  min-height: 0;
  overflow: hidden;
  border-bottom: 1px solid var(--border);
}

.fea-preview {
  flex: 1;
  overflow: auto;
  background: var(--surface-ground);
}

.fea-editor {
  flex: 1;
  min-height: 0;
}

/* ── Disconnected overlay ── */
.fea-disconnected {
  position: absolute;
  inset: 0;
  z-index: 50;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.55);
  backdrop-filter: blur(3px);
}
.fea-disconnected__card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding: 2.5rem 3rem;
  background: var(--surface-card);
  border: 1px solid var(--surface-border);
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0,0,0,0.4);
  text-align: center;
}
.fea-disconnected__icon {
  font-size: 2.5rem;
  color: var(--p-red-400);
  margin-bottom: 0.5rem;
}
.fea-disconnected__title {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--text-color);
  margin: 0;
  font-family: Lato, sans-serif;
}
.fea-disconnected__sub {
  font-size: 0.875rem;
  color: var(--text-color-secondary);
  margin: 0;
  font-family: Lato, sans-serif;
}
.fea-disconnected__btn {
  margin-top: 0.75rem;
  padding: 8px 24px;
  font-size: 0.9rem;
  font-weight: 600;
  font-family: Lato, sans-serif;
  background: var(--p-primary-500);
  color: #fff;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: background 0.15s, opacity 0.15s;
  display: flex;
  align-items: center;
}
.fea-disconnected__btn:hover:not(:disabled) { background: var(--p-primary-600); }
.fea-disconnected__btn:disabled { opacity: 0.6; cursor: not-allowed; }
</style>
