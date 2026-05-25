<template>
  <div class="conditions-table-wrap">
    <h3 class="section-title">Conditions</h3>
    <table class="conditions-table">
      <thead>
        <tr>
          <th>Type</th>
          <th>Status</th>
          <th>Reason</th>
          <th>Message</th>
          <th>Last Transition</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="c in conditions" :key="c.type" :class="rowClass(c)">
          <td class="col-type">{{ c.type }}</td>
          <td class="col-status">
            <span :class="['status-dot', statusDotClass(c)]" />
            {{ c.status }}
          </td>
          <td class="col-reason">{{ c.reason || '—' }}</td>
          <td class="col-message">{{ c.message || '—' }}</td>
          <td class="col-time">{{ fmtTime(c.lastTransitionTime) }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
interface Condition {
  type: string
  status: string
  reason?: string
  message?: string
  lastTransitionTime?: string
}

defineProps<{ conditions: Condition[] }>()

function statusDotClass(c: Condition) {
  if (c.status === 'True')  return 'dot-true'
  if (c.status === 'False') return 'dot-false'
  return 'dot-unknown'
}

function rowClass(c: Condition) {
  if (c.status === 'False' && c.type !== 'Degraded') return 'row-warn'
  if (c.status === 'True'  && c.type === 'Degraded') return 'row-warn'
  return ''
}

function fmtTime(iso?: string) {
  if (!iso) return '—'
  return new Date(iso).toLocaleString(undefined, { dateStyle: 'short', timeStyle: 'short' })
}
</script>

<style scoped>
.conditions-table-wrap { overflow-x: auto; }

.section-title {
  font-size: var(--fs-sm);
  font-weight: 600;
  color: var(--text-color-secondary);
  text-transform: uppercase;
  letter-spacing: 0.06em;
  margin: 0 0 0.6rem;
}

.conditions-table {
  width: 100%;
  border-collapse: collapse;
  font-size: var(--fs-sm);
}

.conditions-table th {
  text-align: left;
  padding: 0.35rem 0.7rem;
  color: var(--text-color-secondary);
  border-bottom: 1px solid var(--surface-border);
  font-weight: 500;
  white-space: nowrap;
}

.conditions-table td {
  padding: 0.35rem 0.7rem;
  border-bottom: 1px solid var(--surface-border);
  color: var(--text-color-secondary);
  vertical-align: top;
}

.conditions-table tr:last-child td { border-bottom: none; }
.conditions-table tr.row-warn td {
  background: color-mix(in srgb, var(--p-amber-500) 8%, transparent);
}

.col-type    { font-weight: 500; color: var(--text-color); white-space: nowrap; }
.col-status  { white-space: nowrap; }
.col-reason  { white-space: nowrap; }
.col-message { max-width: 320px; }
.col-time    { white-space: nowrap; font-size: var(--fs-sm); color: var(--text-color-secondary); }

.status-dot {
  display: inline-block;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  margin-right: 0.35em;
  vertical-align: middle;
}
.dot-true    { background: var(--p-green-500); }
.dot-false   { background: var(--p-amber-500); }
.dot-unknown { background: var(--text-color-secondary); }
</style>
