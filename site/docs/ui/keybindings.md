# Keyboard Shortcuts

Keyboard shortcuts are available throughout the UI to speed up common workflows. All shortcuts are suppressed when focus is inside an input, textarea, or the Monaco code editor.

## Global

Available on every page.

| Key | Action |
| --- | ------ |
| `g` then `h` | Go to Home (explorer dashboard) |
| `g` then `s` | Go to Scopes |
| `r` | Refresh explorer list |
| `b` | Toggle sidebar |
| `d` | Toggle dark / light mode |
| `/` | Focus filter search *(dashboard only)* |
| `?` | Show keyboard shortcuts help |

> `g` chords have a 300 ms window — press the second key within 300 ms of `g`.

## Explorer detail

Available on `/explorers/:ns/:name`.

| Key | Action | Condition |
| --- | ------ | --------- |
| `f` | Browse files | Explorer is **Running** |
| `w` | Wake / Connect | Explorer is **ScaledToZero** |
| `x` | Disconnect / Sleep | Explorer is **Running** |
| `r` | Refresh detail | Always |

## File browser

Available on `/explorers/:ns/:name/files`.

| Key | Action | Condition |
| --- | ------ | --------- |
| `Ctrl` / `⌘` + `A` | Select all files | Always |
| `Delete` / `Backspace` | Delete selected | Items selected, not read-only |
| `Ctrl` / `⌘` + `S` | Save current file | Editor open *(built into Monaco)* |
| `Escape` | Close context menu or modal | Always |

## Built-in browser shortcuts

These are provided by the browser or OS and cannot be overridden.

| Key | Action |
| --- | ------ |
| `Tab` | Move focus between interactive elements |
| `Enter` / `Space` | Activate focused button or link |
