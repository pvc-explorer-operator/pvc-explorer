<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'
import AppLayout from '@/layout/AppLayout.vue'

const route = useRoute()

// Force PrimeVue form-field tokens to match the topbar picker style.
// Must run after mount so PrimeVue's dynamic <style> injection has already happened.
onMounted(() => {
  const r = document.documentElement.style
  r.setProperty('--p-form-field-background',           'rgba(255,255,255,0.08)')
  r.setProperty('--p-form-field-disabled-background',  'rgba(255,255,255,0.04)')
  r.setProperty('--p-form-field-filled-background',    'rgba(255,255,255,0.08)')
  r.setProperty('--p-form-field-border-color',         'rgba(255,255,255,0.2)')
  r.setProperty('--p-form-field-hover-border-color',   'rgba(255,255,255,0.35)')
  r.setProperty('--p-form-field-color',                '#e2e8f0')
  r.setProperty('--p-form-field-placeholder-color',    '#94a3b8')
  r.setProperty('--p-form-field-disabled-color',       '#64748b')
  r.setProperty('--p-form-field-border-radius',        '6px')
  r.setProperty('--p-form-field-shadow',               'none')
})
</script>

<template>
  <app-layout v-if="route.meta.requiresAuth" />
  <router-view v-else />
</template>

<style>
/* Override PrimeVue form fields to match the topbar picker style */
.p-checkbox-box,
.p-radiobutton-box,
.p-inputtext,
.p-textarea,
.p-select,
.p-multiselect,
.p-datepicker-input,
.p-inputnumber-input {
  background: rgba(255, 255, 255, 0.08) !important;
  border-color: rgba(255, 255, 255, 0.2) !important;
  color: #e2e8f0 !important;
  border-radius: 6px !important;
  box-shadow: none !important;
}

.p-inputtext::placeholder,
.p-textarea::placeholder {
  color: #94a3b8 !important;
}

.p-checkbox-box:hover,
.p-radiobutton-box:hover,
.p-inputtext:hover,
.p-select:hover,
.p-multiselect:hover {
  border-color: rgba(255, 255, 255, 0.35) !important;
}

.p-inputtext:focus,
.p-textarea:focus {
  border-color: rgba(255, 255, 255, 0.5) !important;
  box-shadow: none !important;
  outline: none !important;
}

.p-checkbox-checked .p-checkbox-box {
  background: var(--p-primary-color) !important;
  border-color: var(--p-primary-color) !important;
}

.p-radiobutton-checked .p-radiobutton-box {
  border-color: var(--p-primary-color) !important;
}
</style>
