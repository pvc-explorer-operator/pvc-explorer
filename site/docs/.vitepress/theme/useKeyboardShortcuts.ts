import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter, useData } from 'vitepress'

export const shortcutsModalOpen = ref(false)

function isTypingTarget(el: EventTarget | null): boolean {
  if (!(el instanceof HTMLElement)) return false
  const tag = el.tagName
  return (
    tag === 'INPUT' ||
    tag === 'TEXTAREA' ||
    tag === 'SELECT' ||
    el.isContentEditable
  )
}

export function useKeyboardShortcuts() {
  const router = useRouter()
  const { page } = useData()

  function getNavLink(rel: 'prev' | 'next'): string | undefined {
    const selector = rel === 'prev'
      ? '.pager-link.prev'
      : '.pager-link.next'
    const anchor = document.querySelector<HTMLAnchorElement>(selector)
    return anchor?.getAttribute('href') ?? undefined
  }

  function onKeydown(e: KeyboardEvent) {
    if (e.metaKey || e.ctrlKey || e.altKey) return
    if (isTypingTarget(e.target)) return

    switch (e.key) {
      case 'ArrowRight': {
        const href = getNavLink('next')
        if (href) router.go(href)
        break
      }
      case 'ArrowLeft': {
        const href = getNavLink('prev')
        if (href) router.go(href)
        break
      }
      case 't': {
        const reduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches
        window.scrollTo({ top: 0, behavior: reduced ? 'instant' : 'smooth' })
        break
      }
      case '?':
        shortcutsModalOpen.value = !shortcutsModalOpen.value
        break
    }
  }

  onMounted(() => document.addEventListener('keydown', onKeydown))
  onUnmounted(() => document.removeEventListener('keydown', onKeydown))
}
