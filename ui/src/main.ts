import { createApp } from 'vue'
import { createPinia } from 'pinia'
import PrimeVue from 'primevue/config'

// Configure Monaco Editor to use local workers (no CDN required)
import { loader } from '@guolao/vue-monaco-editor'
import * as monaco from 'monaco-editor'
import editorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker'
import jsonWorker from 'monaco-editor/esm/vs/language/json/json.worker?worker'
import cssWorker from 'monaco-editor/esm/vs/language/css/css.worker?worker'
import htmlWorker from 'monaco-editor/esm/vs/language/html/html.worker?worker'
import tsWorker from 'monaco-editor/esm/vs/language/typescript/ts.worker?worker'

self.MonacoEnvironment = {
  getWorker(_: string, label: string) {
    if (label === 'json') return new jsonWorker()
    if (label === 'css' || label === 'scss' || label === 'less') return new cssWorker()
    if (label === 'html' || label === 'handlebars' || label === 'razor') return new htmlWorker()
    if (label === 'typescript' || label === 'javascript') return new tsWorker()
    return new editorWorker()
  },
}
loader.config({ monaco })
import ConfirmationService from 'primevue/confirmationservice'
import Lara from '@primevue/themes/lara'
import { definePreset } from '@primevue/themes'

const SakaiSky = definePreset(Lara, {
  semantic: {
    primary: {
      50:  '{sky.50}',
      100: '{sky.100}',
      200: '{sky.200}',
      300: '{sky.300}',
      400: '{sky.400}',
      500: '{sky.500}',
      600: '{sky.600}',
      700: '{sky.700}',
      800: '{sky.800}',
      900: '{sky.900}',
      950: '{sky.950}',
    },
    colorScheme: {
      light: {
        text: {
          color: '{surface.900}',
          mutedColor: '{surface.600}',
          hoverMutedColor: '{surface.700}',
        },
        content: {
          color: '{surface.900}',
        },
      },
    },
  },
})
import 'primeicons/primeicons.css'
import './layout/layout.css'
import './theme/variables.css'
import './style.css'
import App from './App.vue'
import router from './router'
import { loadTheme } from './theme/themeLoader'
import { useAuthStore } from './stores/authStore'

loadTheme()

const app = createApp(App)
const pinia = createPinia()
app.use(pinia)
app.use(router)
app.use(PrimeVue, { theme: { preset: SakaiSky, options: { darkModeSelector: '.app-dark' } } })
app.use(ConfirmationService)

const auth = useAuthStore()
auth.init().finally(() => {
  app.mount('#app')
})
