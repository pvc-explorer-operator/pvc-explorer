<template>
  <div class="not-found-wrap">
    <svg class="pvc-error-icon" viewBox="0 0 44 44" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
      <rect x="7" y="5" width="30" height="34" rx="4" stroke="currentColor"/>
      <path d="M7 13h30" stroke="currentColor"/>
      <rect x="11" y="18" width="8" height="4" rx="1" stroke="currentColor" stroke-width="1.3"/>
      <rect x="11" y="25" width="8" height="4" rx="1" stroke="currentColor" stroke-width="1.3"/>
      <rect x="11" y="32" width="8" height="4" rx="1" stroke="currentColor" stroke-width="1.3"/>
      <circle cx="28" cy="27" r="6" stroke="#EF4444" stroke-width="1.8" opacity="0.7"/>
      <path d="M26 25l4 4M30 25l-4 4" stroke="#EF4444" stroke-width="1.8" opacity="0.7"/>
    </svg>

    <h1 class="nf-title">404</h1>
    <p class="nf-subtitle">PVC Not Found</p>

    <p class="nf-quote">
      But if you don't change your direction, and if you keep looking, you may end up where you are heading.
    </p>

    <a href="/" class="nf-home">
      <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M20 12H4M10 18l-6-6 6-6"/>
      </svg>
      Take me home
    </a>
  </div>
</template>

<style scoped>
.not-found-wrap {
  text-align: center;
  padding: 4rem 1.5rem 6rem;
}

.pvc-error-icon {
  display: block;
  width: 56px;
  height: 56px;
  margin: 0 auto 1rem;
  color: var(--vp-c-text-3);
  animation: nf-float 3s ease-in-out infinite;
}

@keyframes nf-float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-6px); }
}

@media (prefers-reduced-motion: reduce) {
  .pvc-error-icon {
    animation: none;
  }
}

.nf-title {
  font-size: 4rem;
  font-weight: 800;
  letter-spacing: -0.04em;
  margin: 0;
  line-height: 1;
  color: var(--vp-c-text-1);
}

.nf-subtitle {
  font-size: 1.1rem;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--vp-c-text-2);
  margin: 0.5rem 0 2rem;
}

.nf-quote {
  max-width: 480px;
  margin: 0 auto 2.5rem;
  color: var(--vp-c-text-2);
  font-style: italic;
  line-height: 1.7;
}

.nf-home {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  font-weight: 600;
  color: var(--vp-c-brand-1);
  text-decoration: none;
  transition: color 0.2s;
}

.nf-home:hover {
  color: var(--vp-c-brand-2);
  text-decoration: none;
}
</style>
