[← Back to README](../README.md)

# Usage Reference

Command-line options, URL scheme, configuration, themes, operation modes, access logging, live reload and math rendering.

## URL Scheme

| Type | URL Format |
|------|-----------|
| Top page | `http://localhost:9080/` |
| Local file | `http://localhost:9080/local/path/to/file.md` |
| Local directory | `http://localhost:9080/local/path/to/dir/` |
| Remote (HTTP) | `http://localhost:9080/http/server/path/to/file.md` |
| Remote (HTTPS) | `http://localhost:9080/https/server/path/to/file.md` |
| GitHub repo | `http://localhost:9080/https/github.com/user/repo/blob/main/README.md` |

### Line Anchor Links

Markdown and text files support line-level linking using the `file:line` syntax:

- `[text](foo.md:12)` — rewritten to a link to `foo.md#L12` (line 12 of `foo.md`)
- `[text](foo.md:12-34)` — rewritten to `foo.md#L12-L34` (lines 12–34)
- `[text](foo.txt:12)` — rewritten to `foo.txt#L12` (line 12 of `foo.txt`)
- `<a href="foo.md:12">text</a>` — same, using raw HTML

When navigating to a line anchor (`#L12` or `#L12-L34`), the page scrolls to the target line and highlights the surrounding content. For text files, individual lines are highlighted; for Markdown files, the containing block element is highlighted. Highlighting is hidden in print output.

## Command Line

```bash
markdown-proxy [options] [file-or-url]
```

When `file-or-url` is provided:
- **Local file**: Opens the file via the proxy (relative paths are resolved automatically)
  - `markdown-proxy README.md`
  - `markdown-proxy ../docs/design.md`
  - `markdown-proxy /absolute/path/to/file.md`
- **Remote URL**: Opens the URL via the proxy
  - `markdown-proxy https://github.com/user/repo`
- If the server is not already running, it is started automatically in the background
- The file is opened in the default browser (`xdg-open` on Linux, `start` on Windows, `open` on macOS)

### Options

> **Note:** Options use a single dash (e.g., `-port`) following Go's `flag` package convention. Double dashes (`--port`) also work.

| Flag | Description | Default |
|------|-------------|---------|
| `-port`, `-p` | Listen port | `9080` |
| `-listen` | Bind address (`127.0.0.1` for local, `0.0.0.0` for remote) | `127.0.0.1` |
| `-theme` | Default CSS theme name (any file in the themes directory) | `github` |
| `-plantuml-server` | PlantUML server URL | (disabled) |
| `-auth-token` | Authentication token (required in remote mode) | |
| `-auth-cookie-max-age` | Authentication cookie max age in days | `30` |
| `-access-log` | Access log file path | |
| `-access-log-max-size` | Max log file size in MB before rotation | `100` |
| `-access-log-max-backups` | Max number of old log files to retain | `3` |
| `-access-log-max-age` | Max days to retain old log files | `28` |
| `-verbose`, `-v` | Enable debug logging to stderr | `false` |
| `-configure` | Interactively create configuration file | |
| `-version` | Show version and exit | |

## Configuration File

Settings can be saved to a configuration file so you don't need to pass flags every time.

### Interactive Setup

```bash
markdown-proxy --configure
```

This walks you through each setting interactively and saves the result.

### File Location

| Platform | Path |
|----------|------|
| Linux | `~/.config/markdown-proxy/config.json` |
| macOS | `~/Library/Application Support/markdown-proxy/config.json` |
| Windows | `%APPDATA%/markdown-proxy/config.json` |

### Supported Settings

The configuration file stores these settings as JSON:

```json
{
  "plantuml-server": "https://www.plantuml.com/plantuml",
  "theme": "dark",
  "port": 9080,
  "listen": "127.0.0.1"
}
```

Command-line flags override configuration file values. Security-sensitive settings (`-auth-token`, `-access-log`, etc.) are not stored in the config file.

## CSS Themes

markdown-proxy supports user-defined CSS themes in addition to the three built-in themes.

