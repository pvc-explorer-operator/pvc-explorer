<template>
  <TransitionGroup
    v-if="explorers.length"
    name="card"
    tag="div"
    class="card-grid"
  >
    <AppCard
      v-for="explorer in explorers"
      :key="`${explorer.namespace}/${explorer.name}`"
      :explorer="explorer"
    />
  </TransitionGroup>
  <div v-else-if="loading" class="card-grid">
    <div v-for="i in 6" :key="i" class="sk-card">
      <div class="sk-card-header">
        <Skeleton shape="circle" size="0.625rem" />
        <Skeleton width="60%" height="1rem" />
        <Skeleton width="3.5rem" height="1.25rem" borderRadius="999px" />
      </div>
      <div class="sk-card-meta">
        <Skeleton width="40%" height="0.75rem" />
        <Skeleton width="55%" height="0.75rem" />
      </div>
      <div class="sk-card-labels">
        <Skeleton v-for="j in 3" :key="j" width="3rem" height="1.25rem" borderRadius="4px" />
      </div>
      <div class="sk-card-actions">
        <Skeleton width="7rem" height="2rem" borderRadius="6px" />
        <Skeleton width="5rem" height="2rem" borderRadius="6px" />
      </div>
    </div>
  </div>
  <div v-else class="card-grid">
    <div class="empty-state">
      <i class="pi pi-inbox empty-icon" />
      <div class="empty-msg">No explorers found.</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import AppCard from './AppCard.vue'
import Skeleton from 'primevue/skeleton'
import type { Explorer } from '../../stores/explorerStore'

defineProps<{ explorers: Explorer[]; loading?: boolean }>()
</script>

<style scoped>
.card-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  grid-auto-rows: 1fr;
  gap: 1.25rem;
  align-items: stretch;
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
  .card-leave-to  { opacity: 0; translate: 0 12px; }
  .card-enter-active,
  .card-leave-active { transition: opacity 0.2s ease, translate 0.2s ease; }
  .card-move      { transition: transform 0.25s ease; }
  .card-leave-active { position: absolute; }
}

/* ── Skeleton card ── */
.sk-card {
  display: flex;
  flex-direction: column;
  gap: 0.7rem;
  padding: 1.3rem 1.5rem 1.1rem;
  background: var(--surface-card);
  border: 1px solid var(--surface-border);
  border-radius: 8px;
}
.sk-card-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}
.sk-card-meta {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}
.sk-card-labels {
  display: flex;
  gap: 0.25rem;
  flex-wrap: wrap;
}
.sk-card-actions {
  display: flex;
  gap: 0.5rem;
  margin-top: 0.25rem;
}
</style>
