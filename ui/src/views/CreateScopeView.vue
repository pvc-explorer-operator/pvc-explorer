<template>
  <div v-if="authStore.isAdmin" class="cs-page">

    <!-- Header -->
    <div class="cs-header">
      <div>
        <div class="cs-info-banner">
          <i class="pi pi-info-circle cs-info-icon" aria-hidden="true"></i>
          <span>Fill in the fields below to generate a <code>PVCExplorerScope</code> manifest you can apply with <code>kubectl apply -f</code>.</span>
        </div>
      </div>
      <button class="cs-cancel-btn" @click="router.push('/scopes')">
        <i class="pi pi-times"></i> Cancel
      </button>
    </div>

    <!-- Two-column layout -->
    <div class="cs-body">

      <!-- ── Left: Form ── -->
      <div class="cs-form-col">

        <!-- Identity -->
        <div class="cs-section">
          <div class="cs-section-title">Identity</div>
          <div class="cs-field">
            <label class="cs-label" for="scope-name">Name <span class="cs-required">*</span></label>
            <input
              id="scope-name"
              v-model="name"
              class="cs-input"
              :class="{ 'cs-input--error': nameError }"
              placeholder="my-scope"
              maxlength="63"
              autocomplete="off"
              spellcheck="false"
            />
            <span v-if="nameError" class="cs-field-error">{{ nameError }}</span>
            <span v-else class="cs-field-hint">Lowercase letters, numbers and hyphens only.</span>
          </div>
        </div>

        <!-- Namespaces -->
        <div class="cs-section">
          <div class="cs-section-title">Namespaces</div>
          <div class="cs-field">
            <label class="cs-label" for="ns-tag-input">Namespace names</label>
            <div class="cs-tag-input" @click="focusTagInput('ns')">
              <span v-for="(ns, i) in nsNames" :key="i" class="cs-tag">
                {{ ns }}<button class="cs-tag-remove" :aria-label="`Remove namespace ${ns}`" @click.stop="nsNames.splice(i, 1)">×</button>
              </span>
              <input
                id="ns-tag-input"
                ref="nsInputEl"
                v-model="nsInputVal"
                class="cs-tag-inner-input"
                placeholder="namespace…"
                @keydown.enter.prevent="addTag(nsNames, 'ns')"
                @keydown.tab.prevent="addTag(nsNames, 'ns')"
                @keydown.comma.prevent="addTag(nsNames, 'ns')"
                @blur="addTag(nsNames, 'ns')"
              />
            </div>
            <span class="cs-field-hint">Press Enter or comma to add. Leave empty to match all namespaces via label selector.</span>
          </div>
          <fieldset class="cs-field">
            <legend class="cs-label">Label selector <span class="cs-optional">(optional)</span></legend>
            <div class="cs-label-pairs">
              <div v-for="(pair, i) in labelPairs" :key="i" class="cs-label-pair">
                <input v-model="pair.key"   class="cs-input cs-input--sm" placeholder="key" spellcheck="false" />
                <span class="cs-label-pair-eq">=</span>
                <input v-model="pair.value" class="cs-input cs-input--sm" placeholder="value" spellcheck="false" />
                <button class="cs-icon-btn cs-icon-btn--danger" @click="labelPairs.splice(i, 1)" title="Remove">×</button>
              </div>
              <button class="cs-add-btn" @click="labelPairs.push({ key: '', value: '' })">
                <i class="pi pi-plus"></i> Add label
              </button>
            </div>
            <span class="cs-field-hint">Namespaces matching <em>any</em> of these labels are also included.</span>
          </fieldset>
        </div>

        <!-- Discovery -->
        <div class="cs-section">
          <div class="cs-section-title">Discovery</div>
          <div class="cs-field">
            <label class="cs-label" id="mode-label">Mode</label>
            <div class="cs-chip-group" role="radiogroup" aria-labelledby="mode-label">
              <button
                v-for="m in ['Auto', 'Explicit']" :key="m"
                class="cs-chip" :class="{ 'cs-chip--active': discoveryMode === m }"
                role="radio" :aria-checked="discoveryMode === m"
                @click="discoveryMode = m as 'Auto' | 'Explicit'"
              >{{ m }}</button>
            </div>
            <span class="cs-field-hint">
              <template v-if="discoveryMode === 'Auto'">All PVCs in registered namespaces, minus excluded patterns.</template>
              <template v-else>Only the PVC names you list explicitly.</template>
            </span>
          </div>
          <div v-if="discoveryMode === 'Explicit'" class="cs-field">
            <label class="cs-label" for="pvc-tag-input">PVC names</label>
            <div class="cs-tag-input" @click="focusTagInput('pvc')">
              <span v-for="(p, i) in pvcNames" :key="i" class="cs-tag">
                {{ p }}<button class="cs-tag-remove" :aria-label="`Remove PVC ${p}`" @click.stop="pvcNames.splice(i, 1)">×</button>
              </span>
              <input
                id="pvc-tag-input"
                ref="pvcInputEl"
                v-model="pvcInputVal"
                class="cs-tag-inner-input"
                placeholder="my-pvc…"
                @keydown.enter.prevent="addTag(pvcNames, 'pvc')"
                @keydown.tab.prevent="addTag(pvcNames, 'pvc')"
                @keydown.comma.prevent="addTag(pvcNames, 'pvc')"
                @blur="addTag(pvcNames, 'pvc')"
              />
            </div>
          </div>
          <div class="cs-field">
            <label class="cs-label" for="excl-tag-input">Exclude PVCs <span class="cs-optional">(glob patterns)</span></label>
            <div class="cs-tag-input" @click="focusTagInput('excl')">
              <span v-for="(p, i) in excludePVCs" :key="i" class="cs-tag cs-tag--warn">
                {{ p }}<button class="cs-tag-remove" :aria-label="`Remove exclude pattern ${p}`" @click.stop="excludePVCs.splice(i, 1)">×</button>
              </span>
              <input
                id="excl-tag-input"
                ref="exclInputEl"
                v-model="exclInputVal"
                class="cs-tag-inner-input"
                placeholder="tmp-*…"
                @keydown.enter.prevent="addTag(excludePVCs, 'excl')"
                @keydown.tab.prevent="addTag(excludePVCs, 'excl')"
                @keydown.comma.prevent="addTag(excludePVCs, 'excl')"
                @blur="addTag(excludePVCs, 'excl')"
              />
            </div>
            <span class="cs-field-hint">Supports glob syntax, e.g. <code>tmp-*</code>, <code>*-backup</code>.</span>
          </div>
        </div>

        <!-- Deletion Policy -->
        <div class="cs-section">
          <div class="cs-section-title" id="deletion-policy-label">Deletion Policy</div>
          <div class="cs-field">
            <div class="cs-chip-group" role="radiogroup" aria-labelledby="deletion-policy-label">
              <button
                v-for="p in deletionPolicies" :key="p.value"
                class="cs-chip" :class="{ 'cs-chip--active': deletionPolicy === p.value }"
                role="radio" :aria-checked="deletionPolicy === p.value"
                @click="deletionPolicy = p.value as 'Cleanup' | 'Orphan'"
              >{{ p.label }}</button>
            </div>
            <span class="cs-field-hint">
              <template v-if="deletionPolicy === 'Cleanup'">Deletes all owned PVCExplorer CRs when this scope is removed. PVCs are never touched.</template>
              <template v-else>Leaves PVCExplorer CRs in place as standalone resources when this scope is removed.</template>
            </span>
          </div>
        </div>

        <!-- Defaults -->
        <div class="cs-section">
          <button class="cs-section-title cs-section-title--toggle" @click="showDefaults = !showDefaults">
            <i class="pi" :class="showDefaults ? 'pi-chevron-down' : 'pi-chevron-right'"></i>
            Defaults
            <span class="cs-optional cs-section-badge">applied to every created explorer</span>
          </button>
          <template v-if="showDefaults">
            <div class="cs-field">
              <label class="cs-label" id="explorer-mode-label">Explorer mode</label>
              <div class="cs-chip-group" role="radiogroup" aria-labelledby="explorer-mode-label">
                <button
                  v-for="m in ['ScaledToZero', 'Running']" :key="m"
                  class="cs-chip" :class="{ 'cs-chip--active': defaultMode === m }"
                  role="radio" :aria-checked="defaultMode === m"
                  @click="defaultMode = m as 'ScaledToZero' | 'Running'"
                >{{ m }}</button>
              </div>
            </div>
            <div class="cs-field">
              <label class="cs-label" for="agent-image">Agent image <span class="cs-optional">(optional)</span></label>
              <input id="agent-image" v-model="image" class="cs-input" placeholder="ghcr.io/org/pvc-explorer-agent:latest" spellcheck="false" />
            </div>
            <div class="cs-field cs-field--row">
              <label class="cs-label cs-label--inline" for="force-rw">Force read-write mount</label>
              <input id="force-rw" type="checkbox" v-model="forceRW" class="cs-checkbox" />
            </div>
            <div class="cs-field-row-2">
              <div class="cs-field">
                <label class="cs-label" for="idle-timeout">Idle timeout</label>
                <input id="idle-timeout" v-model="idleTimeout" class="cs-input" placeholder="10m" spellcheck="false" />
              </div>
              <div class="cs-field">
                <label class="cs-label" for="startup-timeout">Startup timeout</label>
                <input id="startup-timeout" v-model="startupTimeout" class="cs-input" placeholder="60s" spellcheck="false" />
              </div>
            </div>
            <div class="cs-field cs-field--row">
              <label class="cs-label cs-label--inline" for="allow-node-affinity">Allow node affinity</label>
              <input id="allow-node-affinity" type="checkbox" v-model="allowNodeAffinity" class="cs-checkbox" />
            </div>
            <div class="cs-field cs-field--row">
              <label class="cs-label cs-label--inline" for="enable-resources">Set resource requests/limits</label>
              <input id="enable-resources" type="checkbox" v-model="enableResources" class="cs-checkbox" />
            </div>
            <template v-if="enableResources">
              <div class="cs-resources-grid">
                <span class="cs-resources-head"></span>
                <span class="cs-resources-head">CPU</span>
                <span class="cs-resources-head">Memory</span>
                <span class="cs-resources-row-label">Requests</span>
                <input v-model="cpuRequest"  class="cs-input cs-input--sm" placeholder="50m" spellcheck="false" />
                <input v-model="memRequest"  class="cs-input cs-input--sm" placeholder="64Mi" spellcheck="false" />
                <span class="cs-resources-row-label">Limits</span>
                <input v-model="cpuLimit"    class="cs-input cs-input--sm" placeholder="200m" spellcheck="false" />
                <input v-model="memLimit"    class="cs-input cs-input--sm" placeholder="256Mi" spellcheck="false" />
              </div>
            </template>
          </template>
        </div>

      </div><!-- /form-col -->

      <!-- ── Right: YAML preview ── -->
      <div class="cs-yaml-col">
        <div class="cs-yaml-card">
          <div class="cs-yaml-header">
            <span class="cs-yaml-title">Generated YAML</span>
            <div class="cs-yaml-actions">
              <button class="cs-action-btn" @click="copyYaml" :class="{ 'cs-action-btn--done': copied }">
                <i class="pi" :class="copied ? 'pi-check' : 'pi-copy'"></i>
                {{ copied ? 'Copied!' : 'Copy' }}
              </button>
              <button class="cs-action-btn cs-action-btn--primary" @click="downloadYaml" :disabled="!!nameError || !name">
                <i class="pi pi-download"></i> Download
              </button>
            </div>
          </div>
          <pre class="cs-yaml-pre"><code v-html="highlightedYaml"></code></pre>
          <div class="cs-yaml-footer">
            Apply with: <code>kubectl apply -f {{ name || 'scope' }}.yaml</code>
          </div>
        </div>
      </div>

    </div><!-- /body -->

  </div>
  <div v-else class="cs-not-admin">Admin access required.</div>
