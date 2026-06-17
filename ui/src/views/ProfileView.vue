<template>
  <div class="profile-page">
    <div class="profile-card">
      <div class="profile-header">
        <span class="profile-avatar">{{ avatarInitial }}</span>
        <div>
          <h1 class="profile-name">{{ auth.username || 'Unknown' }}</h1>
          <span class="profile-role" :class="roleClass">{{ auth.role }}</span>
        </div>
      </div>

      <div class="profile-fields">
        <div class="profile-field" v-if="auth.email">
          <label>Email</label>
          <span>{{ auth.email }}</span>
        </div>

        <div class="profile-field">
          <label>Role</label>
          <Tag :value="auth.role" :severity="roleSeverity" />
        </div>

        <div class="profile-field" v-if="auth.subject">
          <label>Subject</label>
          <code class="profile-subject">{{ auth.subject }}</code>
        </div>
      </div>

      <div class="profile-section" v-if="auth.groups.length > 0">
        <label>Groups</label>
        <div class="profile-groups">
          <Tag v-for="group in auth.groups" :key="group" :value="group" severity="info" />
        </div>
      </div>

      <div class="profile-section" v-else>
        <label>Groups</label>
        <span class="profile-empty">No group memberships</span>
      </div>

      <div class="profile-actions">
        <Button label="Back to Dashboard" icon="pi pi-arrow-left" severity="secondary" @click="router.push('/')" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/authStore'
import Tag from 'primevue/tag'
import Button from 'primevue/button'

const router = useRouter()
const auth = useAuthStore()

const avatarInitial = computed(() => auth.username?.charAt(0)?.toUpperCase() || '?')

const roleClass = computed(() => ({
  'role-admin': auth.role === 'admin',
  'role-user': auth.role === 'user',
  'role-viewer': auth.role === 'viewer',
}))

const roleSeverity = computed(() => {
  switch (auth.role) {
    case 'admin': return 'danger'
    case 'user': return 'warn'
    default: return 'secondary'
  }
})
</script>

<style scoped>
.profile-page {
  padding: 2rem;
  display: flex;
  justify-content: center;
}

.profile-card {
  width: 100%;
  max-width: 40rem;
  background: var(--surface-card);
  border: 1px solid var(--surface-border);
  border-radius: 12px;
  padding: 2rem;
}

.profile-header {
  display: flex;
  align-items: center;
  gap: 1.25rem;
  margin-bottom: 2rem;
  padding-bottom: 1.5rem;
  border-bottom: 1px solid var(--surface-border);
}

.profile-avatar {
  width: 4rem;
  height: 4rem;
  border-radius: 50%;
  background: var(--primary-color);
  color: var(--primary-color-text);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  font-weight: 700;
  flex-shrink: 0;
}

.profile-name {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--text-color);
}

.profile-role {
  font-size: 0.8rem;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.role-admin { color: var(--red-500); }
.role-user { color: var(--orange-500); }
.role-viewer { color: var(--text-color-secondary); }

.profile-fields {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.profile-field {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.profile-field label {
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--text-color-secondary);
}

.profile-field span {
  color: var(--text-color);
  font-size: 0.95rem;
}

.profile-subject {
  font-family: monospace;
  font-size: 0.85rem;
  background: var(--surface-hover);
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  color: var(--text-color-secondary);
  word-break: break-all;
}

.profile-section {
  margin-bottom: 1.5rem;
}

.profile-section > label {
  display: block;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--text-color-secondary);
  margin-bottom: 0.5rem;
}

.profile-groups {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.profile-empty {
  color: var(--text-color-secondary);
  font-size: 0.9rem;
  font-style: italic;
}

.profile-actions {
  padding-top: 1rem;
  border-top: 1px solid var(--surface-border);
}
</style>
