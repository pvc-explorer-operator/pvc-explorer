<template>
  <main v-if="scope" class="sd-page">

    <!-- Header -->
    <div class="sd-header detail-section">
      <div>
        <div class="sd-badges">
          <span class="sd-badge sd-badge--ns">
            <i class="pi pi-server"></i> {{ scope.status?.namespaceCount ?? 0 }} namespace{{ scope.status?.namespaceCount !== 1 ? 's' : '' }}
          </span>
          <span class="sd-badge sd-badge--exp">
            <i class="pi pi-database"></i> {{ scope.status?.explorerCount ?? 0 }} explorer{{ scope.status?.explorerCount !== 1 ? 's' : '' }}
          </span>
          <span class="sd-badge" :class="readyCondition ? 'sd-badge--ok' : 'sd-badge--warn'">
            <i class="pi" :class="readyCondition ? 'pi-check-circle' : 'pi-clock'"></i>
            {{ readyCondition ? 'Ready' : 'Reconciling' }}
          </span>
        </div>
      </div>

    </div>

    <!-- Body: 2-col -->
    <div class="sd-body">

      <!-- ── Left: Info ── -->
      <div class="sd-info-col">

        <!-- Overview -->
        <div class="sd-section detail-section">
          <div class="sd-section-title">Overview</div>
          <div class="sd-rows">
            <div class="sd-row">
              <span class="sd-row-label">Created</span>
              <span class="sd-row-value">{{ scope.metadata.creationTimestamp ? new Date(scope.metadata.creationTimestamp).toLocaleString() : '—' }}</span>
            </div>
            <div class="sd-row">
              <span class="sd-row-label">Deletion policy</span>
              <span class="sd-row-chip">{{ scope.spec.deletionPolicy || 'Cleanup' }}</span>
            </div>
            <div class="sd-row">
              <span class="sd-row-label">Discovery mode</span>
              <span class="sd-row-chip">{{ scope.spec.discovery?.mode || 'Auto' }}</span>
            </div>
            <div class="sd-row" v-if="scope.metadata.finalizers?.length">
              <span class="sd-row-label">Finalizers</span>
              <span class="sd-row-value mono">{{ scope.metadata.finalizers.join(', ') }}</span>
            </div>
          </div>
        </div>

        <!-- Namespaces -->
        <div class="sd-section detail-section">
          <div class="sd-section-title">Namespaces</div>
          <template v-if="scope.spec.namespaces?.names?.length">
            <div class="sd-row-label sd-sub-label">Explicit names</div>
            <div class="sd-tag-row">
              <span v-for="ns in scope.spec.namespaces.names" :key="ns" class="sd-tag">{{ ns }}</span>
            </div>
          </template>
          <template v-if="scope.spec.namespaces?.labelSelector?.matchLabels">
            <div class="sd-row-label sd-sub-label" style="margin-top: 0.6rem">Label selector</div>
            <div class="sd-tag-row">
              <span
                v-for="(v, k) in scope.spec.namespaces.labelSelector.matchLabels" :key="k"
                class="sd-tag sd-tag--label"
              >{{ k }}={{ v }}</span>
            </div>
          </template>
          <div v-if="!scope.spec.namespaces?.names?.length && !scope.spec.namespaces?.labelSelector" class="sd-empty">No namespace selector defined.</div>
        </div>

        <!-- Discovery: excludePVCs / pvcNames -->
        <div v-if="scope.spec.discovery?.excludePVCs?.length || scope.spec.discovery?.pvcNames?.length" class="sd-section detail-section">
          <div class="sd-section-title">Discovery</div>
          <template v-if="scope.spec.discovery?.pvcNames?.length">
            <div class="sd-row-label sd-sub-label">Explicit PVCs</div>
            <div class="sd-tag-row">
              <span v-for="p in scope.spec.discovery.pvcNames" :key="p" class="sd-tag">{{ p }}</span>
            </div>
          </template>
          <template v-if="scope.spec.discovery?.excludePVCs?.length">
            <div class="sd-row-label sd-sub-label" style="margin-top: 0.6rem">Excluded patterns</div>
            <div class="sd-tag-row">
              <span v-for="p in scope.spec.discovery.excludePVCs" :key="p" class="sd-tag sd-tag--warn">{{ p }}</span>
            </div>
          </template>
        </div>

        <!-- Defaults -->
        <div v-if="scope.spec.defaults" class="sd-section detail-section">
          <div class="sd-section-title">Defaults</div>
          <div class="sd-rows">
            <div class="sd-row" v-if="scope.spec.defaults.mode">
              <span class="sd-row-label">Mode</span>
              <span class="sd-row-chip">{{ scope.spec.defaults.mode }}</span>
            </div>
            <div class="sd-row" v-if="scope.spec.defaults.image">
              <span class="sd-row-label">Image</span>
              <span class="sd-row-value mono">{{ scope.spec.defaults.image }}</span>
            </div>
            <div class="sd-row">
              <span class="sd-row-label">Force RW</span>
              <span class="sd-row-bool" :class="scope.spec.defaults.forceRw !== false ? 'sd-row-bool--true' : 'sd-row-bool--false'">
                {{ scope.spec.defaults.forceRw !== false ? 'true' : 'false' }}
              </span>
            </div>
            <template v-if="scope.spec.defaults.scaling">
              <div class="sd-row">
                <span class="sd-row-label">Idle timeout</span>
                <span class="sd-row-value mono">{{ scope.spec.defaults.scaling.idleTimeout || '—' }}</span>
              </div>
              <div class="sd-row">
                <span class="sd-row-label">Startup timeout</span>
                <span class="sd-row-value mono">{{ scope.spec.defaults.scaling.startupTimeout || '—' }}</span>
              </div>
            </template>
            <template v-if="scope.spec.defaults.mountStrategy">
              <div class="sd-row" v-if="scope.spec.defaults.mountStrategy.allowNodeAffinity !== undefined">
                <span class="sd-row-label">Node affinity</span>
                <span class="sd-row-bool" :class="scope.spec.defaults.mountStrategy.allowNodeAffinity ? 'sd-row-bool--true' : 'sd-row-bool--false'">
                  {{ scope.spec.defaults.mountStrategy.allowNodeAffinity ? 'allowed' : 'disabled' }}
                </span>
              </div>
              <div class="sd-row" v-if="scope.spec.defaults.mountStrategy.fallbackOnConflict">
                <span class="sd-row-label">Fallback policy</span>
                <span class="sd-row-chip">{{ scope.spec.defaults.mountStrategy.fallbackOnConflict }}</span>
              </div>
            </template>
            <template v-if="scope.spec.defaults.resources">
              <div class="sd-row" v-if="scope.spec.defaults.resources.requests">
                <span class="sd-row-label">CPU request</span>
                <span class="sd-row-value mono">{{ scope.spec.defaults.resources.requests.cpu || '—' }}</span>
              </div>
              <div class="sd-row" v-if="scope.spec.defaults.resources.requests">
                <span class="sd-row-label">Memory request</span>
                <span class="sd-row-value mono">{{ scope.spec.defaults.resources.requests.memory || '—' }}</span>
              </div>
              <div class="sd-row" v-if="scope.spec.defaults.resources.limits">
                <span class="sd-row-label">CPU limit</span>
                <span class="sd-row-value mono">{{ scope.spec.defaults.resources.limits.cpu || '—' }}</span>
              </div>
              <div class="sd-row" v-if="scope.spec.defaults.resources.limits">
                <span class="sd-row-label">Memory limit</span>
                <span class="sd-row-value mono">{{ scope.spec.defaults.resources.limits.memory || '—' }}</span>
              </div>
            </template>
          </div>
        </div>

        <!-- Conditions -->
        <div v-if="scope.status?.conditions?.length" class="sd-section detail-section">
          <div class="sd-section-title">Conditions</div>
          <div class="sd-conditions">
            <div v-for="cond in scope.status.conditions" :key="cond.type" class="sd-condition">
              <span class="sd-condition-dot" :class="cond.status === 'True' ? 'sd-condition-dot--ok' : 'sd-condition-dot--warn'"></span>
              <div class="sd-condition-body">
                <div class="sd-condition-type">{{ cond.type }} <span class="sd-condition-status">{{ cond.status }}</span></div>
                <div class="sd-condition-reason" v-if="cond.reason">{{ cond.reason }}</div>
                <div class="sd-condition-msg" v-if="cond.message">{{ cond.message }}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- Explorers -->
        <div class="sd-section detail-section">
          <div class="sd-section-title">Explorers ({{ explorers.length }})</div>
          <div v-if="explorers.length" class="sd-explorer-list">
            <div
              v-for="e in explorers" :key="e.namespace + '/' + e.name"
              class="sd-explorer-item"
              @click="router.push(`/explorers/${e.namespace}/${e.name}`)"
            >
              <span class="phase-dot" :class="phaseDotClass(e.phase)"></span>
              <span class="sd-explorer-name">{{ e.name }}</span>
              <span class="sd-explorer-ns">{{ e.namespace }}</span>
              <span class="sd-explorer-pvc">{{ e.pvcName }}</span>
              <i class="pi pi-chevron-right sd-explorer-arrow"></i>
            </div>
          </div>
          <div v-else class="sd-empty">No explorers managed by this scope yet.</div>
        </div>

      </div><!-- /info-col -->

      <!-- ── Right: YAML ── -->
      <div class="sd-yaml-col">
        <div class="sd-yaml-card detail-section">
          <div class="sd-yaml-header">
            <span class="sd-yaml-title">Manifest</span>
            <div class="sd-yaml-actions">
              <button class="sd-action-btn" @click="copyYaml" :class="{ 'sd-action-btn--done': copied }">
                <i class="pi" :class="copied ? 'pi-check' : 'pi-copy'"></i>
                {{ copied ? 'Copied!' : 'Copy' }}
              </button>
              <button class="sd-action-btn sd-action-btn--primary" @click="downloadYaml">
                <i class="pi pi-download"></i> Download
              </button>
            </div>
          </div>
          <pre class="sd-yaml-pre"><code v-html="highlightedYaml"></code></pre>
          <div class="sd-yaml-footer">
            <code>kubectl apply -f {{ scope.metadata.name }}.yaml</code>
          </div>
        </div>
      </div>

    </div><!-- /body -->

  </main>
  <div v-else-if="loading" class="sd-page">
    <div class="sd-header">
      <div class="sd-badges">
        <Skeleton width="8rem" height="1.5rem" borderRadius="20px" />
        <Skeleton width="9rem" height="1.5rem" borderRadius="20px" />
        <Skeleton width="5rem" height="1.5rem" borderRadius="20px" />
      </div>
    </div>
    <div class="sd-body">
      <div class="sd-info-col">
        <div class="sd-section">
          <Skeleton width="4rem" height="0.72rem" />
          <div class="sd-rows" style="margin-top:0.75rem">
            <div class="sd-row"><Skeleton width="6rem" height="0.8125rem" /></div>
            <div class="sd-row"><Skeleton width="8rem" height="0.8125rem" /></div>
            <div class="sd-row"><Skeleton width="5rem" height="0.8125rem" /></div>
          </div>
        </div>
        <div class="sd-section">
          <Skeleton width="5.5rem" height="0.72rem" />
          <div class="sd-tag-row" style="margin-top:0.75rem">
            <Skeleton v-for="j in 5" :key="j" width="4rem" height="1.25rem" borderRadius="4px" />
          </div>
        </div>
        <div class="sd-section">
          <Skeleton width="4.5rem" height="0.72rem" />
          <div class="sd-explorer-list" style="margin-top:0.75rem">
            <div v-for="j in 4" :key="j" class="sd-explorer-item">
              <Skeleton shape="circle" size="0.5rem" />
              <Skeleton width="30%" height="0.875rem" />
              <Skeleton width="25%" height="0.78rem" />
              <Skeleton width="35%" height="0.78rem" />
            </div>
          </div>
        </div>
      </div>
      <div class="sd-yaml-col">
        <div class="sd-yaml-card">
          <div class="sd-yaml-header">
            <Skeleton width="4rem" height="0.72rem" />
            <div style="display:flex;gap:0.5rem">
              <Skeleton width="3.5rem" height="1.5rem" borderRadius="5px" />
              <Skeleton width="4rem" height="1.5rem" borderRadius="5px" />
            </div>
          </div>
          <div style="padding:1rem">
            <Skeleton v-for="j in 8" :key="j" width="60%" height="0.78rem" style="margin-bottom:0.5rem" />
          </div>
        </div>
      </div>
    </div>
  </div>
  <div v-else class="sd-loading">Scope not found.</div>

  <!-- keep old template open tag removed -->
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useExplorerStore } from '../stores/explorerStore'
import { highlightYaml } from '../utils/yamlHighlight'
import Skeleton from 'primevue/skeleton'

