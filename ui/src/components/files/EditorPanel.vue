<template>
  <div class="ep-root">

    <!-- ── Single-pane layout ─────────────────────────────────────── -->
    <template v-if="!splitView">
      <!-- Tab bar -->
      <div class="ep-tabs" v-if="tabs.length">
        <button
          v-for="tab in tabs"
          :key="tab.path"
          class="ep-tab"
          :class="{ 'ep-tab--active': tab.path === activeTabPath }"
          @click="activateTab(tab.path)"
          @mousedown.middle.prevent="closeTab(tab.path)"
        >
          <span class="ep-tab__dot" v-if="tab.dirty" title="Unsaved changes" />
          <span class="ep-tab__name">{{ tab.name }}</span>
          <span role="button" tabindex="0" class="ep-tab__close" @click.stop="closeTab(tab.path)" @keydown.enter.stop="closeTab(tab.path)" @keydown.space.stop="closeTab(tab.path)" title="Close tab">×</span>
        </button>
        <div class="ep-tabs__spacer" />
        <!-- Settings hamburger -->
        <div ref="settingsWrapEl" class="ep-settings-wrap">
          <button
            ref="settingsBtnEl"
            class="ep-toolbar-btn"
            style="anchor-name: --ep-settings"
            title="Editor settings"
            :aria-expanded="settingsOpen"
            aria-controls="ep-settings-menu"
            @click.stop="toggleSettings"
          >
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="3" y1="6"  x2="21" y2="6"/>
              <line x1="3" y1="12" x2="21" y2="12"/>
              <line x1="3" y1="18" x2="21" y2="18"/>
            </svg>
          </button>
        </div>
        <!-- Split editor button — only when 2+ tabs are open -->
        <button v-if="tabs.length >= 2" class="ep-toolbar-btn" title="Split editor" @click="toggleSplit">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="2" y="3" width="9" height="18" rx="1"/><rect x="13" y="3" width="9" height="18" rx="1"/>
          </svg>
        </button>
      </div>

      <template v-if="activeTab">
        <div class="ep-breadcrumb">
          <span class="ep-breadcrumb__path">{{ activeTab.path }}</span>
          <span v-if="activeTab.dirty" class="ep-breadcrumb__unsaved">● unsaved</span>
          <button
            v-if="!readonly && activeTab.dirty"
            class="ep-save-btn"
            :disabled="saving"
            @click="saveTab(activeTabPath!)"
            title="Save (Ctrl+S)"
          >{{ saving ? 'Saving…' : 'Save' }}</button>
          <span v-if="savedMsg" class="ep-breadcrumb__saved">✓ saved</span>
          <span v-if="saveError" class="ep-save-error">{{ saveError }}</span>
        </div>
        <div class="ep-editor-wrap">
          <VueMonacoEditor
            v-show="t.path === activeTabPath"
            v-for="t in tabs"
            :key="t.path"
            v-model:value="t.content"
            :language="t.language"
            :theme="monacoTheme"
            :options="editorOptions"
            style="height: 100%; width: 100%"
            @change="markDirty(t.path)"
            @mount="(editor, m) => onEditorMount(editor, m, t.path)"
          />
        </div>
      </template>

      <div v-else class="ep-empty">
        <span>Open a file to start editing</span>
      </div>
    </template>

    <!-- ── Split-pane layout ──────────────────────────────────────── -->
    <div v-else class="ep-split-wrap">

      <!-- Left pane -->
      <div class="ep-pane" @click.capture="focusedPane = 'left'">
        <div class="ep-pane-tabs">
          <button
            v-for="t in tabs"
            :key="t.path"
            class="ep-tab ep-tab--sm"
            :class="{ 'ep-tab--active': t.path === leftTabPath }"
            @click="leftTabPath = t.path"
            @mousedown.middle.prevent="closeTab(t.path)"
          >
            <span class="ep-tab__dot" v-if="t.dirty" />
            <span class="ep-tab__name">{{ t.name }}</span>
            <span role="button" tabindex="0" class="ep-tab__close" @click.stop="closeTab(t.path)" @keydown.enter.stop="closeTab(t.path)" @keydown.space.stop="closeTab(t.path)">×</span>
          </button>
          <div class="ep-tabs__spacer" />
          <button class="ep-toolbar-btn" title="Close split" @click="toggleSplit">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
        <div class="ep-pane-bar" v-if="leftTab">
          <span class="ep-breadcrumb__path">{{ leftTab.path }}</span>
          <span v-if="leftTab.dirty" class="ep-breadcrumb__unsaved">● unsaved</span>
          <button v-if="!readonly && leftTab.dirty" class="ep-save-btn" :disabled="saving" @click="saveTab(leftTabPath!)">
            {{ saving ? 'Saving…' : 'Save' }}
          </button>
        </div>
        <div class="ep-editor-wrap">
          <VueMonacoEditor
            v-show="t.path === leftTabPath"
            v-for="t in tabs"
            :key="'L:' + t.path"
            v-model:value="t.content"
            :language="t.language"
            :theme="monacoTheme"
            :options="splitEditorOptions"
            style="height: 100%; width: 100%"
            @change="markDirty(t.path)"
            @mount="(editor, m) => onEditorMount(editor, m, 'L:' + t.path)"
          />
        </div>
      </div>

      <div class="ep-pane-divider" />

      <!-- Right pane -->
      <div class="ep-pane" @click.capture="focusedPane = 'right'">
        <div class="ep-pane-tabs">
          <button
            v-for="t in tabs"
            :key="t.path"
            class="ep-tab ep-tab--sm"
            :class="{ 'ep-tab--active': t.path === rightTabPath }"
            @click="rightTabPath = t.path"
            @mousedown.middle.prevent="closeTab(t.path)"
          >
            <span class="ep-tab__dot" v-if="t.dirty" />
            <span class="ep-tab__name">{{ t.name }}</span>
            <span role="button" tabindex="0" class="ep-tab__close" @click.stop="closeTab(t.path)" @keydown.enter.stop="closeTab(t.path)" @keydown.space.stop="closeTab(t.path)">×</span>
          </button>
        </div>
        <div class="ep-pane-bar" v-if="rightTab">
          <span class="ep-breadcrumb__path">{{ rightTab.path }}</span>
          <span v-if="rightTab.dirty" class="ep-breadcrumb__unsaved">● unsaved</span>
          <button v-if="!readonly && rightTab.dirty" class="ep-save-btn" :disabled="saving" @click="saveTab(rightTabPath!)">
            {{ saving ? 'Saving…' : 'Save' }}
          </button>
        </div>
        <div class="ep-editor-wrap">
          <VueMonacoEditor
            v-show="t.path === rightTabPath"
            v-for="t in tabs"
            :key="'R:' + t.path"
            v-model:value="t.content"
            :language="t.language"
            :theme="monacoTheme"
            :options="splitEditorOptions"
            style="height: 100%; width: 100%"
            @change="markDirty(t.path)"
            @mount="(editor, m) => onEditorMount(editor, m, 'R:' + t.path)"
          />
        </div>
      </div>

    </div>

  </div>

  <!-- Settings menu — teleported to body to escape overflow:hidden tab bar -->
  <Teleport to="body">
    <div
      v-if="settingsOpen"
      ref="settingsMenuEl"
      id="ep-settings-menu"
      class="ep-settings-menu"
      :style="anchorSupported ? {} : { top: menuY + 'px', right: menuRight + 'px' }"
      @click.stop
    >
      <div class="ep-settings-title">Editor settings</div>

      <!-- Theme -->
      <div class="ep-settings-row ep-settings-row--col">
        <span>Theme</span>
        <div class="ep-settings-btn-group">
          <button
            v-for="opt in themeOptions" :key="opt.value"
            class="ep-settings-chip"
            :class="{ 'ep-settings-chip--active': editorSettings.theme === opt.value }"
            @click="editorSettings.theme = opt.value"
          >{{ opt.label }}</button>
        </div>
      </div>

      <!-- Font size -->
      <div class="ep-settings-row ep-settings-row--col">
        <span>Font size</span>
        <div class="ep-settings-btn-group">
          <button
            v-for="sz in fontSizeOptions" :key="sz"
            class="ep-settings-chip"
            :class="{ 'ep-settings-chip--active': editorSettings.fontSize === sz }"
            @click="editorSettings.fontSize = sz"
          >{{ sz }}</button>
        </div>
      </div>

      <!-- Tab size -->
      <div class="ep-settings-row">
        <span>Tab size</span>
        <div class="ep-settings-btn-group">
          <button
            v-for="sz in [2, 4]" :key="sz"
            class="ep-settings-chip"
            :class="{ 'ep-settings-chip--active': editorSettings.tabSize === sz }"
            @click="editorSettings.tabSize = sz"
          >{{ sz }}</button>
        </div>
      </div>

      <div class="ep-settings-sep" />

      <!-- Whitespace -->
      <div class="ep-settings-row ep-settings-row--col">
        <span>Render whitespace</span>
        <div class="ep-settings-btn-group">
          <button
            v-for="opt in whitespaceOptions" :key="opt.value"
            class="ep-settings-chip"
            :class="{ 'ep-settings-chip--active': editorSettings.renderWhitespace === opt.value }"
            @click="editorSettings.renderWhitespace = opt.value"
          >{{ opt.label }}</button>
        </div>
      </div>

      <div class="ep-settings-sep" />

      <label class="ep-settings-row">
        <span>Word wrap</span>
        <input type="checkbox" :checked="editorSettings.wordWrap === 'on'" @change="toggleWordWrap" />
      </label>
      <label class="ep-settings-row">
        <span>Line numbers</span>
        <input type="checkbox" :checked="editorSettings.lineNumbers === 'on'" @change="toggleLineNumbers" />
      </label>
      <label class="ep-settings-row">
        <span>Minimap</span>
        <input type="checkbox" v-model="editorSettings.minimap" />
      </label>
    </div>
  </Teleport>

