<template>
  <aside class="filter-sidebar">
    <!-- Search -->
    <div class="sidebar-section">
      <span class="section-header">
        <span>Search</span>
      </span>
      <div class="filter-search-wrap">
        <IconField>
          <InputIcon class="pi pi-search" />
          <label for="filter-search" class="sr-only">Search explorers</label>
          <InputText
            id="filter-search"
            v-model="search"
            placeholder="Search name, PVC..."
            class="w-full filter-search-input"
            @input="emit_()"
          />
        </IconField>
      </div>
    </div>

    <Divider class="sidebar-divider" />

    <!-- Phase filter -->
    <div class="sidebar-section">
      <button class="section-header" @click="open.phase = !open.phase">
        <span>Phase</span>
        <i :class="['pi', open.phase ? 'pi-chevron-down' : 'pi-chevron-right']" />
      </button>
      <div v-if="open.phase" class="filter-options">
        <label
          v-for="opt in phaseOptions"
          :key="opt.value"
          class="filter-option"
        >
          <Checkbox v-model="phases" :value="opt.value" @change="emit_()" />
          <span class="filter-dot" :style="{ background: opt.color }" />
          <Tag :value="opt.label" :severity="opt.severity" rounded class="opt-tag" />
          <span class="opt-count">{{ phaseCounts[opt.value] ?? 0 }}</span>
        </label>
      </div>
    </div>

    <Divider class="sidebar-divider" />

    <!-- Namespace filter -->
    <div class="sidebar-section">
      <button class="section-header" @click="open.namespace = !open.namespace">
        <span>Namespace</span>
        <i :class="['pi', open.namespace ? 'pi-chevron-down' : 'pi-chevron-right']" />
      </button>
      <div v-if="open.namespace" class="filter-options">
        <label
          v-for="ns in namespaceOptions"
          :key="ns"
          class="filter-option"
        >
          <Checkbox v-model="namespaces" :value="ns" @change="emit_()" />
          <span class="filter-dot" :style="{ background: stringToColor(ns) }" />
          <span class="opt-label">{{ ns }}</span>
          <span class="opt-count">{{ namespaceCounts[ns] ?? 0 }}</span>
        </label>
      </div>
    </div>

    <Divider class="sidebar-divider" />

    <!-- Scope filter -->
    <div class="sidebar-section">
      <button class="section-header" @click="open.scope = !open.scope">
        <span>Scope</span>
        <i :class="['pi', open.scope ? 'pi-chevron-down' : 'pi-chevron-right']" />
      </button>
      <div v-if="open.scope" class="filter-options">
        <label
          v-for="sc in scopeOptions"
          :key="sc"
          class="filter-option"
        >
          <Checkbox v-model="scopes" :value="sc" @change="emit_()" />
          <span class="filter-dot" :style="{ background: stringToColor(sc) }" />
          <span class="opt-label">{{ sc }}</span>
          <span class="opt-count">{{ scopeCounts[sc] ?? 0 }}</span>
        </label>
        <div v-if="!scopeOptions.length" class="text-muted-color text-sm">None</div>
      </div>
    </div>

    <Divider class="sidebar-divider" />

    <!-- Mount state filter -->
    <div class="sidebar-section">
      <button class="section-header" @click="open.mount = !open.mount">
        <span>Mount State</span>
        <i :class="['pi', open.mount ? 'pi-chevron-down' : 'pi-chevron-right']" />
      </button>
      <div v-if="open.mount" class="filter-options">
        <label
          v-for="ms in mountStateOptions"
          :key="ms"
          class="filter-option"
        >
          <Checkbox v-model="mountStates" :value="ms" @change="emit_()" />
          <span class="filter-dot" :style="{ background: mountStateColor(ms) }" />
          <Tag :value="ms" :severity="mountSeverity(ms)" rounded class="opt-tag" />
          <span class="opt-count">{{ mountCounts[ms] ?? 0 }}</span>
        </label>
      </div>
    </div>

    <Divider class="sidebar-divider" />

    <!-- Access Mode filter -->
    <div class="sidebar-section">
      <button class="section-header" @click="open.accessMode = !open.accessMode">
        <span>Access Mode</span>
        <i :class="['pi', open.accessMode ? 'pi-chevron-down' : 'pi-chevron-right']" />
      </button>
      <div v-if="open.accessMode" class="filter-options">
        <label
          v-for="am in accessModeOptions"
          :key="am"
          class="filter-option"
        >
          <Checkbox v-model="accessModes" :value="am" @change="emit_()" />
          <span class="filter-dot" :style="{ background: accessModeColor(am) }" />
          <span class="opt-label">{{ am }}</span>
          <span class="opt-count">{{ accessModeCounts[am] ?? 0 }}</span>
        </label>
        <div v-if="!accessModeOptions.length" class="text-muted-color text-sm">None</div>
      </div>
    </div>

    <Divider class="sidebar-divider" />

    <!-- Consumers filter -->
    <div class="sidebar-section">
      <button class="section-header" @click="open.consumers = !open.consumers">
        <span>Consumers</span>
        <i :class="['pi', open.consumers ? 'pi-chevron-down' : 'pi-chevron-right']" />
      </button>
      <div v-if="open.consumers" class="filter-options">
        <label class="filter-option">
          <RadioButton v-model="consumers" value="" @change="emit_()" />
          <span class="opt-label">Any</span>
        </label>
        <label class="filter-option">
          <RadioButton v-model="consumers" value="has" @change="emit_()" />
          <span class="filter-dot" style="background:#22c55e" />
          <span class="opt-label">Has consumers</span>
          <span class="opt-count">{{ hasConsumerCount }}</span>
        </label>
        <label class="filter-option">
          <RadioButton v-model="consumers" value="none" @change="emit_()" />
          <span class="filter-dot" style="background:#94a3b8" />
          <span class="opt-label">No consumers</span>
          <span class="opt-count">{{ noConsumerCount }}</span>
        </label>
      </div>
    </div>

    <Divider class="sidebar-divider" />

    <!-- Created filter -->
    <div class="sidebar-section">
      <button class="section-header" @click="open.created = !open.created">
        <span>Created</span>
        <i :class="['pi', open.created ? 'pi-chevron-down' : 'pi-chevron-right']" />
      </button>
      <div v-if="open.created" class="filter-options">
        <label class="filter-option">
          <RadioButton v-model="created" value="" @change="emit_()" />
          <span class="opt-label">Any time</span>
        </label>
        <label class="filter-option">
          <RadioButton v-model="created" value="24h" @change="emit_()" />
          <span class="opt-label">Last 24 hours</span>
          <span class="opt-count">{{ createdCounts['24h'] }}</span>
        </label>
        <label class="filter-option">
          <RadioButton v-model="created" value="7d" @change="emit_()" />
          <span class="opt-label">Last 7 days</span>
          <span class="opt-count">{{ createdCounts['7d'] }}</span>
        </label>
        <label class="filter-option">
          <RadioButton v-model="created" value="30d" @change="emit_()" />
          <span class="opt-label">Last 30 days</span>
          <span class="opt-count">{{ createdCounts['30d'] }}</span>
        </label>
        <label class="filter-option">
          <RadioButton v-model="created" value="older" @change="emit_()" />
          <span class="opt-label">Older than 30d</span>
          <span class="opt-count">{{ createdCounts['older'] }}</span>
        </label>
      </div>
    </div>

    <Divider class="sidebar-divider" />

    <!-- Labels filter -->
    <div class="sidebar-section">
      <button class="section-header" @click="open.labels = !open.labels">
        <span>Labels</span>
        <i :class="['pi', open.labels ? 'pi-chevron-down' : 'pi-chevron-right']" />
      </button>
      <div v-if="open.labels" class="filter-options">
        <div class="label-input-wrap">
          <input
            v-model="labelInput"
            placeholder="key=value"
            class="label-input"
            @keydown.enter.prevent="addLabel"
            @keydown.backspace="removeLastLabel"
          />
        </div>
        <div v-if="labels.length" class="label-chips">
          <span
            v-for="l in labels"
            :key="l"
            class="label-tag"
            :style="{ background: stringToColor(l) + '25', color: stringToColor(l), borderColor: stringToColor(l) + '40' }"
          >
            {{ l }}
            <button class="label-remove" @click="removeLabel(l)">&times;</button>
          </span>
        </div>
      </div>
    </div>

    <!-- Active filter summary -->
    <div v-if="activeCount > 0" class="clear-row">
      <span class="shown-count">{{ shown }}/{{ total }} shown</span>
      <Button label="Clear" text size="small" severity="secondary" @click="clearAll" />
    </div>
    <div v-else class="shown-row">
      <span class="shown-count">{{ shown }} agents</span>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import Checkbox from 'primevue/checkbox'
