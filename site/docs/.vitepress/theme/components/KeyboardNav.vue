<template>
  <dialog ref="dialogRef" class="kb-dialog">
    <div class="kb-modal">
      <div class="kb-modal-header">
        <span class="kb-modal-title">Keyboard Shortcuts</span>
        <button class="kb-close" @click="closeHelp" aria-label="Close shortcuts">&times;</button>
      </div>
      <div class="kb-modal-body">
        <div v-for="group in keyGroups" :key="group.label" class="kb-group">
          <div class="kb-group-label">{{ group.label }}</div>
          <div v-for="item in group.items" :key="item.keys.join()" class="kb-row">
            <span class="kb-keys">
              <kbd v-for="k in item.keys" :key="k">{{ k }}</kbd>
            </span>
            <span class="kb-desc">{{ item.desc }}</span>
          </div>
        </div>
      </div>
    </div>
  </dialog>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue"

interface KeyItem {
  keys: string[]
  desc: string
}

interface KeyGroup {
  label: string
  items: KeyItem[]
}

const keyGroups: KeyGroup[] = [
  {
    label: "Navigation",
    items: [
      { keys: ["J"], desc: "Next heading" },
      { keys: ["K"], desc: "Previous heading" },
      { keys: ["G", "G"], desc: "Scroll to top" },
      { keys: ["Shift", "G"], desc: "Scroll to bottom" },
      { keys: ["N"], desc: "Next page" },
      { keys: ["P"], desc: "Previous page" },
    ],
  },
  {
    label: "Page",
    items: [
      { keys: ["?"], desc: "Open keyboard shortcuts" },
      { keys: ["Esc"], desc: "Close modal" },
    ],
  },
]

const dialogRef = ref<HTMLDialogElement | null>(null)

function isEditable(el: EventTarget | null): boolean {
  if (!el) return true
  const tag = (el as HTMLElement).tagName
  return tag === "INPUT" || tag === "TEXTAREA" || (el as HTMLElement).isContentEditable
}

function smoothBehavior(): ScrollBehavior {
  return window.matchMedia("(prefers-reduced-motion: reduce)").matches ? "auto" : "smooth"
}

function getHeadings(): Element[] {
  return Array.from(document.querySelectorAll<HTMLElement>("[id]:is(h1,h2,h3,h4,h5,h6)"))
}

function jumpToHeading(dir: "next" | "prev") {
  const headings = getHeadings()
  if (!headings.length) return

  if (dir === "next") {
    for (const h of headings) {
      if (h.getBoundingClientRect().top > 10) {
        smoothScrollTo(h)
        return
      }
    }
  } else {
    for (let i = headings.length - 1; i >= 0; i--) {
      if (headings[i].getBoundingClientRect().top < -5) {
        smoothScrollTo(headings[i])
        return
      }
    }
  }
}

function smoothScrollTo(el: Element) {
  el.scrollIntoView({ behavior: smoothBehavior(), block: "start" })
  try {
    (el as HTMLElement).focus({ preventScroll: true })
  } catch {
    // non-focusable element, ignore
  }
}

function closeHelp() {
  dialogRef.value?.close()
}

// g,g double-tap tracking
let gTimer: ReturnType<typeof setTimeout> | null = null
let gPending = false

