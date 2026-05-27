<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useExplorerStore } from '@/stores/explorerStore'
import AppLayout from '@/layout/AppLayout.vue'

const route = useRoute()
const explorerStore = useExplorerStore()

// Set PrimeVue form-field tokens to match the topbar picker style.
// Guarded so dark-mode values only apply in dark mode — light mode
// lets the Lara preset handle its own (readable) defaults.
const DARK_FORM_FIELD = [
  ['--p-form-field-background',           'rgba(255,255,255,0.08)'],
  ['--p-form-field-disabled-background',  'rgba(255,255,255,0.04)'],
  ['--p-form-field-filled-background',    'rgba(255,255,255,0.08)'],
  ['--p-form-field-border-color',         'rgba(255,255,255,0.2)'],
  ['--p-form-field-hover-border-color',   'rgba(255,255,255,0.35)'],
  ['--p-form-field-color',                '#e2e8f0'],
  ['--p-form-field-placeholder-color',    '#94a3b8'],
  ['--p-form-field-disabled-color',       '#64748b'],
  ['--p-form-field-border-radius',        '6px'],
  ['--p-form-field-shadow',               'none'],
] as const

const formFieldObserver = ref<MutationObserver | null>(null)

function syncFormFieldTokens() {
  const s = document.documentElement.style
  const isDark = document.documentElement.classList.contains('app-dark')
  for (const [prop, darkVal] of DARK_FORM_FIELD) {
    if (isDark) {
      s.setProperty(prop, darkVal)
    } else {
      s.removeProperty(prop)
    }
  }
}

onMounted(() => {
  syncFormFieldTokens()
  // Re-sync when the user toggles dark/light mode
  const mo = new MutationObserver(syncFormFieldTokens)
  mo.observe(document.documentElement, { attributes: true, attributeFilter: ['class'] })
  formFieldObserver.value = mo
})
onUnmounted(() => {
  formFieldObserver.value?.disconnect()
  explorerStore.teardown()
})
</script>

<template>
  <app-layout v-if="route.meta.requiresAuth" />
  <router-view v-else />
</template>

<style>
/* Override PrimeVue form fields to match the topbar picker style (dark mode only).
   Light mode uses PrimeVue Lara default tokens for proper contrast. */
:root.app-dark .p-checkbox-box,
:root.app-dark .p-radiobutton-box,
:root.app-dark .p-inputtext,
:root.app-dark .p-textarea,
:root.app-dark .p-select,
:root.app-dark .p-multiselect,
:root.app-dark .p-datepicker-input,
:root.app-dark .p-inputnumber-input {
  background: rgba(255, 255, 255, 0.08) !important;
  border-color: rgba(255, 255, 255, 0.2) !important;
  color: #e2e8f0 !important;
  border-radius: 6px !important;
  box-shadow: none !important;
}

:root.app-dark .p-inputtext::placeholder,
:root.app-dark .p-textarea::placeholder {
  color: #94a3b8 !important;
}

:root.app-dark .p-checkbox-box:hover,
:root.app-dark .p-radiobutton-box:hover,
:root.app-dark .p-inputtext:hover,
:root.app-dark .p-select:hover,
:root.app-dark .p-multiselect:hover {
  border-color: rgba(255, 255, 255, 0.35) !important;
}

:root.app-dark .p-inputtext:focus,
:root.app-dark .p-textarea:focus {
  border-color: rgba(255, 255, 255, 0.5) !important;
  box-shadow: none !important;
  outline: none !important;
}

:root.app-dark .p-checkbox-checked .p-checkbox-box {
  background: var(--p-primary-color) !important;
  border-color: var(--p-primary-color) !important;
}

:root.app-dark .p-radiobutton-checked .p-radiobutton-box {
  border-color: var(--p-primary-color) !important;
}
</style>
