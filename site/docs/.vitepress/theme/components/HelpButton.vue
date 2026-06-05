<template>
  <button
    class="help-btn"
    aria-label="Keyboard shortcuts"
    title="Keyboard shortcuts"
    @click="toggle"
  >
    <svg class="help-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <circle cx="12" cy="12" r="10" />
      <path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3" />
      <line x1="12" y1="17" x2="12.01" y2="17" />
    </svg>
  </button>
</template>

<script setup lang="ts">
function toggle() {
  document.dispatchEvent(new CustomEvent("opencode:toggle-keyboard-help"))
}
</script>

<style scoped>
.help-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: transparent;
  padding: 0 8px;
  line-height: var(--vp-nav-height);
  cursor: pointer;
  color: var(--vp-c-text-2);
  transition: color 0.2s;
}

.help-btn:hover {
  color: var(--vp-c-brand-1);
}

.help-icon {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
}

@media (max-width: 768px) {
  .help-btn {
    display: none;
  }
}
</style>
