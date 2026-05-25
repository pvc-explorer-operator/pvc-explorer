<template>
  <nav class="fe-crumb" aria-label="Directory navigation">
    <button class="fe-crumb__seg" @click="emit('navigate', '')" title="Root">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
        <polyline points="9 22 9 12 15 12 15 22"/>
      </svg>
    </button>
    <template v-for="(seg, i) in segments" :key="i">
      <span class="fe-crumb__sep" aria-hidden="true">/</span>
      <button
        class="fe-crumb__seg"
        :class="{ 'fe-crumb__seg--active': i === segments.length - 1 }"
        :disabled="i === segments.length - 1"
        @click="emit('navigate', segments.slice(0, i + 1).join('/'))"
      >{{ seg }}</button>
    </template>
  </nav>
</template>

<script setup lang="ts">
import { computed } from 'vue'
const props = defineProps<{ path: string }>()
const emit = defineEmits<{ (e: 'navigate', path: string): void }>()
const segments = computed(() => props.path ? props.path.split('/') : [])
</script>

<style scoped>
.fe-crumb {
  display: flex;
  align-items: center;
  gap: 1px;
  flex: 1;
  overflow: hidden;
}
.fe-crumb__seg {
  background: transparent;
  border: none;
  cursor: pointer;
  color: var(--muted, #5a6490);
  padding: 3px 6px;
  border-radius: 4px;
  font-family: Lato, sans-serif;
  font-size: 12px;
  white-space: nowrap;
  transition: color .15s, background .15s;
  display: flex;
  align-items: center;
}
.fe-crumb__seg svg { width: 13px; height: 13px; }
.fe-crumb__seg:hover:not(:disabled) { color: var(--text, #dde3f8); background: var(--surface2, #1c2038); }
.fe-crumb__seg--active { color: var(--text, #dde3f8); cursor: default; }
.fe-crumb__seg:disabled { opacity: 1; }
.fe-crumb__sep { color: var(--border, #252a42); font-size: 14px; padding: 0 1px; flex-shrink: 0; }
</style>
