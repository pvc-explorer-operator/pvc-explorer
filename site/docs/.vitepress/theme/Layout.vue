<template>
  <Layout>
    <template #layout-top>
      <ZoomInit />
    </template>
    <template #nav-bar-content-after>
      <GitHubStarNav />
    </template>
    <template #not-found>
      <NotFound />
    </template>
  </Layout>
</template>

<script setup lang="ts">
import DefaultTheme from "vitepress/theme"
import GitHubStarNav from "./components/GitHubStarNav.vue"
import ZoomInit from "./components/ZoomInit.vue"
import NotFound from "./components/NotFound.vue"

const { Layout } = DefaultTheme
</script>
