import { computed, reactive } from 'vue'

const THEME_STORAGE_KEY = 'pvc-theme'

function hasStoredPreference(): boolean {
  return localStorage.getItem(THEME_STORAGE_KEY) !== null
}

function loadThemePreference(): boolean {
  const stored = localStorage.getItem(THEME_STORAGE_KEY)
  if (stored !== null) return stored === 'dark'
  return window.matchMedia('(prefers-color-scheme: dark)').matches
}

function applyTheme(dark: boolean) {
  document.documentElement.classList.toggle('app-dark', dark)
}

const initialDark = loadThemePreference()
applyTheme(initialDark)

const layoutConfig = reactive({
  menuMode: 'static' as 'static' | 'overlay',
  darkTheme: initialDark,
})

// React to OS preference changes only when the user hasn't set a manual preference
if (!hasStoredPreference()) {
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)')
  prefersDark.addEventListener('change', (e) => {
    document.documentElement.classList.toggle('app-dark', e.matches)
    layoutConfig.darkTheme = e.matches
  })
}

const layoutState = reactive({
  staticMenuInactive: false,
  overlayMenuActive: false,
  mobileMenuActive: false,
  menuHoverActive: false,
  activePath: null as string | null,
})

export function useLayout() {
  const toggleMenu = () => {
    if (isDesktop()) {
      if (layoutConfig.menuMode === 'static') {
        layoutState.staticMenuInactive = !layoutState.staticMenuInactive
      } else {
        layoutState.overlayMenuActive = !layoutState.overlayMenuActive
      }
    } else {
      layoutState.mobileMenuActive = !layoutState.mobileMenuActive
    }
  }

  const hideMobileMenu = () => {
    layoutState.mobileMenuActive = false
    layoutState.overlayMenuActive = false
  }

  const toggleDarkMode = () => {
    if (!document.startViewTransition) {
      executeDarkModeToggle()
      return
    }
    document.startViewTransition(() => executeDarkModeToggle())
  }

  const executeDarkModeToggle = () => {
    const newDark = !layoutConfig.darkTheme
    layoutConfig.darkTheme = newDark
    document.documentElement.classList.toggle('app-dark')
    localStorage.setItem(THEME_STORAGE_KEY, newDark ? 'dark' : 'light')
    const meta = document.querySelector('meta[name="color-scheme"]')
    if (meta) meta.setAttribute('content', newDark ? 'dark' : 'light')
  }

  const isDesktop = () => window.innerWidth > 991
  const isDarkTheme = computed(() => layoutConfig.darkTheme)
  const hasOpenOverlay = computed(() => layoutState.overlayMenuActive || layoutState.mobileMenuActive)

  return {
    layoutConfig,
    layoutState,
    isDarkTheme,
    hasOpenOverlay,
    toggleMenu,
    toggleDarkMode,
    hideMobileMenu,
    isDesktop,
  }
}
