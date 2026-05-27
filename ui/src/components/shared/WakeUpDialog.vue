<template>
  <Dialog
    v-model:visible="open"
    :header="`Starting ${explorer.name}`"
    :modal="true"
    :closable="false"
    :draggable="false"
    style="min-width: 320px; max-width: 95vw"
    @hide="onDialogHide"
  >
    <div class="dialog-body">
      <div v-if="!error" class="spinner" />
      <p class="status-msg" :class="{ 'is-error': !!error }">{{ message }}</p>
    </div>

    <template #footer>
      <Button v-if="!error" severity="secondary" icon="pi pi-times" label="Cancel" rounded @click="cancel" />
      <Button v-else severity="primary" icon="pi pi-check" label="Close" rounded @click="close" />
    </template>
  </Dialog>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import type { Explorer } from '../../stores/explorerStore'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import { useWebSocket } from '../../composables/useWebSocket'

const props = defineProps<{ explorer: Explorer }>()
const emit = defineEmits<{ (e: 'close'): void }>()

const router = useRouter()
const open = ref(true)
const error = ref(false)
const message = ref('Sending wake request...')

let cancelled = false
let deadlineTimer: ReturnType<typeof setTimeout> | null = null
let cycleTimer: ReturnType<typeof setTimeout> | null = null

const MESSAGES = ['Creating agent pod', 'Mounting PVC', 'Waiting for agent to be ready']
let msgIdx = 0

function cycleMessage() {
  msgIdx = (msgIdx + 1) % MESSAGES.length
  message.value = MESSAGES[msgIdx]
  cycleTimer = setTimeout(cycleMessage, 4000)
}

function onReady() {
  if (cancelled) return
  open.value = false
  emit('close')
  router.push(`/explorers/${props.explorer.namespace}/${props.explorer.name}/files`)
}

const { connect: wsConnect, disconnect: wsDisconnect } = useWebSocket({
  onAgentReady(p) {
    if (p.namespace === props.explorer.namespace && p.name === props.explorer.name) {
      onReady()
    }
  },
})

async function run() {
  const base = `/api/v1/explorers/${props.explorer.namespace}/${props.explorer.name}`
  const wakeRes = await fetch(`${base}/wake`, { method: 'POST' })
  if (!wakeRes.ok) throw new Error('Wake request failed')
  message.value = MESSAGES[0]
  cycleTimer = setTimeout(cycleMessage, 4000)

  deadlineTimer = setTimeout(() => {
    if (!cancelled && open.value) {
      error.value = true
      message.value = 'Timed out waiting for agent to start'
      if (cycleTimer) clearTimeout(cycleTimer)
    }
  }, 120_000)
}

onMounted(() => {
  wsConnect()
  run().catch(err => {
    error.value = true
    message.value = err instanceof Error ? err.message : 'Failed to start agent'
    if (cycleTimer) clearTimeout(cycleTimer)
  })
})

onUnmounted(() => {
  cancelled = true
  wsDisconnect()
  if (cycleTimer) clearTimeout(cycleTimer)
  if (deadlineTimer) clearTimeout(deadlineTimer)
})

let emitted = false

function cancel() {
  cancelled = true
  close()
}

function close() {
  if (emitted) return
  emitted = true
  open.value = false
  emit('close')
}

function onDialogHide() {
  // Called by PrimeVue after Escape or programmatic close.
  // Ensures the parent is notified and the wake operation is cancelled
  // even when PrimeVue closes the dialog via Escape key.
  cancelled = true
  close()
}
</script>

<style scoped>
.dialog-body {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  padding: 0.5rem 0;
}
.spinner {
  width: 36px;
  height: 36px;
  border: 4px solid var(--surface-hover);
  border-top: 4px solid var(--primary-color);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}
@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
.status-msg { color: var(--text-color-secondary); font-size: 0.95rem; text-align: center; }
.status-msg.is-error { color: var(--p-red-500); }

@media (prefers-reduced-motion: reduce) {
  :deep(.p-dialog-mask),
  :deep(.p-dialog) {
    animation: none;
    transition: none;
  }
}
</style>