</template>

<script setup lang="ts">
import { ref, computed, reactive, watch, onMounted, onUnmounted } from 'vue'
import { VueMonacoEditor } from '@guolao/vue-monaco-editor'
import { useLayout } from '@/layout/composables/layout'
import { detectLang } from '../../utils/lang'

interface Tab {
  path: string
  name: string
  content: string
  originalContent: string
  language: string
  dirty: boolean
}

interface OpenFileArgs {
  path: string
  name: string
  content: string
}

const props = withDefaults(defineProps<{
  readonly?: boolean
}>(), { readonly: false })

const emit = defineEmits<{
  (e: 'save', path: string, content: string): void
}>()

// ── Dark mode detection ───────────────────────────────────────────────────────
const { isDarkTheme } = useLayout()

const monacoTheme = computed(() => {
  if (editorSettings.theme === 'dark')  return 'vs-dark'
  if (editorSettings.theme === 'light') return 'vs'
  return isDarkTheme.value ? 'vs-dark' : 'vs'
})

// ── Editor settings (hamburger menu) ─────────────────────────────────────────
const settingsOpen   = ref(false)
const settingsBtnEl  = ref<HTMLElement | null>(null)
const settingsWrapEl = ref<HTMLElement | null>(null)
const settingsMenuEl = ref<HTMLElement | null>(null)
const anchorSupported = CSS.supports('anchor-name', '--x')
const menuY     = ref(0)
const menuRight = ref(0)