</template>

<script setup lang="ts">
import { ref, computed, nextTick, onMounted, onUnmounted, type Ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/authStore'
import { highlightYaml } from '../utils/yamlHighlight'

const authStore = useAuthStore()
const router = useRouter()

// ── Identity ─────────────────────────────────────────────────────────
const name = ref('')

const nameError = computed(() => {
  const v = name.value
  if (!v) return ''
  if (!/^[a-z0-9]([a-z0-9-]*[a-z0-9])?$/.test(v))
    return 'Must start/end with a letter or digit and contain only lowercase letters, digits and hyphens.'
  if (v.length > 63) return 'Max 63 characters.'
  return ''
})

// ── Namespaces ────────────────────────────────────────────────────────
const nsNames    = ref<string[]>([])
const nsInputVal = ref('')
const nsInputEl  = ref<HTMLInputElement | null>(null)

interface LabelPair { key: string; value: string }
const labelPairs = ref<LabelPair[]>([])

// ── Discovery ─────────────────────────────────────────────────────────
const discoveryMode = ref<'Auto' | 'Explicit'>('Auto')
const pvcNames      = ref<string[]>([])
const pvcInputVal   = ref('')
const pvcInputEl    = ref<HTMLInputElement | null>(null)
const excludePVCs   = ref<string[]>([])
const exclInputVal  = ref('')
const exclInputEl   = ref<HTMLInputElement | null>(null)

// ── Deletion Policy ───────────────────────────────────────────────────
const deletionPolicy  = ref<'Cleanup' | 'Orphan'>('Cleanup')
const deletionPolicies = [
  { label: 'Cleanup', value: 'Cleanup' },
  { label: 'Orphan',  value: 'Orphan'  },
]

// ── Defaults ──────────────────────────────────────────────────────────
const showDefaults     = ref(false)
const defaultMode      = ref<'ScaledToZero' | 'Running'>('ScaledToZero')
const image            = ref('')
const forceRW          = ref(true)
const idleTimeout      = ref('10m')
const startupTimeout   = ref('60s')
const allowNodeAffinity= ref(true)
const enableResources  = ref(false)
const cpuRequest       = ref('50m')
const memRequest       = ref('64Mi')
const cpuLimit         = ref('200m')
const memLimit         = ref('256Mi')

// ── Tag input helpers ─────────────────────────────────────────────────
type TagKey = 'ns' | 'pvc' | 'excl'
const inputValMap: Record<TagKey, Ref<string>> = {
  ns:   nsInputVal,
  pvc:  pvcInputVal,
  excl: exclInputVal,
}
const inputElMap: Record<TagKey, Ref<HTMLInputElement | null>> = {
  ns:   nsInputEl,
  pvc:  pvcInputEl,
  excl: exclInputEl,
}

function addTag(list: string[], key: TagKey) {
  const valRef = inputValMap[key]
  const raw = valRef.value.trim().replace(/,$/, '')
  if (raw && !list.includes(raw)) list.push(raw)
  valRef.value = ''
}

function focusTagInput(key: TagKey) {
  nextTick(() => inputElMap[key].value?.focus())
}

// ── YAML generation ───────────────────────────────────────────────────
const highlightedYaml = computed(() => highlightYaml(generatedYaml.value))

const generatedYaml = computed(() => {
  const n = name.value || '<name>'
  const lines: string[] = []

  lines.push('apiVersion: pvcexplorer.io/v1alpha1')
  lines.push('kind: PVCExplorerScope')
  lines.push('metadata:')
  lines.push(`  name: ${n}`)
  lines.push('spec:')

  // namespaces
  const validPairs = labelPairs.value.filter(p => p.key)
  const hasNsNames = nsNames.value.length > 0
  const hasLabelSel = validPairs.length > 0
  lines.push('  namespaces:')
  if (hasNsNames) {
    lines.push('    names:')
    for (const ns of nsNames.value) { lines.push(`      - ${ns}`) }
  }
  if (hasLabelSel) {
    lines.push('    labelSelector:')
    lines.push('      matchLabels:')
    for (const p of validPairs) { lines.push(`        ${p.key}: "${p.value}"`) }
  }
  if (!hasNsNames && !hasLabelSel) {
    lines.push('    names: []')
  }

  // discovery
  lines.push('  discovery:')
  lines.push(`    mode: ${discoveryMode.value}`)
  if (discoveryMode.value === 'Explicit' && pvcNames.value.length) {
    lines.push('    pvcNames:')
    for (const p of pvcNames.value) { lines.push(`      - ${p}`) }
  }
  if (excludePVCs.value.length) {
    lines.push('    excludePVCs:')
    for (const p of excludePVCs.value) { lines.push(`      - "${p}"`) }
  }

  // deletion policy
  lines.push(`  deletionPolicy: ${deletionPolicy.value}`)

  // defaults
  lines.push('  defaults:')
  lines.push(`    mode: ${defaultMode.value}`)
  if (image.value) lines.push(`    image: "${image.value}"`)
  lines.push(`    forceRW: ${forceRW.value}`)
  lines.push('    scaling:')
  lines.push(`      idleTimeout: "${idleTimeout.value || '10m'}"`)
  lines.push(`      startupTimeout: "${startupTimeout.value || '60s'}"`)
  lines.push('    mountStrategy:')
  lines.push(`      allowNodeAffinity: ${allowNodeAffinity.value}`)
  lines.push('      fallbackOnConflict: Pending')
  if (enableResources.value) {
    lines.push('    resources:')
    if (cpuRequest.value || memRequest.value) {
      lines.push('      requests:')
      if (cpuRequest.value) lines.push(`        cpu: "${cpuRequest.value}"`)
      if (memRequest.value) lines.push(`        memory: "${memRequest.value}"`)
    }
    if (cpuLimit.value || memLimit.value) {
      lines.push('      limits:')
      if (cpuLimit.value) lines.push(`        cpu: "${cpuLimit.value}"`)
      if (memLimit.value) lines.push(`        memory: "${memLimit.value}"`)
    }
  }

  return lines.join('\n')
})

// ── Copy / Download ───────────────────────────────────────────────────
const copied = ref(false)

function onEscape(e: KeyboardEvent) { if (e.key === 'Escape') router.push('/scopes') }
onMounted(() => document.addEventListener('keydown', onEscape))
onUnmounted(() => document.removeEventListener('keydown', onEscape))

async function copyYaml() {
  await navigator.clipboard.writeText(generatedYaml.value)
  copied.value = true
  setTimeout(() => { copied.value = false }, 2000)
}

function downloadYaml() {
  const blob = new Blob([generatedYaml.value], { type: 'text/yaml' })
  const url  = URL.createObjectURL(blob)
  const a    = document.createElement('a')
  a.href     = url
  a.download = `${name.value || 'scope'}.yaml`
  a.click()
  URL.revokeObjectURL(url)
}
</script>

<style scoped>
/* ── Page ── */
.cs-page {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  width: 100%;
}

/* ── Header ── */
.cs-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1rem;
}
.cs-title {
  font-size: 1.35rem;
  font-weight: 700;
  color: var(--text-color);
  margin: 0 0 0.25rem;
}
.cs-info-banner {
  display: flex;
  align-items: flex-start;
  gap: 0.55rem;
  padding: 0.6rem 0.9rem;
  border-radius: 8px;
  background: color-mix(in srgb, var(--p-primary-400) 12%, transparent);
  border: 1px solid color-mix(in srgb, var(--p-primary-400) 35%, transparent);
  font-size: 0.84rem;
  color: var(--text-color);
  line-height: 1.5;
}
.cs-info-icon {
  color: var(--p-primary-500);
  font-size: 0.95rem;
  flex-shrink: 0;
  margin-top: 0.1rem;
}
.cs-info-banner code {
  font-family: 'JetBrains Mono', monospace;
  background: color-mix(in srgb, var(--p-primary-400) 18%, transparent);
  padding: 1px 5px;
  border-radius: 3px;
  font-size: 0.8rem;
}
.cs-cancel-btn {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  padding: 0.4rem 1rem;
  border: 1px solid var(--surface-border);
  background: transparent;
  color: var(--text-color-secondary);
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.875rem;
  white-space: nowrap;
  transition: background 0.12s, color 0.12s;
}
.cs-cancel-btn:hover { background: var(--surface-hover); color: var(--text-color); }

