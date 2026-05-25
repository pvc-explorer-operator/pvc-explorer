<template>
  <div class="card-grid">
    <template v-if="explorers.length">
      <AppCard
        v-for="explorer in explorers"
        :key="`${explorer.namespace}/${explorer.name}`"
        :explorer="explorer"
      />
    </template>
    <template v-else>
      <div class="empty-state">
        <i class="pi pi-inbox empty-icon" />
        <div class="empty-msg">No explorers found.</div>
      </div>
    </template>
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
</style>
