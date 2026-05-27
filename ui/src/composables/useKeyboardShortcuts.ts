import { onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useLayout } from '@/layout/composables/layout'
import { useExplorerStore } from '@/stores/explorerStore'
import { shortcutsModalOpen } from '@/composables/useShortcutsModal'
import { searchDialogOpen } from '@/composables/useSearchDialog'

function isTypingTarget(el: EventTarget | null): boolean {
  if (!(el instanceof HTMLElement)) return false
  return (
    el.tagName === 'INPUT' ||
    el.tagName === 'TEXTAREA' ||
    el.tagName === 'SELECT' ||
    el.isContentEditable
  )
}

export function useKeyboardShortcuts() {
  const router = useRouter()
  const route = useRoute()
  const { toggleMenu, toggleDarkMode } = useLayout()
  const store = useExplorerStore()

  let gPending = false
  let gTimer: ReturnType<typeof setTimeout> | null = null

  function clearGChord() {
    gPending = false
    if (gTimer) { clearTimeout(gTimer); gTimer = null }
  }

  function onKeydown(e: KeyboardEvent) {
    // Cmd+K / Ctrl+K — open search dialog (allow through modifier guard)
    if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'k') {
      e.preventDefault()
      searchDialogOpen.value = !searchDialogOpen.value
      return
    }

    if (e.metaKey || e.ctrlKey || e.altKey) return
    if (isTypingTarget(e.target)) return

    if (gPending) {
      clearGChord()
      switch (e.key) {
        case 'h': router.push('/'); break
        case 's': router.push('/scopes'); break
      }
      e.preventDefault()
      return
    }

    switch (e.key) {
      case 'g':
        gPending = true
        gTimer = setTimeout(clearGChord, 300)
        e.preventDefault()
        break
      case 'r':
        store.fetchExplorers()
        break
      case 'b':
        toggleMenu()
        break
      case 'd':
        toggleDarkMode()
        break
      case '/':
        if (route.path === '/') {
          const input = document.querySelector<HTMLInputElement>('.filter-search-input')
          if (input) { input.focus(); e.preventDefault() }
        }
        break
      case '?':
        shortcutsModalOpen.value = !shortcutsModalOpen.value
        break
    }
  }

  onMounted(() => document.addEventListener('keydown', onKeydown))
  onUnmounted(() => { document.removeEventListener('keydown', onKeydown); clearGChord() })
}
