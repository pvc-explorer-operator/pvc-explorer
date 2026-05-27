<template>
  <a
    :class="['app-card', { 'card-in-use': !!explorer.consumerCount }]"
    :href="`/explorers/${explorer.namespace}/${explorer.name}`"
    @click.prevent="goToDetails"
  >
    <!-- Header -->
    <div class="flex justify-between items-center mb-4">
      <div class="flex items-center gap-2 card-name-wrap">
        <span :class="['phase-dot', `dot-${phaseCss}`]" aria-hidden="true" />
        <span class="explorer-name font-semibold" :title="explorer.name">{{ explorer.name }}</span>
      </div>
      <div class="flex items-center gap-2 card-tags-wrap">
        <Tag v-if="explorer.consumerCount" value="In use" severity="warn" rounded />
        <Tag
          :value="explorer.phase"
          :severity="phaseSeverity(explorer.phase)"
          rounded
        />
      </div>
    </div>

    <!-- Waking progress indicator -->
    <ProgressBar
      v-if="explorer.phase === 'Waking'"
      mode="indeterminate"
      class="waking-bar"
    />

    <!-- Metadata rows -->
    <div class="card-body">
      <div class="flex justify-between">
        <span class="text-muted-color text-sm">Namespace</span>
        <span class="text-surface-900 text-sm card-val" :title="explorer.namespace">{{ explorer.namespace }}</span>
      </div>
      <div class="flex justify-between">
        <span class="text-muted-color text-sm">PVC</span>
        <span class="text-surface-900 text-sm mono card-val" :title="explorer.pvcName">{{ explorer.pvcName }}</span>
      </div>
      <div v-if="explorer.scope" class="flex justify-between">
        <span class="text-muted-color text-sm">Scope</span>
        <span class="text-surface-900 text-sm card-val" :title="explorer.scope">{{ explorer.scope }}</span>
      </div>
      <div class="flex justify-between">
        <span class="text-muted-color text-sm">Mount</span>
        <Tag
          :value="mountLabel"
          :severity="mountSeverity"
          rounded
        />
      </div>
      <div v-if="explorer.accessMode || explorer.mode" class="flex justify-between">
        <span class="text-muted-color text-sm">Access</span>
        <Tag
          :value="explorer.accessMode || explorer.mode"
          severity="secondary"
          rounded
        />
      </div>
      <div v-if="explorer.consumerCount" class="flex justify-between">
        <span class="text-muted-color text-sm">Consumers</span>
        <Badge :value="explorer.consumerCount" severity="warning" />
      </div>
      <div v-else-if="explorer.phase !== 'Running'" class="flex justify-between">
        <span class="text-muted-color text-sm">Status</span>
        <Tag value="Not attached" severity="secondary" rounded />
      </div>
      <div v-if="idleDisplay && explorer.phase === 'Running'" class="flex justify-between">
        <span class="text-muted-color text-sm">Idle timeout</span>
        <span class="text-surface-900 text-sm">{{ idleDisplay }}</span>
      </div>
      <div v-if="explorer.createdAt" class="flex justify-between">
        <span class="text-muted-color text-sm">Created</span>
        <span class="text-surface-900 text-sm" :title="explorer.createdAt">{{ relativeTime(explorer.createdAt) }}</span>
      </div>
    </div>

    <!-- Labels -->
    <div v-if="explorer.labels?.length" class="flex gap-2 flex-wrap mt-2">
      <Chip
        v-for="label in explorer.labels"
        :key="label"
        :label="label"
        class="label-chip"
      />
    </div>

    <!-- Actions -->
    <div class="card-actions" @click.stop>
      <Button
        v-if="explorer.phase === 'Running'"
        label="Browse"
        icon="pi pi-folder-open"
        severity="success"
        size="small"
        @click="goToFiles"
      />
      <Button
        v-if="explorer.phase === 'Running'"
        label="Disconnect"
        icon="pi pi-power-off"
        severity="warn"
        size="small"
        :loading="disconnecting"
        @click="doDisconnect"
      />
      <Button
        v-else-if="explorer.phase === 'ScaledToZero'"
        label="Connect"
        icon="pi pi-power-off"
        size="small"
        :loading="connecting"
        @click="doConnect"
      />
      <Button
        v-else-if="explorer.phase === 'Waking'"
        label="Waking..."
        icon="pi pi-spin pi-spinner"
        severity="info"
        size="small"
        disabled
      />
      <Button
        v-else
        label="Details"
        icon="pi pi-info-circle"
        severity="secondary"
        size="small"
        @click="goToDetails"
      />
      <Button
        icon="pi pi-ellipsis-v"
        severity="secondary"
        text
        size="small"
        class="ml-auto"
        @click="goToDetails"
        :title="'View details'"
      />
    </div>
  </a>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import Tag from 'primevue/tag'
import Chip from 'primevue/chip'
import Button from 'primevue/button'
import Badge from 'primevue/badge'
import ProgressBar from 'primevue/progressbar'
import type { Explorer } from '../../stores/explorerStore'
import { useExplorerStore } from '../../stores/explorerStore'

const props = defineProps<{ explorer: Explorer }>()
const router = useRouter()
const explorerStore = useExplorerStore()