const route   = useRoute()
const router  = useRouter()
const store   = useExplorerStore()
const name    = route.params.name as string

// ── Types ─────────────────────────────────────────────────────────────
interface Scope {
  metadata: {
    name: string
    creationTimestamp?: string
    finalizers?: string[]
    [k: string]: any
  }
  spec: {
    namespaces?: {
      names?: string[]
      labelSelector?: { matchLabels?: Record<string, string>; matchExpressions?: any[] }
    }
    discovery?: { mode?: string; pvcNames?: string[]; excludePVCs?: string[] }
    deletionPolicy?: string
    defaults?: {
      mode?: string
      image?: string
      forceRw?: boolean
      scaling?: { idleTimeout?: string; startupTimeout?: string }
      mountStrategy?: { allowNodeAffinity?: boolean; fallbackOnConflict?: string }
      resources?: { requests?: { cpu?: string; memory?: string }; limits?: { cpu?: string; memory?: string } }
    }
    [k: string]: any
  }
  status?: {
    namespaceCount?: number
    explorerCount?: number
    observedGeneration?: number
    conditions?: Array<{ type: string; status: string; reason?: string; message?: string }>
    [k: string]: any
  }
}

// ── State ─────────────────────────────────────────────────────────────
const scope       = ref<Scope | null>(null)
const loading     = ref(true)
const copied      = ref(false)

