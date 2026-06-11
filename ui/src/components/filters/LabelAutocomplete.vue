<template>
  <div class="autocomplete">
    <div class="chips" @click="focusInput">
      <LabelChip
        v-for="label in modelValue"
        :key="label"
        :label="label"
        removable
        @remove="remove"
      />
      <input
        ref="inputRef"
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
        @input="open = true"
        @blur="onBlur"
      />
    </div>
    <ul v-if="open && suggestions.length" :id="listId" class="dropdown" role="listbox">
      <li
        v-for="(s, i) in suggestions"
        :key="s"
        :id="`${listId}-${i}`"
        role="option"
        :style="suggestionStyle(s, i)"
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
import { stringToColor } from '@/composables/useFilterColors'
import LabelChip from './LabelChip.vue'

const props = defineProps<{ modelValue: string[] }>()
const emit = defineEmits<{ (e: 'update:modelValue', v: string[]): void }>()

const inputRef = ref<HTMLInputElement | null>(null)
const allLabels = ref<string[]>([])
const typed = ref('')
const open = ref(false)
const activeIndex = ref(-1)
const listId = 'label-suggestions'

function focusInput() { inputRef.value?.focus() }

function suggestionStyle(label: string, i: number) {
  const key = label.split('=')[0] || label
  const isActive = i === activeIndex.value
  const color = stringToColor(key)
  return {
    background: isActive ? color + '25' : 'transparent',
    color,
    borderLeft: `3px solid ${color}`,
  }
}

const suggestions = computed(() =>
  allLabels.value.filter(l => {
    if (props.modelValue.includes(l)) return false
    if (!typed.value.length) return true
    return l.includes(typed.value)
  })
)

function flattenLabels(raw: unknown): string[] {
  if (Array.isArray(raw)) return raw as string[]
  if (raw && typeof raw === 'object') {
    const out: string[] = []
    for (const [k, vs] of Object.entries(raw)) {
      if (Array.isArray(vs)) for (const v of vs) out.push(`${k}=${v}`)
    }
    return out
  }
  return []
}

onMounted(async () => {
  try {
    const res = await apiFetch('/api/v1/labels')
    if (res.ok) allLabels.value = flattenLabels(await res.json())
  } catch { /* no labels available */ }
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
