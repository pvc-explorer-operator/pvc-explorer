import { defineConfig, type Plugin } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'
import { mockWsPlugin } from './dev-mock/wsPlugin'

export default defineConfig(async ({ command }) => {
  const plugins: Plugin[] = [vue()]
  if (command === 'serve') {
    plugins.push(mockWsPlugin())
  }
  return {
    plugins,
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
      },
    },
  }
})
