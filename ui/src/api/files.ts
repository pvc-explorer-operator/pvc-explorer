export interface FileEntry {
  name: string
  size: number
  isDir: boolean
  modTime: string
}

export interface AgentConfig {
  readonly: boolean
  forceRW: boolean
  pvc: string
  namespace: string
  pod: string
  cluster: string
  version: string
}

/**
 * Factory that creates all file-API functions bound to a specific agent proxy base URL.
 * The base URL is typically /api/v1/explorers/{ns}/{name}/proxy/api
 */
export function createFileApi(proxyBase: string) {
  function endpoint(rel: string) { return `${proxyBase}/${rel}` }

  async function fetchFiles(path: string): Promise<{ entries: FileEntry[] }> {
    const q = path ? `?path=${encodeURIComponent(path)}` : ''
    const res = await fetch(endpoint(`files${q}`))
    if (!res.ok) throw new Error(`List failed: ${res.statusText}`)
    const data = await res.json()
    return { entries: data.entries ?? [] }
  }

  async function fetchContent(path: string): Promise<string> {
    const res = await fetch(endpoint(`download?path=${encodeURIComponent(path)}`))
    if (!res.ok) throw new Error(`Read failed: ${res.statusText}`)
    return res.text()
  }

  async function saveFile(path: string, content: string): Promise<void> {
    const res = await fetch(endpoint(`edit?path=${encodeURIComponent(path)}`), {
      method: 'PUT',
      body: content,
    })
    if (!res.ok) throw new Error(`Save failed: ${res.statusText}`)
  }

  async function deleteFile(path: string): Promise<void> {
    const res = await fetch(endpoint(`files?path=${encodeURIComponent(path)}`), { method: 'DELETE' })
    if (!res.ok) throw new Error(`Delete failed: ${res.statusText}`)
  }

  async function renameFile(from: string, to: string): Promise<void> {
    const res = await fetch(endpoint('rename'), {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ from, to }),
    })
    if (!res.ok) throw new Error(`Rename failed: ${res.statusText}`)
  }

  async function uploadFiles(path: string, files: File[]): Promise<void> {
    for (const file of files) {
      const form = new FormData()
      form.append('file', file)
      const q = path ? `?path=${encodeURIComponent(path)}` : ''
      const res = await fetch(endpoint(`upload${q}`), { method: 'POST', body: form })
      if (!res.ok) throw new Error(`Upload failed: ${res.statusText}`)
    }
  }

  async function createFile(path: string): Promise<void> {
    const res = await fetch(endpoint(`edit?path=${encodeURIComponent(path)}`), {
      method: 'PUT',
      body: '',
    })
    if (!res.ok) throw new Error(`Create failed: ${res.statusText}`)
  }

  async function fetchConfig(): Promise<AgentConfig> {
    const res = await fetch(endpoint('config'))
    if (!res.ok) throw new Error(`Config failed: ${res.statusText}`)
    return res.json()
  }

  function downloadUrl(path: string): string {
    return endpoint(`download?path=${encodeURIComponent(path)}`)
  }

  return { fetchFiles, fetchContent, saveFile, deleteFile, renameFile, uploadFiles, createFile, fetchConfig, downloadUrl }
}
