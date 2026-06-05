import { defineConfig } from "vitepress"
import { withMermaid } from "vitepress-plugin-mermaid"

const base = process.env.VITEPRESS_BASE ?? "/"

export default withMermaid(
  defineConfig({
    title: "PVC Explorer Operator",
    description:
      "Suite of open-source Kubernetes tools for browsing, managing, and exploring PersistentVolumeClaims — safely, scoped, and on demand.",

    base,
    cleanUrls: true,
    lastUpdated: true,
    ignoreDeadLinks: true,

    head: [
      [
        "link",
        {
          rel: "preconnect",
          href: "https://fonts.googleapis.com",
        },
      ],
      [
        "link",
        {
          rel: "preconnect",
          href: "https://fonts.gstatic.com",
          crossorigin: "",
        },
      ],
      [
        "link",
        {
          rel: "stylesheet",
          href: "https://fonts.googleapis.com/css2?family=IBM+Plex+Sans:wght@300;400;500;600;700&family=JetBrains+Mono:wght@400;500;600;700&display=swap",
        },
      ],
      [
        "link",
        {
          rel: "icon",
          type: "image/svg+xml",
          href: `${base}images/branding/logo-icon.svg`,
        },
      ],
      [
        "meta",
        { property: "og:title", content: "PVC Explorer Operator" },
      ],
      [
        "meta",
        {
          property: "og:description",
          content:
            "Suite of open-source Kubernetes tools for browsing, managing, and exploring PersistentVolumeClaims — safely, scoped, and on demand.",
        },
      ],
    ],

    themeConfig: {
      logo: {
        light: "/images/branding/logo-icon.svg",
        dark: "/images/branding/logo-icon-darkbg.svg",
        alt: "PVC Explorer Operator",
      },

      nav: [
        { text: "Overview", link: "/overview" },
        { text: "Operator", link: "/guide/getting-started" },
        { text: "Agent", link: "/pvc-explorer-agent/overview" },
        { text: "CLI", link: "/kubectl-pvc-explorer/overview" },
      ],

      sidebar: [
        {
          text: "Welcome",
          collapsed: false,
          items: [
            { text: "Home", link: "/" },
            { text: "Overview", link: "/overview" },
            { text: "Architecture", link: "/architecture" },
          ],
        },
        {
          text: "pvc-explorer (Operator)",
          collapsed: false,
          items: [
            { text: "Overview", link: "/guide/getting-started" },
            { text: "Install", link: "/install" },
            { text: "Run Local", link: "/guide/local-run" },
            { text: "Core Concepts", link: "/user-guide/core-concepts" },
            { text: "Examples", link: "/user-guide/examples" },
            {
              text: "How-To Guides",
              collapsed: true,
              items: [
                { text: "Login", link: "/user-guide/how-to/login" },
                {
                  text: "Connect to PVC",
                  link: "/user-guide/how-to/connect-to-pvc",
                },
                {
                  text: "Create Scope",
                  link: "/user-guide/how-to/create-scope",
                },
              ],
            },
            {
              text: "Operations & Security",
              collapsed: true,
              items: [
                { text: "Operations", link: "/operations" },
                { text: "Security", link: "/operator-guide/security" },
                { text: "Scope Examples", link: "/operator-guide/scope-examples" },
                { text: "Compliance", link: "/compliance-security" },
              ],
            },
            {
              text: "API Reference",
              collapsed: true,
              items: [
                { text: "REST API", link: "/api/rest" },
                { text: "WebSocket", link: "/api/websocket" },
                { text: "CRDs", link: "/api/crds" },
                { text: "PVCExplorer CRD", link: "/api/crds/pvcexplorer" },
                {
                  text: "PVCExplorerScope CRD",
                  link: "/api/crds/pvcexplorerscope",
                },
              ],
            },
            {
              text: "Development",
              collapsed: true,
              items: [
                { text: "Development Guide", link: "/development" },
                { text: "UI Reference", link: "/ui/index" },
                {
                  text: "UI Components",
                  link: "/ui/components",
                },
                { text: "UI Display", link: "/ui/display" },
                { text: "UI Flows", link: "/ui/flows" },
                { text: "Keybindings", link: "/ui/keybindings" },
              ],
            },
          ],
        },
        {
          text: "pvc-explorer-agent",
          collapsed: false,
          items: [
            {
              text: "Overview",
              link: "/pvc-explorer-agent/overview",
            },
          ],
        },
        {
          text: "kubectl-pvc-explorer (CLI)",
          collapsed: false,
          items: [
            {
              text: "Overview",
              link: "/kubectl-pvc-explorer/overview",
            },
          ],
        },
        {
          text: "Reference",
          collapsed: true,
          items: [
            { text: "Releases", link: "/releases" },
            { text: "ADRs", link: "/adrs" },
            {
              text: "Contributor Guide",
              collapsed: true,
              items: [
                {
                  text: "Kubebuilder",
                  link: "/contributor-guide/kubebuilder",
                },
                {
                  text: "Signing Commits",
                  link: "/contributor-guide/signing-commits",
                },
                {
                  text: "Vue Interface",
                  link: "/contributor-guide/vue-interface",
                },
              ],
            },
          ],
        },
      ],

      search: {
        provider: "local",
      },

      footer: {
        message:
          "Apache 2.0 licensed. Built with the Kubebuilder framework.",
        copyright:
          "PVC Explorer Operator — open source Kubernetes tooling.",
      },


    },
  })
)
