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
        v-model="typed"
        class="chip-input"
        placeholder="key=value"
        role="combobox"
        aria-autocomplete="list"
        :aria-expanded="open && suggestions.length > 0"
        :aria-controls="listId"
        :aria-activedescendant="activeIndex >= 0 ? `${listId}-${activeIndex}` : undefined"
        @keydown="onKeydown"
        @keydown.backspace="onBackspace"
        @focus="open = true"
        @blur="onBlur"
      />
    </div>
    <ul v-if="open && suggestions.length" :id="listId" class="dropdown" role="listbox">
      <li
        v-for="(s, i) in suggestions"
        :key="s"
        :id="`${listId}-${i}`"
        role="option"
        :aria-selected="i === activeIndex"
        @mousedown.prevent="addSuggestion(s)"
        @mouseenter="activeIndex = i"
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
const activeIndex = ref(-1)
const listId = 'label-suggestions'

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
  activeIndex.value = -1
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
  setTimeout(() => { open.value = false; activeIndex.value = -1 }, 150)
}

function moveActive(dir: 1 | -1) {
  if (!open.value || !suggestions.value.length) return
  const max = suggestions.value.length - 1
  activeIndex.value = Math.min(Math.max(activeIndex.value + dir, 0), max)
}

function onEscape() {
  if (open.value) {
    open.value = false
    activeIndex.value = -1
  }
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'ArrowDown') { e.preventDefault(); moveActive(1) }
  else if (e.key === 'ArrowUp') { e.preventDefault(); moveActive(-1) }
  else if (e.key === 'Escape') { onEscape() }
  else if (e.key === 'Enter') {
    if (activeIndex.value >= 0) {
      e.preventDefault()
      addSuggestion(suggestions.value[activeIndex.value])
    } else {
      addTyped()
    }
  }
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
