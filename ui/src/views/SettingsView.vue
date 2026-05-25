<template>
  <div v-if="authStore.isAdmin">
    <div class="card" style="max-width: 28rem">
      <div class="flex justify-between mb-4">
        <span class="text-muted-color">Username:</span>
        <span class="text-surface-900 font-medium">{{ authStore.username || '-' }}</span>
      </div>
      <div class="flex justify-between mb-4">
        <span class="text-muted-color">Role:</span>
        <span class="text-surface-900 font-medium">{{ authStore.isAdmin ? 'Admin' : 'User' }}</span>
      </div>
      <div class="flex justify-between mb-4">
        <span class="text-muted-color">Version:</span>
        <span class="text-surface-900 font-medium">{{ version || '-' }}</span>
      </div>
      <div class="flex justify-between mb-4">
        <span class="text-muted-color">About:</span>
        <a class="settings-link" href="https://github.com/pvc-explorer-operator/pvc-explorer" target="_blank" rel="noopener">GitHub</a>
      </div>
      <div class="flex justify-between mb-4">
        <span class="text-muted-color">Dark Mode:</span>
        <InputSwitch :modelValue="isDarkTheme" @update:modelValue="toggleDarkMode" />
      </div>
    </div>

    <Button severity="danger" label="Sign Out" icon="pi pi-sign-out" rounded @click="signOut" />
  </div>
  <Message v-else severity="error" :closable="false">Admin access required.</Message>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/authStore'
import Button from 'primevue/button'
import Message from 'primevue/message'
import InputSwitch from 'primevue/inputswitch'
import { useLayout } from '../layout/composables/layout'

const authStore = useAuthStore()
const router = useRouter()
const version = ref('')
const { isDarkTheme, toggleDarkMode } = useLayout()

onMounted(async () => {
  const res = await fetch('/api/version')
  if (res.ok) {
    version.value = await res.text()
  }
})

function signOut() {
  authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.settings-link {
  color: var(--primary-color);
  text-decoration: underline;
}
</style>
