<template>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from "vue"

let observer: MutationObserver | null = null

function tryZoom() {
  const svgs = document.querySelectorAll<HTMLElement>(".mermaid svg")
  if (!svgs.length) return

  svgs.forEach((svg) => {
    if (svg.dataset.zoomAttached) return
    svg.dataset.zoomAttached = "true"
    svg.style.cursor = "zoom-in"

    svg.addEventListener("click", function onClick() {
      const clone = this.cloneNode(true) as HTMLElement
      clone.style.maxWidth = "90vw"
      clone.style.maxHeight = "90vh"
      clone.style.width = "auto"
      clone.style.height = "auto"

      const overlay = document.createElement("div")
      overlay.style.cssText = `
        position: fixed; inset: 0; z-index: 9999;
        background: ${getComputedStyle(document.documentElement).getPropertyValue("--vp-c-bg").trim() || "#fff"};
        display: flex; align-items: center; justify-content: center;
        cursor: zoom-out; padding: 48px;
        opacity: 0; transition: opacity 0.2s;
      `
      overlay.appendChild(clone)
      document.body.appendChild(overlay)

      requestAnimationFrame(() => overlay.style.opacity = "1")

      overlay.addEventListener("click", () => {
        overlay.style.opacity = "0"
        setTimeout(() => overlay.remove(), 200)
      })

      document.addEventListener("keydown", function onKey(e) {
        if (e.key === "Escape") {
          overlay.style.opacity = "0"
          setTimeout(() => overlay.remove(), 200)
          document.removeEventListener("keydown", onKey)
        }
      })
    })
  })
}

onMounted(() => {
  // Watch for dynamically rendered mermaid SVGs
  observer = new MutationObserver(tryZoom)
  observer.observe(document.body, { childList: true, subtree: true })
  tryZoom()
})

onUnmounted(() => observer?.disconnect())
</script>