/* ── Body (2-col) ── */
.cs-body {
  display: flex;
  gap: 1.5rem;
  align-items: flex-start;
}
.cs-form-col {
  flex: 3;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
.cs-yaml-col {
  width: 42%;
  min-width: 380px;
  flex-shrink: 0;
  position: sticky;
  top: 1rem;
}

/* ── Section ── */
.cs-section {
  background: var(--surface-card);
  border: 1px solid var(--surface-border);
  border-radius: 8px;
  padding: 1rem 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
.cs-section-title {
  font-size: 0.75rem;
  font-weight: 700;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  color: var(--text-color-secondary);
  margin: 0;
}
.cs-section-title--toggle {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  border: none;
  background: transparent;
  cursor: pointer;
  padding: 0;
  color: var(--text-color-secondary);
  font-size: 0.75rem;
  font-weight: 700;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  width: 100%;
  text-align: left;
}
.cs-section-badge {
  font-size: 0.7rem;
  font-weight: 400;
  text-transform: none;
  letter-spacing: 0;
  margin-left: auto;
}

/* ── Field ── */
.cs-field {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}
fieldset.cs-field {
  border: none;
  padding: 0;
  margin: 0;
}
.cs-field--row {
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
}
.cs-field-row-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.75rem;
}
.cs-label {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--text-color);
}
.cs-label--inline { margin: 0; }
.cs-required { color: var(--p-red-400); }
.cs-optional { font-weight: 400; color: var(--text-color-secondary); font-size: 0.75rem; }
.cs-field-hint {
  font-size: 0.75rem;
  color: var(--text-color-secondary);
}
.cs-field-hint code {
  font-family: 'JetBrains Mono', monospace;
  background: var(--surface-hover);
  padding: 0px 4px;
  border-radius: 3px;
}
.cs-field-error {
  font-size: 0.75rem;
  color: var(--p-red-400);
}

