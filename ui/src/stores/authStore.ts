import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const username = ref<string | null>(null)
  const role = ref<'admin' | 'viewer' | null>(null)
  const isAuthenticated = computed(() => role.value !== null)
  const isAdmin = computed(() => role.value === 'admin')
  const devAuthBypassEnabled = import.meta.env.VITE_DEV_AUTH_BYPASS !== 'false'

  let _initPromise: Promise<void> | null = null

  function ready(): Promise<void> {
    return _initPromise ?? Promise.resolve()
  }

  function setAuth(u: string, r: 'admin' | 'viewer') {
    username.value = u
    role.value = r
  }

  function clearAuth() {
    username.value = null
    role.value = null
  }

  function init(): Promise<void> {
    _initPromise = (async () => {
      // DEV AUTH BYPASS: enabled by default in dev, can be disabled via VITE_DEV_AUTH_BYPASS=false.
      if (import.meta.env.DEV && devAuthBypassEnabled) {
        if (!isAuthenticated.value) {
          setAuth('devuser', 'admin')
        }
        return
      }
      try {
        const res = await fetch('/api/v1/auth/me')
        if (res.ok) {
          const data = await res.json()
          setAuth(data.username, data.role)
        }
      } catch (_) {
      }
    })()
    return _initPromise
  }

  async function login(username: string, password: string): Promise<void> {
    const res = await fetch('/api/v1/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password })
    });
    if (res.ok) {
      const data = await res.json();
      setAuth(data.username, data.role);
    } else {
      throw new Error('Invalid credentials');
    }
  }

  async function logout(): Promise<void> {
    await fetch('/api/v1/auth/logout', { method: 'POST' });
    clearAuth();
  }

  return { username, role, isAuthenticated, isAdmin, setAuth, clearAuth, init, ready, login, logout }
})
