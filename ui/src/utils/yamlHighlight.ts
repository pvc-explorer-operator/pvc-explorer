/**
 * Lightweight YAML syntax highlighter.
 * Returns HTML with <span> tags and CSS class annotations.
 */

type Token = { text: string; class?: string }

function tokenizeLine(line: string): Token[] {
  const tokens: Token[] = []
  let i = 0

  while (i < line.length) {
    // comment
    if (line[i] === '#') {
      tokens.push({ text: line.slice(i), class: 'sy-comment' })
      return tokens
    }

    // quoted string (single or double)
    if (line[i] === '"' || line[i] === "'") {
      const quote = line[i]
      const end = line.indexOf(quote, i + 1)
      if (end === -1) {
        tokens.push({ text: line.slice(i), class: 'sy-string' })
        return tokens
      }
      tokens.push({ text: line.slice(i, end + 1), class: 'sy-string' })
      i = end + 1
      continue
    }

    // key: value separator — color the key part
    if (i === 0 || line[i - 1] === ' ') {
      const colonIdx = line.indexOf(':', i)
      if (colonIdx > i) {
        const afterColon = colonIdx + 1
        // make sure it's a key: (followed by space, end, or comment)
        if (afterColon >= line.length || line[afterColon] === ' ' || line[afterColon] === '#' || line[afterColon] === '\t') {
          tokens.push({ text: line.slice(i, colonIdx), class: 'sy-key' })
          tokens.push({ text: ':' })
          i = colonIdx + 1
          continue
        }
      }
    }

    // list marker
    if (line.slice(i).match(/^-\s/)) {
      tokens.push({ text: '-', class: 'sy-list' })
      i += 1
      continue
    }

    // special values
    const wordEnd = line.slice(i).search(/[\s,#\]})"\']|$/)
    const word = line.slice(i, i + wordEnd)
    if (['true', 'false', 'True', 'False', 'TRUE', 'FALSE', 'null', 'Null', 'NULL', '~'].includes(word)) {
      tokens.push({ text: word, class: 'sy-bool' })
      i += wordEnd
      continue
    }
    if (/^\d+(\.\d+)?$/.test(word)) {
      tokens.push({ text: word, class: 'sy-number' })
      i += wordEnd
      continue
    }

    // default: regular text
    const nextSpecial = line.slice(i).search(/["'#]/)
    if (nextSpecial === -1) {
      tokens.push({ text: line.slice(i) })
      return tokens
    }
    tokens.push({ text: line.slice(i, i + nextSpecial) })
    i += nextSpecial
  }

  return tokens
}

export function highlightYaml(yaml: string): string {
  return yaml
    .split('\n')
    .map(line => {
      const tokens = tokenizeLine(line)
      const html = tokens
        .map(t => (t.class ? `<span class="${t.class}">${escapeHtml(t.text)}</span>` : escapeHtml(t.text)))
        .join('')
      return `<span class="sy-line">${html}</span>`
    })
    .join('\n')
}

function escapeHtml(s: string): string {
  return s
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
}
