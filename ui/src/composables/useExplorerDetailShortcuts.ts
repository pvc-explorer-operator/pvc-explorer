import { onMounted, onUnmounted, type Ref } from 'vue'
import type { Explorer } from '@/stores/explorerStore'

interface ExplorerDetailHandlers {
  explorer: Ref<Explorer | null>
  goToFiles: () => void
  wake: () => void
  doDisconnect: () => void
  refresh: () => void
}

function isTypingTarget(el: EventTarget | null): boolean {
  if (!(el instanceof HTMLElement)) return false
  return (
    el.tagName === 'INPUT' ||
    el.tagName === 'TEXTAREA' ||
    el.tagName === 'SELECT' ||
    el.isContentEditable
  )
}

export function useExplorerDetailShortcuts(handlers: ExplorerDetailHandlers) {
  function onKeydown(e: KeyboardEvent) {
    if (e.metaKey || e.ctrlKey || e.altKey) return
    if (isTypingTarget(e.target)) return

    const phase = handlers.explorer.value?.phase

    switch (e.key) {
      case 'f':
        if (phase === 'Running') handlers.goToFiles()
        break
      case 'w':
        if (phase === 'ScaledToZero') handlers.wake()
        break
      case 'x':
        if (phase === 'Running') handlers.doDisconnect()
        break
      case 'r':
        handlers.refresh()
        break
    }
  }

  onMounted(() => document.addEventListener('keydown', onKeydown))
  onUnmounted(() => document.removeEventListener('keydown', onKeydown))
}