import RadioButton from 'primevue/radiobutton'
import Divider from 'primevue/divider'
import Tag from 'primevue/tag'
import Button from 'primevue/button'
import IconField from 'primevue/iconfield'
import InputIcon from 'primevue/inputicon'
import InputText from 'primevue/inputtext'
import type { Explorer } from '../../stores/explorerStore'
import {
  stringToColor,
  mountStateColor,
  accessModeColor,
  type Filters,
  type ConsumerFilter,
  type CreatedFilter,
} from '../../composables/useFilterColors'

export type { Filters } from '../../composables/useFilterColors'

const props = defineProps<{ explorers: Explorer[]; shown: number; total: number }>()
const emit = defineEmits<{ (e: 'update:filters', v: Filters): void }>()

const search = ref('')
const phases = ref<string[]>([])
const namespaces = ref<string[]>([])
const mountStates = ref<string[]>([])
const scopes = ref<string[]>([])
const accessModes = ref<string[]>([])
const consumers = ref<ConsumerFilter>('')
const created = ref<CreatedFilter>('')
const labels = ref<string[]>([])
const labelInput = ref('')

const open = ref({ phase: true, namespace: true, scope: false, mount: false, accessMode: false, consumers: false, created: false, labels: false })

const phaseOptions = [
  { value: 'Running',      label: 'Running',        severity: 'success',   color: '#22c55e' },
  { value: 'ScaledToZero', label: 'Scaled to Zero',  severity: 'secondary', color: '#94a3b8' },
  { value: 'Waking',       label: 'Waking',          severity: 'info',      color: '#6366f1' },
  { value: 'Pending',      label: 'Pending',         severity: 'warn',      color: '#eab308' },
  { value: 'Failed',       label: 'Failed',          severity: 'danger',    color: '#ef4444' },
  { value: 'InUse',        label: 'In use',          severity: 'warn',      color: '#f59e0b' },
]

