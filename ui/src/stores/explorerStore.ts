import { defineStore } from 'pinia';
import { ref } from 'vue';
import { apiFetch } from '@/composables/useAuth'

let countdownTimer: ReturnType<typeof setInterval> | null = null;

function ensureCountdown(idleRemaining: { value: Record<string, number> }) {
  if (countdownTimer) return;
  countdownTimer = setInterval(() => {
    const updated: Record<string, number> = {};
    let hasActive = false;
    for (const [key, val] of Object.entries(idleRemaining.value)) {
      const next = val > 0 ? val - 1 : 0;
      updated[key] = next;
      if (next > 0) hasActive = true;
    }
    idleRemaining.value = updated;
    if (!hasActive && countdownTimer) {
      clearInterval(countdownTimer);
      countdownTimer = null;
    }
  }, 1000);
}

export interface ExplorerCondition {
  type: string;
  status: string;
  reason?: string;
  message?: string;
  lastTransitionTime?: string;
}

export interface ConsumerInfo {
  podName: string;
  ownerKind?: string;
  ownerName?: string;
  nodeName?: string;
  mountReadOnly?: boolean;
}

export interface Explorer {
  name: string;
  namespace: string;
  phase: string;
  mountState: string;
  pvcName: string;
  scope?: string;
  accessMode?: string;
  mode?: string;
  labels?: string[];
  createdAt?: string;
  idleTimeout?: string;
  consumerCount?: number;
  conditions?: ExplorerCondition[];
  consumers?: ConsumerInfo[];
}

export interface Scope {
  name: string;
  phase: string;
  namespaceCount: number;
  explorerCount: number;
}

function explorerFromK8s(raw: Record<string, unknown>): Explorer {
  const meta = (raw.metadata ?? {}) as Record<string, unknown>;
  const spec = (raw.spec ?? {}) as Record<string, unknown>;
  const status = (raw.status ?? {}) as Record<string, unknown>;
  const mount = (status.mount ?? {}) as Record<string, unknown>;
  const labels = (meta.labels ?? {}) as Record<string, string>;
  const scaling = (spec.scaling ?? {}) as Record<string, unknown>;
  const consumers = (mount.consumers ?? []) as unknown[];
  return {
    name: meta.name as string,
    namespace: meta.namespace as string,
    phase: status.phase as string ?? '',
    mountState: mount.strategy as string ?? '',
    pvcName: spec.pvcName as string ?? '',
    scope: labels['pvcexplorer.io/scope'],
    accessMode: mount.accessMode as string,
    mode: status.mode as string,
    labels: Object.entries(labels)
      .filter(([k]) => !k.startsWith('pvcexplorer.io/'))
      .map(([k, v]) => `${k}=${v}`),
    createdAt: meta.creationTimestamp as string,
    idleTimeout: scaling.idleTimeout as string,
    consumerCount: consumers.length,
    conditions: status.conditions as ExplorerCondition[] | undefined,
    consumers: consumers as ConsumerInfo[],
  };
}

function scopeFromK8s(raw: Record<string, unknown>): Scope {
  const meta = (raw.metadata ?? {}) as Record<string, unknown>;
  const status = (raw.status ?? {}) as Record<string, unknown>;
  const conditions = (status.conditions ?? []) as Array<Record<string, unknown>>;
  const ready = conditions.find(c => c.type === 'Ready');
  return {
    name: meta.name as string,
    phase: ready ? (ready.reason as string ?? '') : 'Unknown',
    namespaceCount: status.namespaceCount as number ?? 0,
    explorerCount: status.explorerCount as number ?? 0,
  };
}

function sortExplorers(list: Explorer[]): Explorer[] {
  return [...list].sort((a, b) => {
    const nsA = a.namespace ?? '';
    const nsB = b.namespace ?? '';
    if (nsA !== nsB) return nsA.localeCompare(nsB);
    return (a.name ?? '').localeCompare(b.name ?? '');
  });
}

