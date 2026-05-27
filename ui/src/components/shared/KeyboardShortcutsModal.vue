<template>
  <dialog v-if="open" ref="dlgEl" class="ks-modal" aria-labelledby="ks-title" @close="onClose">
    <div class="ks-header">
      <span id="ks-title" class="ks-title">Keyboard Shortcuts</span>
      <button class="ks-close" aria-label="Close" @click="close">✕</button>
    </div>

    <div class="ks-body">
      <section>
        <h3 class="ks-section-title">Global</h3>
        <table class="ks-table">
          <tbody>
            <tr><td><kbd>g</kbd> <kbd>h</kbd></td><td>Go to Home</td></tr>
            <tr><td><kbd>g</kbd> <kbd>s</kbd></td><td>Go to Scopes</td></tr>
            <tr><td><kbd>r</kbd></td><td>Refresh explorer list</td></tr>
            <tr><td><kbd>b</kbd></td><td>Toggle sidebar</td></tr>
            <tr><td><kbd>d</kbd></td><td>Toggle dark / light mode</td></tr>
            <tr><td><kbd>/</kbd></td><td>Focus filter search <span class="ks-note">(dashboard)</span></td></tr>
            <tr><td><kbd>⌘</kbd>/<kbd>Ctrl</kbd> <kbd>K</kbd></td><td>Open page search</td></tr>
            <tr><td><kbd>?</kbd></td><td>Toggle this help</td></tr>
          </tbody>
        </table>
      </section>

      <section>
        <h3 class="ks-section-title">Explorer detail</h3>
        <table class="ks-table">
          <tbody>
            <tr><td><kbd>f</kbd></td><td>Browse files <span class="ks-note">(Running)</span></td></tr>
            <tr><td><kbd>w</kbd></td><td>Wake / Connect <span class="ks-note">(ScaledToZero)</span></td></tr>
            <tr><td><kbd>x</kbd></td><td>Disconnect <span class="ks-note">(Running)</span></td></tr>
            <tr><td><kbd>r</kbd></td><td>Refresh</td></tr>
          </tbody>
        </table>
      </section>

      <section>
        <h3 class="ks-section-title">File browser</h3>
        <table class="ks-table">
          <tbody>
            <tr><td><kbd>Ctrl</kbd>/<kbd>⌘</kbd> <kbd>A</kbd></td><td>Select all files</td></tr>
            <tr><td><kbd>Delete</kbd> / <kbd>⌫</kbd></td><td>Delete selected</td></tr>
          </tbody>
        </table>
      </section>

      <div class="ks-tour-row">
        <button class="p-button p-button-sm p-button-outlined" @click="startTour">Show welcome tour</button>
      </div>
    </div>
  </dialog>
</template>

<script setup lang="ts">
import { watch, ref, nextTick } from 'vue'
import { shortcutsModalOpen } from '@/composables/useShortcutsModal'

const emit = defineEmits<{ requestTour: [] }>()

const open = shortcutsModalOpen

function startTour() {
  close()
  emit('requestTour')
}
const dlgEl = ref<HTMLDialogElement | null>(null)

watch(open, (val) => {
  if (val) nextTick(() => dlgEl.value?.showModal())
})

function close() { open.value = false }

function onClose() { open.value = false }
</script>

<style scoped>
.ks-modal {
  background: var(--surface-card);
  border: 1px solid var(--surface-border);
  border-radius: 12px;
  box-shadow: 0 24px 64px color-mix(in srgb, #000 40%, transparent);
  width: min(560px, calc(100vw - 2rem));
  max-height: calc(100vh - 4rem);
  display: flex;
  flex-direction: column;
  overflow: hidden;

  /* Transition */
  opacity: 0;
  transform: scale(0.95);
  transition: opacity 0.2s ease, transform 0.2s ease,
              display 0.2s allow-discrete, overlay 0.2s allow-discrete;
  transition-behavior: allow-discrete;
}
.ks-modal[open] {
  opacity: 1;
  transform: scale(1);
  @starting-style { opacity: 0; transform: scale(0.95); }
}
.ks-modal::backdrop {
  background: color-mix(in srgb, var(--surface-900, #0f1117) 0%, transparent);
  backdrop-filter: blur(0px);
  transition: display 0.2s allow-discrete, overlay 0.2s allow-discrete,
              background 0.2s ease, backdrop-filter 0.2s ease;
  transition-behavior: allow-discrete;
}
.ks-modal[open]::backdrop {
  background: color-mix(in srgb, var(--surface-900, #0f1117) 60%, transparent);
  backdrop-filter: blur(4px);
  @starting-style { background: color-mix(in srgb, var(--surface-900, #0f1117) 0%, transparent); backdrop-filter: blur(0px); }
}
@media (prefers-reduced-motion: reduce) {
  .ks-modal { transform: none; transition-duration: 0.1s; }
}
.ks-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 1.25rem 0.75rem;
  border-bottom: 1px solid var(--surface-border);
}

.ks-title {
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--text-color);
}

.ks-close {
  background: none;
  border: none;
  cursor: pointer;
  color: var(--text-color-secondary);
  font-size: 0.875rem;
  padding: 0.25rem 0.4rem;
  border-radius: 4px;
  line-height: 1;
  transition: background 0.15s, color 0.15s;
}
.ks-close:hover {
  background: var(--surface-hover);
  color: var(--text-color);
}

.ks-body {
  overflow-y: auto;
  padding: 1rem 1.25rem 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.ks-section-title {
  font-size: 0.75rem;
  font-weight: 600;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  color: var(--text-color-secondary);
  margin: 0 0 0.5rem;
}

.ks-table {
  width: 100%;
  border-collapse: collapse;
}
.ks-table td {
  padding: 0.3rem 0.5rem;
  font-size: 0.875rem;
  color: var(--text-color);
  vertical-align: middle;
}
.ks-table td:first-child {
  white-space: nowrap;
  width: 1%;
  padding-right: 1.25rem;
}
.ks-table tr:hover td {
  background: var(--surface-hover);
}
.ks-table tr:hover td:first-child { border-radius: 4px 0 0 4px; }
.ks-table tr:hover td:last-child  { border-radius: 0 4px 4px 0; }

kbd {
  display: inline-block;
  padding: 0.1em 0.45em;
  font-size: 0.75rem;
  font-family: ui-monospace, 'Cascadia Code', 'Fira Mono', monospace;
  color: var(--text-color);
  background: var(--surface-ground);
  border: 1px solid var(--surface-border);
  border-bottom-width: 2px;
  border-radius: 4px;
  line-height: 1.4;
}

.ks-note {
  font-size: 0.75rem;
  color: var(--text-color-secondary);
}

.ks-tour-row {
  display: flex;
  justify-content: center;
  padding-top: 0.5rem;
  border-top: 1px solid var(--surface-border, #e2e8f0);
}
</style>