/* ------------------- Derived option lists ------------------- */

const namespaceOptions = computed(() => {
  const seen = new Set<string>()
  for (const e of props.explorers) if (e.namespace) seen.add(e.namespace)
  return [...seen].sort()
})

const scopeOptions = computed(() => {
  const seen = new Set<string>()
  for (const e of props.explorers) if (e.scope) seen.add(e.scope)
  return [...seen].sort()
})

const mountStateOptions = computed(() => {
  const seen = new Set<string>()
  for (const e of props.explorers) if (e.mountState) seen.add(e.mountState)
  return [...seen].sort()
})

const accessModeOptions = computed(() => {
  const seen = new Set<string>()
  for (const e of props.explorers) {
    const m = e.accessMode || e.mode
    if (m) seen.add(m)
  }
  return [...seen].sort()
})

/* ------------------- Counts ------------------- */

const phaseCounts = computed(() => {
  const c: Record<string, number> = {}
  for (const e of props.explorers) c[e.phase] = (c[e.phase] ?? 0) + 1
  c.InUse = props.explorers.filter(e => (e.consumerCount ?? 0) > 0).length
  return c
})

const namespaceCounts = computed(() => {
  const c: Record<string, number> = {}
  for (const e of props.explorers) c[e.namespace] = (c[e.namespace] ?? 0) + 1
  return c
})

const scopeCounts = computed(() => {
  const c: Record<string, number> = {}
  for (const e of props.explorers) {
    if (e.scope) c[e.scope] = (c[e.scope] ?? 0) + 1
  }
  return c
})

const mountCounts = computed(() => {
  const c: Record<string, number> = {}
  for (const e of props.explorers) if (e.mountState) c[e.mountState] = (c[e.mountState] ?? 0) + 1
  return c
})

const accessModeCounts = computed(() => {
  const c: Record<string, number> = {}
  for (const e of props.explorers) {
    const m = e.accessMode || e.mode
    if (m) c[m] = (c[m] ?? 0) + 1
  }
  return c
})

