<template>
  <div class="arch-wrapper">
    <!-- Zone: User Space -->
    <div class="arch-zone zone-user">
      <div class="zone-label">User Space</div>
      <div class="zone-row">
        <div class="arch-card card-sm">
          <div class="card-icon-wrap">
            <svg class="card-icon-svg" viewBox="0 0 40 40" fill="none">
              <rect x="4" y="10" width="32" height="22" rx="3" stroke="currentColor" stroke-width="1.8"/>
              <path d="M4 16h32" stroke="currentColor" stroke-width="1.8"/>
              <circle cx="8" cy="13" r="1" fill="currentColor"/>
              <circle cx="11" cy="13" r="1" fill="currentColor"/>
              <circle cx="14" cy="13" r="1" fill="currentColor"/>
              <path d="M14 22l-4 4 4 4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M22 22l4 4-4 4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </div>
          <div class="card-text">
            <div class="card-label">kubectl pvc</div>
            <div class="card-sub-label">CLI plugin · kubeconfig</div>
          </div>
        </div>

        <div class="arch-card card-sm">
          <div class="card-icon-wrap">
            <svg class="card-icon-svg" viewBox="0 0 40 40" fill="none">
              <rect x="5" y="8" width="30" height="24" rx="3" stroke="currentColor" stroke-width="1.8"/>
              <path d="M5 14h30" stroke="currentColor" stroke-width="1.8"/>
              <circle cx="9" cy="11" r="1.2" fill="currentColor"/>
              <circle cx="12.5" cy="11" r="1.2" fill="currentColor"/>
              <circle cx="16" cy="11" r="1.2" fill="currentColor"/>
              <path d="M20 22h8M20 26h6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
              <path d="M5 24a4 4 0 016 0" stroke="currentColor" stroke-width="1.3" stroke-linecap="round"/>
              <circle cx="8" cy="28" r="2" stroke="currentColor" stroke-width="1.3"/>
            </svg>
          </div>
          <div class="card-text">
            <div class="card-label">Web Dashboard</div>
            <div class="card-sub-label">Vue UI</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Cross-zone connector -->
    <div class="arch-cross">
      <div class="cross-track">
        <div class="cross-dot"></div>
      </div>
      <span class="cross-label">Auth</span>
    </div>

    <!-- Zone: Control Plane -->
    <div class="arch-zone zone-control">
      <div class="zone-label zone-label-light">Control Plane</div>
      <div class="zone-row zone-row-single">
        <div class="arch-card card-lg">
          <div class="card-lg-header">
            <span class="card-badge badge-operator">Operator</span>
            <div class="card-icon-wrap">
              <svg class="card-icon-svg icon-operator" viewBox="0 0 44 44" fill="none">
                <circle cx="22" cy="22" r="16" stroke="currentColor" stroke-width="1.8"/>
                <circle cx="22" cy="22" r="8" stroke="currentColor" stroke-width="1.8"/>
                <path d="M22 2v6M22 36v6M2 22h6M36 22h6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                <path d="M8.5 8.5l4.2 4.2M31.3 31.3l4.2 4.2M31.3 8.5l-4.2 4.2M8.5 31.3l4.2-4.2" stroke="currentColor" stroke-width="1.3" stroke-linecap="round"/>
                <circle cx="22" cy="22" r="3" fill="currentColor" opacity="0.3"/>
              </svg>
            </div>
            <div class="card-label card-label-lg">pvc-explorer</div>
          </div>
          <ul class="card-list">
            <li>Validates PVCs and consumers</li>
            <li>Creates and scales agent pods</li>
            <li>Enforces RBAC and mount guardrails</li>
          </ul>
        </div>
      </div>
    </div>

    <!-- Cross-zone connector -->
    <div class="arch-cross">
      <div class="cross-track">
        <div class="cross-dot"></div>
      </div>
      <span class="cross-label">creates &amp; manages</span>
    </div>

    <!-- Zone: Data Plane -->
    <div class="arch-zone zone-data">
      <div class="zone-label zone-label-light">Data Plane</div>
      <div class="zone-row">
        <div class="arch-card card-lg">
          <div class="card-lg-header">
            <span class="card-badge badge-agent">Pod</span>
            <div class="card-icon-wrap">
              <svg class="card-icon-svg icon-agent" viewBox="0 0 44 44" fill="none">
                <rect x="6" y="8" width="32" height="28" rx="4" stroke="currentColor" stroke-width="1.8"/>
                <path d="M6 16h32" stroke="currentColor" stroke-width="1.8"/>
                <path d="M16 8v28M28 8v28" stroke="currentColor" stroke-width="1.3" opacity="0.4"/>
                <rect x="10" y="20" width="4" height="3" rx="1" fill="currentColor" opacity="0.4"/>
                <rect x="10" y="26" width="4" height="3" rx="1" fill="currentColor" opacity="0.4"/>
                <rect x="10" y="32" width="4" height="2" rx="1" fill="currentColor" opacity="0.4"/>
                <path d="M30 24h4M30 28h4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
              </svg>
            </div>
            <div class="card-label card-label-lg">pvc-explorer-agent</div>
          </div>
          <ul class="card-list">
            <li>Mounts the target PVC</li>
            <li>Serves REST file-browser API</li>
            <li>Detects conflicts → read-only</li>
          </ul>
        </div>
        <div class="arch-connector conn-lg">
          <svg viewBox="0 0 36 16" fill="none" class="conn-arrow">
            <path d="M2 8h28M28 4l4 4-4 4" stroke="currentColor" stroke-width="1.3" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          <span class="conn-label-lg">mounts &amp; serves</span>
        </div>
        <div class="arch-card card-sm card-pvc">
          <div class="card-icon-wrap">
            <svg class="card-icon-svg icon-pvc" viewBox="0 0 44 44" fill="none">
              <rect x="7" y="5" width="30" height="34" rx="4" stroke="currentColor" stroke-width="1.8"/>
              <path d="M7 13h30" stroke="currentColor" stroke-width="1.8"/>
              <rect x="11" y="18" width="8" height="4" rx="1" stroke="currentColor" stroke-width="1.3"/>
              <rect x="11" y="25" width="8" height="4" rx="1" stroke="currentColor" stroke-width="1.3"/>
              <rect x="11" y="32" width="8" height="4" rx="1" stroke="currentColor" stroke-width="1.3"/>
              <path d="M25 20h6M25 27h6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
              <path d="M37 5l-7 7M37 5l-7-7" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" opacity="0.4"/>
            </svg>
          </div>
          <div class="card-text">
            <div class="card-label">PVC</div>
            <div class="card-sub-label">Kubernetes storage</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
