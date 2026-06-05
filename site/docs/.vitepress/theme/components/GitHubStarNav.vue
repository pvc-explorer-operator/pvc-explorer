<template>
  <a
    :href="repoUrl"
    class="gh-star-nav"
    target="_blank"
    rel="noopener noreferrer"
    :aria-label="stars ? `${stars} stars on GitHub` : 'GitHub'"
  >
    <svg class="gh-icon" viewBox="0 0 24 24" fill="currentColor">
      <path d="M12 2C6.477 2 2 6.477 2 12c0 4.42 2.865 8.166 6.839 9.489.5.092.682-.217.682-.482 0-.237-.008-.866-.013-1.7-2.782.604-3.369-1.34-3.369-1.34-.454-1.156-1.11-1.462-1.11-1.462-.908-.62.069-.608.069-.608 1.003.07 1.531 1.03 1.531 1.03.892 1.529 2.341 1.087 2.91.831.092-.646.35-1.086.636-1.336-2.22-.253-4.555-1.11-4.555-4.943 0-1.091.39-1.984 1.029-2.683-.103-.253-.446-1.27.098-2.647 0 0 .84-.269 2.75 1.025A9.578 9.578 0 0 1 12 6.836a9.59 9.59 0 0 1 2.504.337c1.909-1.294 2.747-1.025 2.747-1.025.546 1.377.202 2.394.1 2.647.64.699 1.028 1.592 1.028 2.683 0 3.842-2.339 4.687-4.566 4.935.359.309.678.919.678 1.852 0 1.336-.012 2.415-.012 2.743 0 .267.18.578.688.48C19.138 20.163 22 16.418 22 12c0-5.523-4.477-10-10-10z"/>
    </svg>
    <svg class="gh-star-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
      <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2" fill="currentColor"/>
    </svg>
    <span v-if="ready" class="gh-star-count">{{ formatted }}</span>
  </a>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from "vue"

const repoUrl = "https://github.com/pvc-explorer-operator/pvc-explorer"
const stars = ref(0)
const ready = ref(false)

const formatted = computed(() => {
  const n = stars.value
  if (n >= 1000) return `${(n / 1000).toFixed(1)}k`
  return String(n)
})

onMounted(async () => {
  const cached = localStorage.getItem("gh-stars-pvc-explorer")
  if (cached) {
    stars.value = parseInt(cached, 10)
    ready.value = true
    return
  }
  try {
    const res = await fetch("https://api.github.com/repos/pvc-explorer-operator/pvc-explorer")
    if (!res.ok) throw new Error("API error")
    const data = await res.json()
    stars.value = data.stargazers_count ?? 0
    localStorage.setItem("gh-stars-pvc-explorer", String(stars.value))
  } catch {
    // fallback: just show the icon
  }
  ready.value = true
})
</script>

<style scoped>
.gh-star-nav {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--vp-c-text-1);
  text-decoration: none;
  padding: 0 12px;
  line-height: var(--vp-nav-height);
  transition: color 0.2s;
  white-space: nowrap;
}

.gh-star-nav:hover {
  color: var(--vp-c-brand-1);
  text-decoration: none;
}

.gh-icon {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
}

.gh-star-icon {
  width: 14px;
  height: 14px;
  flex-shrink: 0;
  color: var(--vp-c-text-2);
}

.gh-star-count {
  font-variant-numeric: tabular-nums;
  color: var(--vp-c-text-2);
  font-size: 0.8rem;
}

.gh-star-nav:hover .gh-star-count {
  color: var(--vp-c-brand-2);
}

@media (max-width: 768px) {
  .gh-star-nav {
    display: none;
  }
}
</style>