const hasConsumerCount = computed(() => props.explorers.filter(e => (e.consumerCount ?? 0) > 0).length)
const noConsumerCount = computed(() => props.explorers.filter(e => !(e.consumerCount ?? 0)).length)

const createdCounts = computed(() => {
  const now = Date.now()
  const c: Record<string, number> = { '24h': 0, '7d': 0, '30d': 0, older: 0 }
  for (const e of props.explorers) {
    if (!e.createdAt) continue
    const age = now - new Date(e.createdAt).getTime()
    if (age < 86_400_000) c['24h']++
    else if (age < 604_800_000) c['7d']++
    else if (age < 2_592_000_000) c['30d']++
    else c['older']++
  }
  return c
})

/* ------------------- Labels ------------------- */

function addLabel() {
  const v = labelInput.value.trim()
  if (v && !labels.value.includes(v)) {
    labels.value = [...labels.value, v]
    emit_()
  }
  labelInput.value = ''
}

function removeLabel(l: string) {
  labels.value = labels.value.filter(x => x !== l)
  emit_()
}

function removeLastLabel() {
  if (!labelInput.value && labels.value.length) {
    labels.value = labels.value.slice(0, -1)
    emit_()
  }
}

/* ------------------- Helpers ------------------- */

function mountSeverity(ms: string) {
  const s = ms.toLowerCase()
  if (s === 'mounted')  return 'success'
  if (s === 'readonly') return 'warning'
  if (s === 'conflict') return 'danger'
  return 'secondary'
}

const activeCount = computed(
  () => phases.value.length + namespaces.value.length + mountStates.value.length
    + scopes.value.length + accessModes.value.length
    + (consumers.value ? 1 : 0) + (created.value ? 1 : 0)
    + labels.value.length + (search.value ? 1 : 0)
)

function emit_() {
  emit('update:filters', {
    search: search.value,
    phases: phases.value,
    namespaces: namespaces.value,
    mountStates: mountStates.value,
    scopes: scopes.value,
    accessModes: accessModes.value,
    consumers: consumers.value,
    created: created.value,
    labels: labels.value,
  })
}

function clearAll() {
  search.value = ''
  phases.value = []
  namespaces.value = []
  mountStates.value = []
  scopes.value = []
  accessModes.value = []
  consumers.value = ''
  created.value = ''
  labels.value = []
  emit_()
}

let debounce: ReturnType<typeof setTimeout> | null = null
watch(search, () => {
  if (debounce) clearTimeout(debounce)
  debounce = setTimeout(emit_, 250)
})
</script>

