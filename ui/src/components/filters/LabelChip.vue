<template>
  <span
    class="label-chip"
    :style="chipStyle"
  >
    {{ label }}
    <button
      v-if="removable"
      class="remove"
      @click.stop="$emit('remove', label)"
      aria-label="Remove"
    >&times;</button>
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { stringToColor } from '../../composables/useFilterColors'

const props = defineProps<{ label: string; removable?: boolean }>()
defineEmits<{ (e: 'remove', label: string): void }>()

const chipStyle = computed(() => {
  const color = stringToColor(props.label.split('=')[0] || props.label)
  return {
    background: color + '18',
    color,
    borderColor: color + '35',
  }
})
</script>

<style scoped>
.label-chip {
  display: inline-flex;
  align-items: center;
  gap: 0.2rem;
  border: 1px solid;
  border-radius: 4px;
  padding: 1px 7px;
  font-size: 0.7rem;
  font-family: 'JetBrains Mono', monospace;
}
.remove {
  background: none;
  border: none;
  cursor: pointer;
  padding: 0;
  line-height: 1;
  font-size: 1em;
  color: inherit;
  opacity: 0.7;
}
.remove:hover { opacity: 1; }
</style>