### Theme Directory

| Platform | Path |
|----------|------|
| Linux | `~/.config/markdown-proxy/themes/` |
| macOS | `~/Library/Application Support/markdown-proxy/themes/` |
| Windows | `%APPDATA%/markdown-proxy/themes/` |

On first launch, the three built-in themes are written here as editable CSS files:

- `github.css` — GitHub-style light theme
- `simple.css` — Serif light theme
- `dark.css` — Dark theme

### Adding a Custom Theme

1. Create a `.css` file in the themes directory, e.g. `~/.config/markdown-proxy/themes/mycolor.css`
2. Write standard CSS (no special scoping required — the file is loaded as the sole active stylesheet)
3. Restart markdown-proxy (or simply open a new page — the theme list is read at startup)
4. Select `mycolor` from the Theme dropdown

**Example:**

```css
body { font-family: "Source Serif Pro", serif; color: #1a1a1a; background: #fffff8; }
.toolbar { background: #f0ece0; border-color: #c8bfa0; }
.home-link, .toolbar-link { color: #8b4513; }
.markdown-body a { color: #8b4513; }
.markdown-body pre { background: #f5f0e8; border: 1px solid #c8bfa0; overflow: auto; }
.markdown-body pre code { background: none; padding: 0; font-size: 100%; }
.copy-btn { background: #fff; color: #444; border-color: #c8bfa0; }
.copy-btn:hover { background: #f0ece0; }
.copy-btn.copied { background: #d4edda; color: #155724; border-color: #c3e6cb; }
.toc-panel { background: #f5f0e8; border-left-color: #c8bfa0; }
.toc-header { border-bottom-color: #c8bfa0; }
.toc-list a.active { border-left-color: #8b4513; background: rgba(139,69,19,0.08); }
```

You can also edit the built-in themes directly — changes take effect on the next request (the file is read from disk on each theme CSS request).

## Operation Modes

### Local Mode (default)

```bash
markdown-proxy
```

The server binds to `127.0.0.1` and all features are available:
- Local file access (`/local/...`)
- Remote file access (`/http/...`, `/https/...`)
- Live reload via SSE (`/_sse`)
- Private network access is allowed for remote file fetching

### Remote Mode

```bash
markdown-proxy -listen 0.0.0.0 -auth-token my-secret-token
```

The server binds to the specified address for network access. For security:
- **Local file access is disabled**: `/local/` and `/_sse` return 403 Forbidden
- **Authentication is required**: `-auth-token` must be specified
- **Private network access is blocked**: SSRF protection prevents fetching from internal IPs
- **Top page** shows URL input only (no local file path input)

Users must authenticate via a login page (`/_login`) by entering the access token. The token is stored in an HttpOnly cookie for the configured duration.

## Access Logging

Access logs record each request in the following format:

```
2026-02-22T15:04:05+09:00 192.168.1.10 GET /https/github.com/user/repo 200 1234 150ms
```

- `-access-log /var/log/mdproxy/access.log`: Log to a file with automatic rotation
- In remote mode without `-access-log`: Logs to stdout by default
- In local mode without `-access-log`: No access logging

Log rotation is handled automatically using configurable size, count, and age limits.

## Live Reload

When viewing local Markdown files or directories (`/local/...`), the browser automatically reloads when the file or directory contents change. This uses Server-Sent Events (SSE) with filesystem notifications (fsnotify).

- **Local files only**: Remote files (`/http/...`, `/https/...`) are not affected
- **Local mode only**: Not available in remote mode
- **No configuration needed**: Works automatically for all local file views
- **Debounced**: Multiple rapid changes are coalesced into a single reload (100ms debounce)

## Math Rendering

Mathematical expressions are rendered using [KaTeX](https://katex.org/). Use standard LaTeX syntax:

- **Inline math**: `$E = mc^2$` renders inline within text
- **Display math**: `$$\int_0^\infty e^{-x} dx = 1$$` renders as a centered block

No configuration needed. Math expressions are automatically detected and rendered.
