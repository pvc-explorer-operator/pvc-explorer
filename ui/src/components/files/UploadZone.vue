<template>
  <label
    class="fe-upload-btn"
    :class="{ 'fe-upload-btn--disabled': readonly || uploading }"
    :title="uploading ? 'Upload in progress…' : 'Upload files'"
  >
    <svg v-if="!uploading" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <polyline points="16 16 12 12 8 16"/>
      <line x1="12" y1="12" x2="12" y2="21"/>
      <path d="M20.39 18.39A5 5 0 0 0 18 9h-1.26A8 8 0 1 0 3 16.3"/>
    </svg>
    <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="fe-upload-spin">
      <path d="M21 12a9 9 0 1 1-6.219-8.56"/>
    </svg>
    {{ uploading ? "Uploading…" : "Upload" }}
    <input
      type="file"
      multiple
      hidden
      :disabled="readonly || uploading"
      @change="onFileChange"
    />
  </label>
</template>

<script setup lang="ts">
const props = defineProps<{ readonly: boolean; uploading: boolean }>()
const emit = defineEmits<{ (e: 'upload', files: File[]): void }>()

function onFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  const files = Array.from(input.files ?? [])
  if (files.length) emit('upload', files)
  input.value = ''
}
</script>

<style scoped>
.fe-upload-btn {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 5px 10px;
  border-radius: 6px;
  border: 1px solid rgba(255,255,255,.2);
  background: rgba(255,255,255,.08);
  color: var(--text-color, #e2e8f0);
  cursor: pointer;
  font-family: Lato, sans-serif;
  font-size: 0.8125rem;
  white-space: nowrap;
  line-height: 1.4;
  transition: background .15s, border-color .15s;
  user-select: none;
}
.fe-upload-btn svg { width: 13px; height: 13px; flex-shrink: 0; }
.fe-upload-btn:hover:not(.fe-upload-btn--disabled) { background: rgba(255,255,255,.15); border-color: rgba(255,255,255,.35); color: #fff; }
.fe-upload-btn--disabled { opacity: .35; cursor: not-allowed; pointer-events: none; }

@keyframes fe-spin { to { transform: rotate(360deg) } }
.fe-upload-spin { animation: fe-spin 0.9s linear infinite; }
</style>