function toggleSettings() {
  if (!anchorSupported && !settingsOpen.value) {
    const rect = settingsBtnEl.value?.getBoundingClientRect()
    if (rect) {
      menuY.value     = rect.bottom + 4
      menuRight.value = window.innerWidth - rect.right
    }
  }
  settingsOpen.value = !settingsOpen.value
}

const editorSettings = reactive({
  wordWrap:         'off' as 'off' | 'on',
  minimap:          true,
  lineNumbers:      'on'  as 'on' | 'off',
  theme:            'auto' as 'auto' | 'light' | 'dark',
  fontSize:         13,
  tabSize:          2,
  renderWhitespace: 'boundary' as 'none' | 'boundary' | 'all',
})

const themeOptions      = [
  { label: 'Auto',  value: 'auto'  as const },
  { label: 'Light', value: 'light' as const },
  { label: 'Dark',  value: 'dark'  as const },
]
const fontSizeOptions   = [12, 13, 14, 16]
const whitespaceOptions = [
  { label: 'None',     value: 'none'     as const },
  { label: 'Boundary', value: 'boundary' as const },
  { label: 'All',      value: 'all'      as const },
]

function toggleWordWrap()    { editorSettings.wordWrap    = editorSettings.wordWrap    === 'on' ? 'off' : 'on' }
function toggleLineNumbers() { editorSettings.lineNumbers = editorSettings.lineNumbers === 'on' ? 'off' : 'on' }

