<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'

interface TourStep {
  target: string
  title: string
  description: string
  placement?: 'bottom' | 'top' | 'left' | 'right' | 'bottom-start'
}

const STORAGE_KEY = 'pvc-explorer-tour-seen'

const steps: TourStep[] = [
  {
    target: '.app-card:first-of-type',
    title: 'Discovered PVCs',
    description: 'Each card represents a PersistentVolumeClaim discovered in your cluster. Click to browse files, connect, or view details.',
    placement: 'bottom',
  },
  {
    target: '.filter-sidebar',
    title: 'Filter Settings',
    description: 'Narrow down explorers by phase, namespace, labels, and more. Active filters show a summary at the top of the list.',
    placement: 'right',
  },
  {
    target: '.view-toggle',
    title: 'View Modes',
    description: 'Switch between Card and List view to browse explorers the way you prefer.',
    placement: 'bottom',
  },
  {
    target: '.layout-topbar',
    title: 'Topbar Controls',
    description: 'Manage auto-refresh, toggle dark mode, and access keyboard shortcuts from here.',
    placement: 'bottom',
  },
]

const active = ref(false)
const currentIndex = ref(0)
const popoverEl = ref<HTMLElement | null>(null)
const targetRect = ref({ top: 0, right: 0, bottom: 0, left: 0, width: 0, height: 0 })
const popoverRect = ref({ width: 320, height: 0 })
const arrowDir = ref<'top' | 'bottom' | 'left' | 'right'>('top')

const currentStep = computed(() => steps[currentIndex.value] ?? null)

function scrollTargetIntoView(sel: string) {
  const el = document.querySelector(sel)
  if (el) el.scrollIntoView({ block: 'center', behavior: 'smooth' })
}

function calcTargetRect(sel: string) {
  const el = document.querySelector(sel)
  if (!el) return null
  const rect = el.getBoundingClientRect()
  return { top: rect.top, right: rect.right, bottom: rect.bottom, left: rect.left, width: rect.width, height: rect.height }
}

function calcPopoverPosition() {
  const step = currentStep.value
  if (!step) return
  const tr = calcTargetRect(step.target)
  if (!tr) return
  targetRect.value = tr

  const pw = popoverRect.value.width
  const gap = 12
  let top = 0, left = 0

  const placement = step.placement ?? 'bottom'

  if (placement === 'bottom') {
    top = tr.bottom + gap
    left = tr.left + tr.width / 2 - pw / 2
    arrowDir.value = 'top'
  } else if (placement === 'bottom-start') {
    top = tr.bottom + gap
    left = tr.left
    arrowDir.value = 'top'
  } else if (placement === 'top') {
    top = tr.top - gap
    left = tr.left + tr.width / 2 - pw / 2
    arrowDir.value = 'bottom'
  } else if (placement === 'right') {
    top = tr.top + tr.height / 2
    left = tr.right + gap
    arrowDir.value = 'left'
  } else if (placement === 'left') {
    top = tr.top + tr.height / 2
    left = tr.left - gap - pw
    arrowDir.value = 'right'
  }

  // Keep popover within viewport
  const vw = window.innerWidth
  const vh = window.innerHeight
  const popupH = popoverEl.value?.offsetHeight ?? 200
  popoverRect.value.height = popupH

  if (left < 12) left = 12
  if (left + pw > vw - 12) left = vw - pw - 12
  if (top < 12) top = 12
  if (top + popupH > vh - 12) top = vh - popupH - 12

  return { top, left }
}

const popoverStyle = computed(() => {
  const pos = calcPopoverPosition()
  if (!pos) return { display: 'none' }
  return {
    top: `${pos.top}px`,
    left: `${pos.left}px`,
  }
})

const highlightStyle = computed(() => {
  const tr = targetRect.value
  if (!tr.width) return { display: 'none' }
  return {
    top: `${tr.top - 4}px`,
    left: `${tr.left - 4}px`,
    width: `${tr.width + 8}px`,
    height: `${tr.height + 8}px`,
  }
})

function next() {
  if (currentIndex.value < steps.length - 1) {
    currentIndex.value++
    scrollTargetIntoView(steps[currentIndex.value].target)
  } else {
    finish()
  }
}

function prev() {
  if (currentIndex.value > 0) {
    currentIndex.value--
    scrollTargetIntoView(steps[currentIndex.value].target)
  }
}

function finish() {
  active.value = false
  localStorage.setItem(STORAGE_KEY, '1')
}

function skip() {
  finish()
}

function start() {
  active.value = true
  currentIndex.value = 0
  scrollTargetIntoView(steps[0].target)
}

function onResize() {
  if (active.value && currentStep.value) {
    calcPopoverPosition()
  }
}

function onKeydown(e: KeyboardEvent) {
  if (!active.value) return
  if (e.key === 'Escape') skip()
  if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); next() }
  if (e.key === 'ArrowLeft') { e.preventDefault(); prev() }
  if (e.key === 'ArrowRight') { e.preventDefault(); next() }
}