// ── Derived ───────────────────────────────────────────────────────────
const readyCondition = computed(() =>
  scope.value?.status?.conditions?.some(c => c.type === 'Ready' && c.status === 'True') ?? false
)

const explorers = computed(() =>
  store.explorers.filter(e => {
    if (!e.labels) return false
    const lbls = e.labels as any
    if (Array.isArray(lbls)) return false
    return lbls['pvcexplorer.io/scope'] === name
  })
)

function phaseDotClass(phase: string) {
  const p = phase?.toLowerCase() ?? ''
  if (p === 'running') return 'dot-running'
  if (p === 'pending') return 'dot-pending'
  if (p === 'failed')  return 'dot-failed'
  return 'dot-warning'
}

// ── YAML generation ───────────────────────────────────────────────────
const highlightedYaml = computed(() => highlightYaml(scopeYaml.value))

const scopeYaml = computed(() => {
  const s = scope.value
  if (!s) return ''
  const lines: string[] = []
  lines.push('apiVersion: pvcexplorer.io/v1alpha1')
  lines.push('kind: PVCExplorerScope')
  lines.push('metadata:')
  lines.push(`  name: ${s.metadata.name}`)
  lines.push('spec:')

  // namespaces
  const ns = s.spec.namespaces
  lines.push('  namespaces:')
  if (ns?.names?.length) {
    lines.push('    names:')
    for (const n of ns.names) { lines.push(`      - ${n}`) }
  }
  if (ns?.labelSelector?.matchLabels && Object.keys(ns.labelSelector.matchLabels).length) {
    lines.push('    labelSelector:')
    lines.push('      matchLabels:')
    for (const [k, v] of Object.entries(ns.labelSelector.matchLabels)) { lines.push(`        ${k}: "${v}"`) }
  }
  if (!ns?.names?.length && !ns?.labelSelector) lines.push('    names: []')

  // discovery
  const disc = s.spec.discovery
  lines.push('  discovery:')
  lines.push(`    mode: ${disc?.mode ?? 'Auto'}`)
  if (disc?.pvcNames?.length) {
    lines.push('    pvcNames:')
    for (const p of disc.pvcNames) { lines.push(`      - ${p}`) }
  }
  if (disc?.excludePVCs?.length) {
    lines.push('    excludePVCs:')
    for (const p of disc.excludePVCs) { lines.push(`      - "${p}"`) }
  }

  // deletion policy
  lines.push(`  deletionPolicy: ${s.spec.deletionPolicy ?? 'Cleanup'}`)

  // defaults
  const d = s.spec.defaults
  if (d) {
    lines.push('  defaults:')
    if (d.mode)  lines.push(`    mode: ${d.mode}`)
    if (d.image) lines.push(`    image: "${d.image}"`)
    if (d.forceRw !== undefined) lines.push(`    forceRW: ${d.forceRw}`)
    if (d.scaling) {
      lines.push('    scaling:')
      if (d.scaling.idleTimeout)   lines.push(`      idleTimeout: "${d.scaling.idleTimeout}"`)
      if (d.scaling.startupTimeout) lines.push(`      startupTimeout: "${d.scaling.startupTimeout}"`)
    }
    if (d.mountStrategy) {
      lines.push('    mountStrategy:')
      if (d.mountStrategy.allowNodeAffinity !== undefined) lines.push(`      allowNodeAffinity: ${d.mountStrategy.allowNodeAffinity}`)
      if (d.mountStrategy.fallbackOnConflict) lines.push(`      fallbackOnConflict: ${d.mountStrategy.fallbackOnConflict}`)
    }
    if (d.resources) {
      lines.push('    resources:')
      if (d.resources.requests && (d.resources.requests.cpu || d.resources.requests.memory)) {
        lines.push('      requests:')
        if (d.resources.requests.cpu)    lines.push(`        cpu: "${d.resources.requests.cpu}"`)
        if (d.resources.requests.memory) lines.push(`        memory: "${d.resources.requests.memory}"`)
      }
      if (d.resources.limits && (d.resources.limits.cpu || d.resources.limits.memory)) {
        lines.push('      limits:')
        if (d.resources.limits.cpu)    lines.push(`        cpu: "${d.resources.limits.cpu}"`)
        if (d.resources.limits.memory) lines.push(`        memory: "${d.resources.limits.memory}"`)
      }
    }
  }

  return lines.join('\n')
})

