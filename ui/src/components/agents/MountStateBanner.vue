<template>
  <span :class="['mount-badge', `state-${state.toLowerCase()}`]">{{ state }}</span>
</template>

<script setup lang="ts">
defineProps<{ state: string }>()
</script>

<style scoped>
.mount-badge {
  display: inline-block;
  padding: 0.1em 0.6em;
  border-radius: 4px;
  font-size: var(--fs-sm);
  font-weight: 600;
  background: var(--surface-hover);
  color: var(--text-color-secondary);
}
.state-mounted   { background: #22c55e; color: #fff; }
.state-readonly  { background: #f59e0b; color: #000; }
.state-conflict  { background: #ef4444;  color: #fff; }
.state-unmounted { background: var(--surface-hover); color: var(--text-color-secondary); }
</style>
