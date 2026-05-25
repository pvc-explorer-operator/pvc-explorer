<template>
  <div class="ab-page">

    <!-- Hero banner -->
    <div class="ab-hero">
      <div class="ab-hero-bg">
        <div class="ab-hero-grid"></div>
        <div class="ab-hero-glow"></div>
      </div>
      <div class="ab-hero-content">
        <div class="ab-hero-text">
          <h1 class="ab-project-name">PVC Explorer</h1>
          <p class="ab-tagline">On-demand browsing of PersistentVolumeClaims via ephemeral agent pods — scaled to zero, woken on demand.</p>
          <div class="ab-version-badge" v-if="version">v{{ version }}</div>
        </div>
      </div>
    </div>

    <!-- Links -->
    <div class="ab-cards">

      <a class="ab-card" href="https://github.com/pvc-explorer-operator/pvc-explorer" target="_blank" rel="noopener">
        <div class="ab-card-icon ab-card-icon--gh">
          <svg viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 .297c-6.63 0-12 5.373-12 12 0 5.303 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61-.546-1.385-1.335-1.755-1.335-1.755-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.417-1.305.76-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 22.092 24 17.592 24 12.297c0-6.627-5.373-12-12-12"/>
          </svg>
        </div>
        <div class="ab-card-body">
          <div class="ab-card-title">GitHub Repository</div>
          <div class="ab-card-desc">Source code, issues, pull requests and releases.</div>
          <div class="ab-card-url">github.com/pvc-explorer-operator/pvc-explorer</div>
        </div>
        <i class="pi pi-arrow-right ab-card-arrow"></i>
      </a>

      <a class="ab-card" href="https://pvc-explorer-operator.github.io/pvc-explorer" target="_blank" rel="noopener">
        <div class="ab-card-icon ab-card-icon--docs">
          <i class="pi pi-book"></i>
        </div>
        <div class="ab-card-body">
          <div class="ab-card-title">Documentation</div>
          <div class="ab-card-desc">Getting started, architecture reference and API guides.</div>
          <div class="ab-card-url">pvc-explorer-operator.github.io/pvc-explorer</div>
        </div>
        <i class="pi pi-arrow-right ab-card-arrow"></i>
      </a>

      <a class="ab-card" href="https://pvc-explorer-operator.github.io/pvc-explorer/getting-started" target="_blank" rel="noopener">
        <div class="ab-card-icon ab-card-icon--start">
          <i class="pi pi-play-circle"></i>
        </div>
        <div class="ab-card-body">
          <div class="ab-card-title">Getting Started</div>
          <div class="ab-card-desc">Install the operator, create your first Scope and browse PVC contents.</div>
          <div class="ab-card-url">pvc-explorer-operator.github.io/pvc-explorer/getting-started</div>
        </div>
        <i class="pi pi-arrow-right ab-card-arrow"></i>
      </a>

      <a class="ab-card" href="https://github.com/pvc-explorer-operator/pvc-explorer/releases" target="_blank" rel="noopener">
        <div class="ab-card-icon ab-card-icon--rel">
          <i class="pi pi-tag"></i>
        </div>
        <div class="ab-card-body">
          <div class="ab-card-title">Releases</div>
          <div class="ab-card-desc">Changelog, version history and container image tags.</div>
          <div class="ab-card-url">github.com/pvc-explorer-operator/pvc-explorer/releases</div>
        </div>
        <i class="pi pi-arrow-right ab-card-arrow"></i>
      </a>

      <a class="ab-card" href="https://github.com/pvc-explorer-operator/pvc-explorer/issues" target="_blank" rel="noopener">
        <div class="ab-card-icon ab-card-icon--issue">
          <i class="pi pi-exclamation-circle"></i>
        </div>
        <div class="ab-card-body">
          <div class="ab-card-title">Issues &amp; Feedback</div>
          <div class="ab-card-desc">Report bugs, request features or ask questions.</div>
          <div class="ab-card-url">github.com/pvc-explorer-operator/pvc-explorer/issues</div>
        </div>
        <i class="pi pi-arrow-right ab-card-arrow"></i>
      </a>

    </div>

    <!-- Meta strip -->
    <div class="ab-meta">
      <div class="ab-meta-item">
        <span class="ab-meta-label">Version</span>
        <span class="ab-meta-value">{{ version || '—' }}</span>
      </div>
      <div class="ab-meta-sep"></div>
      <div class="ab-meta-item">
        <span class="ab-meta-label">License</span>
        <a class="ab-meta-link" href="https://github.com/pvc-explorer-operator/pvc-explorer/blob/main/LICENSE" target="_blank" rel="noopener">Apache 2.0</a>
      </div>

    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const version = ref('')

onMounted(async () => {
  try {
    const res = await fetch('/api/version')
    if (res.ok) version.value = (await res.text()).trim()
  } catch { /* ignore */ }
})
</script>

<style scoped>
.ab-page {
  display: flex;
  flex-direction: column;
  gap: 2rem;
  width: 100%;
  padding: 0 1.5rem;
}

