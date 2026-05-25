import type { Plugin } from 'vite'
import type { IncomingMessage, ServerResponse } from 'node:http'
import { WebSocketServer } from 'ws'
import { mockExplorers, mockScopes } from './data'

// ---------------------------------------------------------------------------
// In-memory mock file system for the agent proxy API
// ---------------------------------------------------------------------------

interface MockFile { content: string; size: number; modTime: string; isDir: boolean }
const mockFs = new Map<string, MockFile>()

function seed(path: string, content: string) {
  mockFs.set(path, { content, size: content.length, modTime: new Date().toISOString(), isDir: false })
}
function seedDir(path: string) {
  mockFs.set(path, { content: '', size: 0, modTime: new Date().toISOString(), isDir: true })
}

seedDir('config')
seedDir('data')
seedDir('logs')
seed('README.md', `# Mock PVC\n\nThis is a **dev-mock** filesystem.\n\nEdit me and press Ctrl+S to save.\n`)
seed('config/app.yaml', `server:\n  port: 8080\n  debug: true\n\ndatabase:\n  host: postgres\n  port: 5432\n  name: myapp\n`)
seed('config/secret.env', `DATABASE_URL=postgres://user:pass@localhost/myapp\nSECRET_KEY=dev-only-secret\n`)
seed('data/sample.json', JSON.stringify({ version: 1, items: [{ id: 1, name: 'foo' }, { id: 2, name: 'bar' }] }, null, 2) + '\n')
seed('data/query.sql', `SELECT id, name, created_at\nFROM items\nWHERE active = true\nORDER BY created_at DESC\nLIMIT 100;\n`)
seed('logs/app.log', `2026-05-23T10:00:00Z INFO  Server started port=8080\n2026-05-23T10:01:00Z INFO  Request GET /health\n2026-05-23T10:02:00Z WARN  Slow query duration=320ms\n`)

function listDir(dir: string): { name: string; isDir: boolean; size: number; modTime: string }[] {
  const prefix = dir ? dir + '/' : ''
  const seen = new Set<string>()
  const result: { name: string; isDir: boolean; size: number; modTime: string }[] = []
  for (const [p] of mockFs.entries()) {
    if (!p.startsWith(prefix)) continue
    const rest = p.slice(prefix.length)
    if (!rest) continue
    const seg = rest.split('/')[0]
    if (seen.has(seg)) continue
    seen.add(seg)
    // Is seg itself a dir entry, or is it a prefix of deeper paths?
    const directKey = prefix + seg
    const entry = mockFs.get(directKey)
    result.push({
      name: seg,
      isDir: entry ? entry.isDir : rest.includes('/'),
      size: entry ? entry.size : 0,
      modTime: entry ? entry.modTime : new Date().toISOString(),
    })
  }
  return result
}

function handleAgentProxy(
  method: string,
  endpoint: string,
  query: URLSearchParams,
  body: string,
  res: ServerResponse<IncomingMessage>,
) {
  const path = query.get('path') ?? ''

  if (endpoint === 'config') {
    res.setHeader('Content-Type', 'application/json')
    res.end(JSON.stringify({ readonly: false, forceRW: false, pvc: 'mock-pvc', namespace: 'demo', pod: 'mock-pod', cluster: 'dev', version: 'dev-mock' }))
    return
  }

  if (endpoint === 'files') {
    if (method === 'GET') {
      res.setHeader('Content-Type', 'application/json')
      res.end(JSON.stringify({ entries: listDir(path) }))
      return
    }
    if (method === 'DELETE') {
      // Delete entry and all children
      for (const k of [...mockFs.keys()]) {
        if (k === path || k.startsWith(path + '/')) mockFs.delete(k)
      }
      res.writeHead(204)
      res.end()
      return
    }
  }

  if (endpoint === 'download' && method === 'GET') {
    const entry = mockFs.get(path)
    if (!entry || entry.isDir) { res.writeHead(404); res.end('not found'); return }
    res.setHeader('Content-Type', 'application/octet-stream')
    res.setHeader('Content-Disposition', `attachment; filename="${path.split('/').pop()}"`)
    res.end(entry.content)
    return
  }

  if (endpoint === 'edit' && (method === 'PUT' || method === 'POST')) {
    mockFs.set(path, { content: body, size: body.length, modTime: new Date().toISOString(), isDir: false })
    res.writeHead(204)
    res.end()
    return
  }

  if (endpoint === 'rename' && method === 'POST') {
    try {
      const { from, to } = JSON.parse(body)
      const entry = mockFs.get(from)
      if (entry) { mockFs.set(to, entry); mockFs.delete(from) }
    } catch { /* ignore */ }
    res.writeHead(204)
    res.end()
    return
  }

  if (endpoint === 'upload' && method === 'POST') {
    // Accept upload but don't parse multipart — just acknowledge
    res.writeHead(204)
    res.end()
    return
  }

  res.writeHead(404)
  res.end('not found')
}