// Close settings panel when clicking outside
function onDocClick(e: MouseEvent) {
  const t = e.target as Node
  if (
    settingsWrapEl.value?.contains(t) ||
    settingsMenuEl.value?.contains(t)
  ) return
  settingsOpen.value = false
}
onMounted(() => document.addEventListener('mousedown', onDocClick))
onUnmounted(() => document.removeEventListener('mousedown', onDocClick))

// Propagate settings changes to all mounted editor instances
watch(editorSettings, () => {
  const opts = {
    wordWrap:         editorSettings.wordWrap,
    lineNumbers:      editorSettings.lineNumbers,
    minimap:          { enabled: editorSettings.minimap },
    fontSize:         editorSettings.fontSize,
    tabSize:          editorSettings.tabSize,
    renderWhitespace: editorSettings.renderWhitespace,
  }
  for (const editor of editorInstances.values()) {
    editor.updateOptions(opts)
  }
})

// ── Tabs state ────────────────────────────────────────────────────────────────
const tabs          = ref<Tab[]>([])
const activeTabPath = ref<string | null>(null)
const saving        = ref(false)
const savedMsg      = ref(false)
const saveError     = ref<string | null>(null)

const activeTab = computed(() => tabs.value.find(t => t.path === activeTabPath.value) ?? null)

// ── Split view state ──────────────────────────────────────────────────────────
const splitView    = ref(false)
const leftTabPath  = ref<string | null>(null)
const rightTabPath = ref<string | null>(null)
const focusedPane  = ref<'left' | 'right'>('left')

const leftTab  = computed(() => tabs.value.find(t => t.path === leftTabPath.value)  ?? null)
const rightTab = computed(() => tabs.value.find(t => t.path === rightTabPath.value) ?? null)

function toggleSplit() {
  splitView.value = !splitView.value
  if (splitView.value) {
    leftTabPath.value  = activeTabPath.value
    rightTabPath.value = tabs.value.find(t => t.path !== activeTabPath.value)?.path ?? activeTabPath.value
  }
}

// ── Editor options ─────────────────────────────────────────────────────────────
const baseOptions = computed(() => ({
  fontSize:             editorSettings.fontSize,
  fontFamily:           "'JetBrains Mono', 'Fira Mono', monospace",
  lineHeight:           Math.round(editorSettings.fontSize * 1.7),
  scrollBeyondLastLine: false,
  wordWrap:             editorSettings.wordWrap,
  lineNumbers:          editorSettings.lineNumbers,
  tabSize:              editorSettings.tabSize,
  insertSpaces:         true,
  renderWhitespace:     editorSettings.renderWhitespace,
  bracketPairColorization: { enabled: true },
}))
const editorOptions = computed(() => ({
  ...baseOptions.value,
  minimap: { enabled: editorSettings.minimap },
  readOnly: props.readonly,
}))
const splitEditorOptions = computed(() => ({
  ...baseOptions.value,
  minimap: { enabled: false },
  readOnly: props.readonly,
}))

// ── Editor instances (keyed by tab path or 'L:path' / 'R:path') ───────────────
const editorInstances = new Map<string, any>()

function onEditorMount(editor: any, monaco: any, key: string) {
  editorInstances.set(key, editor)
  editor.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyS, () => {
    if (props.readonly) return
    if (key.startsWith('L:'))      saveTab(leftTabPath.value!)
    else if (key.startsWith('R:')) saveTab(rightTabPath.value!)
    else                           saveTab(activeTabPath.value!)
  })
}

// ── Tab management ─────────────────────────────────────────────────────────────
function openFile({ path, name, content }: OpenFileArgs) {
  const existing = tabs.value.find(t => t.path === path)
  if (existing) {
    if (!existing.dirty) {
      existing.content         = content
      existing.originalContent = content
    }
    _activateInPane(path)
    return
  }
  tabs.value.push({ path, name, content, originalContent: content, language: detectLang(name), dirty: false })
  _activateInPane(path)
}