const connecting = ref(false)
const disconnecting = ref(false)

const idleKey = computed(() => `${props.explorer.namespace}/${props.explorer.name}`)

const idleDisplay = computed(() => {
  const remaining = explorerStore.idleRemaining[idleKey.value]
  if (remaining !== undefined) {
    if (remaining >= 60) return `${Math.floor(remaining / 60)}m ${remaining % 60}s remaining`
    return `${remaining}s remaining`
  }
  return props.explorer.idleTimeout ?? null
})

const phaseCss = computed(() => props.explorer.phase.toLowerCase().replace(/\s/g, '-'))

function goToFiles() {
  router.push(`/explorers/${props.explorer.namespace}/${props.explorer.name}/files`)
}
function goToDetails() {
  router.push(`/explorers/${props.explorer.namespace}/${props.explorer.name}`)
}

async function pollPhase(target: string, timeout = 60000): Promise<void> {
  const deadline = Date.now() + timeout
  while (Date.now() < deadline) {
    try {
      const e = await explorerStore.fetchExplorer(props.explorer.namespace, props.explorer.name)
      if (e.phase === target) return
    } catch {
      // retry
    }
    await new Promise(r => setTimeout(r, 3000))
  }
  throw new Error(`Timed out waiting for phase ${target}`)
}

async function doConnect() {
  connecting.value = true
  try {
    await explorerStore.wakeExplorer(props.explorer.namespace, props.explorer.name)
    await pollPhase('Running')
  } catch {
    // fetchExplorer + upsertExplorer already updated the store
    // if polling timed out, still show current phase (Waking or whatever)
  }
  connecting.value = false
}

async function doDisconnect() {
  disconnecting.value = true
  try {
    await explorerStore.sleepExplorer(props.explorer.namespace, props.explorer.name)
    await explorerStore.fetchExplorer(props.explorer.namespace, props.explorer.name)
    if (props.explorer.phase !== 'ScaledToZero') {
      await pollPhase('ScaledToZero')
    }
  } catch {
    await explorerStore.fetchExplorer(props.explorer.namespace, props.explorer.name).catch(() => {})
  }
  disconnecting.value = false
}

function phaseSeverity(phase: string) {
  if (phase === 'Running') return 'success'
  if (phase === 'ScaledToZero') return 'secondary'
  if (phase === 'Waking') return 'info'
  if (phase === 'Failed') return 'danger'
  if (phase === 'Pending') return 'warn'
  return 'secondary'
}

const mountLabel = computed(() => {
  if (props.explorer.consumerCount) return 'Read-only'
  return 'Read-write'
})

const mountSeverity = computed(() => {
  if (props.explorer.consumerCount) return 'warning'
  return 'success'
})

function relativeTime(iso: string): string {
  const diff = Date.now() - new Date(iso).getTime()
  const sec = Math.floor(diff / 1000)
  if (sec < 60) return `${sec}s ago`
  const min = Math.floor(sec / 60)
  if (min < 60) return `${min}m ago`
  const hr = Math.floor(min / 60)
  if (hr < 24) return `${hr}h ago`
  const d = Math.floor(hr / 24)
  return `${d}d ago`
}
</script>

<style scoped>
.app-card {
  content-visibility: auto;
  contain-intrinsic-size: auto 280px auto 160px;
  background: var(--surface-card);
  border: 1px solid var(--surface-border);
  border-radius: var(--content-border-radius);
  padding: 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  cursor: pointer;
  transition: box-shadow 0.2s, border-color 0.2s;
}
@supports not (content-visibility: auto) {
  .app-card { contain: layout style paint; }
}
.app-card:hover {
  border-color: var(--primary-color);
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.06);
}
.card-in-use {
  border-left: 4px solid var(--p-orange-400, #f59e0b);
}
.phase-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
  background: var(--text-color-secondary);
}
.dot-running      { background: var(--p-green-500); }
.dot-scaledtozero { background: var(--text-color-secondary); }
.dot-waking       { background: var(--primary-color); animation: pulse 1s ease-in-out infinite; }
.dot-failed       { background: var(--p-red-500); }
@keyframes pulse {
  0%, 100% { opacity: 1; }
  50%       { opacity: 0.4; }
}
/* Left side of header: must shrink so the name can truncate */
.card-name-wrap {
  min-width: 0;
  flex: 1;
  overflow: hidden;
}
/* Right side of header: tags must never shrink or wrap */
.card-tags-wrap {
  flex-shrink: 0;
  margin-left: 0.5rem;
}
.explorer-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.waking-bar {
  height: 3px;
  border-radius: 2px;
}
.card-body {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
  flex: 1;
}
.mono {
  font-family: var(--font-mono, monospace);
}
/* Truncate long text values in card body rows */
.card-val {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  text-align: right;
  min-width: 0;
}
.label-chip {
  font-size: var(--fs-sm);
  height: 1.5rem;
  padding: 0 0.5rem;
}
.card-actions {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  padding-top: 0.5rem;
  border-top: 1px solid var(--surface-border);
  margin-top: auto;
}
.ml-auto { margin-left: auto; }
</style>
