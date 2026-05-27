<template>
  <div v-if="authStore.isAdmin">
    <div class="card" style="max-width: 28rem">
      <form @submit.prevent="onSubmit" class="flex flex-column gap-4">
        <div class="flex flex-column gap-2">
          <label for="name" class="font-medium text-sm">Name</label>
          <InputText id="name" v-model="name" type="text" required maxlength="63" autocomplete="off" />
        </div>
        <div class="flex flex-column gap-2">
          <label for="namespace" class="font-medium text-sm">Namespace</label>
          <Dropdown id="namespace" v-model="namespace_" :options="namespaces" placeholder="Select namespace" />
        </div>
        <div class="flex flex-column gap-2">
          <label for="pvc" class="font-medium text-sm">PVC</label>
          <Dropdown id="pvc" v-model="pvc" :options="pvcs" placeholder="Select PVC" />
        </div>
        <div class="flex flex-column gap-2">
          <label for="idleTimeout" class="font-medium text-sm">Idle Timeout (minutes)</label>
          <InputNumber id="idleTimeout" v-model="idleTimeout" :min="1" :max="1440" />
        </div>
        <div class="flex align-items-center gap-2">
          <Checkbox id="forceRW" v-model="forceReadWrite" binary />
          <label for="forceRW" class="font-medium text-sm">Force ReadWrite</label>
        </div>
        <div class="flex gap-2">
          <Button type="submit" severity="primary" icon="pi pi-check" label="Create" rounded :disabled="loading" :loading="loading" />
          <Button type="button" severity="secondary" icon="pi pi-times" label="Cancel" rounded @click="router.push('/explorers')" />
        </div>
        <div v-if="error" class="error-msg">{{ error }}</div>
      </form>
    </div>
  </div>
  <div v-else class="not-admin">Admin access required.</div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/authStore'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Dropdown from 'primevue/dropdown'
import Checkbox from 'primevue/checkbox'
import Button from 'primevue/button'

const authStore = useAuthStore()
const router = useRouter()

const name = ref('')
const namespace_ = ref('')
const pvc = ref('')
const idleTimeout = ref<number | null>(null)
const forceReadWrite = ref(false)
const loading = ref(false)
const error = ref('')

const namespaces = ref<string[]>([])
const pvcs = ref<string[]>([])

onMounted(fetchNamespaces)
watch(namespace_, fetchPVCs)

async function fetchNamespaces() {
  const res = await fetch('/api/v1/namespaces')
  if (res.ok) {
    const data: { name: string }[] = await res.json()
    namespaces.value = data.map(n => n.name)
  }
}

async function fetchPVCs() {
  pvcs.value = []
  if (!namespace_.value) return
  const res = await fetch(`/api/v1/namespaces/${encodeURIComponent(namespace_.value)}/pvcs`)
  if (res.ok) {
    const data: { name: string }[] = await res.json()
    pvcs.value = data.map(p => p.name)
  }
}

async function onSubmit() {
  error.value = ''
  loading.value = true
  const body = {
    apiVersion: 'pvcexplorer.io/v1',
    kind: 'Explorer',
    metadata: { name: name.value, namespace: namespace_.value },
    spec: {
      pvcName: pvc.value,
      idleTimeout: idleTimeout.value,
      forceReadWrite: forceReadWrite.value
    }
  }
  const res = await fetch('/api/v1/explorers', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body)
  })
  if (res.ok) {
    loading.value = false
    router.push('/explorers')
  } else {
    error.value = 'Failed to create agent.'
    loading.value = false
  }
}
</script>

<style scoped>
.error-msg {
  color: var(--p-red-500);
  margin-top: 0.5rem;
  font-size: 0.95rem;
}
.not-admin {
  color: var(--p-red-500);
  font-size: 1.1rem;
  text-align: center;
  margin-top: 2.5rem;
}
</style>
