<template>
  <div class="sc-grid">
    <template v-if="scopes.length">
      <ScopeCard v-for="scope in scopes" :key="scope.name" :scope="scope" />
    </template>
    <div v-else class="empty-state">
      <i class="pi pi-inbox empty-icon" />
      <div class="empty-msg">No scopes found.</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import ScopeCard from './ScopeCard.vue'
import type { Scope } from '../../stores/explorerStore'

defineProps<{ scopes: Scope[] }>()
</script>

<style scoped>
.sc-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
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
  grid-column: 1 / -1;
}
.empty-icon {
  font-size: 2.5rem;
  margin-bottom: 0.7rem;
}
.empty-msg {
  font-size: 1.08rem;
}
</style>
