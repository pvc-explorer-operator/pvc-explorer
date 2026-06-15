import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export type UserRole = 'admin' | 'user' | 'viewer'

export interface UserProfile {
  username: string
  role: UserRole
  email: string
  groups: string[]
  subject: string
}

export const useAuthStore = defineStore('auth', () => {
  const username = ref<string | null>(null)
  const role = ref<UserRole | null>(null)
  const email = ref<string>('')
  const groups = ref<string[]>([])
  const subject = ref<string>('')
  const isAuthenticated = computed(() => role.value !== null)
  const isAdmin = computed(() => role.value === 'admin')
  const oidcEnabled = ref(false)
  const devAuthBypassEnabled = import.meta.env.VITE_DEV_AUTH_BYPASS !== 'false'

  let _initPromise: Promise<void> | null = null

  function ready(): Promise<void> {
    return _initPromise ?? Promise.resolve()
  }

  function setAuth(data: Partial<UserProfile> & { username: string; role: UserRole }) {
    username.value = data.username
    role.value = data.role
    email.value = data.email ?? ''
    groups.value = data.groups ?? []
    subject.value = data.subject ?? ''
  }

  function clearAuth() {
    username.value = null
    role.value = null
    email.value = ''
    groups.value = []
    subject.value = ''
  }

  function init(): Promise<void> {
    _initPromise = (async () => {
      // DEV AUTH BYPASS: enabled by default in dev, can be disabled via VITE_DEV_AUTH_BYPASS=false.
      if (import.meta.env.DEV && devAuthBypassEnabled) {
        if (!isAuthenticated.value) {
          setAuth({ username: 'devuser', role: 'admin' })
        }
        return
      }

      // Check if OIDC is enabled
      try {
        const configRes = await fetch('/api/v1/auth/config')
        if (configRes.ok) {
          const config = await configRes.json()
          oidcEnabled.value = config.oidcEnabled
        }
      } catch (_) {
      }

      try {
        const res = await fetch('/api/v1/auth/me')
        if (res.ok) {
          const data = await res.json()
          setAuth(data)
        }
      } catch (_) {
      }
    })()
    return _initPromise
  }

  async function login(usr: string, pwd: string): Promise<void> {
    const res = await fetch('/api/v1/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: usr, password: pwd })
    });
    if (res.ok) {
      const data = await res.json();
      setAuth(data);
    } else {
      throw new Error('Invalid credentials');
    }
  }

  function loginWithOIDC(): void {
    window.location.href = '/api/v1/auth/oidc/start'
  }

  async function logout(): Promise<void> {
    await fetch('/api/v1/auth/logout', { method: 'POST' });
    clearAuth();
  }

  return { username, role, email, groups, subject, isAuthenticated, isAdmin, oidcEnabled, setAuth, clearAuth, init, ready, login, loginWithOIDC, logout }
})
