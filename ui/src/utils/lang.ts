const LANG_MAP: Record<string, string> = {
  yaml: 'yaml', yml: 'yaml',
  json: 'json',
  js: 'javascript', ts: 'typescript', jsx: 'javascript', tsx: 'typescript',
  py: 'python',
  sh: 'shell', bash: 'shell', env: 'shell',
  go: 'go',
  sql: 'sql',
  md: 'markdown', markdown: 'markdown',
  conf: 'ini', ini: 'ini', toml: 'ini',
  xml: 'xml', html: 'html', css: 'css', scss: 'scss',
  txt: 'plaintext', log: 'plaintext', csv: 'plaintext',
}

const READONLY_EXTS = new Set([
  'parquet', 'gz', 'tar', 'zip', 'bz2', 'tgz', '7z', 'rar',
  'bin', 'exe', 'so', 'dylib',
  'png', 'jpg', 'jpeg', 'gif', 'svg', 'webp', 'ico',
  'pdf',
  'mp3', 'wav', 'ogg', 'flac',
  'mp4', 'mov', 'avi', 'mkv',
])

export function detectLang(filename: string): string {
  const ext = filename.split('.').pop()?.toLowerCase() ?? ''
  return LANG_MAP[ext] ?? 'plaintext'
}

export function isEditable(filename: string): boolean {
  const ext = filename.split('.').pop()?.toLowerCase() ?? ''
  return !READONLY_EXTS.has(ext)
}
