<template>
  <Layout>
    <template #layout-top>
      <div class="alpha-banner">
        <svg class="alpha-banner-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M9 3H5a2 2 0 0 0-2 2v4m6-6h10a2 2 0 0 1 2 2v4M9 3v18m0 0h10a2 2 0 0 0 2-2V9M9 21H5a2 2 0 0 1-2-2V9m0 0h18" />
        </svg>
        <span class="alpha-banner-text">Alpha Preview</span>
        <span class="alpha-banner-dot">·</span>
        <span class="alpha-banner-sub">Documentation in progress</span>
      </div>
      <ZoomInit />
      <KeyboardNav />
    </template>
    <template #nav-bar-content-after>
      <GitHubStarNav />
      <HelpButton />
    </template>
    <template #not-found>
      <NotFound />
    </template>
  </Layout>
</template>

<script setup lang="ts">
import DefaultTheme from "vitepress/theme"
import GitHubStarNav from "./components/GitHubStarNav.vue"
import HelpButton from "./components/HelpButton.vue"
import ZoomInit from "./components/ZoomInit.vue"
import KeyboardNav from "./components/KeyboardNav.vue"
import NotFound from "./components/NotFound.vue"

const { Layout } = DefaultTheme
</script>
