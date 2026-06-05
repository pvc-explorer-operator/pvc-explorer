import DefaultTheme from "vitepress/theme"
import Layout from "./Layout.vue"
import ArchitectureDiagram from "./components/ArchitectureDiagram.vue"
import "./custom.css"

export default {
  extends: DefaultTheme,
  Layout,
  enhanceApp({ app }) {
    app.component("ArchitectureDiagram", ArchitectureDiagram)
  },
}
