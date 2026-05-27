<template>
  <a
    class="sc-card"
    :href="`/scopes/${encodeURIComponent(scope.name)}`"
    @click.prevent="$router.push(`/scopes/${encodeURIComponent(scope.name)}`)"
  >
    <div class="sc-card-top">
      <span class="sc-card-name">{{ scope.name }}</span>
      <span :class="['sc-card-badge', phaseClass]">{{ scope.phase }}</span>
    </div>
    <div class="sc-card-stats">
      <div class="sc-stat">
        <i class="pi pi-server"></i>
        <span>{{ scope.namespaceCount }} namespace{{ scope.namespaceCount !== 1 ? 's' : '' }}</span>
      </div>
      <div class="sc-stat">
        <i class="pi pi-database"></i>
        <span>{{ scope.explorerCount }} explorer{{ scope.explorerCount !== 1 ? 's' : '' }}</span>
      </div>
    </div>
  </a>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Scope } from '../../stores/explorerStore'

const props = defineProps<{ scope: Scope }>()

const phaseClass = computed(() => {
  const p = props.scope.phase?.toLowerCase() ?? ''
  if (p === 'ready' || p === 'true') return 'sc-phase-ok'
  if (p === 'failed') return 'sc-phase-fail'
  return 'sc-phase-warn'
})
</script>

<style scoped>
.sc-card {
  content-visibility: auto;
  contain-intrinsic-size: auto 280px auto 120px;
  background: var(--surface-card);
  border: 1px solid var(--surface-border);
  border-radius: 10px;
  padding: 1.25rem;
  cursor: pointer;
  transition: border-color 0.15s, box-shadow 0.15s, transform 0.15s;
  box-shadow: 0 1px 3px rgba(0,0,0,0.04);
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
@supports not (content-visibility: auto) {
  .sc-card { contain: layout style paint; }
}
.sc-card:hover {
  border-color: var(--p-primary-400);
  box-shadow: 0 4px 16px rgba(0,0,0,0.1);
  transform: translateY(-2px);
}
.sc-card-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.5rem;
}
.sc-card-name {
  font-size: 1rem;
  font-weight: 700;
  color: var(--text-color);
  font-family: 'JetBrains Mono', monospace;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.sc-card-badge {
  font-size: 0.7rem;
  font-weight: 600;
  padding: 2px 10px;
  border-radius: 20px;
  white-space: nowrap;
  flex-shrink: 0;
}
.sc-phase-ok   { background: rgba(34,197,94,0.12); color: #22c55e; border: 1px solid rgba(34,197,94,0.25); }
.sc-phase-warn { background: rgba(245,158,11,0.12); color: #f59e0b; border: 1px solid rgba(245,158,11,0.25); }
.sc-phase-fail { background: rgba(239,68,68,0.12);  color: #ef4444; border: 1px solid rgba(239,68,68,0.25); }
.sc-card-stats {
  display: flex;
  gap: 1.25rem;
}
.sc-stat {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  font-size: 0.82rem;
  color: var(--text-color-secondary);
}
.sc-stat i {
  font-size: 0.85rem;
}
</style>
