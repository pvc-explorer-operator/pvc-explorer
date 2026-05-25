const PALETTE = [
  '#6366f1', '#8b5cf6', '#a855f7', '#d946ef',
  '#ec4899', '#f43f5e', '#fb7185', '#f97316',
  '#14b8a6', '#06b6d4', '#0ea5e9', '#22c55e',
  '#84cc16', '#eab308', '#f59e0b',
]

export function stringToColor(s: string): string {
  let hash = 5381
  for (let i = 0; i < s.length; i++) {
    hash = ((hash << 5) + hash) + s.charCodeAt(i)
  }
  return PALETTE[Math.abs(hash) % PALETTE.length]
}

export function mountStateColor(state: string): string {
  const s = state.toLowerCase()
  if (s === 'mounted') return '#22c55e'
  if (s === 'readonly') return '#f59e0b'
  if (s === 'conflict') return '#ef4444'
  return '#94a3b8'
}

export function accessModeColor(mode: string): string {
  const m = mode.toLowerCase()
  if (m.includes('once')) return '#6366f1'
  if (m.includes('many')) return '#14b8a6'
  return '#8b5cf6'
}

export type ConsumerFilter = '' | 'has' | 'none'
export type CreatedFilter = '' | '24h' | '7d' | '30d' | 'older'

export interface Filters {
  search: string
  phases: string[]
  namespaces: string[]
  mountStates: string[]
  scopes: string[]
  accessModes: string[]
  consumers: ConsumerFilter
  created: CreatedFilter
  labels: string[]
}