// ── Copy / Download ───────────────────────────────────────────────────
async function copyYaml() {
  await navigator.clipboard.writeText(scopeYaml.value)
  copied.value = true
  setTimeout(() => { copied.value = false }, 2000)
}
function downloadYaml() {
  const blob = new Blob([scopeYaml.value], { type: 'text/yaml' })
  const url  = URL.createObjectURL(blob)
  const a    = document.createElement('a')
  a.href = url; a.download = `${name}.yaml`; a.click()
  URL.revokeObjectURL(url)
}

// ── Fetch ────────────────────────────────────────────────────────────
async function fetchScope() {
  loading.value = true
  const res = await fetch(`/api/v1/scopes/${encodeURIComponent(name)}`)
  scope.value = res.ok ? await res.json() : null
  loading.value = false
}

onMounted(() => { store.fetchExplorers(); fetchScope() })
</script>

<style scoped>
/* ── Page ── */
.sd-page {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  width: 100%;
}

/* ── Header ── */
.sd-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1rem;
  flex-wrap: wrap;
}
.sd-title {
  font-size: 1.4rem;
  font-weight: 700;
  color: var(--text-color);
  margin: 0 0 0.5rem;
  font-family: 'JetBrains Mono', monospace;
}
.sd-badges {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}
.sd-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.3rem;
  padding: 2px 10px;
  border-radius: 20px;
  font-size: 0.75rem;
  font-weight: 600;
  background: var(--surface-hover);
  color: var(--text-color-secondary);
  border: 1px solid var(--surface-border);
}
.sd-badge--ns  { background: rgba(79,142,247,0.1);  color: #0ea5e9; border-color: rgba(79,142,247,0.25); }
.sd-badge--exp { background: rgba(168,85,247,0.1);  color: #a855f7; border-color: rgba(168,85,247,0.25); }
.sd-badge--ok  { background: rgba(34,197,94,0.1);   color: #22c55e;  border-color: rgba(34,197,94,0.25); }
.sd-badge--warn{ background: rgba(245,158,11,0.1);  color: #f59e0b;  border-color: rgba(245,158,11,0.25); }

/* ── Body ── */
.sd-body {
  display: flex;
  gap: 1.5rem;
  align-items: flex-start;
}
.sd-info-col {
  flex: 3;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
.sd-yaml-col {
  width: 42%;
  min-width: 380px;
  flex-shrink: 0;
  position: sticky;
  top: 1rem;
}

/* ── Section ── */
.sd-section {
  background: var(--surface-card);
  border: 1px solid var(--surface-border);
  border-radius: 8px;
  padding: 1rem 1.25rem;
}
.sd-section-title {
  font-size: 0.72rem;
  font-weight: 700;
  letter-spacing: 0.07em;
  text-transform: uppercase;
  color: var(--text-color-secondary);
  margin: 0 0 0.75rem;
}

/* ── Info rows ── */
.sd-rows { display: flex; flex-direction: column; gap: 0.4rem; }
.sd-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 3px 0;
  border-bottom: 1px solid var(--surface-border);
}
.sd-row:last-child { border-bottom: none; }
.sd-row-label {
  font-size: 0.8125rem;
  color: var(--text-color-secondary);
  white-space: nowrap;
}
.sd-sub-label {
  font-size: 0.72rem;
  font-weight: 600;
  letter-spacing: 0.04em;
  color: var(--text-color-secondary);
  margin-bottom: 0.3rem;
}
.sd-row-value {
  font-size: 0.8125rem;
  color: var(--text-color);
  text-align: right;
}
.sd-row-value.mono { font-family: 'JetBrains Mono', monospace; font-size: 0.78rem; }
.sd-row-chip {
  display: inline-block;
  padding: 1px 8px;
  background: var(--surface-hover);
  color: var(--text-color);
  border: 1px solid var(--surface-border);
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 600;
}
.sd-row-bool {
  font-size: 0.75rem;
  font-weight: 700;
  font-family: 'JetBrains Mono', monospace;
}
.sd-row-bool--true  { color: #22c55e; }
.sd-row-bool--false { color: var(--p-red-400); }

/* ── Tag row ── */
.sd-tag-row { display: flex; flex-wrap: wrap; gap: 5px; }
.sd-tag {
  display: inline-flex; align-items: center;
  padding: 2px 9px;
  background: rgba(79,142,247,0.1);
  color: #0ea5e9;
  border-radius: 4px;
  font-size: 0.78rem;
  font-family: 'JetBrains Mono', monospace;
}
.sd-tag--label {
  background: rgba(168,85,247,0.1);
  color: #a855f7;
}
.sd-tag--warn {
  background: rgba(245,158,11,0.1);
  color: #f59e0b;
}
.sd-empty {
  font-size: 0.82rem;
  color: var(--text-color-secondary);
  padding: 0.25rem 0;
}

/* ── Conditions ── */
.sd-conditions { display: flex; flex-direction: column; gap: 0.6rem; }
.sd-condition {
  display: flex; align-items: flex-start; gap: 0.6rem;
  padding: 0.5rem 0.6rem;
  background: var(--surface-ground);
  border-radius: 6px;
  border: 1px solid var(--surface-border);
}
.sd-condition-dot {
  width: 8px; height: 8px; border-radius: 50%;
  flex-shrink: 0; margin-top: 4px;
}
.sd-condition-dot--ok   { background: #22c55e; }
.sd-condition-dot--warn { background: #f59e0b; }
.sd-condition-type   { font-size: 0.82rem; font-weight: 700; color: var(--text-color); }
.sd-condition-status { font-size: 0.75rem; font-weight: 400; color: var(--text-color-secondary); margin-left: 0.4rem; }
.sd-condition-reason { font-size: 0.75rem; color: var(--text-color-secondary); margin-top: 1px; }
.sd-condition-msg    { font-size: 0.75rem; color: var(--text-color-secondary); margin-top: 2px; }

/* ── Explorers list ── */
.sd-explorer-list { display: flex; flex-direction: column; gap: 0.4rem; }
.sd-explorer-item {
  display: flex; align-items: center; gap: 0.6rem;
  padding: 0.5rem 0.75rem;
  border: 1px solid var(--surface-border);
  border-radius: 6px;
  cursor: pointer;
  transition: border-color 0.12s, background 0.12s;
}
.sd-explorer-item:hover { border-color: var(--p-primary-400); background: var(--surface-hover); }
.sd-explorer-name { font-size: 0.875rem; font-weight: 600; color: var(--text-color); flex-shrink: 0; }
.sd-explorer-ns   { font-size: 0.78rem; color: var(--text-color-secondary); flex-shrink: 0; }
.sd-explorer-pvc  { font-size: 0.78rem; color: var(--text-color-secondary); font-family: 'JetBrains Mono', monospace; flex: 1; min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.sd-explorer-arrow { margin-left: auto; color: var(--text-color-secondary); font-size: 0.7rem; flex-shrink: 0; }

/* ── Phase dots ── */
.phase-dot { display: inline-block; width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; background: var(--text-color-secondary); }
.dot-running { background: #22c55e; }
.dot-pending { background: #f59e0b; }
.dot-failed  { background: var(--p-red-400); }
.dot-warning { background: #f59e0b; }

/* ── YAML panel ── */
.sd-yaml-card {
  background: var(--surface-card);
  border: 1px solid var(--surface-border);
  border-radius: 8px;
  overflow: hidden;
  display: flex; flex-direction: column;
}
.sd-yaml-header {
  display: flex; align-items: center; justify-content: space-between;
  padding: 0.6rem 1rem;
  border-bottom: 1px solid var(--surface-border);
  background: var(--surface-card);
}
.sd-yaml-title {
  font-size: 0.72rem; font-weight: 700; letter-spacing: 0.05em; text-transform: uppercase; color: var(--text-color-secondary);
}
.sd-yaml-actions { display: flex; gap: 0.5rem; }
.sd-action-btn {
  display: inline-flex; align-items: center; gap: 0.35rem;
  padding: 3px 12px;
  border: 1px solid var(--surface-border);
  background: transparent; color: var(--text-color-secondary);
  border-radius: 5px; cursor: pointer; font-size: 0.8rem; font-family: inherit;
  transition: background 0.1s, color 0.1s;
}
.sd-action-btn:hover { background: var(--surface-hover); color: var(--text-color); }
.sd-action-btn--primary { background: var(--p-primary-500); color: #fff; border-color: var(--p-primary-500); }
.sd-action-btn--primary:hover { background: var(--p-primary-400); border-color: var(--p-primary-400); color: #fff; }
.sd-action-btn--done { color: #22c55e; border-color: #22c55e; }
.sd-yaml-pre {
  margin: 0; padding: 1rem;
  background: var(--surface-card);
  font-family: 'JetBrains Mono', 'Fira Mono', monospace;
  font-size: 0.78rem; line-height: 1.6;
  color: var(--text-color);
  overflow-x: auto;
  max-height: calc(100vh - 200px); overflow-y: auto;
  white-space: pre;
}
.sd-yaml-footer {
  padding: 0.5rem 1rem;
  font-size: 0.75rem; color: var(--text-color-secondary);
  border-top: 1px solid var(--surface-border);
  background: var(--surface-card);
}
.sd-yaml-footer code { font-family: 'JetBrains Mono', monospace; font-size: 0.75rem; }
.sd-loading {
  display: flex; align-items: center; gap: 0.5rem;
  padding: 3rem; justify-content: center;
  color: var(--text-color-secondary); font-size: 1rem;
}

@media (max-width: 860px) {
  .sd-body { flex-direction: column; }
  .sd-yaml-col { width: 100%; position: static; }
  .sd-yaml-pre { max-height: 300px; }
}

@media (prefers-reduced-motion: no-preference) {
  @supports (animation-timeline: scroll()) {
    .detail-section {
      animation: sd-fade-in-up linear both;
      animation-timeline: view();
      animation-range: entry 0% entry 25%;
    }
    @keyframes sd-fade-in-up {
      from { opacity: 0; translate: 0 20px; }
      to   { opacity: 1; translate: 0; }
    }
  }
}
</style>
