<template>
  <div class="fe-preview" v-if="file">
    <div class="fe-preview__header">
      <i :class="icon" class="fe-preview__icon" />
      <span class="fe-preview__name">{{ file.name }}</span>
      <span class="fe-preview__size">{{ fmtSize(file.size) }}</span>
      <Button label="Download" text size="small" icon="pi pi-download"
              @click="emit('download', file)" class="ml-auto" />
    </div>
    <img v-if="isImage" :src="downloadUrl" :alt="file.name" class="fe-preview__img" />
    <embed v-else-if="isPdf" :src="downloadUrl" type="application/pdf" class="fe-preview__pdf" />
    <div v-else-if="tooLarge" class="fe-preview__msg">
      File is too large to preview ({{ fmtSize(file.size) }}).
    </div>
    <div v-else class="fe-preview__msg">
      No preview available for this file type.
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import Button from 'primevue/button'
import type { FileEntry } from '../../api/files'

const props = defineProps<{
  file: FileEntry | null
  downloadUrl: string
  tooLarge?: boolean
}>()
const emit = defineEmits<{ (e: 'download', file: FileEntry): void }>()

const ext = computed(() => props.file?.name.split('.').pop()?.toLowerCase() ?? '')
const isImage = computed(() => ['png','jpg','jpeg','gif','svg','webp'].includes(ext.value))
const isPdf   = computed(() => ext.value === 'pdf')
const icon    = computed(() => isImage.value ? 'pi pi-image' : isPdf.value ? 'pi pi-file-pdf' : 'pi pi-file')

function fmtSize(bytes: number) {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1_048_576) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / 1_048_576).toFixed(1)} MB`
}
</script>

<style scoped>
.fe-preview {
  border-top: 1px solid var(--surface-border);
  padding: 12px;
}
.fe-preview__header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
}
.fe-preview__icon { font-size: 1.1rem; color: var(--primary-400); }
.fe-preview__name { font-weight: 600; }
.fe-preview__size { font-size: 0.8rem; color: var(--text-color-secondary); }
.fe-preview__img {
  max-width: 100%;
  max-height: 420px;
  display: block;
  margin: 0 auto;
}
.fe-preview__pdf { width: 100%; height: 500px; border: none; }
.fe-preview__msg { color: var(--text-color-secondary); font-size: 0.9rem; }
</style>