/* ── Inputs ── */
.cs-input {
  height: 2rem;
  padding: 0 0.6rem;
  border: 1px solid var(--surface-border);
  background: var(--surface-ground);
  color: var(--text-color);
  border-radius: 6px;
  font-size: 0.875rem;
  font-family: inherit;
  outline: none;
  transition: border-color 0.12s;
}
.cs-input:focus { border-color: var(--p-primary-500); }
.cs-input--error { border-color: var(--p-red-400) !important; }
.cs-input--sm { height: 1.75rem; font-size: 0.8125rem; }
.cs-checkbox {
  width: 1rem;
  height: 1rem;
  accent-color: var(--p-primary-500);
  cursor: pointer;
}

/* ── Tag input ── */
.cs-tag-input {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 4px;
  min-height: 2rem;
  padding: 3px 6px;
  border: 1px solid var(--surface-border);
  background: var(--surface-ground);
  border-radius: 6px;
  cursor: text;
  transition: border-color 0.12s;
}
.cs-tag-input:focus-within { border-color: var(--p-primary-500); }
.cs-tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 1px 8px 1px 8px;
  background: rgba(79,142,247,0.15);
  color: #0ea5e9;
  border-radius: 4px;
  font-size: 0.8125rem;
  white-space: nowrap;
}
.cs-tag--warn {
  background: rgba(245, 158, 11, 0.15);
  color: #f59e0b;
}
.cs-tag-remove {
  background: transparent;
  border: none;
  color: inherit;
  cursor: pointer;
  font-size: 0.9rem;
  line-height: 1;
  padding: 0;
  opacity: 0.7;
}
.cs-tag-remove:hover { opacity: 1; }
.cs-tag-inner-input {
  border: none;
  background: transparent;
  color: var(--text-color);
  font-size: 0.875rem;
  font-family: inherit;
  outline: none;
  min-width: 80px;
  flex: 1;
  padding: 0;
}

