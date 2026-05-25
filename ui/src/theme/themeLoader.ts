interface ThemeConfig {
  appName: string
  logoUrl?: string
  primaryColor?: string
  darkModeDefault?: boolean
}

const APP_NAME_DEFAULT = 'PVC Explorer'

export async function loadTheme(): Promise<ThemeConfig> {
  let theme: ThemeConfig = { appName: APP_NAME_DEFAULT }

  try {
    const res = await fetch('/api/v1/theme')
    if (res.ok) {
      theme = await res.json()
    }
  } catch {
    // intentionally empty — defaults already set above
  }

  applyTheme(theme)
  return theme
}

function applyTheme(theme: ThemeConfig) {
  const root = document.documentElement

  if (theme.primaryColor) {
    root.style.setProperty('--color-primary', theme.primaryColor)
    root.style.setProperty('--primary-color', theme.primaryColor)
  }

  if (theme.appName) {
    document.title = theme.appName
  }

}

export function getAppName(theme?: ThemeConfig): string {
  return theme?.appName || APP_NAME_DEFAULT
}