</script>

<style scoped>
.arch-wrapper {
  display: flex;
  flex-direction: column;
  align-items: stretch;
  gap: 0;
  padding: 1rem 0;
  max-width: 820px;
  margin: 0 auto;
}

/* ---- Zones ---- */
.arch-zone {
  border-radius: 12px;
  padding: 1rem 1.25rem 1rem;
  position: relative;
}

.zone-user {
  background: linear-gradient(135deg, rgba(0, 207, 255, 0.22), rgba(0, 150, 200, 0.08));
  border: 1px solid rgba(0, 207, 255, 0.35);
}

:root.dark .zone-user {
  background: linear-gradient(135deg, rgba(0, 207, 255, 0.28), rgba(0, 150, 200, 0.12));
  border-color: rgba(0, 207, 255, 0.45);
}

.zone-control {
  background: linear-gradient(135deg, rgba(42, 95, 196, 0.22), rgba(26, 58, 122, 0.08));
  border: 1px solid rgba(42, 95, 196, 0.35);
}

:root.dark .zone-control {
  background: linear-gradient(135deg, rgba(42, 95, 196, 0.28), rgba(26, 58, 122, 0.12));
  border-color: rgba(42, 95, 196, 0.45);
}

.zone-data {
  background: linear-gradient(135deg, rgba(202, 138, 4, 0.22), rgba(180, 120, 0, 0.08));
  border: 1px solid rgba(202, 138, 4, 0.35);
}