<style scoped>
.filter-sidebar {
  display: flex;
  flex-direction: column;
  gap: 0;
  width: 100%;
}
.sidebar-section { padding: 0.25rem 0; }
.sidebar-divider { margin: 0.4rem 0; }
.filter-sidebar :deep(.sidebar-divider .p-divider-content),
.filter-sidebar :deep(.sidebar-divider::before) {
  border-color: #2d3748 !important;
}
.filter-sidebar :deep(.p-divider.sidebar-divider) {
  border-top-color: #2d3748 !important;
}
.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  background: none;
  border: none;
  cursor: pointer;
  padding: 0.25rem 0;
  font-size: var(--fs-base);
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  color: #94a3b8;
}
.section-header:hover { color: #e2e8f0; }
.filter-options {
  display: flex;
  flex-direction: column;
  gap: 0.3rem;
  margin-top: 0.4rem;
}
.filter-option {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  cursor: pointer;
  font-size: var(--fs-body);
  font-weight: normal;
  color: #94a3b8;
}
.filter-option:hover { color: #e2e8f0; }
.filter-sidebar :deep(.p-checkbox .p-checkbox-box),
.filter-sidebar :deep(.p-radiobutton .p-radiobutton-box) {
  background: rgba(255, 255, 255, 0.08) !important;
  border: 1.5px solid rgba(255, 255, 255, 0.2) !important;
  border-radius: 6px !important;
  box-shadow: none !important;
  transition: background 0.15s, border-color 0.15s;
}
.filter-sidebar :deep(.p-checkbox-checked .p-checkbox-box),
.filter-sidebar :deep(.p-radiobutton-checked .p-radiobutton-box) {
  background: var(--p-primary-color) !important;
  border-color: var(--p-primary-color) !important;
}
.filter-sidebar :deep(.p-checkbox .p-checkbox-box:hover),
.filter-sidebar :deep(.p-checkbox .p-checkbox-box.p-focus),
.filter-sidebar :deep(.p-radiobutton .p-radiobutton-box:hover),
.filter-sidebar :deep(.p-radiobutton .p-radiobutton-box.p-focus) {
  background: rgba(255, 255, 255, 0.14) !important;
  border-color: rgba(255, 255, 255, 0.35) !important;
}
.filter-sidebar :deep(.p-checkbox .p-checkbox-icon),
.filter-sidebar :deep(.p-radiobutton .p-radiobutton-icon) {
  color: #fff !important;
}

/* Filter search box (InputText in filter sidebar) */
.filter-sidebar :deep(.p-inputtext),
.filter-sidebar input[type="text"],
.filter-sidebar .label-input {
  background: rgba(255, 255, 255, 0.08) !important;
  color: #e2e8f0 !important;
  border: 1.5px solid rgba(255, 255, 255, 0.2) !important;
  border-radius: 6px !important;
  box-shadow: none !important;
  transition: background 0.15s, border-color 0.15s, color 0.15s;
}
.filter-sidebar :deep(.p-inputtext::placeholder),
.filter-sidebar input[type="text"]::placeholder {
  color: #94a3b8 !important;
}
.filter-sidebar :deep(.p-inputtext:focus),
.filter-sidebar input[type="text"]:focus,
.filter-sidebar .label-input:focus {
  background: rgba(255, 255, 255, 0.12) !important;
  border-color: rgba(255, 255, 255, 0.4) !important;
  color: #fff !important;
}

.filter-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.opt-tag { font-size: 0.75em; }

/* Pin Tag severity colors so they look consistent on the dark sidebar
   regardless of whether light or dark mode is active globally.        */
.filter-sidebar :deep(.opt-tag.p-tag-success)    { background: #16a34a !important; color: #fff !important; }
.filter-sidebar :deep(.opt-tag.p-tag-secondary)  { background: #475569 !important; color: #e2e8f0 !important; }
.filter-sidebar :deep(.opt-tag.p-tag-info)       { background: #4f46e5 !important; color: #fff !important; }
.filter-sidebar :deep(.opt-tag.p-tag-warn),
.filter-sidebar :deep(.opt-tag.p-tag-warning)    { background: #d97706 !important; color: #fff !important; }
.filter-sidebar :deep(.opt-tag.p-tag-danger)     { background: #dc2626 !important; color: #fff !important; }
.opt-label {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.opt-count {
  margin-left: auto;
  font-size: var(--fs-xs);
  color: var(--text-color-secondary);
}

/* Labels filter */
.label-input-wrap {
  display: flex;
  align-items: center;
  background: rgba(255, 255, 255, 0.08);
  border: 1.5px solid rgba(255, 255, 255, 0.2);
  border-radius: 6px;
  padding: 0.3rem 0.5rem;
  transition: background 0.15s, border-color 0.15s;
}
.label-input-wrap:focus-within {
  background: rgba(255, 255, 255, 0.12);
  border-color: rgba(255, 255, 255, 0.4);
}
.filter-search-wrap {
  margin-top: 0.4rem;
}
.label-input {
  border: none;
  background: transparent;
  color: #e2e8f0;
  font-size: var(--fs-sm);
  outline: none;
  width: 100%;
}
.label-input::placeholder { color: #94a3b8; }
.label-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 0.3rem;
  margin-top: 0.4rem;
}
.label-tag {
  display: inline-flex;
  align-items: center;
  gap: 0.2rem;
  font-size: var(--fs-sm);
  border: 1px solid;
  border-radius: 4px;
  padding: 0.1em 0.35em;
}
.label-remove {
  background: none;
  border: none;
  cursor: pointer;
  padding: 0;
  line-height: 1;
  font-size: 1.1em;
  color: inherit;
  opacity: 0.7;
}
.label-remove:hover { opacity: 1; }

.clear-row, .shown-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: auto;
  padding-top: 0.75rem;
}
.shown-count {
  font-size: var(--fs-sm);
  color: var(--text-color-secondary);
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
