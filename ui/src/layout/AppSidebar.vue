<script setup lang="ts">
import { useLayout } from '@/layout/composables/layout'
import { useExplorerStore } from '@/stores/explorerStore'
import { useAuthStore } from '@/stores/authStore'
import { onBeforeUnmount, ref, watch, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import AppMenu from './AppMenu.vue'
import FilterSidebar from '@/components/filters/FilterSidebar.vue'
import type { Filters } from '@/components/filters/FilterSidebar.vue'

const { layoutState, isDesktop, hasOpenOverlay, toggleMenu } = useLayout()
const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const sidebarRef = ref<HTMLElement | null>(null)
let outsideClickListener: ((e: MouseEvent) => void) | null = null

const explorerStore = useExplorerStore()
const filters = ref<Filters>({ search: '', phases: [], namespaces: [], mountStates: [], scopes: [], accessModes: [], consumers: '', created: '', labels: [] })
const isDashboard = computed(() => route.path === '/')

// Expose filters to HomeView via store — simplest approach: store the filters
// HomeView will watch explorerStore.sidebarFilters
watch(filters, (v) => { explorerStore.setSidebarFilters(v as any) }, { deep: true })

const filteredCount = computed(() => {
  const list = explorerStore.explorers
  const f = filters.value
  let count = 0
  for (const e of list) {
    if (f.phases.length && !f.phases.includes(e.phase)) continue
    if (f.namespaces.length && !f.namespaces.includes(e.namespace)) continue
    if (f.mountStates.length && !f.mountStates.includes(e.mountState)) continue
    if (f.scopes.length && (!e.scope || !f.scopes.includes(e.scope))) continue
    if (f.accessModes.length) {
      const m = e.accessMode || e.mode
      if (!m || !f.accessModes.includes(m)) continue
    }
    if (f.consumers === 'has' && !(e.consumerCount ?? 0)) continue
    if (f.consumers === 'none' && (e.consumerCount ?? 0) > 0) continue
    if (f.created) {
      if (!e.createdAt) continue
      const age = Date.now() - new Date(e.createdAt).getTime()
      if (f.created === '24h' && age >= 86_400_000) continue
      if (f.created === '7d' && age >= 604_800_000) continue
      if (f.created === '30d' && age >= 2_592_000_000) continue
      if (f.created === 'older' && age < 2_592_000_000) continue
    }
    if (f.labels.length) {
      if (!e.labels?.length) continue
      let ok = true
      for (const l of f.labels) { if (!e.labels.includes(l)) { ok = false; break } }
      if (!ok) continue
    }
    if (f.search) {
      const q = f.search.toLowerCase()
      if (!e.name.toLowerCase().includes(q) && !e.namespace.toLowerCase().includes(q) && !e.pvcName.toLowerCase().includes(q)) continue
    }
    count++
  }
  return count
})

watch(
  () => route.path,
  () => {
    layoutState.overlayMenuActive = false
    layoutState.mobileMenuActive = false
    layoutState.menuHoverActive = false
  },
)

watch(hasOpenOverlay, (val) => {
  if (isDesktop()) {
    if (val) bindOutsideClickListener()
    else unbindOutsideClickListener()
  }
})

const bindOutsideClickListener = () => {
  if (!outsideClickListener) {
    outsideClickListener = (event: MouseEvent) => {
      if (isOutsideClicked(event)) layoutState.overlayMenuActive = false
    }
    document.addEventListener('click', outsideClickListener)
  }
}

const unbindOutsideClickListener = () => {
  if (outsideClickListener) {
    document.removeEventListener('click', outsideClickListener)
    outsideClickListener = null
  }
}

const isOutsideClicked = (event: MouseEvent) => {
  const btn = document.querySelector('.layout-sidebar-toggle')
  return !(
    sidebarRef.value?.isSameNode(event.target as Node) ||
    sidebarRef.value?.contains(event.target as Node) ||
    btn?.isSameNode(event.target as Node) ||
    btn?.contains(event.target as Node)
  )
}

onBeforeUnmount(unbindOutsideClickListener)

const avatarInitial = computed(() => auth.username?.charAt(0).toUpperCase() ?? '?')

const handleLogout = async () => {
  await auth.logout()
  router.push('/login')
}

const expandIfCollapsed = () => {
  if (isDesktop() && layoutState.staticMenuInactive) {
    layoutState.staticMenuInactive = false
  }
}

const showFilters = ref(true)

const toggleFilters = () => {
  expandIfCollapsed()
  showFilters.value = !showFilters.value
}
</script>

<template>
  <div ref="sidebarRef" class="layout-sidebar">

    <!-- ① Logo + Hamburger at top -->
    <div class="layout-sidebar-header">
      <router-link to="/" class="layout-sidebar-logo">
        <img src="/logo-icon.svg" class="layout-sidebar-logo-icon layout-sidebar-logo-svg" alt="pvc-explorer icon" />
        <span class="layout-sidebar-logo-text">pvc-explorer</span>
      </router-link>
      <button class="layout-sidebar-toggle" @click="toggleMenu" title="Toggle sidebar">
        <i class="pi pi-bars"></i>
      </button>
    </div>

    <!-- ② Nav items -->
    <div class="layout-sidebar-scroll">
      <AppMenu />

      <!-- Filter toggle button (dashboard only) -->
      <template v-if="isDashboard">
        <button
          class="sidebar-filter-btn"
          :class="{ active: showFilters }"
          @click="toggleFilters"
          title="Filters"
        >
          <i class="pi pi-fw pi-filter"></i>
          <span class="layout-menuitem-text">Filters</span>
          <i class="pi pi-chevron-down sidebar-filter-chevron" :class="{ rotated: showFilters }"></i>
        </button>

        <div v-if="showFilters" class="layout-sidebar-filters" @click.capture="expandIfCollapsed">
          <FilterSidebar
            :explorers="explorerStore.explorers"
            :shown="filteredCount"
            :total="explorerStore.explorers.length"
            @update:filters="filters = $event"
          />
        </div>
      </template>
    </div>

    <!-- ③ Footer: user + logo at bottom -->
    <div class="layout-sidebar-footer">
      <template v-if="auth.isAuthenticated">
        <span class="layout-sidebar-footer-avatar">{{ avatarInitial }}</span>
        <div class="layout-sidebar-footer-user">
          <div class="layout-sidebar-footer-name">{{ auth.username }}</div>
        </div>
        <button class="layout-sidebar-footer-logout" @click="handleLogout" title="Logout">
          <i class="pi pi-sign-out"></i>
        </button>
      </template>
    </div>

  </div>
</template>
