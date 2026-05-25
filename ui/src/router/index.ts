import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '../stores/authStore'

const routes: RouteRecordRaw[] = [
  { path: '/login', name: 'Login', component: () => import('../views/LoginView.vue') },
  { path: '/', name: 'Home', component: () => import('../views/HomeView.vue'), meta: { requiresAuth: true } },
  { path: '/scopes', name: 'ScopeList', component: () => import('../views/ScopeListView.vue'), meta: { requiresAuth: true } },
  { path: '/scopes/:name', name: 'ScopeDetail', component: () => import('../views/ScopeDetailView.vue'), meta: { requiresAuth: true } },
  { path: '/explorers/:ns/:name', name: 'AgentDetail', component: () => import('../views/AgentDetailView.vue'), meta: { requiresAuth: true } },
  { path: '/explorers/:ns/:name/files', name: 'FileBrowser', component: () => import('../views/FileBrowserView.vue'), meta: { requiresAuth: true } },
  { path: '/scopes/create', name: 'CreateScope', component: () => import('../views/CreateScopeView.vue'), meta: { requiresAuth: true, adminOnly: true } },
  { path: '/explorers/create', name: 'CreateAgent', component: () => import('../views/CreateAgentView.vue'), meta: { requiresAuth: true, adminOnly: true } },
  { path: '/settings', name: 'Settings', component: () => import('../views/SettingsView.vue'), meta: { requiresAuth: true, adminOnly: true } },
  { path: '/about',    name: 'About',    component: () => import('../views/AboutView.vue'),    meta: { requiresAuth: true } },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach(async (to, _from, next) => {
  const auth = useAuthStore()
  await auth.ready()
  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    next({ path: '/login' })
  } else if (to.meta.adminOnly && !auth.isAdmin) {
    next({ path: '/' })
  } else {
    next()
  }
})

export default router