function _activateInPane(path: string) {
  activeTabPath.value = path
  if (splitView.value) {
    if (focusedPane.value === 'right') rightTabPath.value = path
    else                               leftTabPath.value  = path
  }
}

function activateTab(path: string) {
  activeTabPath.value = path
}

function closeTab(path: string) {
  const idx = tabs.value.findIndex(t => t.path === path)
  if (idx === -1) return
  tabs.value.splice(idx, 1)
  const fallback = tabs.value[Math.min(idx, tabs.value.length - 1)]?.path ?? null
  if (activeTabPath.value  === path) activeTabPath.value  = fallback
  if (leftTabPath.value    === path) leftTabPath.value    = fallback
  if (rightTabPath.value   === path) {
    rightTabPath.value = tabs.value.find(t => t.path !== leftTabPath.value)?.path ?? leftTabPath.value
  }
  if (tabs.value.length < 2) splitView.value = false
}

function markDirty(path: string) {
  const tab = tabs.value.find(t => t.path === path)
  if (tab) tab.dirty = true
}

async function saveTab(path: string) {
  const tab = tabs.value.find(t => t.path === path)
  if (!tab || props.readonly) return
  saving.value    = true
  saveError.value = null
  try {
    emit('save', tab.path, tab.content)
    tab.originalContent = tab.content
    tab.dirty           = false
    savedMsg.value      = true
    setTimeout(() => { savedMsg.value = false }, 2000)
  } catch (e: any) {
    saveError.value = e?.message ?? 'Save failed'
  } finally {
    saving.value = false
  }
}

defineExpose({ openFile, closeTab })
</script>

<style scoped>
/* ── Root ──────────────────────────────────────────────────────────── */
.ep-root {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--bg);
  overflow: hidden;
}

/* ── Tab bars ──────────────────────────────────────────────────────── */
.ep-tabs,
.ep-pane-tabs {
  display: flex;
  align-items: center;
  overflow-x: auto;
  background: var(--surface);
  border-bottom: 1px solid var(--border);
  flex-shrink: 0;
}

.ep-tabs__spacer { flex: 1; min-width: 4px; }

.ep-tab {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 6px 14px 6px 10px;
  border: none;
  border-bottom: 2px solid transparent;
  background: transparent;
  color: var(--muted);
  cursor: pointer;
  font-size: 12px;
  font-family: var(--font);
  white-space: nowrap;
  transition: background 0.12s, color 0.12s;
}
.ep-tab:hover { background: var(--surface2); color: var(--text); }
.ep-tab--active { color: var(--text); background: var(--bg); border-bottom-color: var(--accent); }
.ep-tab--sm { padding: 4px 10px 4px 8px; font-size: 11px; }