export const useExplorerStore = defineStore('explorer', () => {
  const explorers = ref<Explorer[]>([]);
  const scopes = ref<Scope[]>([]);
  const sidebarFilters = ref({ search: '', phases: [] as string[], namespaces: [] as string[], mountStates: [] as string[] });
  const idleRemaining = ref<Record<string, number>>({});

  function setSidebarFilters(f: typeof sidebarFilters.value) {
    sidebarFilters.value = { ...f };
  }

  function setSnapshot(rawExplorers: unknown[], rawScopes: unknown[]) {
    explorers.value = sortExplorers((rawExplorers ?? []).map(e => explorerFromK8s(e as Record<string, unknown>)));
    scopes.value = (rawScopes ?? []).map(s => scopeFromK8s(s as Record<string, unknown>));
  }

  function upsertExplorer(raw: unknown) {
    const e = explorerFromK8s(raw as Record<string, unknown>);
    const idx = explorers.value.findIndex(x => x.name === e.name && x.namespace === e.namespace);
    if (idx !== -1) {
      explorers.value[idx] = e;
    } else {
      explorers.value = sortExplorers([...explorers.value, e]);
    }
  }

  function removeExplorer(namespace: string, name: string) {
    explorers.value = explorers.value.filter(x => !(x.namespace === namespace && x.name === name));
    delete idleRemaining.value[`${namespace}/${name}`];
  }

  function setIdleRemaining(namespace: string, name: string, seconds: number) {
    idleRemaining.value = { ...idleRemaining.value, [`${namespace}/${name}`]: seconds };
    ensureCountdown(idleRemaining);
  }

  function upsertScope(raw: unknown) {
    const s = scopeFromK8s(raw as Record<string, unknown>);
    const idx = scopes.value.findIndex(x => x.name === s.name);
    if (idx !== -1) {
      scopes.value[idx] = s;
    } else {
      scopes.value.push(s);
    }
  }

  function removeScope(name: string) {
    scopes.value = scopes.value.filter(x => x.name !== name);
  }

  async function fetchExplorers(filters?: Record<string, string>): Promise<void> {
    let url = '/api/v1/explorers';
    if (filters && Object.keys(filters).length > 0) {
      url += '?' + new URLSearchParams(filters).toString();
    }
    const res = await apiFetch(url);
    if (!res.ok) throw new Error('Failed to fetch explorers');
    const data = await res.json();
    explorers.value = sortExplorers((data as unknown[]).map(e => explorerFromK8s(e as Record<string, unknown>)));
  }

  async function fetchScopes(): Promise<void> {
    const res = await apiFetch('/api/v1/scopes');
    if (!res.ok) throw new Error('Failed to fetch scopes');
    const data = await res.json();
    scopes.value = (data as unknown[]).map(s => scopeFromK8s(s as Record<string, unknown>));
  }

  async function fetchExplorer(ns: string, name: string): Promise<Explorer> {
    const res = await apiFetch(`/api/v1/explorers/${encodeURIComponent(ns)}/${encodeURIComponent(name)}`);
    if (!res.ok) throw new Error('Failed to fetch explorer');
    const raw = await res.json();
    const e = explorerFromK8s(raw);
    upsertExplorer(raw);
    return e;
  }

  function updatePhase(ns: string, name: string, phase: string) {
    const idx = explorers.value.findIndex(x => x.name === name && x.namespace === ns)
    if (idx !== -1) {
      explorers.value[idx] = { ...explorers.value[idx], phase }
    }
  }

  async function wakeExplorer(ns: string, name: string): Promise<void> {
    const res = await apiFetch(`/api/v1/explorers/${encodeURIComponent(ns)}/${encodeURIComponent(name)}/wake`, {
      method: 'POST',
    });
    if (!res.ok) throw new Error('Failed to wake explorer');
  }

  async function sleepExplorer(ns: string, name: string): Promise<void> {
    const res = await apiFetch(`/api/v1/explorers/${encodeURIComponent(ns)}/${encodeURIComponent(name)}/sleep`, {
      method: 'POST',
    });
    if (!res.ok) throw new Error('Failed to sleep explorer');
  }

  function teardown() {
    if (countdownTimer) {
      clearInterval(countdownTimer);
      countdownTimer = null;
    }
  }

  return {
    explorers,
    scopes,
    sidebarFilters,
    idleRemaining,
    setSidebarFilters,
    setSnapshot,
    upsertExplorer,
    removeExplorer,
    upsertScope,
    removeScope,
    setIdleRemaining,
    fetchExplorers,
    fetchScopes,
    fetchExplorer,
    wakeExplorer,
    sleepExplorer,
    updatePhase,
    teardown,
  };
});
