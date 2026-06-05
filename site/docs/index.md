---
layout: home

title: PVC Explorer Operator

tagline: Browse, manage, and explore PVCs — safely, scoped, and on demand

hero:
  text: Open-source Kubernetes tools for PVC access
  tagline: A suite of tools for platform teams to inspect, debug, and explore PersistentVolumeClaims without disrupting running workloads.
  image:
    light: /images/branding/logo-icon.svg
    dark: /images/branding/logo-icon-darkbg.svg
    alt: PVC Explorer icon
  actions:
    - theme: brand
      text: Get Started
      link: /overview
    - theme: alt
      text: Live Demo
      link: https://pvc-explorer-operator.github.io/demo/

features:
  - title: pvc-explorer (Operator)
    details: Kubernetes controller that manages ephemeral agent pods for browsing PVCs. Scale-to-zero, safe read-only fallback, namespace-scoped access.
    link: /guide/getting-started
  - title: pvc-explorer-agent
    details: Lightweight HTTP file-browser that mounts a PVC and exposes its contents over a REST API. Embedded Vue UI, read-only conflict detection.
    link: /pvc-explorer-agent/overview
  - title: kubectl-pvc-explorer
    details: kubectl plugin for exploring PVCs from the terminal. List, cat, cp, exec, and mount commands — no UI needed.
    link: /kubectl-pvc-explorer/overview
  - title: Local-first Workflow
    details: Docs, UI mock mode, and kind-based testing all run locally with short feedback loops.

---

<div class="home-content">


<div class="section-relationships" id="project-relationships">

<div class="section-label" role="heading" aria-level="2">
  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <path d="M18 8a6 6 0 0 1-6 6"/>
    <path d="M6 16a6 6 0 0 1 6-6"/>
    <circle cx="12" cy="4" r="2"/>
    <circle cx="6" cy="20" r="2"/>
    <circle cx="18" cy="20" r="2"/>
    <path d="m12 8 3 3-3 3-3-3 3-3"/>
  </svg>
  System Topology
</div>

<ArchitectureDiagram />

<div class="pr-cards">

<div class="pr-card">
  <div class="pr-card-icon">
    <svg viewBox="0 0 44 44" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
      <circle cx="22" cy="22" r="16" stroke="currentColor"/>
      <circle cx="22" cy="22" r="8" stroke="currentColor"/>
      <path d="M22 2v6M22 36v6M2 22h6M36 22h6" stroke="currentColor" stroke-width="1.5"/>
      <path d="M8.5 8.5l4.2 4.2M31.3 31.3l4.2 4.2M31.3 8.5l-4.2 4.2M8.5 31.3l4.2-4.2" stroke="currentColor" stroke-width="1.3"/>
      <circle cx="22" cy="22" r="3" fill="currentColor" opacity="0.3" stroke="none"/>
    </svg>
  </div>
  <div class="pr-card-body">
    <div class="pr-card-title">pvc-explorer</div>
    <div class="pr-card-desc">Manages everything — install this first</div>
  </div>
</div>

<div class="pr-card">
  <div class="pr-card-icon pr-card-icon--agent">
    <svg viewBox="0 0 44 44" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
      <rect x="6" y="8" width="32" height="28" rx="4" stroke="currentColor"/>
      <path d="M6 16h32" stroke="currentColor"/>
      <path d="M16 8v28M28 8v28" stroke="currentColor" stroke-width="1.3" opacity="0.4"/>
      <rect x="10" y="20" width="4" height="3" rx="1" fill="currentColor" opacity="0.4" stroke="none"/>
      <rect x="10" y="26" width="4" height="3" rx="1" fill="currentColor" opacity="0.4" stroke="none"/>
      <rect x="10" y="32" width="4" height="2" rx="1" fill="currentColor" opacity="0.4" stroke="none"/>
      <path d="M30 24h4M30 28h4" stroke="currentColor" stroke-width="1.5"/>
    </svg>
  </div>
  <div class="pr-card-body">
    <div class="pr-card-title">pvc-explorer-agent</div>
    <div class="pr-card-desc">Container image deployed by the operator (no manual setup needed)</div>
  </div>
