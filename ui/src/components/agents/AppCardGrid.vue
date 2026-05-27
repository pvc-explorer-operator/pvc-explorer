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
  <div v-else class="card-grid">
    <div class="empty-state">
      <i class="pi pi-inbox empty-icon" />
      <div class="empty-msg">No explorers found.</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import AppCard from './AppCard.vue'
import type { Explorer } from '../../stores/explorerStore'

defineProps<{ explorers: Explorer[] }>()
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
</style>