function cancelG() {
  gPending = false
  if (gTimer !== null) {
    clearTimeout(gTimer)
    gTimer = null
  }
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === "?" && !e.ctrlKey && !e.metaKey && !e.altKey) {
    e.preventDefault()
    if (dialogRef.value?.open) {
      dialogRef.value.close()
    } else {
      dialogRef.value?.showModal()
    }
    return
  }

  if (dialogRef.value?.open) return

  if (isEditable(e.target)) return

  if (e.key === "g" && !e.shiftKey && !e.ctrlKey && !e.metaKey && !e.altKey) {
    if (gPending) {
      e.preventDefault()
      cancelG()
      window.scrollTo({ top: 0, behavior: smoothBehavior() })
      return
    }
    gPending = true
    gTimer = window.setTimeout(() => { gPending = false }, 500)
    return
  }

  if (gPending) cancelG()

  if (e.key === "G") {
    e.preventDefault()
    window.scrollTo({ top: document.documentElement.scrollHeight, behavior: smoothBehavior() })
    return
  }

  if (e.key === "j" || e.key === "J") {
    if (!e.ctrlKey && !e.metaKey && !e.altKey) {
      e.preventDefault()
      jumpToHeading("next")
    }
    return
  }

  if (e.key === "k" || e.key === "K") {
    if (!e.ctrlKey && !e.metaKey && !e.altKey) {
      e.preventDefault()
      jumpToHeading("prev")
    }
    return
  }

  if (e.key === "n" || e.key === "N") {
    if (!e.ctrlKey && !e.metaKey && !e.altKey) {
      e.preventDefault()
      document.querySelector<HTMLAnchorElement>(".pager-link.next")?.click()
    }
    return
  }

  if (e.key === "p" || e.key === "P") {
    if (!e.ctrlKey && !e.metaKey && !e.altKey) {
      e.preventDefault()
      document.querySelector<HTMLAnchorElement>(".pager-link.prev")?.click()
    }
    return
  }
}

function onDialogClick(e: MouseEvent) {
  if (e.target === dialogRef.value) {
    closeHelp()
  }
}

onMounted(() => {
  document.addEventListener("keydown", onKeydown, { capture: true })
  document.addEventListener("ui:toggle-keyboard-help", toggleHelp)
  dialogRef.value?.addEventListener("click", onDialogClick)
})

onUnmounted(() => {
  document.removeEventListener("keydown", onKeydown, { capture: true })
  document.removeEventListener("ui:toggle-keyboard-help", toggleHelp)
  cancelG()
})

function toggleHelp() {
  if (dialogRef.value?.open) {
    dialogRef.value.close()
  } else {
    dialogRef.value?.showModal()
  }
}
</script>

<style scoped>
.kb-dialog {
  border: none;
  border-radius: 12px;
  padding: 0;
  max-width: 400px;
  width: calc(100vw - 2rem);
  box-shadow: 0 16px 48px rgba(0, 0, 0, 0.18);
  color: var(--vp-c-text-1);
  background: var(--vp-c-bg);
}

.kb-dialog::backdrop {
  background: rgba(0, 0, 0, 0.35);
}

.kb-modal {
  padding: var(--space-lg);
}

.kb-modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-lg);
}

.kb-modal-title {
  font-size: 1.05rem;
  font-weight: 700;
}

.kb-close {
  border: none;
  background: transparent;
  font-size: 1.25rem;
  cursor: pointer;
  color: var(--vp-c-text-3);
  padding: 4px 8px;
  border-radius: 6px;
  line-height: 1;
}

.kb-close:hover {
  color: var(--vp-c-text-1);
  background: var(--vp-c-bg-soft);
}

.kb-group {
  margin-bottom: var(--space-lg);
}

.kb-group:last-child {
  margin-bottom: 0;
}

.kb-group-label {
  font-size: 0.7rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--vp-c-text-3);
  margin-bottom: var(--space-sm);
}

.kb-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 0;
  gap: var(--space-md);
}

.kb-keys {
  display: flex;
  gap: 4px;
  align-items: center;
}

.kb-desc {
  font-size: 0.85rem;
  color: var(--vp-c-text-2);
  text-align: right;
}

kbd {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 26px;
  height: 24px;
  padding: 0 7px;
  font-family: var(--vp-font-family-mono);
  font-size: 0.7rem;
  font-weight: 600;
  color: var(--vp-c-text-1);
  background: var(--vp-c-bg-soft);
  border: 1px solid var(--vp-c-border);
  border-radius: 6px;
  box-shadow: 0 1px 0 var(--vp-c-border);
}
</style>