export function mockWsPlugin(): Plugin {
  return {
    name: 'mock-ws',
    configureServer(server) {
      const wss = new WebSocketServer({ noServer: true })

      server.httpServer!.on('upgrade', (req, socket, head) => {
        if (req.url?.startsWith('/ws/v1/events')) {
          wss.handleUpgrade(req, socket as any, head, (ws) => {
            // Send snapshot immediately
            ws.send(JSON.stringify({
              id: 'mock-1',
              type: 'snapshot',
              serverTime: new Date().toISOString(),
              payload: {
                explorers: mockExplorers,
                scopes: mockScopes,
              },
            }))

            // Respond to pings
            ws.on('message', (data) => {
              try {
                const msg = JSON.parse(data.toString())
                if (msg.type === 'pong') return
              } catch { /* ignore */ }
            })

            // Keep alive with pings every 15s
            const ping = setInterval(() => {
              if (ws.readyState === ws.OPEN) {
                ws.send(JSON.stringify({ id: '', type: 'ping', serverTime: new Date().toISOString(), payload: {} }))
              }
            }, 15_000)
            ws.on('close', () => clearInterval(ping))
          })
        }
      })

      // Mock /api/v1/scopes
      server.middlewares.use('/api/v1/scopes', (req, res, next) => {
        const match = req.url?.match(/^\/([^/?]+)/)
        if (match) {
          const name = decodeURIComponent(match[1])
          const raw = mockScopes.find(s => (s.metadata as any).name === name)
          if (raw) {
            res.setHeader('Content-Type', 'application/json')
            res.end(JSON.stringify(raw))
            return
          }
        }
        if (!req.url || req.url === '/' || req.url?.startsWith('?')) {
          res.setHeader('Content-Type', 'application/json')
          res.end(JSON.stringify(mockScopes))
          return
        }
        next()
      })

      // Mock /api/version for Settings view
      server.middlewares.use('/api/version', (req, res, next) => {
        if (req.method === 'GET' && (!req.url || req.url === '/' || req.url.startsWith('?'))) {
          res.setHeader('Content-Type', 'text/plain; charset=utf-8')
          res.end('dev-mock')
          return
        }
        next()
      })

      // Mock REST endpoints needed by detail views
      server.middlewares.use('/api/v1/explorers', (req, res, next) => {
        const url = req.url ?? '/'

        // Agent proxy: /:ns/:name/proxy/api/:endpoint  or  /:ns/:name/heartbeat
        const proxyMatch = url.match(/^\/([^/]+)\/([^/]+)\/(proxy\/api\/([^?]*)|(heartbeat))(.*)?$/)
        if (proxyMatch) {
          const endpoint = proxyMatch[4] ?? ''
          const isHeartbeat = !!proxyMatch[5]
          const queryStr = proxyMatch[6] ?? ''
          const query = new URLSearchParams(queryStr.startsWith('?') ? queryStr.slice(1) : queryStr)
          const method = (req.method ?? 'GET').toUpperCase()

          if (isHeartbeat) {
            res.setHeader('Content-Type', 'application/json')
            res.end(JSON.stringify({ remainingSeconds: 600, phase: 'Running' }))
            return
          }

          // Collect body
          const chunks: Buffer[] = []
          req.on('data', (c: Buffer) => chunks.push(c))
          req.on('end', () => {
            const body = Buffer.concat(chunks).toString('utf-8')
            handleAgentProxy(method, endpoint, query, body, res)
          })
          return
        }

        // Individual explorer: GET /api/v1/explorers/:ns/:name
        const match = url.match(/^\/([^/]+)\/([^/]+)$/)
        if (match) {
          const [, ns, name] = match
          const raw = mockExplorers.find(
            e => (e.metadata as any).namespace === ns && (e.metadata as any).name === name
          )
          if (raw) {
            res.setHeader('Content-Type', 'application/json')
            res.end(JSON.stringify(raw))
            return
          }
        }
        // List: GET /api/v1/explorers
        if (!url || url === '/' || url?.startsWith('?')) {
          res.setHeader('Content-Type', 'application/json')
          res.end(JSON.stringify(mockExplorers))
          return
        }
        next()
      })
    },
  }
}
