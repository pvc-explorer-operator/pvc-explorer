<template>
  <span class="label-chip">
    {{ label }}
    <button class="remove" @click.stop="$emit('remove', label)" aria-label="Remove">×</button>
  </span>
</template>

<script setup lang="ts">
defineProps<{ label: string }>()
defineEmits<{ (e: 'remove', label: string): void }>()
</script>

<style scoped>
.label-chip {
  display: inline-flex;
  align-items: center;
  gap: 0.2rem;
  background: var(--surface-hover);
  color: var(--text-color-secondary);
  border-radius: 4px;
  padding: 0.1em 0.4em 0.1em 0.55em;
  font-size: 0.875rem;
}
.remove {
  background: none;
  border: none;
  color: var(--text-color-secondary);
  cursor: pointer;
  padding: 0;
  line-height: 1;
  font-size: 1.1em;
}
.remove:hover { color: var(--p-red-500); }
</style>
