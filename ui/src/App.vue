<script setup lang="ts">
import { onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { useExplorerStore } from '@/stores/explorerStore'
import AppLayout from '@/layout/AppLayout.vue'

const route = useRoute()
const explorerStore = useExplorerStore()

onUnmounted(() => {
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
