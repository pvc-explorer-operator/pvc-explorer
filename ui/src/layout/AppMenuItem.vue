<script setup lang="ts">
import { useLayout } from '@/layout/composables/layout'
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const { layoutState, isDesktop } = useLayout()

const props = defineProps<{
  item: {
    label?: string
    icon?: string
    to?: string
    url?: string
    target?: string
    path?: string
    items?: any[]
    separator?: boolean
    visible?: unknown
    disabled?: boolean
    class?: string
    command?: (args: { originalEvent: Event; item: any }) => void
  }
  root?: boolean
  parentPath?: string
}>()

const fullPath = computed(() =>
  props.item.path
    ? props.parentPath
      ? props.parentPath + props.item.path
      : props.item.path
    : null
)

const isActive = computed(() => {
  if (props.item.path) return layoutState.activePath?.startsWith(fullPath.value ?? '') ?? false
  return layoutState.activePath === props.item.to
})

const itemClick = (event: MouseEvent, item: typeof props.item) => {
  if (item.disabled) { event.preventDefault(); return }
  // Auto-expand sidebar when collapsed (ArgoCD behaviour)
  if (isDesktop() && layoutState.staticMenuInactive) {
    layoutState.staticMenuInactive = false
  }
  if (item.command) item.command({ originalEvent: event, item })
  if (item.items) {
    if (isActive.value) {
      layoutState.activePath = layoutState.activePath!.replace(item.path!, '')
    } else {
      layoutState.activePath = fullPath.value
      layoutState.menuHoverActive = true
    }
  } else {
    layoutState.overlayMenuActive = false
    layoutState.mobileMenuActive = false
    layoutState.menuHoverActive = false
  }
}

const onMouseEnter = () => {
  if (isDesktop() && props.root && props.item.items && layoutState.menuHoverActive) {
    layoutState.activePath = fullPath.value
  }
}
</script>

<template>
  <li v-if="(item.visible !== false)" :class="{ 'layout-root-menuitem': root, 'active-menuitem': isActive }">
    <div v-if="root && (item.visible !== false)" class="layout-menuitem-root-text">{{ item.label }}</div>
    <a
      v-if="(!item.to || item.items) && (item.visible !== false)"
      :href="item.url"
      @click="itemClick($event, item)"
      :class="item.class"
      :target="item.target"
      tabindex="0"
      @mouseenter="onMouseEnter"
    >
      <i :class="item.icon" class="layout-menuitem-icon" />
      <span class="layout-menuitem-text">{{ item.label }}</span>
      <i class="pi pi-fw pi-angle-down layout-submenu-toggler" v-if="item.items" />
    </a>
    <router-link
      v-if="item.to && !item.items && (item.visible !== false)"
      @click="itemClick($event, item)"
      exactActiveClass="active-route"
      :class="item.class"
      tabindex="0"
      :to="item.to"
      :aria-current="route.path === item.to ? 'page' : undefined"
      @mouseenter="onMouseEnter"
    >
      <i :class="item.icon" class="layout-menuitem-icon" />
      <span class="layout-menuitem-text">{{ item.label }}</span>
    </router-link>
    <Transition v-if="item.items && (item.visible !== false)" name="layout-submenu">
      <ul v-show="root ? true : isActive" class="layout-submenu">
        <app-menu-item
          v-for="child in item.items"
          :key="child.label"
          :item="child"
          :root="false"
          :parentPath="fullPath ?? undefined"
        />
      </ul>
    </Transition>
  </li>
</template>
