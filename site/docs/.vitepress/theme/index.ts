import DefaultTheme from 'vitepress/theme'
import { h, onMounted, ref, type Ref } from 'vue'
import './custom.css'
import KeyboardShortcutsModal from './KeyboardShortcutsModal.vue'
import { useKeyboardShortcuts, shortcutsModalOpen } from './useKeyboardShortcuts'

export default {
  extends: DefaultTheme,
  Layout() {
    return h(DefaultTheme.Layout, null, {
      'layout-bottom': () => h(KeyboardShortcutsModal),
      'nav-bar-content-after': () => h(HelpButton),
    })
  },
  setup() {
    useKeyboardShortcuts()

    onMounted(() => {
      const el = document.createElement('p')
      el.style.cssText = 'position:absolute;width:1px;height:1px;padding:0;margin:-1px;overflow:hidden;clip:rect(0,0,0,0);white-space:nowrap;border:0;'
      el.textContent = 'Press ? to open keyboard shortcuts.'
      el.setAttribute('aria-live', 'polite')
      document.body.appendChild(el)
      setTimeout(() => el.remove(), 3000)
    })
  },
}

/* ── Small ? help button in the nav bar ── */
function HelpButton() {
  const open = shortcutsModalOpen as Ref<boolean>
  return h('button', {
    class: 'kb-nav-btn',
    type: 'button',
    title: 'Keyboard shortcuts',
    'aria-label': 'Toggle keyboard shortcuts',
    onClick: () => { open.value = !open.value },
    innerHTML: '<span style="font-weight:700;font-size:1.05rem;line-height:1">?</span>',
  })
}
