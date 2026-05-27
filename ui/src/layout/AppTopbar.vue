<script setup lang="ts">
import { ref, computed, inject, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useLayout } from '@/layout/composables/layout'
import { useExplorerStore } from '@/stores/explorerStore'
import { shortcutsModalOpen } from '@/composables/useShortcutsModal'

const { toggleMenu, toggleDarkMode, isDarkTheme } = useLayout()
const explorerStore = useExplorerStore()
const startTour = inject<() => void>('startTour', () => {})
const route  = useRoute()
const router = useRouter()

// ── Breadcrumb nav ────────────────────────────────────────────────────────────
const navInfo = computed(() => {
  const { name, params } = route
  const ns = params.ns as string
  const n  = params.name as string

  switch (name) {
    case 'Home':         return { back: null,                      segments: ['Explorers'] }
    case 'ScopeList':    return { back: '/',                        segments: ['Scopes'] }
    case 'ScopeDetail':  return { back: '/scopes',                  segments: ['Scopes', n] }
    case 'CreateScope':  return { back: '/scopes',                  segments: ['Scopes', 'Create'] }
    case 'CreateAgent':  return { back: '/',                        segments: ['Explorers', 'Create'] }
    case 'AgentDetail':  return { back: '/',                        segments: [ns, n] }
    case 'FileBrowser':  return { back: `/explorers/${ns}/${n}`,   segments: [n, 'Files'] }
    case 'Settings':     return { back: '/',                        segments: ['Settings'] }
    case 'About':        return { back: '/',                        segments: ['About'] }
    default:             return { back: null,                       segments: [] }
  }
})

const handleRefresh = async () => {
  await explorerStore.fetchExplorers()
  const now = new Date()
  lastRefreshed.value = now.toLocaleTimeString()
}

/* ── Auto-poll ── */
const pollOptions = [
  { label: 'Off', value: 0 },
  { label: '15s', value: 15_000 },
  { label: '30s', value: 30_000 },
  { label: '1m', value: 60_000 },
  { label: '2m', value: 120_000 },
]
const pollInterval = ref(30_000)
const lastRefreshed = ref('')
let pollTimer: ReturnType<typeof setInterval> | null = null

function startPoll() {
  stopPoll()
  if (!pollInterval.value) return
  pollTimer = setInterval(() => {
    explorerStore.fetchExplorers().then(() => {
      const now = new Date()
      lastRefreshed.value = now.toLocaleTimeString()
    }).catch(() => {})
  }, pollInterval.value)
}

function stopPoll() {
  if (pollTimer) {
    clearInterval(pollTimer)
    pollTimer = null
  }
}

function onPollChange() {
  startPoll()
}

onMounted(() => {
  explorerStore.fetchExplorers().then(() => {
    lastRefreshed.value = new Date().toLocaleTimeString()
  }).catch(() => {})
  startPoll()
})

onUnmounted(() => {
  stopPoll()
})
</script>

