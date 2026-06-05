<template>
  <Layout>
    <template #layout-top>
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
