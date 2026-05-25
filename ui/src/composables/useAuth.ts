import { useAuthStore } from '@/stores/authStore'
import router from '@/router'

export async function apiFetch(input: RequestInfo, init?: RequestInit): Promise<Response> {
  const res = await fetch(input, init)
  if (res.status === 401) {
    const auth = useAuthStore()
    auth.clearAuth()
    router.push('/login')
  }
  return res
}

export function useAuth() {}
