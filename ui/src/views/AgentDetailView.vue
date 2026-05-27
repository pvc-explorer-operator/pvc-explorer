<template>
  <main v-if="explorer" class="detail-view">
    <div class="detail-card detail-section">
      <div class="detail-header">
        <span :class="['phase-dot', `dot-${(explorer.phase || '').toLowerCase()}`]" aria-hidden="true" />
        <h1 class="explorer-name">{{ explorer.name }}</h1>
        <Tag v-if="explorer.phase" :value="explorer.phase" :severity="phaseSeverity(explorer.phase)" rounded />
        <MountStateBanner :state="explorer.mountState" />
      </div>
      <div class="detail-meta">
        <div><span class="meta-label">Namespace:</span> <span class="meta-value">{{ explorer.namespace }}</span></div>
        <div><span class="meta-label">PVC:</span> <span class="meta-value">{{ explorer.pvcName }}</span></div>
        <div v-if="explorer.accessMode || explorer.mode"><span class="meta-label">Mode:</span> <Tag :value="explorer.accessMode || explorer.mode" severity="info" rounded /></div>
      </div>
      <div v-if="explorer.labels?.length" class="labels-row">
        <Chip v-for="label in explorer.labels" :key="label" :label="label" class="label-chip" />
      </div>
      <div class="detail-actions">
        <Button v-if="explorer.phase === 'Running'" severity="success" icon="pi pi-folder-open" label="Browse Files" rounded @click="goToFiles" />
        <Button v-if="explorer.phase === 'Running'" severity="warn" icon="pi pi-power-off" label="Disconnect" rounded :loading="disconnecting" @click="doDisconnect" />
        <Button v-else-if="explorer.phase === 'ScaledToZero'" severity="primary" icon="pi pi-plug" label="Connect" rounded @click="wake" />
        <Button v-else-if="explorer.phase === 'Waking'" severity="primary" icon="pi pi-spin pi-spinner" label="Waking..." rounded disabled />
        <Button v-else severity="secondary" icon="pi pi-ban" label="Unavailable" rounded disabled />
        <Button severity="secondary" icon="pi pi-refresh" label="Refresh" rounded @click="refresh" />
      </div>
    </div>
    <MountStateBanner v-if="explorer.mountState && explorer.mountState !== 'Mounted'" :state="explorer.mountState" style="margin-top:1rem;" class="detail-section" />
    <ConditionsTable v-if="explorer.conditions" :conditions="explorer.conditions" style="margin-top:1.2rem;" class="detail-section" />
    <ConsumerList v-if="explorer.consumers" :consumers="explorer.consumers" style="margin-top:1.2rem;" class="detail-section" />
    <WakeUpDialog v-if="showWakeDialog" :explorer="explorer" @close="onWakeDialogClose" />
  </main>
  <div v-else-if="loading" class="detail-view">
    <div class="detail-card sk-agent-card">
      <div class="detail-header">
        <Skeleton shape="circle" size="0.625rem" />
        <Skeleton width="50%" height="1.25rem" />
        <Skeleton width="4rem" height="1.25rem" borderRadius="999px" />
      </div>
      <div class="detail-meta">
        <div><Skeleton width="8rem" height="0.875rem" /></div>
        <div><Skeleton width="10rem" height="0.875rem" /></div>
        <div><Skeleton width="6rem" height="0.875rem" /></div>
      </div>
      <div class="labels-row">
        <Skeleton v-for="j in 4" :key="j" width="5rem" height="1.25rem" borderRadius="4px" />
      </div>
      <div class="detail-actions">
        <Skeleton width="7.5rem" height="2.25rem" borderRadius="6px" />
        <Skeleton width="6rem" height="2.25rem" borderRadius="6px" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useExplorerStore } from '../stores/explorerStore'
import type { Explorer } from '../stores/explorerStore'
import MountStateBanner from '../components/agents/MountStateBanner.vue'
import WakeUpDialog from '../components/shared/WakeUpDialog.vue'
import ConditionsTable from '../components/agents/ConditionsTable.vue'
import ConsumerList from '../components/agents/ConsumerList.vue'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Chip from 'primevue/chip'
import Skeleton from 'primevue/skeleton'
import { useExplorerDetailShortcuts } from '../composables/useExplorerDetailShortcuts'

