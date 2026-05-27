<template>
  <span v-if="props.remainingSeconds !== null" :class="['idle-timer', { warning: props.remainingSeconds < 60 }]">
    <i class="pi pi-clock" style="margin-right:0.25rem"></i> Idle: {{ formatted }}
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{ remainingSeconds: number | null }>()

const formatted = computed(() => {
  const s = Math.max(0, props.remainingSeconds ?? 0)
  const m = Math.floor(s / 60)
  const sec = s % 60
  return m > 0 ? `${m}m ${sec}s` : `${sec}s`
})
</script>

<style scoped>
.idle-timer { font-size: var(--fs-body); color: var(--text-color-secondary); }
.idle-timer.warning { color: var(--p-red-500); font-weight: 600; }
</style>