</div>

<div class="pr-card">
  <div class="pr-card-icon pr-card-icon--cli">
    <svg viewBox="0 0 40 40" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
      <rect x="4" y="10" width="32" height="22" rx="3" stroke="currentColor"/>
      <path d="M4 16h32" stroke="currentColor"/>
      <circle cx="8" cy="13" r="1" fill="currentColor" stroke="none"/>
      <circle cx="11" cy="13" r="1" fill="currentColor" stroke="none"/>
      <circle cx="14" cy="13" r="1" fill="currentColor" stroke="none"/>
      <path d="M14 22l-4 4 4 4" stroke="currentColor" stroke-width="1.5"/>
      <path d="M22 22l4 4-4 4" stroke="currentColor" stroke-width="1.5"/>
    </svg>
  </div>
  <div class="pr-card-body">
    <div class="pr-card-title">kubectl-pvc-explorer</div>
    <div class="pr-card-desc">Optional — use it if you prefer the terminal over the web UI</div>
  </div>
</div>

</div>

</div>

<div class="section-involved">

<div class="section-label" role="heading" aria-level="2">
  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/>
  </svg>
  Get Involved
</div>

<div class="involved-grid">

<a href="https://github.com/pvc-explorer-operator/pvc-explorer" class="involved-card" target="_blank" rel="noopener noreferrer">
  <div class="involved-card-icon">
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <path d="M12 2C6.477 2 2 6.477 2 12c0 4.42 2.865 8.166 6.839 9.489.5.092.682-.217.682-.482 0-.237-.008-.866-.013-1.7-2.782.604-3.369-1.34-3.369-1.34-.454-1.156-1.11-1.462-1.11-1.462-.908-.62.069-.608.069-.608 1.003.07 1.531 1.03 1.531 1.03.892 1.529 2.341 1.087 2.91.831.092-.646.35-1.086.636-1.336-2.22-.253-4.555-1.11-4.555-4.943 0-1.091.39-1.984 1.029-2.683-.103-.253-.446-1.27.098-2.647 0 0 .84-.269 2.75 1.025A9.578 9.578 0 0 1 12 6.836a9.59 9.59 0 0 1 2.504.337c1.909-1.294 2.747-1.025 2.747-1.025.546 1.377.202 2.394.1 2.647.64.699 1.028 1.592 1.028 2.683 0 3.842-2.339 4.687-4.566 4.935.359.309.678.919.678 1.852 0 1.336-.012 2.415-.012 2.743 0 .267.18.578.688.48C19.138 20.163 22 16.418 22 12c0-5.523-4.477-10-10-10z"/>
    </svg>
  </div>
  <div class="involved-card-body">
    <div class="involved-card-title">GitHub</div>
    <div class="involved-card-desc">Star the repo, open issues, submit PRs</div>
  </div>
  <svg class="involved-card-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <path d="M7 17l9.2-9.2M17 17V7H7"/>
  </svg>
</a>

<a href="https://pvc-explorer-operator.github.io/demo/" class="involved-card" target="_blank" rel="noopener noreferrer">
  <div class="involved-card-icon">
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <circle cx="12" cy="12" r="10"/>
      <polygon points="10 8 16 12 10 16 10 8" fill="currentColor" stroke="none"/>
    </svg>
  </div>
  <div class="involved-card-body">
    <div class="involved-card-title">Live Demo</div>
    <div class="involved-card-desc">Try the UI without installing anything</div>
  </div>
  <svg class="involved-card-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <path d="M7 17l9.2-9.2M17 17V7H7"/>
  </svg>
</a>

<a href="/contributor-guide/kubebuilder" class="involved-card">
  <div class="involved-card-icon">
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"/>
      <path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"/>
      <path d="M8 7h8M8 11h6"/>
    </svg>
  </div>
  <div class="involved-card-body">
    <div class="involved-card-title">Contributing</div>
    <div class="involved-card-desc">Read the contributor guides</div>
  </div>
  <svg class="involved-card-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <path d="M7 17l9.2-9.2M17 17V7H7"/>
  </svg>
</a>

</div>

</div>

</div>
