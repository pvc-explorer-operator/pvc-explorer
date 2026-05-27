import { onMounted, onUnmounted, type Ref } from 'vue'

interface FileBrowserShortcutHandlers {
  selectedNames: Ref<string[]>
  readonly: Ref<boolean>
  selectAll: () => void
  clearSelection: () => void
  deleteSelected: () => void
  downloadSelected: () => void
  openNewFile: () => void
  openNewFolder: () => void
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

export function useFileBrowserShortcuts(handlers: FileBrowserShortcutHandlers) {
  function onKeydown(e: KeyboardEvent) {
    if (isTypingTarget(e.target)) return

    const ctrl = e.metaKey || e.ctrlKey

    if (ctrl && e.key === 'a') {
      e.preventDefault()
      handlers.selectAll()
      return
    }

    if (e.altKey || e.shiftKey) return

    switch (e.key) {
      case 'Delete':
      case 'Backspace':
        if (handlers.selectedNames.value.length > 0 && !handlers.readonly.value) {
          handlers.deleteSelected()
        }
        break
      case ' ':
        e.preventDefault()
        break
    }
  }

  onMounted(() => document.addEventListener('keydown', onKeydown))
  onUnmounted(() => document.removeEventListener('keydown', onKeydown))
}
