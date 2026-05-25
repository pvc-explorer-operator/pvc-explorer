<template>
  <div class="autocomplete">
    <div class="chips">
      <LabelFilter
        v-for="label in modelValue"
        :key="label"
        :label="label"
        @remove="remove"
      />
      <input
        ref="inputEl"
        v-model="typed"
        class="chip-input"
        placeholder="key=value"
        @keydown.enter.prevent="addTyped"
        @keydown.backspace="onBackspace"
        @focus="open = true"
        @blur="onBlur"
      />
    </div>
    <ul v-if="open && suggestions.length" class="dropdown">
      <li
        v-for="s in suggestions"
        :key="s"
        @mousedown.prevent="addSuggestion(s)"
      >{{ s }}</li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { apiFetch } from '@/composables/useAuth'
import { ref, computed, onMounted } from 'vue'
import LabelFilter from './LabelFilter.vue'

const props = defineProps<{ modelValue: string[] }>()
const emit = defineEmits<{ (e: 'update:modelValue', v: string[]): void }>()

const allLabels = ref<string[]>([])
const typed = ref('')
const open = ref(false)


const suggestions = computed(() =>
  typed.value.length
    ? allLabels.value.filter(l => l.includes(typed.value) && !props.modelValue.includes(l))
    : []
)

onMounted(async () => {
  const res = await apiFetch('/api/v1/labels')
  if (res.ok) allLabels.value = await res.json()
})

function add(label: string) {
  const trimmed = label.trim()
  if (trimmed && !props.modelValue.includes(trimmed)) {
    emit('update:modelValue', [...props.modelValue, trimmed])
  }
  typed.value = ''
  open.value = false
}

function addTyped() { add(typed.value) }
function addSuggestion(s: string) { add(s) }

function remove(label: string) {
  emit('update:modelValue', props.modelValue.filter(l => l !== label))
}

function onBackspace() {
  if (!typed.value && props.modelValue.length) {
    emit('update:modelValue', props.modelValue.slice(0, -1))
  }
}

function onBlur() {
  setTimeout(() => { open.value = false }, 150)
}
</script>

<style scoped>
.autocomplete { position: relative; }
.chips {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 0.25rem;
  background: var(--surface-hover);
  border: 1px solid var(--surface-border);
  border-radius: 5px;
  padding: 0.25rem 0.5rem;
  min-height: 2rem;
  cursor: text;
}
.chip-input {
  border: none;
  background: transparent;
  color: var(--text-color);
  font-size: 0.9rem;
  outline: none;
  min-width: 80px;
  flex: 1;
}
.dropdown {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: var(--surface-card);
  border: 1px solid var(--surface-border);
  border-radius: 5px;
  margin: 2px 0 0;
  padding: 0;
  list-style: none;
  z-index: 100;
  max-height: 160px;
  overflow-y: auto;
}
.dropdown li {
  padding: 0.35rem 0.7rem;
  font-size: 0.9rem;
  color: var(--text-color-secondary);
  cursor: pointer;
}
.dropdown li:hover { background: var(--surface-hover); }
</style>