const route = useRoute()
const router = useRouter()
const store = useExplorerStore()
const explorer = ref<Explorer | null>(null)
const loading = ref(true)
const showWakeDialog = ref(false)
const disconnecting = ref(false)

async function fetchExplorer() {
  loading.value = true
  const ns = route.params.ns as string
  const name = route.params.name as string
  try {
    explorer.value = await store.fetchExplorer(ns, name)
  } catch {
    explorer.value = null
  } finally {
    loading.value = false
  }
}

function phaseSeverity(phase: string) {
  if (phase === 'Running') return 'success'
  if (phase === 'ScaledToZero') return 'warning'
  if (phase === 'Waking') return 'info'
  if (phase === 'Failed') return 'danger'
  if (phase === 'Pending') return 'warn'
  return 'secondary'
}

function refresh() {
  fetchExplorer()
}

function goToFiles() {
  if (!explorer.value) return
  router.push(`/explorers/${explorer.value.namespace}/${explorer.value.name}/files`)
}

function wake() {
  showWakeDialog.value = true
}

function onWakeDialogClose() {
  showWakeDialog.value = false
  fetchExplorer()
}

async function doDisconnect() {
  if (!explorer.value) return
  disconnecting.value = true
  try {
    await store.sleepExplorer(explorer.value.namespace, explorer.value.name)
    await pollPhase('ScaledToZero')
  } catch {
    // fetchExplorer already updated the store
  }
  disconnecting.value = false
}

async function pollPhase(target: string, timeout = 60000): Promise<void> {
  if (!explorer.value) return
  const deadline = Date.now() + timeout
  while (Date.now() < deadline) {
    await new Promise(r => setTimeout(r, 3000))
    try {
      const e = await store.fetchExplorer(explorer.value.namespace, explorer.value.name)
      if (e.phase === target) return
    } catch {
      // retry
    }
  }
}

onMounted(fetchExplorer)
watch(() => [route.params.ns, route.params.name], fetchExplorer)

useExplorerDetailShortcuts({ explorer, goToFiles, wake, doDisconnect, refresh })
</script>

<style scoped>
.detail-view {
  max-width: 700px;
  margin: 2.5rem auto 0 auto;
  padding: 1.5rem 1rem 3rem;
}
.detail-card {
  background: var(--surface-card);
  border: 1px solid var(--surface-border);
  border-radius: 8px;
  padding: 1.3rem 1.5rem 1.1rem 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 0.7rem;
}
.detail-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}
.phase-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
  background: var(--text-color-secondary);
}
.dot-running         { background: var(--p-green-500); }
.dot-scaledtozero    { background: var(--p-amber-500); }
.dot-waking          { background: var(--primary-color); }
.dot-pending         { background: var(--p-amber-500); }
.dot-failed          { background: var(--p-red-500); }
.explorer-name {
  flex: 1;
  font-weight: 600;
  color: var(--text-color);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.detail-meta {
  display: flex;
  gap: 1.5rem;
  font-size: 0.98rem;
  flex-wrap: wrap;
}
.meta-label { color: var(--text-color-secondary); }
.meta-value { color: var(--text-color-secondary); }
.labels-row { display: flex; flex-wrap: wrap; gap: 0.25rem; }
.label-chip {
  background: var(--surface-hover);
  color: var(--text-color-secondary);
  border-radius: 4px;
  padding: 0.1em 0.5em;
  font-size: 0.875rem;
}
.detail-actions {
  margin-top: 0.5rem;
  display: flex;
  gap: 0.7rem;
}

.sk-agent-card .detail-header,
.sk-agent-card .detail-meta,
.sk-agent-card .labels-row,
.sk-agent-card .detail-actions {
  pointer-events: none;
}

@media (prefers-reduced-motion: no-preference) {
  @supports (animation-timeline: scroll()) {
    .detail-section {
      animation: ag-fade-in-up linear both;
      animation-timeline: view();
      animation-range: entry 0% entry 25%;
    }
    @keyframes ag-fade-in-up {
      from { opacity: 0; translate: 0 20px; }
      to   { opacity: 1; translate: 0; }
    }
  }
}
</style>
