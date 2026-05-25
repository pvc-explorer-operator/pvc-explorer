<template>
  <div>
    <div class="flex gap-4 items-center flex-wrap">
      <!-- Search -->
      <IconField>
        <InputIcon class="pi pi-search" />
        <InputText v-model="search" placeholder="Search name, PVC..." class="filter-search" @input="emit_()" />
      </IconField>

      <!-- Phase filter chips -->
      <div v-for="opt in phaseOptions" :key="opt.value" class="flex items-center gap-2">
        <Checkbox v-model="phases" :value="opt.value" :inputId="'phase-' + opt.value" @change="emit_()" />
        <label :for="'phase-' + opt.value">
          <Tag :value="opt.label" :severity="opt.severity" rounded />
        </label>
      </div>

      <!-- Namespace select -->
      <Dropdown v-model="namespaceFilter" :options="namespaceOptions" placeholder="Namespace" class="filter-dropdown" @change="onNamespaceChange" />

      <!-- Count & Clear -->
      <span class="text-muted-color text-sm ml-auto">{{ shown }}/{{ total }} shown</span>

      <Button v-if="activeCount > 0" label="Clear" text size="small" severity="secondary" @click="clearAll" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { apiFetch } from '@/composables/useAuth'
import Checkbox from 'primevue/checkbox'
import Dropdown from 'primevue/dropdown'
import Tag from 'primevue/tag'
import Button from 'primevue/button'
import IconField from 'primevue/iconfield'
import InputIcon from 'primevue/inputicon'
import InputText from 'primevue/inputtext'
import type { Explorer } from '../../stores/explorerStore'

export interface Filters {
  search: string
  phases: string[]
  namespaces: string[]
  mountStates: string[]
}

const props = defineProps<{ explorers: Explorer[]; shown: number; total: number }>()
const emit = defineEmits<{ (e: 'update:filters', v: Filters): void }>()

const search = ref('')
const phases = ref<string[]>([])
const namespaces = ref<string[]>([])
const mountStates = ref<string[]>([])
const namespaceFilter = ref<string | null>(null)
const namespaceOptions = ref<string[]>([])

const phaseOptions = [
  { value: 'Running',      label: 'Running',      severity: 'success'   },
  { value: 'ScaledToZero', label: 'Scaled to Zero', severity: 'secondary' },
  { value: 'Waking',       label: 'Waking',       severity: 'info'      },
  { value: 'Pending',      label: 'Pending',      severity: 'warn'      },
  { value: 'Failed',       label: 'Failed',       severity: 'danger'    },
]

const activeCount = computed(
  () => phases.value.length + namespaces.value.length + mountStates.value.length + (search.value ? 1 : 0)
)

function emit_() {
  emit('update:filters', {
    search: search.value,
    phases: phases.value,
    namespaces: namespaces.value,
    mountStates: mountStates.value,
  })
}

function onNamespaceChange() {
  namespaces.value = namespaceFilter.value ? [namespaceFilter.value] : []
  emit_()
}

function clearAll() {
  search.value = ''
  phases.value = []
  namespaces.value = []
  mountStates.value = []
  namespaceFilter.value = null
  emit_()
}

let debounce: ReturnType<typeof setTimeout> | null = null
watch(search, () => {
  if (debounce) clearTimeout(debounce)
  debounce = setTimeout(emit_, 250)
})

import { onMounted } from 'vue'
onMounted(async () => {
  const res = await apiFetch('/api/v1/namespaces')
  if (res.ok) {
    const data: { name: string }[] = await res.json()
    namespaceOptions.value = data.map(n => n.name)
  }
})
</script>

<style scoped>
.filter-search {
  min-width: 200px;
}
.filter-dropdown {
  min-width: 160px;
}
</style>
