<template>
  <FileExplorerApp
    :fetch-files="api.fetchFiles"
    :fetch-content="api.fetchContent"
    :save-file="api.saveFile"
    :delete-file="api.deleteFile"
    :rename-file="api.renameFile"
    :upload-files="api.uploadFiles"
    :create-file="api.createFile"
    :download-url="api.downloadUrl"
    :readonly="config.readonly"
    :explorer-label="`${ns} / ${name}`"
    :remaining-seconds="remainingSeconds"
    :idle-warning="idleWarning"
    @heartbeat="heartbeat"
  />
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import FileExplorerApp from '../components/files/FileExplorerApp.vue'
import { createFileApi } from '../api/files'
import { useWebSocket } from '../composables/useWebSocket'
import type { AgentConfig } from '../api/files'

/* ── Route params ──────────────────────────────────────────────────────────── */
const route = useRoute()
const router = useRouter()
const ns   = route.params.ns   as string
const name = route.params.name as string
const base = `/api/v1/explorers/${ns}/${name}`

/* ── API layer (all fetch calls live here, not in child components) ────────── */
const api = createFileApi(`${base}/proxy/api`)

/* ── Agent config + idle state ─────────────────────────────────────────────── */
const config           = ref<AgentConfig>({ readonly: false, forceRW: false, pvc: '', namespace: '', pod: '', cluster: '', version: '' })
const remainingSeconds = ref(600)
const idleWarning      = ref(false)
let secondsTimer: ReturnType<typeof setInterval> | null = null

/* ── WebSocket (idle ticks, warnings, expiry) ──────────────────────────────── */
const { connect: wsConnect, disconnect: wsDisconnect } = useWebSocket({
  onIdleTick(p) {
    if (p.namespace === ns && p.name === name) {
      remainingSeconds.value = p.remainingSeconds
      if (p.remainingSeconds > 60) idleWarning.value = false
    }
  },
  onIdleWarning(p) {
    if (p.namespace === ns && p.name === name) {
      remainingSeconds.value = p.remainingSeconds
      idleWarning.value = true
    }
  },
  onIdleExpired(p) {
    if (p.namespace === ns && p.name === name) router.push('/')
  },
})

/* ── Heartbeat ─────────────────────────────────────────────────────────────── */
async function heartbeat() {
  const res = await fetch(`${base}/heartbeat`, { method: 'POST' })
  if (!res.ok) return
  const data: { remainingSeconds: number; phase: string } = await res.json()
  remainingSeconds.value = data.remainingSeconds
  if (data.phase !== 'Running') router.push('/')
}

/* ── Lifecycle ─────────────────────────────────────────────────────────────── */
onMounted(async () => {
  config.value = await api.fetchConfig().catch(() => config.value)
  await heartbeat()
  wsConnect()
  setInterval(heartbeat, 3 * 60 * 1000)
  // Smooth 1-second countdown for the idle timer UI
  secondsTimer = setInterval(() => {
    if (remainingSeconds.value > 0) remainingSeconds.value--
  }, 1000)
})

onUnmounted(() => {
  wsDisconnect()
  if (secondsTimer) clearInterval(secondsTimer)
})
</script>