.ep-tab__dot {
  width: 7px; height: 7px; border-radius: 50%;
  background: var(--warn, #f0a500); flex-shrink: 0;
}
.ep-tab__name { max-width: 130px; overflow: hidden; text-overflow: ellipsis; }
.ep-tab__close {
  background: transparent; border: none; color: var(--muted);
  cursor: pointer; font-size: 14px; line-height: 1; padding: 0 2px;
}
.ep-tab__close:hover { color: var(--text); }

/* ── Toolbar button (split / close-split / settings) ──────────────── */
.ep-toolbar-btn {
  flex-shrink: 0;
  display: flex; align-items: center; justify-content: center;
  width: 28px; height: 28px;
  background: transparent; border: none; cursor: pointer;
  color: var(--muted); border-radius: 4px; margin: 0 2px;
  transition: background 0.12s, color 0.12s;
}
.ep-toolbar-btn svg { width: 14px; height: 14px; }
.ep-toolbar-btn:hover { background: var(--surface2); color: var(--text); }

/* ── Breadcrumb / pane bar ────────────────────────────────────────── */
.ep-breadcrumb,
.ep-pane-bar {
  display: flex; align-items: center; gap: 8px;
  padding: 4px 12px;
  background: var(--surface2);
  font-size: 12px;
  color: var(--muted);
  flex-shrink: 0;
  border-bottom: 1px solid var(--border);
}
.ep-breadcrumb__path { font-family: monospace; color: var(--text); }
.ep-breadcrumb__unsaved { color: var(--warn, #f0a500); }
.ep-breadcrumb__saved  { color: #73c991; }
.ep-save-error         { color: var(--danger); }

.ep-save-btn {
  margin-left: auto;
  padding: 2px 10px;
  border: 1px solid var(--border);
  background: var(--surface);
  color: var(--text);
  border-radius: 3px;
  cursor: pointer;
  font-size: 12px;
  font-family: var(--font);
  transition: background 0.12s, border-color 0.12s;
}
.ep-save-btn:hover:not(:disabled) { background: var(--accent); color: #fff; border-color: var(--accent); }
.ep-save-btn:disabled { opacity: 0.5; cursor: not-allowed; }

/* ── Editor wrap ──────────────────────────────────────────────────── */
.ep-editor-wrap {
  flex: 1;
  min-height: 0;
  position: relative;
}

/* ── Empty state ──────────────────────────────────────────────────── */
.ep-empty {
  flex: 1;
  display: flex; align-items: center; justify-content: center;
  color: var(--muted);
  font-size: 14px;
  font-family: var(--font);
}

/* ── Split-pane layout ────────────────────────────────────────────── */
.ep-split-wrap {
  flex: 1;
  display: flex;
  min-height: 0;
  overflow: hidden;
}
.ep-pane {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  overflow: hidden;
}
.ep-pane-divider {
  width: 1px;
  background: var(--border);
  flex-shrink: 0;
}

/* ── Settings wrap (anchor for button only) ───────────────────────── */
.ep-settings-wrap {
  flex-shrink: 0;
  display: flex;
  align-items: center;
}
</style>

<!-- Settings menu styles are NOT scoped — needed because the menu is teleported to body -->
<style>
.ep-settings-menu {
  position: fixed;
  position-anchor: --ep-settings;
  position-area: block-end span-inline-end;
  position-try-fallbacks: flip-block, flip-inline;
  margin-block-start: 4px;
  z-index: 9999;
  min-width: 230px;
  background: var(--surface-card);
  border: 1px solid var(--surface-border);
  border-radius: 8px;
  padding: 8px 0 10px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.18);
  animation: ep-settings-pop .12s ease;
  font-family: Lato, sans-serif;
}
@keyframes ep-settings-pop {
  from { opacity: 0; transform: scale(.95) translateY(-4px) }
  to   { opacity: 1; transform: none }
}
.ep-settings-title {
  padding: 4px 14px 8px;
  font-size: 11px;
  font-weight: 600;
  letter-spacing: .06em;
  text-transform: uppercase;
  color: var(--text-color-secondary);
  font-family: Lato, sans-serif;
}
.ep-settings-sep {
  height: 1px;
  background: var(--surface-border);
  margin: 4px 0;
}
.ep-settings-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 14px;
  font-size: 13px;
  color: var(--text-color);
  cursor: pointer;
  gap: 12px;
  font-family: Lato, sans-serif;
}
.ep-settings-row:hover { background: var(--surface-hover); }
.ep-settings-row--col {
  flex-direction: column;
  align-items: flex-start;
  gap: 8px;
  cursor: default;
}
.ep-settings-row--col:hover { background: transparent; }
.ep-settings-row input[type="checkbox"] {
  accent-color: var(--p-primary-500);
  cursor: pointer;
  width: 15px; height: 15px;
  flex-shrink: 0;
}
/* ── Chip button group ─────────────────────────────────────────────── */
.ep-settings-btn-group {
  display: flex;
  gap: 4px;
}
.ep-settings-chip {
  padding: 2px 10px;
  border: 1px solid var(--surface-border);
  background: transparent;
  color: var(--text-color-secondary);
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  font-family: Lato, sans-serif;
  transition: background 0.1s, color 0.1s, border-color 0.1s;
}
.ep-settings-chip:hover {
  background: var(--surface-hover);
  color: var(--text-color);
}
.ep-settings-chip--active {
  background: var(--p-primary-500);
  color: #fff;
  border-color: var(--p-primary-500);
}
</style>
