<script setup lang="ts">
import { useAuthStore } from '@/stores/authStore'
import { computed } from 'vue'
import AppMenuItem from './AppMenuItem.vue'

const auth = useAuthStore()

const model = computed(() => [
  {
    label: 'Home',
    items: [
      { label: 'Dashboard', icon: 'pi pi-fw pi-home',          to: '/' },
      { label: 'Scopes',    icon: 'pi pi-fw pi-server',        to: '/scopes' },
      { label: 'About',     icon: 'pi pi-fw pi-info-circle',   to: '/about' },
    ],
  },
  ...(auth.isAdmin
    ? [{
        label: 'Admin',
        items: [
          { label: 'Settings', icon: 'pi pi-fw pi-cog', to: '/settings' },
        ],
      }]
    : []),
])
</script>

<template>
  <ul class="layout-menu">
    <template v-for="(item, i) in model" :key="i">
      <app-menu-item v-if="!(item as any).separator" :item="item" :index="i" :root="true" />
      <li v-if="(item as any).separator" class="menu-separator"></li>
    </template>
  </ul>
</template>

<style lang="scss" scoped></style>
