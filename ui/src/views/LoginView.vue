<template>
  <main class="login-bg" :class="{ 'is-dark': isDark }">
    <canvas ref="bgCanvas" class="login-canvas" aria-hidden="true" />

    <div class="card login-card">
      <!-- Logo / title -->
      <div class="login-hero">
        <img src="/logo-icon.svg" class="login-logo" alt="PVC Explorer icon" />
        <h1 class="login-title">PVC Explorer</h1>
        <span class="login-sub">Sign in to manage your PVC agents</span>
      </div>

      <form @submit.prevent="onSubmit" class="login-form">
        <div class="login-field">
          <label for="username">Username</label>
          <InputText id="username" v-model="username" autocomplete="username" placeholder="Username" class="w-full" />
        </div>

        <div class="login-field">
          <label for="password">Password</label>
          <Password id="password" v-model="password" :feedback="false" toggleMask inputClass="w-full" class="w-full" autocomplete="current-password" placeholder="Password" />
        </div>

        <Message v-if="error" severity="error" :closable="false">{{ error }}</Message>

        <Button type="submit" label="Sign In" icon="pi pi-sign-in" :loading="loading" :disabled="!canSubmit" class="w-full" fluid />
      </form>
    </div>
  </main>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/authStore'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import Button from 'primevue/button'
import Message from 'primevue/message'

const username = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)
const router = useRouter()
const authStore = useAuthStore()
const bgCanvas = ref<HTMLCanvasElement | null>(null)

const isDark = ref(document.documentElement.classList.contains('app-dark'))

const canSubmit = computed(
  () => username.value.trim().length > 0 && password.value.length > 0 && !loading.value
)

async function onSubmit() {
  if (!canSubmit.value) return
  error.value = ''
  loading.value = true
  try {
    await authStore.login(username.value.trim(), password.value)
    router.push('/')
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Login failed'
  } finally {
    loading.value = false
  }
}

// ── Background network animation ────────────────────────────────────────────

interface Node {
  x: number; y: number
  vx: number; vy: number
  r: number; phase: number
}

const NODES = 32
const LINK_DIST = 160
let raf = 0
let nodes: Node[] = []

function palette() {
  return isDark.value
    ? { bg1: '#050f1a', bg2: '#0c1f35', node: '#0ea5e9', line: '#38bdf8', glow: 18 }
    : { bg1: '#e0f2fe', bg2: '#f0f9ff', node: '#0284c7', line: '#7dd3fc', glow: 0 }
}

function initNodes(w: number, h: number) {
  nodes = Array.from({ length: NODES }, () => ({
    x: Math.random() * w,
    y: Math.random() * h,
    vx: (Math.random() - 0.5) * 0.35,
    vy: (Math.random() - 0.5) * 0.35,
    r: 3 + Math.random() * 3.5,
    phase: Math.random() * Math.PI * 2,
  }))
}

function draw(canvas: HTMLCanvasElement, t: number) {
  const ctx = canvas.getContext('2d')!
  const w = canvas.width
  const h = canvas.height
  const p = palette()

  // background gradient
  const grad = ctx.createLinearGradient(0, 0, w, h)
  grad.addColorStop(0, p.bg1)
  grad.addColorStop(1, p.bg2)
  ctx.fillStyle = grad
  ctx.fillRect(0, 0, w, h)

  // move nodes
  for (const n of nodes) {
    n.x += n.vx; n.y += n.vy
    if (n.x < -20) n.x = w + 20
    if (n.x > w + 20) n.x = -20
    if (n.y < -20) n.y = h + 20
    if (n.y > h + 20) n.y = -20
  }

  // draw edges
  for (let i = 0; i < nodes.length; i++) {
    for (let j = i + 1; j < nodes.length; j++) {
      const a = nodes[i], b = nodes[j]
      const dx = a.x - b.x, dy = a.y - b.y
      const dist = Math.sqrt(dx * dx + dy * dy)
      if (dist > LINK_DIST) continue
      const alpha = (1 - dist / LINK_DIST) * 0.55
      ctx.beginPath()
      ctx.moveTo(a.x, a.y)
      ctx.lineTo(b.x, b.y)
      ctx.strokeStyle = p.line
      ctx.globalAlpha = alpha
      ctx.lineWidth = 1
      ctx.stroke()
    }
  }

  // draw nodes
  for (const n of nodes) {
    const pulse = 0.85 + 0.15 * Math.sin(t / 1200 + n.phase)
    const r = n.r * pulse
    ctx.globalAlpha = 1

    if (p.glow > 0) {
      ctx.shadowColor = p.node
      ctx.shadowBlur = p.glow
    }

    ctx.beginPath()
    ctx.arc(n.x, n.y, r, 0, Math.PI * 2)
    ctx.fillStyle = p.node
    ctx.fill()

    // inner highlight
    ctx.beginPath()
    ctx.arc(n.x, n.y, r * 0.45, 0, Math.PI * 2)
    ctx.fillStyle = isDark.value ? '#bae6fd' : '#e0f2fe'
    ctx.fill()

    ctx.shadowBlur = 0
  }

  ctx.globalAlpha = 1
}

function resize(canvas: HTMLCanvasElement) {
  canvas.width = window.innerWidth
  canvas.height = window.innerHeight
}

function loop(t: number) {
  const canvas = bgCanvas.value
  if (!canvas) return
  draw(canvas, t)
  raf = requestAnimationFrame(loop)
}

function onResize() {
  const canvas = bgCanvas.value
  if (!canvas) return
  resize(canvas)
  initNodes(canvas.width, canvas.height)
}

onMounted(() => {
  const canvas = bgCanvas.value!
  resize(canvas)
  initNodes(canvas.width, canvas.height)
  raf = requestAnimationFrame(loop)
  window.addEventListener('resize', onResize)
})

onUnmounted(() => {
  cancelAnimationFrame(raf)
  window.removeEventListener('resize', onResize)
})

// observe app-dark class changes so isDark ref stays in sync
const mo = new MutationObserver(() => {
  isDark.value = document.documentElement.classList.contains('app-dark')
})
onMounted(() => mo.observe(document.documentElement, { attributes: true, attributeFilter: ['class'] }))
onUnmounted(() => mo.disconnect())
</script>

<style scoped>
.login-bg {
  position: relative;
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.login-canvas {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.login-card {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 26rem;
  backdrop-filter: blur(6px);
  background: color-mix(in srgb, var(--surface-card) 88%, transparent) !important;
  border: 1px solid color-mix(in srgb, var(--surface-border) 60%, transparent);
}

.login-logo {
  width: clamp(7.15rem, 13.2vw, 9.35rem);
  height: clamp(7.15rem, 13.2vw, 9.35rem);
  object-fit: contain;
  margin-bottom: 0.25rem;
  filter: drop-shadow(0 0 14px #38bdf877);
}

.login-hero {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 2rem;
  text-align: center;
}

.login-title {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--text-color);
}

.login-sub {
  font-size: 0.875rem;
  color: var(--text-color-secondary);
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.login-field {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}

.login-field label {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--text-color);
}

.w-full {
  width: 100%;
}
</style>