:root.dark .zone-data {
  background: linear-gradient(135deg, rgba(202, 138, 4, 0.28), rgba(180, 120, 0, 0.12));
  border-color: rgba(202, 138, 4, 0.45);
}

.zone-label {
  font-size: 0.65rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--vp-c-text-3);
  margin-bottom: 0.8rem;
}

.zone-label-light {
  color: var(--vp-c-text-2);
}

.zone-user .zone-label { color: #0099CC; }
.zone-control .zone-label { color: var(--vp-c-brand-1); }
.zone-data .zone-label { color: #B8860B; }

:root.dark .zone-user .zone-label { color: #00CFFF; }
:root.dark .zone-control .zone-label { color: #5A8BFF; }
:root.dark .zone-data .zone-label { color: #FACC15; }

/* ---- Zone rows ---- */
.zone-row {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.6rem;
}

.zone-row-single {
  justify-content: center;
}

/* ---- Cards ---- */
.arch-card {
  background: var(--vp-c-bg);
  border: 1px solid var(--vp-c-border);
  border-radius: 12px;
  padding: 1.25rem 1.25rem 1rem;
  text-align: center;
  transition: transform 0.2s ease, box-shadow 0.2s ease, border-color 0.2s;
}

.arch-card:hover {
  transform: translateY(-3px);
}

.card-sm {
  flex: 1 1 200px;
  max-width: 220px;
  display: flex;
  align-items: center;
  gap: 0.7rem;
  text-align: left;
  padding: 0.8rem 1rem;
}

.card-lg {
  flex: 0 1 320px;
  max-width: 480px;
  text-align: left;
  padding: 1.1rem 1.25rem 1.1rem;
}

.card-pvc {
  border-color: rgba(202, 138, 4, 0.25);
}

.card-pvc:hover {
  border-color: rgba(202, 138, 4, 0.5);
  box-shadow: 0 4px 16px rgba(202, 138, 4, 0.08);
}

/* ---- SVG Icons ---- */
.card-icon-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
}

.card-sm .card-icon-wrap {
  flex-shrink: 0;
}

.card-lg .card-icon-wrap {
  flex-shrink: 0;
}

.card-lg .card-icon-svg {
  width: 30px;
  height: 30px;
}

.card-icon-svg {
  width: 32px;
  height: 32px;
  color: var(--vp-c-text-3);
}

.card-sm .card-icon-svg {
  width: 28px;
  height: 28px;
}

/* Large card header: icon + title in a row */
.card-lg-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.4rem;
}

.card-lg-header .card-label {
  font-size: 0.9rem;
}

.icon-operator {
  color: var(--vp-c-brand-1);
}

:root.dark .icon-operator {
  color: #60a5fa;
}

.icon-agent {
  color: #0099CC;
}

:root.dark .icon-agent {
  color: #00CFFF;
}

.icon-pvc {
  color: #ca8a04;
}

:root.dark .icon-pvc {
  color: #facc15;
}

/* Text group (for small cards: label + sub-label stacked) */
.card-text {
  display: flex;
  flex-direction: column;
}

/* Labels */
.card-label {
  font-weight: 700;
  font-size: 0.85rem;
  color: var(--vp-c-text-1);
  letter-spacing: -0.01em;
  line-height: 1.35;
}

.card-label-lg {
  font-size: 0.95rem;
}

.card-sub-label {
  font-size: 0.7rem;
  color: var(--vp-c-text-3);
  margin-top: 0.15rem;
  line-height: 1.3;
}

/* Badge */
.card-badge {
  display: inline-block;
  font-size: 0.58rem;
  font-weight: 700;
  padding: 0.15rem 0.5rem;
  border-radius: 999px;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  white-space: nowrap;
}

.badge-operator {
  background: var(--vp-c-brand-soft);
  color: var(--vp-c-brand-1);
}

.badge-agent {
  background: rgba(0, 207, 255, 0.1);
  color: #0099CC;
}

:root.dark .badge-agent {
  background: rgba(0, 207, 255, 0.15);
  color: #00CFFF;
}



/* List */
.card-list {
  list-style: none;
  padding: 0;
  margin: 0.2rem 0 0;
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
}

.card-list li {
  font-size: 0.7rem;
  line-height: 1.4;
  color: var(--vp-c-text-2);
  padding-left: 1rem;
  position: relative;
}

.card-list li::before {
  content: "";
  position: absolute;
  left: 0;
  top: 0.42em;
  width: 4px;
  height: 4px;
  border-radius: 50%;
  background: var(--vp-c-brand-1);
}

.zone-data .card-list li::before {
  background: #CA8A04;
}

:root.dark .zone-data .card-list li::before {
  background: #FACC15;
}

.zone-user .card-list li::before {
  background: #0099CC;
}

:root.dark .zone-user .card-list li::before {
  background: #00CFFF;
}

.zone-control .card-list li::before {
  background: var(--vp-c-brand-1);
}

/* ---- Connectors (within zone) ---- */
.arch-connector {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.15rem;
  flex: 0 0 48px;
}

.conn-sm .conn-arrow {
  width: 30px;
  height: 14px;
  color: var(--vp-c-brand-1);
}

.conn-lg .conn-arrow {
  width: 32px;
  height: 14px;
  color: var(--vp-c-brand-1);
}

.conn-label-sm,
.conn-label-lg {
  font-size: 0.58rem;
  color: var(--vp-c-text-3);
  white-space: nowrap;
  font-weight: 500;
  text-align: center;
}

/* ---- Cross-zone connector ---- */
.arch-cross {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.25rem 0 0.25rem 2rem;
  position: relative;
  height: 2rem;
}

.cross-track {
  width: 2px;
  height: 100%;
  background: var(--vp-c-border);
  position: relative;
  overflow: hidden;
}

.cross-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #00CFFF;
  position: absolute;
  left: -2px;
  top: 30%;
  opacity: 0.6;
}

@media (prefers-reduced-motion: no-preference) {
  .cross-dot {
    animation: flow 1.8s ease-in-out infinite;
  }
}

:root.dark .cross-dot {
  background: #00CFFF;
}

@keyframes flow {
  0% { top: -6px; opacity: 0; }
  30% { opacity: 1; }
  70% { opacity: 1; }
  100% { top: 100%; opacity: 0; }
}

.cross-label {
  font-size: 0.6rem;
  font-weight: 600;
  color: var(--vp-c-text-3);
  text-transform: uppercase;
  letter-spacing: 0.04em;
  white-space: nowrap;
}

/* ---- Hover effects per card ---- */
.card-sm:hover {
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.06);
  border-color: var(--vp-c-brand-1);
}

:root.dark .card-sm:hover {
  box-shadow: 0 6px 24px rgba(0, 0, 0, 0.25);
}

.card-lg:hover {
  box-shadow: 0 8px 28px rgba(0, 0, 0, 0.07);
  border-color: var(--vp-c-brand-1);
}

:root.dark .card-lg:hover {
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.25);
}

/* ---- Responsive ---- */
@media (max-width: 780px) {
  .arch-wrapper {
    padding: 0.5rem 0;
  }

  .arch-zone {
    padding: 1rem 0.75rem 0.75rem;
  }

  .zone-row {
    flex-direction: column;
    gap: 0.5rem;
  }

  .card-sm,
  .card-lg {
    max-width: 100%;
    width: 100%;
    flex: none;
  }

  .arch-connector {
    flex-direction: row;
    flex: 0 0 auto;
    height: 28px;
  }

  .conn-sm .conn-arrow,
  .conn-lg .conn-arrow {
    transform: rotate(90deg);
    width: 14px;
    height: 24px;
  }

  .conn-label-sm,
  .conn-label-lg {
    order: -1;
  }

  .arch-cross {
    padding: 0.25rem 0;
    height: 2rem;
    gap: 0.5rem;
  }
}
</style>