/* ── Chip group ── */
.cs-chip-group {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}
.cs-chip {
  padding: 3px 14px;
  border: 1px solid var(--surface-border);
  background: transparent;
  color: var(--text-color-secondary);
  border-radius: 20px;
  cursor: pointer;
  font-size: 0.8125rem;
  font-family: inherit;
  transition: background 0.1s, color 0.1s, border-color 0.1s;
}
.cs-chip:hover { background: var(--surface-hover); color: var(--text-color); }
.cs-chip--active {
  background: var(--p-primary-500);
  color: #fff;
  border-color: var(--p-primary-500);
}

/* ── Label pairs ── */
.cs-label-pairs {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}
.cs-label-pair {
  display: flex;
  align-items: center;
  gap: 0.4rem;
}
.cs-label-pair-eq {
  font-size: 0.9rem;
  color: var(--text-color-secondary);
}
.cs-add-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.3rem;
  padding: 3px 10px;
  border: 1px dashed var(--surface-border);
  background: transparent;
  color: var(--text-color-secondary);
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.8rem;
  font-family: inherit;
  align-self: flex-start;
  transition: border-color 0.12s, color 0.12s;
}
.cs-add-btn:hover { border-color: var(--p-primary-500); color: var(--p-primary-500); }
.cs-icon-btn {
  background: transparent;
  border: none;
  cursor: pointer;
  font-size: 1rem;
  padding: 0 4px;
  color: var(--text-color-secondary);
}
.cs-icon-btn--danger:hover { color: var(--p-red-400); }

