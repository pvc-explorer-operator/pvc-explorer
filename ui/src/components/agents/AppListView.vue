<template>
  <div class="list-view">
    <template v-if="explorers.length">
      <table class="explorer-table">
        <thead>
          <tr>
            <th>Status</th>
            <th>Name</th>
            <th>Namespace</th>
            <th>PVC</th>
            <th>Mount</th>
            <th>Consumers</th>
            <th></th>
          </tr>
        </thead>
        <TransitionGroup tag="tbody" name="card">
          <tr
            v-for="explorer in explorers"
            :key="`${explorer.namespace}/${explorer.name}`"
            :class="{ 'row-in-use': !!explorer.consumerCount }"
            @click="goToDetails(explorer)"
          >
            <td>
              <div class="flex items-center gap-2">
                <span :class="['phase-dot', `dot-${phaseCss(explorer.phase)}`]" />
                <Tag
                  :value="explorer.phase"
                  :severity="phaseSeverity(explorer.phase)"
                  rounded
                />
              </div>
            </td>
            <td class="font-semibold">
              {{ explorer.name }}
              <span v-if="explorer.labels?.length" class="label-list">
                <Chip v-for="label in explorer.labels" :key="label" :label="label" class="label-chip ml-1" size="small" />
              </span>
            </td>
            <td class="text-muted-color">{{ explorer.namespace }}</td>
            <td class="text-muted-color">{{ explorer.pvcName }}</td>
            <td>
              <span v-if="explorer.mountState" class="mount-label">{{ mountLabel(explorer.mountState) }}</span>
              <span v-else class="text-muted-color">—</span>
            </td>
            <td>
              <span v-if="explorer.consumerCount" class="consumer-badge">{{ explorer.consumerCount }}</span>
              <span v-else class="text-muted-color">—</span>
            </td>
            <td class="text-right">
              <i class="pi pi-chevron-right text-muted-color"></i>
            </td>
          </tr>
        </TransitionGroup>
      </table>
    </template>
    <div v-else class="empty-state">
      <i class="pi pi-inbox empty-icon" />
      <div class="empty-msg">No explorers found.</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import Tag from 'primevue/tag'
import type { Explorer } from '../../stores/explorerStore'

defineProps<{ explorers: Explorer[] }>()

const router = useRouter()

function goToDetails(explorer: Explorer) {
  router.push(`/explorers/${encodeURIComponent(explorer.namespace)}/${encodeURIComponent(explorer.name)}`)
}

function phaseCss(phase: string) {
  return phase.toLowerCase().replace(/\s/g, '-')
}

function phaseSeverity(phase: string) {
  if (phase === 'Running') return 'success'
  if (phase === 'ScaledToZero') return 'secondary'
  if (phase === 'Waking') return 'info'
  if (phase === 'Failed') return 'danger'
  if (phase === 'Pending') return 'warn'
  return 'secondary'
}

function mountLabel(strategy: string) {
  if (strategy === 'Direct') return 'Read-write'
  if (strategy === 'NodeAffinity') return 'Read-only'
  return strategy
}
</script>

<style scoped>
.explorer-table {
  width: 100%;
  border-collapse: collapse;
}
.explorer-table th {
  text-align: left;
  padding: 0.65rem 0.75rem;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--text-color-secondary);
  border-bottom: 1px solid var(--surface-border);
}
.explorer-table td {
  padding: 0.65rem 0.75rem;
  font-size: 0.875rem;
  border-bottom: 1px solid var(--surface-border);
  vertical-align: middle;
}
.explorer-table tbody tr {
  cursor: pointer;
  transition: background-color 0.1s;
}
.explorer-table tbody tr:hover {
  background-color: var(--surface-hover);
}
.explorer-table tbody tr.row-in-use {
  border-left: 3px solid var(--orange-500);
}
.phase-dot {
  display: inline-block;
  width: 0.6rem;
  height: 0.6rem;
  border-radius: 50%;
  flex-shrink: 0;
}
.dot-running { background: var(--green-500); }
.dot-scaled-to-zero { background: var(--surface-400); }
.dot-waking { background: var(--indigo-500); }
.dot-failed { background: var(--red-500); }
.dot-pending { background: var(--yellow-500); }

.consumer-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 1.5rem;
  height: 1.5rem;
  padding: 0 0.35rem;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 600;
  background: var(--orange-100);
  color: var(--orange-800);
}
.dark .consumer-badge {
  background: color-mix(in srgb, var(--orange-500) 20%, transparent);
  color: var(--orange-300);
}

.mount-label {
  font-size: 0.8125rem;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 220px;
  color: var(--text-color-secondary);
  opacity: 0.85;
}
.empty-icon {
  font-size: 2.5rem;
  margin-bottom: 0.7rem;
}
.empty-msg {
  font-size: 1.08rem;
}

@media (prefers-reduced-motion: no-preference) {
  .card-enter-from,
  .card-leave-to  { opacity: 0; }
  .card-enter-active,
  .card-leave-active { transition: opacity 0.2s ease; }
  .card-move      { transition: transform 0.25s ease; }
}
</style>