<template>
  <div class="layout-topbar">
    <!-- Hamburger for mobile only (visibility controlled by CSS media query) -->
    <button
      type="button"
      class="layout-topbar-hamburger"
      @click="toggleMenu"
      title="Open menu"
      aria-label="Open menu"
    >
      <i class="pi pi-bars" aria-hidden="true"></i>
    </button>
    <!-- Breadcrumb navigation -->
    <div class="topbar-nav">
      <button
        v-if="navInfo.back"
        class="topbar-back-btn"
        title="Go back"
        aria-label="Go back"
        @click="router.push(navInfo.back)"
      >
        <i class="pi pi-arrow-left" aria-hidden="true" />
      </button>
      <span v-if="navInfo.segments.length" class="topbar-breadcrumb">
        <template v-for="(seg, i) in navInfo.segments" :key="i">
          <span v-if="i > 0" class="topbar-sep" aria-hidden="true">›</span>
          <span
            :class="i === navInfo.segments.length - 1 ? 'topbar-crumb topbar-crumb--current' : 'topbar-crumb topbar-crumb--parent'"
            :aria-current="i === navInfo.segments.length - 1 ? 'page' : undefined"
          >{{ seg }}</span>
        </template>
      </span>
    </div>

    <div class="layout-topbar-actions">
      <div class="topbar-poll-group">
        <select v-model="pollInterval" class="topbar-poll-picker" @change="onPollChange" title="Auto-refresh interval">
          <option v-for="opt in pollOptions" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
        </select>
        <button type="button" class="layout-topbar-action" @click="handleRefresh" title="Refresh now" aria-label="Refresh now">
          <i class="pi pi-refresh" aria-hidden="true"></i>
        </button>
        <span class="topbar-timestamp"><span class="topbar-sync-label">Last sync </span>{{ lastRefreshed || '—' }}</span>
      </div>
      <div class="topbar-divider"></div>
      <button type="button" class="layout-topbar-action" @click="toggleDarkMode" :title="isDarkTheme ? 'Light mode' : 'Dark mode'" :aria-label="isDarkTheme ? 'Switch to light mode' : 'Switch to dark mode'">
        <i :class="['pi', isDarkTheme ? 'pi-sun' : 'pi-moon']" aria-hidden="true"></i>
      </button>
      <button type="button" class="layout-topbar-action" @click="startTour" title="Start welcome tour" aria-label="Start welcome tour">
        <i class="pi pi-info-circle" aria-hidden="true"></i>
      </button>
      <button type="button" class="layout-topbar-action" @click="shortcutsModalOpen = !shortcutsModalOpen" title="Keyboard shortcuts" aria-label="Toggle keyboard shortcuts">
        <i class="pi pi-question-circle" aria-hidden="true"></i>
      </button>
    </div>
  </div>
</template>

<style scoped>
.topbar-nav {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  flex: 1;
  min-width: 0;
  overflow: hidden;
  margin-left: 0.75rem;
}

.topbar-back-btn {
  display: inline-flex;
  justify-content: center;
  align-items: center;
  flex-shrink: 0;
  width: 2rem;
  height: 2rem;
  border-radius: 50%;
  background: rgba(255,255,255,0.08);
  border: none;
  color: var(--sidebar-text, #e2e8f0);
  cursor: pointer;
  transition: background 0.15s;
}
.topbar-back-btn:hover { background: rgba(255,255,255,0.18); }
.topbar-back-btn .pi { font-size: 0.8rem; }

.topbar-breadcrumb {
  display: flex;
  align-items: center;
  gap: 0.3rem;
  overflow: hidden;
}
.topbar-crumb {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.topbar-crumb--parent {
  color: rgba(255,255,255,0.38);
  font-size: 0.8125rem;
}
.topbar-crumb--current {
  color: var(--sidebar-text, #e2e8f0);
  font-size: 0.875rem;
  font-weight: 600;
}
.topbar-sep {
  color: rgba(255,255,255,0.22);
  font-size: 0.75rem;
  flex-shrink: 0;
}

.topbar-poll-group {
  display: flex;
  align-items: center;
  gap: 0.35rem;
}
.topbar-poll-picker {
  height: 1.9rem;
  padding: 0 0.4rem;
  border: 1px solid rgba(255,255,255,0.2);
  border-radius: 6px;
  background: rgba(255,255,255,0.08);
  color: #e2e8f0;
  font-size: 0.8rem;
  cursor: pointer;
  outline: none;
  width: 5rem;
}
.topbar-poll-picker:focus {
  border-color: var(--primary-color);
}
.topbar-timestamp {
  font-size: 0.72rem;
  color: #94a3b8;
  white-space: nowrap;
  line-height: 1;
}
.topbar-sync-label {
  color: #64748b;
  margin-right: 0.15rem;
}
.topbar-divider {
  width: 1px;
  height: 1.5rem;
  background: rgba(255,255,255,0.15);
  margin: 0 0.25rem;
}
</style>