/* ── Resources grid ── */
.cs-resources-grid {
  display: grid;
  grid-template-columns: 6rem 1fr 1fr;
  gap: 0.4rem 0.6rem;
  align-items: center;
}
.cs-resources-head {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--text-color-secondary);
  text-align: center;
}
.cs-resources-row-label {
  font-size: 0.8rem;
  color: var(--text-color-secondary);
}

/* ── YAML panel ── */
.cs-yaml-card {
  background: var(--surface-card);
  border: 1px solid var(--surface-border);
  border-radius: 8px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}
.cs-yaml-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.6rem 1rem;
  border-bottom: 1px solid var(--surface-border);
  background: var(--surface-card);
}
.cs-yaml-title {
  font-size: 0.75rem;
  font-weight: 700;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  color: var(--text-color-secondary);
}
.cs-yaml-actions { display: flex; gap: 0.5rem; }
.cs-action-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  padding: 3px 12px;
  border: 1px solid var(--surface-border);
  background: transparent;
  color: var(--text-color-secondary);
  border-radius: 5px;
  cursor: pointer;
  font-size: 0.8rem;
  font-family: inherit;
  transition: background 0.1s, color 0.1s, border-color 0.1s;
}
.cs-action-btn:hover { background: var(--surface-hover); color: var(--text-color); }
.cs-action-btn--primary {
  background: var(--p-primary-500);
  color: #fff;
  border-color: var(--p-primary-500);
}
.cs-action-btn--primary:hover:not(:disabled) {
  background: var(--p-primary-400);
  border-color: var(--p-primary-400);
  color: #fff;
}
.cs-action-btn--primary:disabled { opacity: 0.45; cursor: not-allowed; }
.cs-action-btn--done { color: #22c55e; border-color: #22c55e; }
.cs-yaml-pre {
  margin: 0;
  padding: 1rem;
  background: var(--surface-card);
  font-family: 'JetBrains Mono', 'Fira Mono', monospace;
  font-size: 0.78rem;
  line-height: 1.6;
  color: var(--text-color);
  overflow-x: auto;
  max-height: calc(100vh - 220px);
  overflow-y: auto;
  white-space: pre;
}
.cs-yaml-footer {
  padding: 0.5rem 1rem;
  font-size: 0.75rem;
  color: var(--text-color-secondary);
  border-top: 1px solid var(--surface-border);
  background: var(--surface-card);
}
.cs-yaml-footer code {
  font-family: 'JetBrains Mono', monospace;
  font-size: 0.75rem;
}

/* ── Misc ── */
.cs-not-admin {
  color: var(--p-red-400);
  font-size: 1.1rem;
  text-align: center;
  margin-top: 2.5rem;
}

/* ── Responsive ── */
@media (max-width: 860px) {
  .cs-body { flex-direction: column; }
  .cs-yaml-col { width: 100%; position: static; }
  .cs-yaml-pre { max-height: 320px; }
}
</style>
