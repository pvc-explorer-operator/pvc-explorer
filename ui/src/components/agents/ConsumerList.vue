<template>
  <div class="consumer-list-wrap">
    <h3 class="section-title">Active Consumers</h3>
    <div v-if="!consumers.length" class="no-consumers">No active consumers</div>
    <table v-else class="consumer-table">
      <thead>
        <tr>
          <th>Pod</th>
          <th>Owner</th>
          <th>Node</th>
          <th>Mount</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="c in consumers" :key="c.podName">
          <td class="col-pod">{{ c.podName }}</td>
          <td class="col-owner">
            <span v-if="c.ownerKind" class="owner-kind">{{ c.ownerKind }}</span>
            {{ c.ownerName || '—' }}
          </td>
          <td class="col-node">{{ c.nodeName || '—' }}</td>
          <td class="col-mount">
            <span :class="['mount-chip', c.mountReadOnly ? 'chip-ro' : 'chip-rw']">
              {{ c.mountReadOnly ? 'RO' : 'RW' }}
            </span>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
interface ConsumerInfo {
  podName: string
  ownerKind?: string
  ownerName?: string
  nodeName?: string
  mountReadOnly?: boolean
}

defineProps<{ consumers: ConsumerInfo[] }>()
</script>

<style scoped>
.consumer-list-wrap { overflow-x: auto; }

.section-title {
  font-size: var(--fs-sm);
  font-weight: 600;
  color: var(--text-color-secondary);
  text-transform: uppercase;
  letter-spacing: 0.06em;
  margin: 0 0 0.6rem;
}

.no-consumers {
  font-size: var(--fs-sm);
  color: var(--text-color-secondary);
  padding: 0.5rem 0;
}

.consumer-table {
  width: 100%;
  border-collapse: collapse;
  font-size: var(--fs-sm);
}

.consumer-table th {
  text-align: left;
  padding: 0.35rem 0.7rem;
  color: var(--text-color-secondary);
  border-bottom: 1px solid var(--surface-border);
  font-weight: 500;
  white-space: nowrap;
}

.consumer-table td {
  padding: 0.35rem 0.7rem;
  border-bottom: 1px solid var(--surface-border);
  color: var(--text-color-secondary);
}

.consumer-table tr:last-child td { border-bottom: none; }

.col-pod  { font-family: monospace; font-size: var(--fs-sm); color: var(--text-color); }
.col-node { font-size: var(--fs-sm); color: var(--text-color-secondary); }

.owner-kind {
  display: inline-block;
  font-size: var(--fs-xs);
  background: var(--surface-hover);
  color: var(--text-color-secondary);
  border-radius: 3px;
  padding: 0.05em 0.4em;
  margin-right: 0.3em;
  vertical-align: middle;
}

.mount-chip {
  display: inline-block;
  font-size: var(--fs-xs);
  font-weight: 600;
  border-radius: 4px;
  padding: 0.1em 0.5em;
}
.chip-rw { background: color-mix(in srgb, #22c55e 15%, transparent); color: #22c55e; }
.chip-ro { background: color-mix(in srgb, #f59e0b 15%, transparent); color: #f59e0b; }
</style>