/* ── Hero banner ── */
.ab-hero {
  position: relative;
  border-radius: 14px;
  overflow: hidden;
  min-height: 160px;
}
.ab-hero-bg {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, #0d1117 0%, #1e2130 60%, #162032 100%);
}
.ab-hero-grid {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(42,51,72,0.35) 1px, transparent 1px),
    linear-gradient(90deg, rgba(42,51,72,0.35) 1px, transparent 1px);
  background-size: 24px 24px;
}
.ab-hero-glow {
  position: absolute;
  top: -40px;
  right: 60px;
  width: 200px;
  height: 200px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(0,207,255,0.12) 0%, transparent 70%);
  pointer-events: none;
}
.ab-hero-content {
  position: relative;
  display: flex;
  align-items: center;
  gap: 1.5rem;
  padding: 2rem 2rem;
  width: 100%;
}
.ab-hero-text {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
  flex: 1;
}
.ab-project-name {
  margin: 0;
  font-size: 1.8rem;
  font-weight: 800;
  color: #fff;
  letter-spacing: -0.02em;
  font-family: 'Courier New', Courier, monospace;
}
.ab-tagline {
  margin: 0;
  font-size: 0.9rem;
  color: #94a3b8;
  max-width: 520px;
  line-height: 1.5;
}
.ab-version-badge {
  display: inline-block;
  width: fit-content;
  margin-top: 0.25rem;
  padding: 2px 10px;
  background: rgba(0,207,255,0.12);
  color: #00cfff;
  border: 1px solid rgba(0,207,255,0.2);
  border-radius: 20px;
  font-size: 0.78rem;
  font-weight: 600;
  font-family: 'JetBrains Mono', 'Courier New', monospace;
}

/* ── Cards ── */
.ab-cards {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1rem;
}
.ab-card {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.25rem 1.25rem;
  background: var(--surface-card);
  border: 1px solid var(--surface-border);
  border-radius: 10px;
  text-decoration: none;
  color: inherit;
  box-shadow: 0 1px 3px rgba(0,0,0,0.04);
  transition: border-color 0.15s, box-shadow 0.15s, transform 0.15s;
}
.ab-card:hover {
  border-color: var(--p-primary-400);
  box-shadow: 0 4px 16px rgba(0,0,0,0.1);
  transform: translateY(-2px);
}
.ab-card-icon {
  width: 42px;
  height: 42px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.2rem;
  flex-shrink: 0;
}
.ab-card-icon svg { width: 22px; height: 22px; }
.ab-card-icon--gh    { background: rgba(36,41,47,0.12);   color: var(--text-color); }
.ab-card-icon--docs  { background: rgba(79,142,247,0.12); color: var(--p-primary-500); }
.ab-card-icon--start { background: rgba(34,197,94,0.12);  color: #22c55e; }
.ab-card-icon--rel   { background: rgba(245,158,11,0.12); color: #f59e0b; }
.ab-card-icon--issue { background: rgba(239,68,68,0.12);  color: #ef4444; }
.ab-card-body {
  flex: 1;
  min-width: 0;
}
.ab-card-title {
  font-size: 0.9rem;
  font-weight: 700;
  color: var(--text-color);
  margin-bottom: 0.15rem;
}
.ab-card-desc {
  font-size: 0.78rem;
  color: var(--text-color-secondary);
  line-height: 1.4;
  margin-bottom: 0.25rem;
}
.ab-card-url {
  font-size: 0.72rem;
  font-family: 'JetBrains Mono', monospace;
  color: var(--p-primary-500);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.ab-card-arrow {
  color: var(--text-color-secondary);
  font-size: 0.75rem;
  flex-shrink: 0;
  opacity: 0;
  transition: opacity 0.15s, transform 0.15s;
}
.ab-card:hover .ab-card-arrow {
  opacity: 1;
  transform: translateX(2px);
}

/* ── Meta strip ── */
.ab-meta {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0.75rem 1.25rem;
  background: var(--surface-card);
  border: 1px solid var(--surface-border);
  border-radius: 8px;
  flex-wrap: wrap;
}
.ab-meta-sep {
  width: 1px;
  height: 1.2rem;
  background: var(--surface-border);
  flex-shrink: 0;
}
.ab-meta-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}
.ab-meta-label {
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--text-color-secondary);
}
.ab-meta-value {
  font-size: 0.85rem;
  color: var(--text-color);
}
.ab-meta-value.mono {
  font-family: 'JetBrains Mono', monospace;
  font-size: 0.78rem;
}
.ab-meta-link {
  font-size: 0.85rem;
  color: var(--p-primary-500);
  text-decoration: none;
}
.ab-meta-link:hover { text-decoration: underline; }

@media (max-width: 700px) {
  .ab-hero-content { padding: 1.5rem 1.25rem; }
  .ab-project-name { font-size: 1.4rem; }
  .ab-cards { grid-template-columns: 1fr; }
  .ab-meta-sep { display: none; }
}
</style>

<style>
/* Light mode overrides (unscoped so html ancestor selector works) */
html:not(.app-dark) .ab-hero-bg {
  background: linear-gradient(135deg, #f1f5f9 0%, #e2e8f0 60%, #f8fafc 100%);
}
html:not(.app-dark) .ab-hero-grid {
  background-image:
    linear-gradient(rgba(148,163,184,0.18) 1px, transparent 1px),
    linear-gradient(90deg, rgba(148,163,184,0.18) 1px, transparent 1px);
}
html:not(.app-dark) .ab-hero-glow {
  background: radial-gradient(circle, rgba(37,99,235,0.08) 0%, transparent 70%);
}
html:not(.app-dark) .ab-project-name {
  color: #0f172a;
}
html:not(.app-dark) .ab-tagline {
  color: #475569;
}
html:not(.app-dark) .ab-version-badge {
  background: rgba(37,99,235,0.1);
  color: #2563eb;
  border-color: rgba(37,99,235,0.2);
}
</style>