// Expose start for external trigger (e.g. help menu button)
defineExpose({ start })

onMounted(() => {
  const seen = localStorage.getItem(STORAGE_KEY)
  if (!seen) {
    // Small delay to allow the view to render
    setTimeout(() => {
      scrollTargetIntoView(steps[0].target)
      setTimeout(() => { start() }, 400)
    }, 600)
  }
  window.addEventListener('resize', onResize)
  window.addEventListener('keydown', onKeydown)
})

onUnmounted(() => {
  window.removeEventListener('resize', onResize)
  window.removeEventListener('keydown', onKeydown)
})
</script>

<template>
  <Teleport to="body">
    <div v-if="active" class="tour-backdrop" @click="skip" />
    <div v-if="active && currentStep" class="tour-highlight" :style="highlightStyle" />
    <div
      v-if="active && currentStep"
      ref="popoverEl"
      class="tour-popover"
      :style="popoverStyle"
      role="dialog"
      aria-labelledby="tour-title"
    >
      <div class="tour-popover-arrow" :class="`arrow-${arrowDir}`" />
      <div class="tour-popover-header">
        <span class="tour-step-dots">
          <span
            v-for="(_, i) in steps"
            :key="i"
            class="tour-dot"
            :class="{ active: i === currentIndex }"
          />
        </span>
        <button class="tour-close" @click="skip" aria-label="Close tour">
          <i class="pi pi-times" aria-hidden="true" />
        </button>
      </div>
      <div class="tour-popover-body">
        <h2 id="tour-title" class="tour-title">{{ currentStep.title }}</h2>
        <p class="tour-desc">{{ currentStep.description }}</p>
      </div>
      <div class="tour-popover-footer">
        <button class="p-button p-button-text p-button-sm" @click="skip">
          Skip
        </button>
        <div class="tour-footer-right">
          <button
            v-if="currentIndex > 0"
            class="p-button p-button-text p-button-sm"
            @click="prev"
          >
            Back
          </button>
          <button class="p-button p-button-sm" @click="next">
            {{ currentIndex < steps.length - 1 ? 'Next' : 'Got it' }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
.tour-backdrop {
  position: fixed;
  inset: 0;
  z-index: 9999;
  background: rgba(0, 0, 0, 0.4);
  cursor: pointer;
}

.tour-highlight {
  position: fixed;
  z-index: 10000;
  border-radius: 8px;
  box-shadow: 0 0 0 4px var(--primary-color, #3b82f6), 0 0 0 9999px rgba(0, 0, 0, 0.35);
  pointer-events: none;
  transition: all 0.25s ease;
}

.tour-popover {
  position: fixed;
  z-index: 10001;
  width: 320px;
  max-width: calc(100vw - 24px);
  background: var(--surface-card, #fff);
  color: var(--text-color, #1e293b);
  border: 1px solid var(--surface-border, #e2e8f0);
  border-radius: 10px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.18);
  transition: left 0.25s ease, top 0.25s ease;
}

.tour-popover-arrow {
  position: absolute;
  width: 12px;
  height: 12px;
  background: var(--surface-card, #fff);
  border: 1px solid var(--surface-border, #e2e8f0);
  transform: rotate(45deg);
  z-index: -1;
}

.arrow-top {
  top: -6px;
  left: calc(50% - 6px);
  border-right: none;
  border-bottom: none;
}

.arrow-bottom {
  bottom: -6px;
  left: calc(50% - 6px);
  border-left: none;
  border-top: none;
}

.arrow-left {
  left: -6px;
  top: calc(50% - 6px);
  border-right: none;
  border-top: none;
}

.arrow-right {
  right: -6px;
  top: calc(50% - 6px);
  border-left: none;
  border-bottom: none;
}

.tour-popover-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.75rem 1rem 0;
}

.tour-step-dots {
  display: flex;
  gap: 6px;
}

.tour-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--surface-border, #cbd5e1);
  transition: background 0.2s;
}

.tour-dot.active {
  background: var(--primary-color, #3b82f6);
  width: 20px;
  border-radius: 4px;
}

.tour-close {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border: none;
  border-radius: 6px;
  background: transparent;
  color: var(--text-color-secondary, #64748b);
  cursor: pointer;
  font-size: 1rem;
}

.tour-close:hover {
  background: var(--surface-hover, #f1f5f9);
}

.tour-popover-body {
  padding: 0.75rem 1rem;
}

.tour-title {
  margin: 0 0 0.5rem;
  font-size: 1.05rem;
  font-weight: 600;
  line-height: 1.3;
}

.tour-desc {
  margin: 0;
  font-size: 0.875rem;
  line-height: 1.5;
  color: var(--text-color-secondary, #64748b);
}

.tour-popover-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.5rem 1rem 0.75rem;
}

.tour-footer-right {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}
</style>
