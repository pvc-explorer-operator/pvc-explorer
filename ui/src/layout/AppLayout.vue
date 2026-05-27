<script setup lang="ts">
import { useLayout } from '@/layout/composables/layout'
import { computed } from 'vue'
import AppSidebar from './AppSidebar.vue'
import AppTopbar from './AppTopbar.vue'
import { useKeyboardShortcuts } from '@/composables/useKeyboardShortcuts'
import KeyboardShortcutsModal from '@/components/shared/KeyboardShortcutsModal.vue'

const { layoutConfig, layoutState, hideMobileMenu } = useLayout()
useKeyboardShortcuts()

const containerClass = computed(() => ({
  'layout-overlay': layoutConfig.menuMode === 'overlay',
  'layout-static': layoutConfig.menuMode === 'static',
  'layout-overlay-active': layoutState.overlayMenuActive,
  'layout-mobile-active': layoutState.mobileMenuActive,
  'layout-static-inactive': layoutState.staticMenuInactive,
}))
</script>

<template>
  <div class="layout-wrapper" :class="containerClass">
    <a href="#main-content" class="skip-link">Skip to content</a>
    <AppTopbar />
    <AppSidebar />
    <div class="layout-main-container" :inert="layoutState.mobileMenuActive || undefined">
      <div id="main-content" class="layout-main">
        <router-view />
      </div>
      <div class="layout-footer">
        <span>pvc-explorer</span>
      </div>
    </div>
    <div class="layout-mask animate-fadein" @click="hideMobileMenu" />
  </div>
  <Toast />
  <KeyboardShortcutsModal />
</template>
