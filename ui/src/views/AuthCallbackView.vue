<template>
  <main class="callback-bg">
    <div class="callback-card">
      <i class="pi pi-spin pi-spinner" style="font-size: 2rem"></i>
      <p>Signing you in...</p>
    </div>
  </main>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/authStore'

const router = useRouter()
const authStore = useAuthStore()

onMounted(async () => {
  await authStore.init()
  router.replace('/')
})
</script>

<style scoped>
.callback-bg {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
}
.callback-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  color: var(--text-color-secondary);
}
</style>
