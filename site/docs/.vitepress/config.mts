import { defineConfig } from 'vitepress'
import { withMermaid } from 'vitepress-plugin-mermaid'
import { fileURLToPath } from 'node:url'
import path from 'node:path'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
// Absolute path to the site/ package root (two levels up from .vitepress/)
const siteRoot = path.resolve(__dirname, '../..')

const isCI = process.env.GITHUB_ACTIONS === 'true'
const base = process.env.DOCS_BASE_PATH || (isCI ? `/${process.env.GITHUB_REPOSITORY?.split('/')[1] ?? ''}/` : '/')

export default withMermaid(defineConfig({
  title: 'PVC Explorer',
  description: 'Documentation for the PVC Explorer operator',
  base,
  lastUpdated: true,
  cleanUrls: true,
  vite: {
    server: {
      fs: {
        allow: [siteRoot],
      },
    },
  },
  head: [
    ['link', { rel: 'preconnect', href: 'https://fonts.googleapis.com' }],
    ['link', { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: '' }],
    ['link', {
      rel: 'preload',
      as: 'style',
      href: 'https://fonts.googleapis.com/css2?family=Lato:wght@300;400;700&display=swap',
      onload: "this.onload=null;this.rel='stylesheet'",
    }],
    ['noscript', {}, '<link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Lato:wght@300;400;700&display=swap">'],
    ['meta', { name: 'referrer', content: 'strict-origin-when-cross-origin' }],
  ],
  markdown: {
    externalLinks: {
      target: '_blank',
      rel: 'noopener noreferrer',
    },
  },
  themeConfig: {
    logo: '/images/branding/logo-icon-darkbg.svg',
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Overview', link: '/overview' },
      { text: 'Install', link: '/install' },
      { text: 'API', link: '/api/' },
      { text: 'GitHub', items: [
        { text: 'Operator', link: 'https://github.com/pvc-explorer-operator/pvc-explorer' },
        { text: 'Agent', link: 'https://github.com/pvc-explorer-operator/pvc-explorer-agent' },
      ]},
    ],
    sidebar: [
      { text: 'Getting Started', link: '/guide/getting-started' },
      { text: 'Overview', link: '/overview' },
      { text: 'Install', link: '/install' },
      { text: 'Run Local', link: '/guide/local-run' },
      {
        text: 'User Guide',
        collapsed: false,
        items: [
          { text: 'Core Concepts', link: '/user-guide/core-concepts' },
          { text: 'Examples', link: '/user-guide/examples' },
          {
            text: 'How-to',
            collapsed: true,
            items: [
              { text: 'Login', link: '/user-guide/how-to/login' },
              { text: 'Connect to PVC', link: '/user-guide/how-to/connect-to-pvc' },
              { text: 'Create Scope', link: '/user-guide/how-to/create-scope' },
            ],
          },
        ],
      },
      {
        text: 'Core',
        collapsed: false,
        items: [
          { text: 'Architecture', link: '/architecture' },
          {
            text: 'API',
            collapsed: true,
            items: [
              { text: 'Overview', link: '/api/' },
              {
                text: 'CRD',
                collapsed: true,
                items: [
                  { text: 'PVCExplorer', link: '/api/crds/pvcexplorer' },
                  { text: 'PVCExplorerScope', link: '/api/crds/pvcexplorerscope' },
                ],
              },
              { text: 'REST', link: '/api/rest' },
              { text: 'WebSocket', link: '/api/websocket' },
            ],
          },
          { text: 'Operations', link: '/operations' },
        ],
      },
      {
        text: 'Operator Guide',
        collapsed: false,
        items: [
          { text: 'Scope Examples', link: '/operator-guide/scope-examples' },
          { text: 'Security', link: '/operator-guide/security' },
        ],
      },
      {
        text: 'Contributor Guide',
        collapsed: true,
        items: [
          { text: 'Kubebuilder', link: '/contributor-guide/kubebuilder' },
          { text: 'Vue Interface', link: '/contributor-guide/vue-interface' },
          { text: 'Signing Commits', link: '/contributor-guide/signing-commits' },
        ],
      },
      {
        text: 'Project',
        collapsed: false,
        items: [
          { text: 'Development', link: '/development' },
          { text: 'Releases', link: '/releases' },
          { text: 'Compliance and Security', link: '/compliance-security' },
          { text: 'ADRs', link: '/adrs' },
        ],
      },
    ],
    search: {
      provider: 'local',
    },
    footer: {
      message: 'Released under Apache-2.0',
      copyright: 'Copyright PVC Explorer contributors',
    },
  },
}))
