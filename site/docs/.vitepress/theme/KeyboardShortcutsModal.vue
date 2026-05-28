<template>
  <dialog
    ref="dialogEl"
    class="kb-modal"
    aria-labelledby="kb-title"
    @click.self="close"
    @cancel.prevent="close"
    @close="close"
  >
    <div class="kb-header">
      <span id="kb-title" class="kb-title">Keyboard Shortcuts</span>
      <button class="kb-close" aria-label="Close" @click="close">✕</button>
    </div>
    <dl class="kb-list">
      <div class="kb-row">
        <dt><kbd>→</kbd></dt>
        <dd>Next page</dd>
      </div>
      <div class="kb-row">
        <dt><kbd>←</kbd></dt>
        <dd>Previous page</dd>
      </div>
      <div class="kb-row">
        <dt><kbd>t</kbd></dt>
        <dd>Scroll to top</dd>
      </div>
      <div class="kb-row">
        <dt><kbd>⌘</kbd><kbd>K</kbd></dt>
        <dd>Open search</dd>
      </div>
      <div class="kb-row">
        <dt><kbd>?</kbd></dt>
        <dd>Toggle this panel</dd>
      </div>
      <div class="kb-row">
        <dt><kbd>Esc</kbd></dt>
        <dd>Close search / menu</dd>
      </div>
    </dl>
  </dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { shortcutsModalOpen } from './useKeyboardShortcuts'

const dialogEl = ref<HTMLDialogElement | null>(null)

// biome-ignore lint/correctness/noUnusedVariables: used in <template> — false positive from Vue SFC
function close() {
  shortcutsModalOpen.value = false
}

watch(shortcutsModalOpen, (open) => {
  if (!dialogEl.value) return
  if (open) dialogEl.value.showModal()
  else dialogEl.value.close()
})
</script>

<style scoped>
/* ── native <dialog> — card surface ── */
.kb-modal {
  padding: 1.25rem 1.5rem 1.5rem;
  border: 1px solid var(--vp-c-divider);
  border-radius: 12px;
  min-width: 280px;
  max-width: 360px;
  width: 90vw;
  background: var(--vp-c-bg-elv);
  box-shadow: 0 20px 48px rgba(0, 0, 0, 0.32);

  /* base: visible states (no transition by default) */
  opacity: 0;
  transform: scale(0.96);
}

@media (prefers-reduced-motion: no-preference) {
  .kb-modal {
    transition: opacity 0.15s ease, transform 0.15s ease,
                display 0.15s ease allow-discrete,
                overlay 0.15s ease allow-discrete;
  }

  .kb-modal::backdrop {
    transition: opacity 0.15s ease;
  }

  @starting-style {
    .kb-modal[open] {
      opacity: 0;
      transform: scale(0.96);
    }
  }
}

.kb-modal[open] {
  opacity: 1;
  transform: scale(1);
}

/* ── backdrop overlay ── */
.kb-modal::backdrop {
  background: rgba(0, 0, 0, 0.45);
  opacity: 0;
}

.kb-modal[open]::backdrop {
  opacity: 1;
}

/* ── header ── */
.kb-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1rem;
}

.kb-title {
  font-size: 0.8rem;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--vp-c-text-3);
}

.kb-close {
  background: none;
  border: none;
  cursor: pointer;
  color: var(--vp-c-text-3);
  font-size: 0.85rem;
  padding: 0.1rem 0.3rem;
  border-radius: 4px;
  line-height: 1;
}

.kb-close:hover {
  color: var(--vp-c-text-1);
  background: var(--vp-c-bg-soft);
}

.kb-close:focus-visible {
  outline: 2px solid var(--vp-c-brand-1);
  outline-offset: 2px;
}

/* ── shortcut list ── */
.kb-list {
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 0.45rem;
}

.kb-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}

.kb-row dt {
  display: flex;
  gap: 0.25rem;
  flex-shrink: 0;
}

.kb-row dd {
  margin: 0;
  font-size: 0.875rem;
  color: var(--vp-c-text-2);
}

kbd {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 1.5rem;
  padding: 0.15rem 0.4rem;
  font-family: var(--vp-font-family-mono);
  font-size: 0.75rem;
  color: var(--vp-c-text-1);
  background: var(--vp-c-bg-soft);
  border: 1px solid var(--vp-c-divider);
  border-bottom-width: 2px;
  border-radius: 4px;
  line-height: 1.4;
}
</style>
