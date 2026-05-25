<template>
  <div class="readonly-banner">Read-Only — this PVC has active consumers. Modifications are disabled.</div>
</template>

<style scoped>
.readonly-banner {
  background: var(--p-amber-100);
  color: #000;
  font-weight: 600;
  border-radius: 6px;
  padding: 0.55em 1rem;
  font-size: var(--fs-sm);
}
</style>
