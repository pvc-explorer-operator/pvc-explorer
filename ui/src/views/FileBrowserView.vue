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
    :disconnected="disconnected"
    :reconnecting="reconnecting"
    @heartbeat="heartbeat"
    @reconnect="reconnect"
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
const remainingSeconds = ref<number | null>(null)
const idleWarning      = ref(false)
const disconnected     = ref(false)
const reconnecting     = ref(false)
let secondsTimer: ReturnType<typeof setInterval> | null = null
// Timestamp of the last heartbeat reset. idle.tick events that arrive within
// a short window after a reset are stale ring-buffer replays and are ignored.
let lastHeartbeatAt = 0
const heartbeatCooldownMs = 15_000 // ignore idle.ticks for 15s after a reset

/* ── WebSocket (idle ticks, warnings, expiry) ──────────────────────────────── */
const { connect: wsConnect, disconnect: wsDisconnect } = useWebSocket({
  onIdleTick(p) {
    if (p.namespace === ns && p.name === name) {
      // Ignore ticks that arrive shortly after a heartbeat reset — they are
      // stale ring-buffer replays of the reset tick (full timeout value) being
      // replayed on WebSocket reconnect, not real decrements from the server.
      if (Date.now() - lastHeartbeatAt < heartbeatCooldownMs) return
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
    if (p.namespace === ns && p.name === name) {
      disconnected.value = true
      remainingSeconds.value = 0
      idleWarning.value = false
      if (secondsTimer) { clearInterval(secondsTimer); secondsTimer = null }
    }
  },
})

/* ── Heartbeat ─────────────────────────────────────────────────────────────── */
async function heartbeat() {
  const res = await fetch(`${base}/heartbeat`, { method: 'POST' })
  if (!res.ok) return
  const data: { remainingSeconds: number; phase: string } = await res.json()
  lastHeartbeatAt = Date.now()
  remainingSeconds.value = data.remainingSeconds
  if (data.phase !== 'Running') router.push('/')
}

/* ── Reconnect after idle expiry ───────────────────────────────────────────── */
async function reconnect() {
  reconnecting.value = true
  try {
    const wakeRes = await fetch(`${base}/wake`, { method: 'POST' })
    if (!wakeRes.ok) return
    // Poll heartbeat until the agent is Running again (up to 60s)
    const deadline = Date.now() + 60_000
    while (Date.now() < deadline) {
      await new Promise(r => setTimeout(r, 3_000))
      const hbRes = await fetch(`${base}/heartbeat`, { method: 'POST' })
      if (hbRes.ok) {
        const data: { remainingSeconds: number; phase: string } = await hbRes.json()
        if (data.phase === 'Running') {
          lastHeartbeatAt = Date.now()
          remainingSeconds.value = data.remainingSeconds
          disconnected.value = false
          idleWarning.value = false
          // Restart the local countdown
          secondsTimer = setInterval(() => {
            if (remainingSeconds.value != null && remainingSeconds.value > 0) remainingSeconds.value--
          }, 1000)
          return
        }
      }
    }
  } finally {
    reconnecting.value = false
  }
}

/* ── Lifecycle ─────────────────────────────────────────────────────────────── */
onMounted(async () => {
  config.value = await api.fetchConfig().catch(() => config.value)
  await heartbeat()
  wsConnect()
  // Smooth 1-second countdown for the idle timer UI
  secondsTimer = setInterval(() => {
    if (remainingSeconds.value != null && remainingSeconds.value > 0) remainingSeconds.value--
  }, 1000)
})

onUnmounted(() => {
  wsDisconnect()
  if (secondsTimer) clearInterval(secondsTimer)
})
</script>
